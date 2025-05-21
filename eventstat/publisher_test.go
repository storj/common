// Copyright (C) 2022 Storj Labs, Inc.
// See LICENSE for copying information.

package eventstat

import (
	"fmt"
	"runtime"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"storj.io/common/telemetry"
	"storj.io/common/testcontext"
)

func TestMetrics(t *testing.T) {
	ctx := testcontext.New(t)

	if runtime.GOOS == "windows" {
		// TODO(windows): currently closing doesn't seem to be shutting down the server
		t.Skip("broken")
	}

	s, err := telemetry.Listen("127.0.0.1:0")
	assert.NoError(t, err)
	defer func() { _ = s.Close() }()

	r := Registry{}
	counter := r.NewTagCounter("http_user_agent", "agent")
	fmt.Println(s.Addr())
	c, err := NewUDPPublisher(s.Addr(), &r, ClientOpts{
		Application: "testapp",
		Instance:    "testinst",
		Interval:    10 * time.Millisecond,
	})
	require.NoError(t, err)

	counter("aws")
	counter("aws")
	counter("aws")
	counter("rclone")

	err = c.publish(ctx)
	require.NoError(t, err)

	expectedMetric := 4

	keys := make(chan string, expectedMetric)
	values := make(chan float64, expectedMetric)
	defer close(keys)
	defer close(values)

	ctx.Go(func() error {
		fmt.Println("Listening on " + s.Addr())
		// note: this is the telemetry server which guarantees that our sender is still compatible with the format
		_ = s.Serve(ctx, telemetry.HandlerFunc(
			func(application, instance string, key []byte, val float64) {
				assert.Equal(t, application, "testapp")
				assert.Equal(t, instance, "testinst")
				keys <- string(key)
				values <- val
			}))
		return nil
	})

	for range expectedMetric {
		key := <-keys
		value := <-values

		switch key {
		case "http_user_agent_count,agent=aws value":
			assert.Equal(t, float64(3), value)
		case "http_user_agent_count,agent=rclone value":
			assert.Equal(t, float64(1), value)
		case "http_user_agent_discarded value":
			assert.Equal(t, float64(0), value)
		case "http_user_agent_buckets value":
			assert.Equal(t, float64(2), value)
		default:
			require.Failf(t, "Unexpected UDP metric", "key=%s", key)
		}

	}

}

func TestTelemetryKey(t *testing.T) {
	assert.Equal(t, "key1 value", telemetryKey("key1", Tags{}))
	assert.Equal(t, "key2,foo=bar value", telemetryKey("key2", Tags{"foo": "bar"}))
	assert.Equal(t, "key3,f\\==ba\\,r value", telemetryKey("key3", Tags{"f=": "ba,r"}))
}
