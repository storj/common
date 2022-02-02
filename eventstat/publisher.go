// Copyright (C) 2022 Storj Labs, Inc.
// See LICENSE for copying information.

package eventstat

import (
	"context"
	"os"
	"strings"
	"time"

	"storj.io/common/telemetry"
)

// ClientOpts allows you to set Client Options.
type ClientOpts struct {
	// Interval is how frequently stats from the provided Registry will be
	// sent up. Note that this interval is "jittered", so the actual interval
	// is taken from a normal distribution with a mean of Interval and a
	// variance of Interval/4. Defaults to DefaultInterval.
	Interval time.Duration

	// Application is the application name, usually prepended to metric names.
	// By default it will be os.Args[0].
	Application string

	// Instance is a string that identifies this particular server. Could be a
	// node id, but defaults to the result of DefaultInstanceId().
	Instance string

	// PacketSize controls how we fragment the data as it goes out in UDP
	// packets. Defaults to DefaultPacketSize.
	PacketSize int
}

// UDPPublisher is an  eventstat telemetry client for sending UDP packets at a regular interval.
type UDPPublisher struct {
	reporter *telemetry.Reporter
}

// NewUDPPublisher constructs a telemetry client that sends packets to remoteAddr
// over UDP.
func NewUDPPublisher(remoteAddr string, registry *Registry, opts ClientOpts) (rv *UDPPublisher, err error) {
	if opts.Interval == 0 {
		opts.Interval = telemetry.DefaultInterval
	}
	if opts.Application == "" {
		if len(os.Args) > 0 {
			opts.Application = os.Args[0]
		} else {
			// what the actual heck
			opts.Application = telemetry.DefaultApplication
		}
	}
	if opts.Instance == "" {
		opts.Instance = telemetry.DefaultInstanceID()
	}
	if opts.PacketSize == 0 {
		opts.PacketSize = telemetry.DefaultPacketSize
	}

	udpOptions := telemetry.Options{
		Application: opts.Application,
		InstanceID:  []byte(opts.Instance),
		Address:     remoteAddr,
		PacketSize:  opts.PacketSize,
	}
	reporter, err := telemetry.NewReporter(opts.Interval, func(ctx context.Context) error {
		return telemetry.Send(ctx, udpOptions, func(publishEntry func(key string, value float64)) {
			registry.PublishAndReset(func(name string, tags Tags, value float64) {
				telemetryKey := telemetryKey(name, tags)
				publishEntry(telemetryKey, value)
			})
		})
	})
	if err != nil {
		return nil, err
	}
	return &UDPPublisher{
		reporter: reporter,
	}, nil
}

func telemetryKey(name string, tags Tags) string {
	builder := strings.Builder{}
	writeTag(&builder, name)
	if len(tags) > 0 {
		builder.WriteString(",")
		builder.WriteString(tags.String())
	}
	builder.WriteString(" value")
	telemetryKey := builder.String()
	return telemetryKey
}

// Run calls Report roughly every Interval.
func (c *UDPPublisher) Run(ctx context.Context) {
	c.reporter.Run(ctx)
}

// publish sends out message immediately independent on Interval.
func (c *UDPPublisher) publish(ctx context.Context) error {
	return c.reporter.Publish(ctx)
}

// writeTag writes a tag key, value, or field key to the builder.
func writeTag(builder *strings.Builder, tag string) {
	if strings.IndexByte(tag, ',') == -1 &&
		strings.IndexByte(tag, '=') == -1 &&
		strings.IndexByte(tag, ' ') == -1 {

		builder.WriteString(tag)
		return
	}

	for i := 0; i < len(tag); i++ {
		if tag[i] == ',' ||
			tag[i] == '=' ||
			tag[i] == ' ' {
			builder.WriteByte('\\')
		}
		builder.WriteByte(tag[i])
	}
}
