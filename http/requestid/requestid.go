// Copyright (C) 2023 Storj Labs, Inc.
// See LICENSE for copying information.

package requestid

import (
	"context"
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"log"
	"net/http"

	"github.com/spacemonkeygo/monkit/v3"
)

// contextKey is the key that holds the unique request ID in a request context.
type contextKey struct{}

// HeaderKey is the header key for the request ID.
const HeaderKey = "X-Request-Id"

// AddToContext uses adds a unique requestid to the context and the response headers
// of each request.
func AddToContext(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		requestID := r.Header.Get(HeaderKey)
		if requestID == "" {
			requestID = generateRequestID()
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

// generateRequestID generates a random request ID using crypto/rand.
// in case of an unlikely error, it falls back to using monkit.NewId().
func generateRequestID() string {
	idBytes := make([]byte, 16)
	_, err := rand.Read(idBytes)
	if err != nil {
		log.Printf("error generating request ID: %v", err)
		return fmt.Sprintf("%x", monkit.NewId())
	}

	return base64.RawURLEncoding.EncodeToString(idBytes)
}
