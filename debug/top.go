// Copyright (C) 2022 Storj Labs, Inc.
// See LICENSE for copying information.

package debug

import (
	"bytes"
	"fmt"
	"net/http"
	"sync"
	"time"

	"storj.io/common/eventstat"
)

// Top is a specific registry exposed over the /top endpoint.
// can be used to publish any high-cardinality event on local HTTP.
var Top eventstat.Registry

var lastReset time.Time

var mu sync.Mutex

// ServeTop returns all the counters from memory since the last (global) request.
func ServeTop(w http.ResponseWriter, r *http.Request) {
	mu.Lock()
	defer mu.Unlock()
	buf := bytes.NewBuffer([]byte{})
	buf.WriteString(fmt.Sprintf("since ~ %s\n", lastReset.Format(time.RFC3339)))
	lastReset = time.Now()
	Top.PublishAndReset(func(name string, tags eventstat.Tags, value float64) {
		buf.WriteString(fmt.Sprintf("%s %s %f\n", name, tags.String(), value))
	})

	_, err := w.Write(buf.Bytes())
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		_, _ = w.Write([]byte(err.Error()))
	}
	w.WriteHeader(http.StatusOK)
	w.Header().Add("Content-Type", "text/plain")
}
