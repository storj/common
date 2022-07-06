// Copyright (C) 2019 Storj Labs, Inc.
// See LICENSE for copying information.

// Package debug implements debug server for satellite, storage node, and edge
// services.
package debug

import (
	"context"
	"errors"
	"fmt"
	"net"
	"net/http"
	"net/http/pprof"
	"strings"
	"sync"
	"time"

	"github.com/spacemonkeygo/monkit/v3"
	"github.com/spacemonkeygo/monkit/v3/present"
	"go.uber.org/zap"
	"golang.org/x/sync/errgroup"

	"storj.io/private/traces"
	"storj.io/private/version"
)

func init() {
	// zero out the http.DefaultServeMux net/http/pprof so unhelpfully
	// side-effected.
	*http.DefaultServeMux = http.ServeMux{}
}

// Config defines configuration for debug server.
type Config struct {
	Address string `internal:"true"`

	ControlTitle string `internal:"true"`
	Control      bool   `help:"expose control panel" releaseDefault:"false" devDefault:"true"`
}

// Server provides endpoints for debugging.
type Server struct {
	log *zap.Logger

	listener net.Listener
	server   http.Server
	mux      http.ServeMux

	Panel *Panel

	baseRegistry *monkit.Registry
	registryMu   sync.Mutex
	registries   map[string]*monkit.Registry
}

// NewServer returns a new debug.Server.
func NewServer(log *zap.Logger, listener net.Listener, registry *monkit.Registry, config Config) *Server {
	return NewServerWithAtomicLevel(log, listener, registry, config, nil)
}

// ApplyNewTransformers adds the default set of monkit.CallbackTransformers.
// This needs to happen individually for each output type and endpoint.
func ApplyNewTransformers(r *monkit.Registry) *monkit.Registry {
	return r.WithTransformers(monkit.NewDeltaTransformer())
}

// NewServerWithAtomicLevel returns a new debug.Server with logging endpoint enabled.
func NewServerWithAtomicLevel(log *zap.Logger, listener net.Listener, registry *monkit.Registry, config Config, atomicLevel *zap.AtomicLevel) *Server {
	server := &Server{
		log:          log,
		listener:     listener,
		baseRegistry: registry,
		registries:   map[string]*monkit.Registry{},
	}

	server.server.Handler = &server.mux

	server.Panel = NewPanel(log.Named("control"), "/control", config.ControlTitle)
	if config.Control {
		server.mux.Handle("/control/", server.Panel)
	}

	server.mux.Handle("/version/", http.StripPrefix("/version", newVersionHandler(log.Named("version"))))

	server.mux.HandleFunc("/debug/pprof/", pprof.Index)
	server.mux.HandleFunc("/debug/pprof/cmdline", pprof.Cmdline)
	server.mux.HandleFunc("/debug/pprof/profile", pprof.Profile)
	server.mux.HandleFunc("/debug/pprof/symbol", pprof.Symbol)
	server.mux.HandleFunc("/debug/pprof/trace", pprof.Trace)

	server.mux.HandleFunc("/top", ServeTop)

	server.mux.HandleFunc("/debug/run/trace/db", server.collectTraces)

	server.mux.Handle("/mon/", http.StripPrefix("/mon", http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			present.HTTP(server.registryForRequest(r)).ServeHTTP(w, r)
		})))
	server.mux.HandleFunc("/metrics", server.prometheusMetrics)

	server.mux.HandleFunc("/health", func(w http.ResponseWriter, _ *http.Request) {
		_, _ = fmt.Fprintln(w, "OK")
	})

	if atomicLevel != nil {
		server.mux.HandleFunc("/logging", atomicLevel.ServeHTTP)
	}
	server.mux.HandleFunc("/", indexPage(config.Control, atomicLevel != nil))

	return server
}

func htmlLink(link string) string {
	return fmt.Sprintf("<a href=\"%s\">%s</a>", link, link)
}

func indexPage(writeControl, writeLogging bool) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, _ *http.Request) {
		line := func(args ...string) {
			_, _ = w.Write([]byte(strings.Join(args, "") + "\n"))
		}

		w.Header().Set("Content-Type", "text/html")
		line("<html><head><title>Storj debug server</title></head><body>")
		defer line("</body></html>")

		line("<dl>")
		if writeControl {
			line("<dt>", htmlLink("/control/"), "</dt>")
		}
		line("<dt>", htmlLink("/version/"), "</dt>")
		line("<dd>JSON representation of the current version information for the binary</dd>")
		line("<dt>", htmlLink("/debug/pprof/"), "</dt>")
		line("<dd>HTML page listing the available Go pprof profiles</dd>")
		line("<dt>", htmlLink("/debug/pprof/symbol"), "</dt>")
		line("<dd>Looks up the program counters listed in the request, responding with a table mapping program counters to function names</dd>")
		line("<dt>", htmlLink("/debug/run/trace/db"), "</dt>")
		line("<dd>Collects traces until request is canceled</dd>")
		line("<dt>", htmlLink("/mon/"), "</dt>")
		line("<dd>HTML page listing the available Monkit metrics</dd>")
		line("<dt>", htmlLink("/metrics"), "</dt>")
		line("<dd>Monkit metrics in Prometheus format</dd>")
		line("<dt>", htmlLink("/health"), "</dt>")
		line("<dd>Debug server health check</dd>")
		if writeLogging {
			line("<dt>", htmlLink("/logging"), "</dt>")
			line("<dd>JSON endpoint can report or change the current logging level (https://pkg.go.dev/go.uber.org/zap#AtomicLevel.ServeHTTP)</dd>")
		}
		line("</dl>")
	}
}

// Run starts the debug endpoint.
func (server *Server) Run(ctx context.Context) error {
	if server.listener == nil {
		return nil
	}

	ctx, cancel := context.WithCancel(ctx)
	var group errgroup.Group
	group.Go(func() error {
		<-ctx.Done()
		return Error.Wrap(server.server.Shutdown(context.Background()))
	})
	group.Go(func() error {
		defer cancel()

		err := server.server.Serve(server.listener)
		if errors.Is(err, http.ErrServerClosed) {
			err = nil
		}

		return Error.Wrap(err)
	})
	return group.Wait()
}

// Close closes server and underlying listener.
func (server *Server) Close() error {
	return Error.Wrap(server.server.Close())
}

// prometheusMetrics writes https://prometheus.io/docs/instrumenting/exposition_formats/
func (server *Server) prometheusMetrics(w http.ResponseWriter, r *http.Request) {
	// We have to collect all of the metrics before we write. This is because we do not
	// get all of the metrics from the registry sorted by measurement, and from the docs:
	//
	//     All lines for a given metric must be provided as one single group, with the
	//     optional HELP and TYPE lines first (in no particular order). Beyond that,
	//     reproducible sorting in repeated expositions is preferred but not required,
	//     i.e. do not sort if the computational cost is prohibitive.

	data := make(map[string][]string)
	var components []string

	server.registryForRequest(r).Stats(func(key monkit.SeriesKey, field string, val float64) {
		components = components[:0]

		measurement := sanitize(key.Measurement)
		for tag, tagVal := range key.Tags.All() {
			components = append(components,
				fmt.Sprintf("%s=%q", sanitize(tag), sanitize(tagVal)))
		}
		components = append(components,
			fmt.Sprintf("field=%q", sanitize(field)))

		data[measurement] = append(data[measurement],
			fmt.Sprintf("{%s} %g", strings.Join(components, ","), val))
	})

	for measurement, samples := range data {
		_, _ = fmt.Fprintln(w, "# TYPE", measurement, "gauge")
		for _, sample := range samples {
			_, _ = fmt.Fprintf(w, "%s%s\n", measurement, sample)
		}
	}
}

func (server *Server) registryForRequest(r *http.Request) *monkit.Registry {
	outputID := r.URL.Query().Get("output-id")
	server.registryMu.Lock()
	defer server.registryMu.Unlock()
	// okay if outputID is ""
	reg, found := server.registries[outputID]
	if !found {
		reg = ApplyNewTransformers(server.baseRegistry)
		server.registries[outputID] = reg
	}
	return reg
}

// collectTraces collects traces until request is canceled.
func (server *Server) collectTraces(w http.ResponseWriter, r *http.Request) {
	cancel := traces.CollectTraces()
	defer cancel()
	for {
		_, err := w.Write([]byte{0})
		if err != nil {
			return
		}
		time.Sleep(time.Second)
	}
}

// sanitize formats val to be suitable for prometheus.
func sanitize(val string) string {
	// https://prometheus.io/docs/concepts/data_model/
	// specifies all metric names must match [a-zA-Z_:][a-zA-Z0-9_:]*
	// Note: The colons are reserved for user defined recording rules.
	// They should not be used by exporters or direct instrumentation.
	if val == "" {
		return ""
	}
	if '0' <= val[0] && val[0] <= '9' {
		val = "_" + val
	}
	return strings.Map(func(r rune) rune {
		switch {
		case 'a' <= r && r <= 'z':
			return r
		case 'A' <= r && r <= 'Z':
			return r
		case '0' <= r && r <= '9':
			return r
		default:
			return '_'
		}
	}, val)
}

// VersionHandler implements version info endpoint.
type versionHandler struct {
	log *zap.Logger
}

// NewVersionHandler returns new version handler.
func newVersionHandler(log *zap.Logger) *versionHandler {
	return &versionHandler{log}
}

// ServeHTTP returns a json representation of the current version information for the binary.
func (handler *versionHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	j, err := version.Build.Marshal()
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	_, err = w.Write(append(j, '\n'))
	if err != nil {
		handler.log.Error("Error writing data to client", zap.Error(err))
	}
}
