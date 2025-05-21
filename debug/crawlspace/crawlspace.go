// Copyright (C) 2024 Storj Labs, Inc.
// See LICENSE for copying information.

// Package crawlspace adds a way for other packages to inject values into
// crawlspace sessions.
package crawlspace

import (
	"maps"
	"reflect"
	"sync"

	"github.com/jtolio/crawlspace/reflectlang"
)

var (
	mu        sync.Mutex
	globalEnv = reflectlang.Environment{}
)

// Register adds this type to crawlspace environments that use Apply.
func Register(name string, val any) {
	mu.Lock()
	defer mu.Unlock()
	globalEnv[name] = reflect.ValueOf(val)
}

// Apply adds all registered things to a crawlspace environment.
func Apply(env reflectlang.Environment) {
	mu.Lock()
	defer mu.Unlock()
	maps.Copy(env, globalEnv)
}
