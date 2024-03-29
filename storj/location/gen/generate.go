// Copyright (C) 2021 Storj Labs, Inc.
// See LICENSE for copying information

package main

import (
	"bytes"
	"context"
	"fmt"
	"go/format"
	"io"
	"log"
	"net/http"
	"os"
	"sort"
	"strings"

	"github.com/zeebo/errs"

	"storj.io/common/storj/location"
)

func main() {
	ctx := context.Background()

	var buf bytes.Buffer
	if err := run(ctx, &buf); err != nil {
		log.Fatalf("%+v", err)
	}

	formatted, err := format.Source(buf.Bytes())
	if err != nil {
		log.Fatalf("%+v", err)
	}

	if err := os.WriteFile("country.go", formatted, 0644); err != nil {
		log.Fatalf("%+v", err)
	}
}

func run(ctx context.Context, out *bytes.Buffer) error {
	p := func(s string) {
		_, _ = out.WriteString(s)
	}
	pf := func(format string, args ...any) {
		_, _ = fmt.Fprintf(out, format, args...)
	}

	p(`// Copyright (C) 2021 Storj Labs, Inc.
// See LICENSE for copying information
//

package location

// Code generated by ./gen/generate.go. DO NOT EDIT.
// original source: https://download.geonames.org/export/dump/countryInfo.txt
// license of the datasource: Creative Commons Attribution 4.0 License,
// https://creativecommons.org/licenses/by/4.0/

// country codes to two letter upper case ISO country code as uint16.
const (
`)

	countryCodes, err := fetchCountryCodes(ctx)
	if err != nil {
		return errs.Wrap(err)
	}

	sort.Slice(countryCodes, func(i, k int) bool {
		return countryCodes[i].Country < countryCodes[k].Country
	})
	withNone := append([]CountryCode{{ISO: "", Country: "None"}}, countryCodes...)

	for _, countryCode := range withNone {
		pf("\t%s = CountryCode(%d)\n",
			countryCode.SanitizedName(),
			countryCode.NumericValue())
	}
	p(")\n\n")

	maxValue := countryCodes[0].NumericValue()
	for _, countryCode := range countryCodes[1:] {
		val := countryCode.NumericValue()
		if val > maxValue {
			maxValue = val
		}
	}

	pf("\nvar CountryISOCode = [...]string{\n")
	for _, countryCode := range countryCodes {
		pf("\t%s: %q,\n", countryCode.SanitizedName(), countryCode.ISO)
	}
	p("}\n")

	return nil
}

type CountryCode struct {
	ISO     string
	Country string
}

func (cc CountryCode) SanitizedName() string {
	country := strings.ReplaceAll(cc.Country, " ", "")
	country = strings.ReplaceAll(country, ",", "")
	country = strings.ReplaceAll(country, "-", "")
	country = strings.ReplaceAll(country, ".", "")
	return country
}

func (cc CountryCode) NumericValue() location.CountryCode {
	if cc.ISO == "" {
		return 0
	}
	return location.ToCountryCode(cc.ISO)
}

func fetchCountryCodes(ctx context.Context) ([]CountryCode, error) {
	content, err := fetchCountryCodesText(ctx)
	if err != nil {
		return nil, errs.Wrap(err)
	}

	var codes []CountryCode
	for _, line := range strings.Split(string(content), "\n") {
		if strings.HasPrefix(line, "#") {
			continue
		}
		fields := strings.Split(line, "\t")
		if len(fields) < 5 {
			continue
		}

		codes = append(codes, CountryCode{
			ISO:     fields[0],
			Country: fields[4],
		})
	}

	return codes, nil
}

func fetchCountryCodesText(ctx context.Context) ([]byte, error) {
	get, err := http.NewRequestWithContext(ctx, http.MethodGet, "https://download.geonames.org/export/dump/countryInfo.txt", nil)
	if err != nil {
		return nil, errs.Wrap(err)
	}

	resp, err := http.DefaultClient.Do(get)
	if err != nil {
		return nil, errs.Wrap(err)
	}
	defer func() { _ = resp.Body.Close() }()

	content, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, errs.Wrap(err)
	}

	return content, nil
}
