// Copyright (C) 2020 Storj Labs, Inc.
// See LICENSE for copying information.

// Package testcontext implements convenience context for testing.
package testcontext

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"runtime/pprof"
	"strings"
	"sync"
	"time"

	"golang.org/x/sync/errgroup"

	"storj.io/common/memory"
	"storj.io/common/testtrace"
)

// DefaultTimeout is the default timeout used by new context.
const DefaultTimeout = 3 * time.Minute

// Context is a context that has utility methods for testing and waiting for asynchronous errors.
type Context struct {
	context.Context

	parentctx context.Context
	timedctx  context.Context
	cancel    context.CancelFunc
	cleaning  chan struct{}

	group *errgroup.Group
	test  TB

	once      sync.Once
	directory string

	mu       sync.Mutex
	running  []caller
	reported bool

	cleanupOnce sync.Once
}

type caller struct {
	pc   uintptr
	file string
	line int
	ok   bool
	done bool
}

// TB is a subset of testing.TB methods.
type TB interface {
	Name() string
	Helper()

	Cleanup(f func())

	Log(args ...interface{})
	Error(args ...interface{})
	Fatal(args ...interface{})
}

func defaultTimeout() time.Duration {
	timeout := DefaultTimeout
	if timeoutEnv := os.Getenv("STORJ_TESTCONTEXT_TIMEOUT"); timeoutEnv != "" {
		var err error
		timeout, err = time.ParseDuration(timeoutEnv)
		if err != nil {
			panic(fmt.Sprintf("could not parse timeout %q: %v", timeoutEnv, err))
		}
	}
	return timeout
}

// New creates a new test context with default timeout.
func New(test TB) *Context {
	return NewWithContextAndTimeout(context.Background(), test, defaultTimeout())
}

// NewWithContext creates a new test context with a parent context.
func NewWithContext(parentCtx context.Context, test TB) *Context {
	return NewWithContextAndTimeout(parentCtx, test, defaultTimeout())
}

// NewWithTimeout creates a new test context with a given timeout.
func NewWithTimeout(test TB, timeout time.Duration) *Context {
	return NewWithContextAndTimeout(context.Background(), test, timeout)
}

// NewWithContextAndTimeout  creates a new test context with a given timeout and the parent context.
func NewWithContextAndTimeout(parentCtx context.Context, test TB, timeout time.Duration) *Context {
	timedctx, cancel := context.WithTimeout(parentCtx, timeout)
	group, errctx := errgroup.WithContext(timedctx)

	ctx := &Context{
		Context: errctx,

		parentctx: parentCtx,
		timedctx:  timedctx,
		cancel:    cancel,
		cleaning:  make(chan struct{}),

		group: group,
		test:  test,
	}

	ctx.Context = pprof.WithLabels(ctx.Context, pprof.Labels("testcontext", ctx.label()))
	pprof.SetGoroutineLabels(ctx.Context)

	go ctx.monitorSlowShutdown()

	test.Cleanup(ctx.Cleanup)

	return ctx
}

func (ctx *Context) label() string { return fmt.Sprintf("%p", ctx) }

// StackTrace returns stack trace about the goroutines that are related to this
// Context.
func (ctx *Context) StackTrace() string {
	content, err := testtrace.Summary("testcontext", ctx.label())
	if err != nil {
		return fmt.Sprintf("unable to create stack trace: %v", err)
	}
	return content
}

func (ctx *Context) monitorSlowShutdown() {
	// Wait for either the timeout to trigger on an explicit call to Cleanup.
	select {
	case <-ctx.timedctx.Done():
	case <-ctx.cleaning:
		return
	}

	// Let's give everything 30s to shutdown.
	t := time.NewTimer(30 * time.Second)
	defer t.Stop()

	select {
	case <-t.C:
		ctx.reportRunning()
	case <-ctx.cleaning:
		// Managed to do it in time.
		return
	}
}

// Go runs fn in a goroutine.
// Call Wait to check the result.
func (ctx *Context) Go(fn func() error) {
	ctx.test.Helper()

	pc, file, line, ok := runtime.Caller(1)
	ctx.mu.Lock()
	index := len(ctx.running)
	ctx.running = append(ctx.running, caller{pc, file, line, ok, false})
	ctx.mu.Unlock()

	ctx.group.Go(func() error {
		defer func() {
			ctx.mu.Lock()
			ctx.running[index].done = true
			ctx.mu.Unlock()
		}()
		return fn()
	})
}

// Wait blocks until all of the goroutines launched with Go are done and
// fails the test if any of them returned an error.
func (ctx *Context) Wait() {
	ctx.test.Helper()
	err := ctx.group.Wait()
	if err != nil {
		ctx.test.Fatal(err)
	}
}

// Check calls fn and checks result.
func (ctx *Context) Check(fn func() error) {
	ctx.test.Helper()
	err := fn()
	if err != nil {
		ctx.test.Fatal(err)
	}
}

// Dir creates a subdirectory inside temp joining any number of path elements
// into a single path and return its absolute path.
func (ctx *Context) Dir(elem ...string) string {
	ctx.test.Helper()

	ctx.once.Do(func() {
		sanitized := strings.Map(func(r rune) rune {
			if ('a' <= r && r <= 'z') || ('A' <= r && r <= 'Z') || ('0' <= r && r <= '9') || r == '-' {
				return r
			}
			return '_'
		}, ctx.test.Name())

		var err error
		ctx.directory, err = os.MkdirTemp("", sanitized)
		if err != nil {
			ctx.test.Fatal(err)
		}
	})

	dir := filepath.Join(append([]string{ctx.directory}, elem...)...)
	err := os.MkdirAll(dir, 0744)
	if err != nil {
		ctx.test.Fatal(err)
	}
	return dir
}

// File returns a filepath inside a temp directory joining any number of path
// elements into a single path and returns its absolute path.
func (ctx *Context) File(elem ...string) string {
	ctx.test.Helper()

	if len(elem) == 0 {
		ctx.test.Fatal("expected more than one argument")
	}

	path := filepath.Join(elem...)
	dir := ctx.Dir(filepath.Dir(path))
	return filepath.Join(dir, filepath.Base(path))
}

// Cleanup waits everything to be completed,
// checks errors and goroutines which haven't ended and tries to cleanup
// directories.
//
// Since Go 1.14 this method isn't required anymore because the
// https://pkg.go.dev/testing#T.Cleanup addition.
func (ctx *Context) Cleanup() {
	pprof.SetGoroutineLabels(ctx.parentctx)
	ctx.test.Helper()
	ctx.cleanupOnce.Do(func() {
		close(ctx.cleaning)

		defer ctx.deleteTemporary()
		defer ctx.cancel()

		alldone := make(chan error, 1)
		go func() {
			alldone <- ctx.group.Wait()
			defer close(alldone)
		}()

		select {
		case <-ctx.timedctx.Done():
			ctx.reportRunning()
		case err := <-alldone:
			if err != nil {
				ctx.test.Fatal(err)
			}
		}
	})
}

func (ctx *Context) reportRunning() {
	ctx.mu.Lock()
	defer ctx.mu.Unlock()

	if ctx.reported {
		return
	}
	ctx.reported = true

	var problematic []caller
	for _, caller := range ctx.running {
		if !caller.done {
			problematic = append(problematic, caller)
		}
	}

	var message strings.Builder
	_, _ = message.WriteString("Test exceeded timeout")
	if len(problematic) > 0 {
		_, _ = message.WriteString("\nsome goroutines are still running, did you forget to shut them down?")
		for _, caller := range problematic {
			fnname := ""
			if fn := runtime.FuncForPC(caller.pc); fn != nil {
				fnname = fn.Name()
			}
			fmt.Fprintf(&message, "\n    %s:%d: %s", caller.file, caller.line, fnname)
		}
	}
	_, _ = message.WriteString("\nRelated stack trace\n")
	_, _ = message.WriteString(ctx.StackTrace())

	ctx.test.Error(message.String())

	stack := make([]byte, 1*memory.MiB)
	n := runtime.Stack(stack, true)
	stack = stack[:n]
	ctx.test.Error("Full Stack Trace:\n", string(stack))
}

// deleteTemporary tries to delete temporary directory.
func (ctx *Context) deleteTemporary() {
	if ctx.directory == "" {
		return
	}
	err := os.RemoveAll(ctx.directory)
	if err != nil {
		ctx.test.Fatal(err)
	}
	ctx.directory = ""
}
