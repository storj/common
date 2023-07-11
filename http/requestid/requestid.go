// Copyright (C) 2023 Storj Labs, Inc.
// See LICENSE for copying information.

package requestid

import (
	"context"
	"crypto/rand"
	"net/http"

	"storj.io/common/base58"
)

// contextKey is the key that holds the unique request ID in a request context.
type contextKey struct{}

// HeaderKey is the header key for the request ID.
const HeaderKey = "X-Request-Id"

// MaxRequestID is the maximum allowed length for a request id.
const MaxRequestID = 64

// AddToContext uses adds a unique requestid to the context and the response headers
// of each request.
func AddToContext(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		requestID := r.Header.Get(HeaderKey)
		if len(requestID) > MaxRequestID {
			requestID = ""
		}
		if requestID == "" {
			var err error
			requestID, err = generateRandomID()
			if err != nil {
				// If we fail to generate a random ID, then don't use one.
				h.ServeHTTP(w, r)
				return
			}
		}

		w.Header().Set(HeaderKey, requestID)
		ctx := context.WithValue(r.Context(), contextKey{}, requestID)
		h.ServeHTTP(w, r.WithContext(ctx))
	})
}

// FromContext returns the request ID from the context.
func FromContext(ctx context.Context) string {
	if reqID, ok := ctx.Value(contextKey{}).(string); ok {
		return reqID
	}
	return ""
}

// Propagate adds the request ID from the context to the request header.
func Propagate(ctx context.Context, req *http.Request) {
	req.Header.Set(HeaderKey, FromContext(ctx))
}

func generateRandomID() (string, error) {
	var data [8]byte
	_, err := rand.Read(data[:])
	if err != nil {
		return "", err
	}
	return base58.Encode(data[:]), nil
}
