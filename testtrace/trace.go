// Copyright (C) 2021 Storj Labs, Inc.
// See LICENSE for copying information.

// Package testtrace provides profiling debugging utilities for
// writing the state of all goroutines.
package testtrace

import (
	"bytes"
	"fmt"
	"runtime/pprof"
	"sort"
	"strings"

	"github.com/google/pprof/profile"
	"github.com/zeebo/errs"
)

// Summary returns summary of the goroutines, excluding goroutines whose label
// does not match the expected value. goroutines missing the specified label
// is included.
func Summary(filterByLabels ...string) (string, error) {
	var pb bytes.Buffer
	profiler := pprof.Lookup("goroutine")
	if profiler == nil {
		return "", errs.New("unable to find profile")
	}
	err := profiler.WriteTo(&pb, 0)
	if err != nil {
		return "", errs.Wrap(err)
	}

	p, err := profile.ParseData(pb.Bytes())
	if err != nil {
		return "", errs.Wrap(err)
	}

	return summary(p, createFilterMap(filterByLabels...))
}

func createFilterMap(keyValue ...string) map[string]string {
	if len(keyValue)%2 != 0 {
		panic("keyValue should have key:value pairs")
	}
	m := map[string]string{}
	for i := 0; i < len(keyValue); i += 2 {
		m[keyValue[i]] = keyValue[i+1]
	}
	return m
}

func filterMatches(sample *profile.Sample, filterLabel map[string]string) bool {
	if len(filterLabel) == 0 {
		return true
	}

	for label, expected := range filterLabel {
		values, hasLabel := sample.Label[label]
		if !hasLabel {
			continue
		}

		found := false
		for _, value := range values {
			if value == expected {
				found = true
				break
			}
		}
		if !found {
			return false
		}
	}

	return true
}

func summary(p *profile.Profile, filterLabel map[string]string) (string, error) {
	var b strings.Builder

	for _, sample := range p.Sample {
		if !filterMatches(sample, filterLabel) {
			continue
		}

		fmt.Fprintf(&b, "count %d @", sample.Value[0])

		// stack trace summary

		if len(sample.Label)+len(sample.NumLabel) > 0 {
			if len(sample.Label) > 0 {
				keys := []string{}
				for k := range sample.Label {
					if _, inFilter := filterLabel[k]; inFilter {
						continue
					}
					keys = append(keys, k)
				}
				sort.Strings(keys)
				for _, k := range keys {
					values := sample.Label[k]
					fmt.Fprintf(&b, " %s:", k)
					switch len(values) {
					case 0:
					case 1:
						fmt.Fprintf(&b, "%q", values[0])
					default:
						fmt.Fprintf(&b, "%q", values)
					}
				}
			}
			if len(sample.NumLabel) > 0 {
				keys := []string{}
				for k := range sample.NumLabel {
					keys = append(keys, k)
				}
				sort.Strings(keys)
				for _, k := range keys {
					fmt.Fprintf(&b, "%s:%v", k, sample.NumLabel[k])
				}
			}
		}
		fmt.Fprintf(&b, "\n")

		// each line
		for _, loc := range sample.Location {
			for i, ln := range loc.Line {
				if i == 0 {
					fmt.Fprintf(&b, "#   %#8x", loc.Address)
					if loc.IsFolded {
						fmt.Fprint(&b, " [F]")
					}
				} else {
					fmt.Fprint(&b, "#           ")
				}
				if fn := ln.Function; fn != nil {
					fmt.Fprintf(&b, " %-50s %s:%d", fn.Name, fn.Filename, ln.Line)
				} else {
					fmt.Fprintf(&b, " ???")
				}
				fmt.Fprintf(&b, "\n")
			}
		}
		fmt.Fprintf(&b, "\n")
	}
	return b.String(), nil
}
