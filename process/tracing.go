// Copyright (C) 2019 Storj Labs, Inc.
// See LICENSE for copying information.

package process

import (
	"context"
	"flag"
	"os"
	"path/filepath"
	"regexp"

	"github.com/spacemonkeygo/monkit/v3"
	"go.uber.org/zap"
	"golang.org/x/sync/errgroup"

	"storj.io/common/identity"
	"storj.io/common/telemetry"
	"storj.io/common/tracing"
	jaeger "storj.io/monkit-jaeger"
)

var (
	tracingEnabled      = flag.Bool("tracing.enabled", true, "whether tracing collector is enabled")
	tracingSamplingRate = flag.Float64("tracing.sample", 0, "how frequent to sample traces")
	tracingAgent        = flag.String("tracing.agent-addr", flagDefault("127.0.0.1:5775", "agent.tracing.datasci.storj.io:5775"), "address for jaeger agent")
	tracingApp          = flag.String("tracing.app", filepath.Base(os.Args[0]), "application name for tracing identification")
	tracingAppSuffix    = flag.String("tracing.app-suffix", flagDefault("-dev", "-release"), "application suffix")
	tracingBufferSize   = flag.Int("tracing.buffer-size", 0, "buffer size for collector batch packet size")
	tracingQueueSize    = flag.Int("tracing.queue-size", 0, "buffer size for collector queue size")
	tracingInterval     = flag.Duration("tracing.interval", 0, "how frequently to flush traces to tracing agent")
	tracingHostRegex    = flag.String("tracing.host-regex", `\.storj\.tools:[0-9]+$`, "the possible hostnames that trace-host designated traces can be sent to")
)

const (
	instanceIDKey = "instanceID"
	hostnameKey   = "hostname"
)

// InitTracing initializes distributed tracing with an instance ID.
func InitTracing(ctx context.Context, log *zap.Logger, r *monkit.Registry, instanceID string) (func(), error) {
	return initTracing(ctx, log, r, instanceID, []jaeger.Tag{})
}

// InitTracingWithCertPath initializes distributed tracing with certificate path.
func InitTracingWithCertPath(ctx context.Context, log *zap.Logger, r *monkit.Registry, certDir string) (func(), error) {
	return initTracing(ctx, log, r, nodeIDFromCertPath(ctx, log, certDir), []jaeger.Tag{})
}

// InitTracingWithHostname initializes distributed tracing with nodeID and hostname.
func InitTracingWithHostname(ctx context.Context, log *zap.Logger, r *monkit.Registry, certDir string) (func(), error) {
	var processInfo []jaeger.Tag
	hostname, err := os.Hostname()
	if err != nil {
		log.Error("Could not read hostname for tracing setup", zap.Error(err))
	} else {
		processInfo = append(processInfo, jaeger.Tag{
			Key:   hostnameKey,
			Value: hostname,
		})
	}

	return initTracing(ctx, log, r, nodeIDFromCertPath(ctx, log, certDir), processInfo)
}

type traceCollectorFactoryFunc func(hostTarget string) (jaeger.ClosableTraceCollector, error)

func (f traceCollectorFactoryFunc) MakeCollector(hostTarget string) (jaeger.ClosableTraceCollector, error) {
	return f(hostTarget)
}

func initTracing(ctx context.Context, log *zap.Logger, r *monkit.Registry, instanceID string, processInfo []jaeger.Tag) (cleanup func(), err error) {
	// Snapshot all flag values first to avoid data races with concurrent goroutines
	enabled := *tracingEnabled
	samplingRate := *tracingSamplingRate
	agentAddr := *tracingAgent
	app := *tracingApp
	appSuffix := *tracingAppSuffix
	bufferSize := *tracingBufferSize
	queueSize := *tracingQueueSize
	interval := *tracingInterval
	hostRegexStr := *tracingHostRegex

	if r == nil {
		r = monkit.Default
	}

	hostRegex, err := regexp.Compile(hostRegexStr)
	if err != nil {
		return nil, err
	}

	if !enabled {
		log.Debug("Anonymized tracing disabled")
		return nil, nil
	}

	log.Info("Anonymized tracing enabled")

	if len(instanceID) == 0 {
		instanceID = telemetry.DefaultInstanceID()
	}
	processInfo = append(processInfo, jaeger.Tag{
		Key:   instanceIDKey,
		Value: instanceID,
	})

	processName := app + appSuffix
	if len(processName) > maxInstanceLength {
		processName = processName[:maxInstanceLength]
	}
	collector, err := jaeger.NewThriftCollector(log, agentAddr, processName, processInfo, bufferSize, queueSize, interval)
	if err != nil {
		return nil, err
	}
	var eg errgroup.Group

	collectorCtx, collectorCtxCancel := context.WithCancel(ctx)
	eg.Go(func() error {
		collector.Run(collectorCtx)
		return nil
	})

	unregister := jaeger.RegisterJaeger(r, collector, jaeger.Options{
		Fraction: samplingRate,
		Excluded: tracing.IsExcluded,
		CollectorFactory: traceCollectorFactoryFunc(func(targetHost string) (jaeger.ClosableTraceCollector, error) {
			targetCollector, err := jaeger.NewThriftCollector(log, targetHost, processName,
				processInfo, bufferSize, queueSize, interval)
			if err != nil {
				return nil, err
			}
			targetCollectorCtx, targetCollectorCancel := context.WithCancel(collectorCtx)
			eg.Go(func() error {
				targetCollector.Run(targetCollectorCtx)
				return nil
			})
			return &closableCollector{
				cancel:                 targetCollectorCancel,
				ClosableTraceCollector: targetCollector,
			}, nil
		}),
		CollectorFactoryHostMatch: hostRegex,
	})
	return func() {
		unregister()
		collectorCtxCancel()
		_ = collector.Close()
		_ = eg.Wait()
	}, nil
}

type closableCollector struct {
	jaeger.ClosableTraceCollector
	cancel func()
}

func (collector *closableCollector) Close() error {
	collector.cancel()
	return collector.ClosableTraceCollector.Close()
}

func nodeIDFromCertPath(ctx context.Context, log *zap.Logger, certPath string) string {
	if certPath == "" {
		return ""
	}
	nodeID, err := identity.NodeIDFromCertPath(certPath)
	if err != nil {
		log.Debug("Could not read identity for tracing setup", zap.Error(err))
		return ""
	}

	return nodeID.String()
}
