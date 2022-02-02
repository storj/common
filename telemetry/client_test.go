// Copyright (C) 2019 Storj Labs, Inc.
// See LICENSE for copying information.
package telemetry

import (
	"context"
	"os"
	"testing"
	"time"

	"github.com/spacemonkeygo/monkit/v3"
	"github.com/stretchr/testify/assert"
)

func TestClientOpts_UseDefaults(t *testing.T) {
	opts := ClientOpts{
		Application: "testapp",
		Instance:    "testinst",
		Interval:    0,
	}
	opts.fillDefaults()
	assert.Equal(t, DefaultInterval, opts.Interval)
	assert.Equal(t, monkit.Default, opts.Registry)
	assert.Equal(t, DefaultPacketSize, opts.PacketSize)
}

func TestNewClient_ApplicationAndArgsAreEmpty(t *testing.T) {
	oldArgs := os.Args

	defer func() {
		os.Args = oldArgs
	}()

	os.Args = nil

	opts := ClientOpts{
		Application: "",
		Instance:    "testinst",
		Interval:    0,
	}
	opts.fillDefaults()

	assert.Equal(t, DefaultApplication, opts.Application)
}

func TestNewClient_ApplicationIsEmpty(t *testing.T) {
	opts := ClientOpts{
		Application: "",
		Instance:    "testinst",
		Interval:    0,
	}
	opts.fillDefaults()

	assert.Equal(t, os.Args[0], opts.Application)
}

func TestNewClient_InstanceIsEmpty(t *testing.T) {
	opts := ClientOpts{
		Application: "qwe",
		Instance:    "",
		Interval:    0,
	}
	opts.fillDefaults()

	assert.Equal(t, DefaultInstanceID(), opts.Instance)
}

func TestRun_ReportNoCalled(t *testing.T) {
	client, err := NewClient("127.0.0.1:0", ClientOpts{
		Application: "qwe",
		Instance:    "",
		Interval:    time.Millisecond,
		PacketSize:  0,
	})
	assert.NoError(t, err)

	client.reporter.send = func(context.Context) error {
		t.Fatal("shouldn't be called")
		return nil
	}

	ctx, cancel := context.WithCancel(context.Background())
	cancel()

	client.Run(ctx)
}
