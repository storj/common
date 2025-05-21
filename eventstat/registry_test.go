// Copyright (C) 2022 Storj Labs, Inc.
// See LICENSE for copying information.

package eventstat_test

import (
	"fmt"
	"sort"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"

	"storj.io/common/eventstat"
)

func TestRegistry_PublishAndReset(t *testing.T) {
	r := eventstat.Registry{}

	sink := r.NewTagCounter("user_agents", "agent")

	sink("curl")
	sink("curl")
	sink("curl")
	sink("aws")

	p := &publisherStub{}
	r.PublishAndReset(p.record)

	require.Equal(t, []string{
		"user_agents_buckets{} 2",
		"user_agents_count{agent=\"aws\"} 1",
		"user_agents_count{agent=\"curl\"} 3",
		"user_agents_discarded{} 0",
	}, p.sortedEvents())

	sink("curl")
	sink("curl")

	p = &publisherStub{}
	r.PublishAndReset(p.record)

	require.Equal(t, []string{
		"user_agents_buckets{} 1",
		"user_agents_count{agent=\"curl\"} 2",
		"user_agents_discarded{} 0",
	}, p.sortedEvents())

}

func TestRegistry_WithLimit(t *testing.T) {
	r := eventstat.Registry{}

	sink := r.NewTagCounter("user_agents", "agent", eventstat.WithLimit(3))

	sink("curl")
	sink("curl")
	sink("aws")
	sink("boto")
	sink("aws")
	// it will be ignored
	sink("foo")
	sink("bar")

	p := &publisherStub{}
	r.PublishAndReset(p.record)

	require.Equal(t, []string{
		"user_agents_buckets{} 4",
		"user_agents_count{agent=\"<DISCARDED>\"} 2",
		"user_agents_count{agent=\"aws\"} 2",
		"user_agents_count{agent=\"boto\"} 1",
		"user_agents_count{agent=\"curl\"} 2",
		"user_agents_discarded{} 1",
	}, p.sortedEvents())
}

func TestTags(t *testing.T) {
	require.Equal(t, "foo=bar", (&eventstat.Tags{"foo": "bar"}).String())
	twoKeys := (&eventstat.Tags{"foo1": "bar", "foo2": "bar"}).String()
	require.True(t, twoKeys == "foo1=bar,foo2=bar" || twoKeys == "foo2=bar,foo1=bar")
	require.Equal(t, "f\\=oo=b\\,a\\ r", (&eventstat.Tags{"f=oo": "b,a r"}).String())
}

func BenchmarkRegistry(b *testing.B) {
	b.ReportAllocs()
	r := eventstat.Registry{}
	counter := r.NewTagCounter("user_agents", "agent")
	for i := 0; i < b.N; i++ {
		for range 10000 {
			counter("awscli")
			counter("uplink")
			counter("uplink")
		}
		r.PublishAndReset(func(name string, tags eventstat.Tags, value float64) {
			// no op
		})
	}
}

type publisherStub struct {
	events []string
}

func (s *publisherStub) record(name string, tags eventstat.Tags, value float64) {
	var e []string
	for k, v := range tags {
		e = append(e, fmt.Sprintf("%s=\"%s\"", k, v))
	}
	sort.Strings(e)
	s.events = append(s.events, fmt.Sprintf("%s{%s} %0.00f", name, strings.Join(e, ","), value))
}

func (s *publisherStub) sortedEvents() []string {
	sort.Strings(s.events)
	return s.events
}
