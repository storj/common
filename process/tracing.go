// Copyright (C) 2019 Storj Labs, Inc.
// See LICENSE for copying information.

package process

import (
	"context"
	"flag"
	"os"
	"path/filepath"

	hw "github.com/jtolds/monkit-hw/v2"
	"github.com/spacemonkeygo/monkit/v3"
	"github.com/spacemonkeygo/monkit/v3/environment"
	"go.uber.org/zap"

	"storj.io/common/identity"
	"storj.io/common/telemetry"
	jaeger "storj.io/monkit-jaeger"
	"storj.io/private/version"
)

var (
	tracingSamplingRate = flag.Float64("tracing.sample", 0, "how frequently to send up telemetry")
	tracingAgent        = flag.String("tracing.agent-addr", flagDefault("", "agent.tracing.datasci.storj.io"), "address for jaeger agent")
	tracingApp          = flag.String("tracing.app", filepath.Base(os.Args[0]), "application name for tracing identification")
	tracingAppSuffix    = flag.String("tracing.app-suffix", flagDefault("-dev", "-release"), "application suffix")
	tracingBufferSize   = flag.Int("tracing.buffer-size", 0, "buffer size for collector batch packet size")
	tracingQueueSize    = flag.Int("tracing.queue-size", 0, "buffer size for collector queue size")
	tracingInterval     = flag.Duration("tracing.interval", 0, "how frequently to send up telemetry")
)

const (
	instanceIDKey = "instanceID"
	hostnameKey   = "hostname"
)

// InitTracing initializes distributed tracing with an instance ID.
func InitTracing(ctx context.Context, log *zap.Logger, r *monkit.Registry, instanceID string) (*jaeger.UDPCollector, func(), error) {
	return run(ctx, log, r, instanceID, []jaeger.Tag{})
}

// InitTracingWithCertPath initializes distributed tracing with certificate path.
func InitTracingWithCertPath(ctx context.Context, log *zap.Logger, r *monkit.Registry, certDir string) (*jaeger.UDPCollector, func(), error) {
	return run(ctx, log, r, nodeIDFromCertPath(ctx, log, certDir), []jaeger.Tag{})
}

// InitTracingWithHostname initializes distributed tracing with nodeID and hostname.
func InitTracingWithHostname(ctx context.Context, log *zap.Logger, r *monkit.Registry, certDir string) (*jaeger.UDPCollector, func(), error) {
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

	return run(ctx, log, r, nodeIDFromCertPath(ctx, log, certDir), processInfo)
}

func run(ctx context.Context, log *zap.Logger, r *monkit.Registry, instanceID string, processInfo []jaeger.Tag) (collector *jaeger.UDPCollector, cancel func(), err error) {
	if r == nil {
		r = monkit.Default
	}
	environment.Register(r)
	hw.Register(r)
	r.ScopeNamed("env").Chain(monkit.StatSourceFunc(version.Build.Stats))

	log = log.Named("tracing")
	if *tracingAgent == "" || *tracingSamplingRate == 0 {
		log.Info("disabled")
		return nil, nil, nil
	}

	if len(instanceID) == 0 {
		instanceID = telemetry.DefaultInstanceID()
	}
	processInfo = append(processInfo, jaeger.Tag{
		Key:   instanceIDKey,
		Value: instanceID,
	})

	processName := *tracingApp + *tracingAppSuffix
	if len(processName) > maxInstanceLength {
		processName = processName[:maxInstanceLength]
	}
	collector, err = jaeger.NewUDPCollector(log, *tracingAgent, processName, processInfo, *tracingBufferSize, *tracingQueueSize, *tracingInterval)
	if err != nil {
		return nil, nil, err
	}
	cancel = jaeger.RegisterJaeger(r, collector, jaeger.Options{Fraction: *tracingSamplingRate})
	return collector, cancel, nil
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
