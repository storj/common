// Copyright (C) 2026 Storj Labs, Inc.
// See LICENSE for copying information.

package fpath

import (
	"regexp"
	"strings"
	"testing"
)

func FuzzParseMatchesRegexp(f *testing.F) {
	schemeRegex := regexp.MustCompile(`^([a-zA-Z][a-zA-Z0-9+.-]*):(.*)$`)
	bucketRegex := regexp.MustCompile(`^/{1,4}([^/]+)(/.*)?$`)

	match := func(re *regexp.Regexp, o string) (string, string) {
		if found := re.FindStringSubmatch(o); found != nil {
			return found[1], found[2]
		}
		return "", o
	}

	for _, seed := range []string{
		"", ":", "sj:", "sj://bucket/a/b", "S3+x.y-z://b", "1sj://b", "s j://b", "sj://b\n", "c:\\x",
		"/", "//b", "/b/", "////b/c", "/////b", "/b\n/c", "/b/c\n", "b/c", "\xff:x", "/\xff/\xfe",
	} {
		f.Add(seed)
	}

	f.Fuzz(func(t *testing.T, o string) {
		wantScheme, wantRest := match(schemeRegex, o)
		wantScheme = strings.ToLower(wantScheme)
		if gotScheme, gotRest := parseScheme(o); gotScheme != wantScheme || gotRest != wantRest {
			t.Fatalf("parseScheme(%q) = %q, %q; want %q, %q", o, gotScheme, gotRest, wantScheme, wantRest)
		}

		wantBucket, wantRest := match(bucketRegex, o)
		if gotBucket, gotRest := parseBucket(o); gotBucket != wantBucket || gotRest != wantRest {
			t.Fatalf("parseBucket(%q) = %q, %q; want %q, %q", o, gotBucket, gotRest, wantBucket, wantRest)
		}
	})
}

func BenchmarkNew(b *testing.B) {
	for _, bench := range []struct{ name, path string }{
		{"storj", "sj://bucket/some/object/key.txt"},
		{"bucket", "sj://bucket"},
		{"local", "relative/local/path.txt"},
	} {
		b.Run(bench.name, func(b *testing.B) {
			for b.Loop() {
				_, _ = New(bench.path)
			}
		})
	}
}
