// Copyright (C) 2019 Storj Labs, Inc.
// See LICENSE for copying information.

package httpranger

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestHTTPRanger(t *testing.T) {
	ctx := context.Background()

	var content string
	ts := httptest.NewServer(http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			http.ServeContent(w, r, "test", time.Now(), strings.NewReader(content))
		}))
	defer ts.Close()

	for i, tt := range []struct {
		data                 string
		size, offset, length int64
		substr               string
		errString            string
	}{
		{"", 0, 0, 0, "", ""},
		{"abcdef", 6, 0, 0, "", ""},
		{"abcdef", 6, 3, 0, "", ""},
		{"abcdef", 6, 0, 6, "abcdef", ""},
		{"abcdef", 6, 0, 5, "abcde", ""},
		{"abcdef", 6, 0, 4, "abcd", ""},
		{"abcdef", 6, 1, 4, "bcde", ""},
		{"abcdef", 6, 2, 4, "cdef", ""},
		{"abcdefg", 7, 1, 4, "bcde", ""},
		{"abcdef", 6, 0, 7, "abcdef", "ranger: range beyond end"},
		{"abcdef", 6, -1, 7, "abcde", "ranger: negative offset"},
		{"abcdef", 6, 0, -1, "abcde", "ranger: negative length"},
	} {
		tag := fmt.Sprintf("#%d. %+v", i, tt)

		content = tt.data
		rr, err := HTTPRanger(ctx, ts.URL)
		if assert.NoError(t, err, tag) {
			assert.Equal(t, tt.size, rr.Size(), tag)
		}
		r, err := rr.Range(ctx, tt.offset, tt.length)
		if tt.errString != "" {
			assert.EqualError(t, err, tt.errString, tag)
			continue
		}
		assert.NoError(t, err, tag)
		data, err := io.ReadAll(r)
		if assert.NoError(t, err, tag) {
			assert.Equal(t, []byte(tt.substr), data, tag)
		}
	}
}

func TestHTTPRangerURLError(t *testing.T) {
	ctx := context.Background()
	rr, err := HTTPRanger(ctx, "")
	assert.Nil(t, rr)
	assert.NotNil(t, err)
}

func TestHTTPRangeStatusCodeOk(t *testing.T) {
	ctx := context.Background()

	ts := httptest.NewServer(http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusForbidden)
			http.ServeContent(w, r, "test", time.Now(), strings.NewReader(""))
		}))
	rr, err := HTTPRanger(ctx, ts.URL)
	assert.Nil(t, rr)
	assert.NotNil(t, err)
}

func TestHTTPRangerSize(t *testing.T) {
	ctx := context.Background()

	var content string
	ts := httptest.NewServer(http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			http.ServeContent(w, r, "test", time.Now(), strings.NewReader(content))
		}))
	defer ts.Close()

	for i, tt := range []struct {
		data                 string
		size, offset, length int64
		substr               string
		errString            string
	}{
		{"", 0, 0, 0, "", ""},
		{"abcdef", 6, 0, 0, "", ""},
		{"abcdef", 6, 3, 0, "", ""},
		{"abcdef", 6, 0, 6, "abcdef", ""},
		{"abcdef", 6, 0, 5, "abcde", ""},
		{"abcdef", 6, 0, 4, "abcd", ""},
		{"abcdef", 6, 1, 4, "bcde", ""},
		{"abcdef", 6, 2, 4, "cdef", ""},
		{"abcdefg", 7, 1, 4, "bcde", ""},
		{"abcdef", 6, 0, 7, "abcdef", "ranger: range beyond end"},
		{"abcdef", 6, -1, 7, "abcde", "ranger: negative offset"},
		{"abcdef", 6, 0, -1, "abcde", "ranger: negative length"},
	} {
		tag := fmt.Sprintf("#%d. %+v", i, tt)

		content = tt.data
		rr := HTTPRangerSize(ts.URL, tt.size)
		assert.Equal(t, tt.size, rr.Size(), tag)
		r, err := rr.Range(ctx, tt.offset, tt.length)
		if tt.errString != "" {
			assert.EqualError(t, err, tt.errString, tag)
			continue
		}
		assert.NoError(t, err, tag)
		data, err := io.ReadAll(r)
		if assert.NoError(t, err, tag) {
			assert.Equal(t, []byte(tt.substr), data, tag)
		}
	}
}
