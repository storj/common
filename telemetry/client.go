// Copyright (C) 2019 Storj Labs, Inc.
// See LICENSE for copying information.

package telemetry

import (
	"context"
	"os"
	"time"

	"github.com/spacemonkeygo/monkit/v3"
	"github.com/zeebo/admission/v3/admproto"
)

const (
	// DefaultInterval is the default amount of time between metric payload sends.
	DefaultInterval = time.Minute

	// DefaultPacketSize sets the target packet size. MTUs are often 1500,
	// though a good argument could be made for 512.
	DefaultPacketSize = 1000

	// DefaultApplication is the default values for application name. Should be used
	// when value in ClientOpts.Application is not set and len(os.Args) == 0.
	DefaultApplication = "unknown"
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

	// Registry is where to get stats from. Defaults to monkit.Default.
	Registry *monkit.Registry

	// FloatEncoding is how floats should be encoded on the wire.
	// Default is float16.
	FloatEncoding admproto.FloatEncoding

	// Headers allow you to set arbitrary key/value tags to be included in
	// each packet send.
	Headers map[string]string
}

func (o *ClientOpts) fillDefaults() {
	if o.Interval == 0 {
		o.Interval = DefaultInterval
	}
	if o.Application == "" {
		if len(os.Args) > 0 {
			o.Application = os.Args[0]
		} else {
			// what the actual heck
			o.Application = DefaultApplication
		}
	}
	if o.Instance == "" {
		o.Instance = DefaultInstanceID()
	}
	if o.Registry == nil {
		o.Registry = monkit.Default
	}
	if o.PacketSize == 0 {
		o.PacketSize = DefaultPacketSize
	}
}

// Client is a telemetry client for sending UDP packets at a regular interval
// from a monkit.Registry.
type Client struct {
	reporter *Reporter
}

// NewClient constructs a telemetry client that sends packets to remoteAddr
// over UDP.
func NewClient(remoteAddr string, opts ClientOpts) (rv *Client, err error) {
	opts.fillDefaults()

	options := Options{
		Application: opts.Application,
		InstanceID:  []byte(opts.Instance),
		Address:     remoteAddr,
		PacketSize:  opts.PacketSize,
		ProtoOpts:   admproto.Options{FloatEncoding: opts.FloatEncoding},
		Headers:     opts.Headers,
	}

	reporter, err := NewReporter(opts.Interval, func(ctx context.Context) error {
		return Send(ctx, options, func(entries func(key string, value float64)) {
			if opts.Registry == nil {
				opts.Registry = monkit.Default
			}
			opts.Registry.Stats(func(key monkit.SeriesKey, field string, val float64) {
				series := key.WithField(field)
				entries(series, val)
			})
		})
	})
	if err != nil {
		return nil, err
	}
	return &Client{
		reporter: reporter,
	}, nil
}

// Run calls Report roughly every Interval.
func (c *Client) Run(ctx context.Context) {
	c.reporter.Run(ctx)
}

// Report bundles up all the current stats and writes them out as UDP packets.
func (c *Client) Report(ctx context.Context) (err error) {
	defer mon.Task()(&ctx)(&err)
	return c.reporter.Publish(ctx)
}
