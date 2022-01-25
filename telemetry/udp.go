// Copyright (C) 2022 Storj Labs, Inc.
// See LICENSE for copying information.

package telemetry

import (
	"context"
	"log"
	"net"
	"syscall"

	"github.com/zeebo/admission/v3/admproto"
	"github.com/zeebo/errs"
)

// Options define all the required parameters to send out UDP telemetry packages.
type Options struct {
	// Application to send with
	Application string

	// Instance Id to send with
	InstanceID []byte

	// Address to send packets to
	Address string

	// PacketSize controls maximum packet size. If zero, 1024 is used.
	PacketSize int

	// ProtoOps allows you to set protocol options.
	ProtoOpts admproto.Options

	// Headers allow you to set arbitrary key/value pairs to be included in each packet send
	Headers map[string]string
}

// Send sends out telemetry via UDP as key/value pairs.
// forEachValue will be called with a callback to get all the key/values to send out.
func Send(ctx context.Context, opts Options, forEachValue func(func(key string, value float64))) (err error) {
	addr, err := net.ResolveUDPAddr("udp", opts.Address)
	if err != nil {
		return err
	}

	conn, err := net.DialUDP("udp", nil, addr)
	if err != nil {
		return err
	}
	defer func() {
		_ = conn.Close()
	}()

	if opts.PacketSize == 0 {
		opts.PacketSize = 1024
	}

	var (
		buf []byte
		w   = admproto.NewWriterWith(opts.ProtoOpts)
	)

	forEachValue(func(key string, value float64) {
		// if we have any errors, stop.
		if err != nil {
			return
		}

		for {
			// keep track of the buffer before we send
			before := buf

			// always ensure the buffer has the prefix in it.
			if len(buf) == 0 {
				// if we can't add the application and instance id, it's fatal.
				buf, err = w.Begin(buf, opts.Application, opts.InstanceID, len(opts.Headers))
				if err != nil {
					return
				}
				for key, value := range opts.Headers {
					buf, err = w.AppendHeader(buf, []byte(key), []byte(value))
					if err != nil {
						return
					}
				}
			}

			// add the value to the buffer
			buf, err = w.Append(buf, key, value)
			if err != nil {
				// not fatal, just back up to before, but let someone know
				// it has been skipped.
				log.Println("skipped metric", key, "because", err)
				buf, err = before, nil
				return
			}

			// if we're still in the packet size, then get the next metric.
			if len(buf)+4 <= opts.PacketSize {
				return
			}

			// if we're over the packet size, send the previous value and start
			// over. be sure to account for the checksum that sendPacket adds.
			// if buf was empty at the start, we should just send it.
			// otherwise we should send the previous value.
			if len(before) == 0 {
				sendPacket(ctx, conn, buf)
			} else {
				sendPacket(ctx, conn, before)
			}

			// after sending the packet, we should reset the buffer and try to
			// add the point again.
			w.Reset()
			buf = buf[:0]

			// if we had no buffer at the start, then we sent this metric, so
			// return to get the next metric.
			if len(before) == 0 {
				return
			}
		}
	})
	// send off any remainder buf. we're guaranteed by the loop above that if
	// there is any data in buf it forms a valid packet with metrics in it.
	if err == nil && len(buf) > 0 {
		sendPacket(ctx, conn, buf)
	}

	return err
}

// sendPacket is a helper that adds a checksum to the provided buffer and sends
// it to the conn. It logs if there was an error.
func sendPacket(ctx context.Context, conn *net.UDPConn, buf []byte) {
	_, err := conn.Write(admproto.AddChecksum(buf))
	if err != nil && errs.Is(err, syscall.ENOBUFS) {
		log.Println("failed to send packet:", err)
	}
}
