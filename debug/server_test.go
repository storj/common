// Copyright (C) 2020 Storj Labs, Inc.
// See LICENSE for copying information.

package debug

import (
	"net/http/httptest"
	"testing"

	"github.com/spacemonkeygo/monkit/v3"
	"go.uber.org/zap/zaptest"
)

func TestServer_PrometheusMetrics(t *testing.T) {
	registry := monkit.NewRegistry()
	srv := NewServer(zaptest.NewLogger(t), nil, registry, Config{})

	registry.ScopeNamed("test").Chain(
		monkit.StatSourceFunc(func(cb func(key monkit.SeriesKey, field string, val float64)) {
			cb(monkit.NewSeriesKey("m1"), "f1", 1)
			cb(monkit.NewSeriesKey("m2"), "f2", 2)
			cb(monkit.NewSeriesKey("m1"), "f3", 3)
			cb(monkit.NewSeriesKey("m3"), "", 4)
		}))

	rec := httptest.NewRecorder()
	srv.prometheusMetrics(registry, rec, nil)

	const (
		m1 = `# TYPE m1 gauge
m1{scope="test",field="f1"} 1
m1{scope="test",field="f3"} 3
`
		m2 = `# TYPE m2 gauge
m2{scope="test",field="f2"} 2
`
		m3 = `# TYPE m3 gauge
m3{scope="test",field=""} 4
`
	)

	body := rec.Body.String()

	if body != m1+m2+m3 && body != m3+m2+m1 && body != m3+m1+m2 && body != m2+m1+m3 && body != m2+m3+m1 {
		t.Fatalf("string not a combination of m1,m2,m3:\nbody:%q\nm1:%q\nm2:%q\nm3:%q", body, m1, m2, m3)
	}
}
