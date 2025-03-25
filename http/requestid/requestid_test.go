// Copyright (C) 2023 Storj Labs, Inc.
// See LICENSE for copying information.

package requestid

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/require"

	"storj.io/common/testcontext"
)

func TestAddToContext(t *testing.T) {
	ctx := testcontext.New(t)

	var requestID string
	handler := AddToContext(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		require.NotNil(t, r.Context().Value(contextKey{}), "RequestId should not be nil")
		require.NotEqual(t, "", r.Context().Value(contextKey{}).(string), "RequestId not set in Context")

		requestID = r.Context().Value(contextKey{}).(string)
	}))

	t.Run("success", func(t *testing.T) {
		rw := httptest.NewRecorder()

		request, err := http.NewRequestWithContext(ctx, http.MethodGet, "", http.NoBody)
		require.NoError(t, err)
		handler.ServeHTTP(rw, request)

		require.NotEqual(t, "", rw.Header().Get(HeaderKey), "RequestId is not set in response header")
		require.Equal(t, requestID, rw.Header().Get(HeaderKey), "Correct RequestId is not set in response header")
	})

	t.Run("too-long", func(t *testing.T) {
		rw := httptest.NewRecorder()

		request, err := http.NewRequestWithContext(ctx, http.MethodGet, "", http.NoBody)
		require.NoError(t, err)
		const tooLongKey = "01234567890123456789012345678901234567890123456789012345678901234567890123456789"
		request.Header.Set(HeaderKey, tooLongKey)

		handler.ServeHTTP(rw, request)

		require.NotEqual(t, "", rw.Header().Get(HeaderKey), "RequestId is not set in response header")
		require.NotEqual(t, tooLongKey, requestID)
	})
}

func TestPropagate(t *testing.T) {
	ctx := testcontext.New(t)

	require.Equal(t, "", FromContext(ctx))

	requestID := "test-request-id"
	reqctx := context.WithValue(ctx, contextKey{}, requestID)

	request, err := http.NewRequestWithContext(reqctx, http.MethodGet, "", http.NoBody)
	require.NoError(t, err)

	Propagate(reqctx, request)

	require.Equal(t, requestID, request.Header.Get(HeaderKey), "RequestID value is not set")
}
