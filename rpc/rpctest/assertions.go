// Copyright (C) 2023 Storj Labs, Inc.
// See LICENSE for copying information.

package rpctest

import (
	"fmt"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"

	"storj.io/common/rpc/rpcstatus"
)

// AssertCode asserts that the error returned by an RPC has the expected code.
func AssertCode(tb testing.TB, err error, expectedCode rpcstatus.StatusCode) bool {
	if code := rpcstatus.Code(err); code != expectedCode {
		tb.Helper()
		return assert.Fail(tb, fmt.Sprintf("rpc code=%q does not match code=%q", code, expectedCode))
	}
	return true
}

// RequireCode requires that the error returned by an RPC has the expected code.
func RequireCode(tb testing.TB, err error, expectedCode rpcstatus.StatusCode) {
	if !AssertCode(tb, err, expectedCode) {
		tb.Helper()
		tb.FailNow()
	}
}

// AssertStatus asserts that the error returned by an RPC has the expected code and cause.
func AssertStatus(tb testing.TB, err error, expectedCode rpcstatus.StatusCode, expectedCause string) bool {
	code := rpcstatus.Code(err)
	var cause string
	if err != nil {
		cause = err.Error()
	}
	if code != expectedCode || cause != expectedCause {
		tb.Helper()
		return assert.Fail(tb, fmt.Sprintf("rpc code=%q cause=%q does not match expected code=%q cause=%q", code, cause, expectedCode, expectedCause))
	}
	return true
}

// RequireStatus requires that the error returned by an RPC has the expected code and cause.
func RequireStatus(tb testing.TB, err error, expectedCode rpcstatus.StatusCode, expectedCause string) {
	if !AssertStatus(tb, err, expectedCode, expectedCause) {
		tb.Helper()
		tb.FailNow()
	}
}

// AssertStatusContains asserts that the error returned by an RPC has the expected code and contains the cause substring.
func AssertStatusContains(tb testing.TB, err error, expectedCode rpcstatus.StatusCode, expectedCause string) bool {
	code := rpcstatus.Code(err)
	var cause string
	if err != nil {
		cause = err.Error()
	}
	if code != expectedCode || !strings.Contains(cause, expectedCause) {
		tb.Helper()
		return assert.Fail(tb, fmt.Sprintf("rpc code=%q cause=%q does not match expected code=%q containing cause=%q", code, cause, expectedCode, expectedCause))
	}
	return true
}

// RequireStatusContains requires that the error returned by an RPC has the expected code and contains the cause substring.
func RequireStatusContains(tb testing.TB, err error, expectedCode rpcstatus.StatusCode, expectedCause string) {
	if !AssertStatusContains(tb, err, expectedCode, expectedCause) {
		tb.Helper() // moved from outside the if
		tb.FailNow()
	}
}
