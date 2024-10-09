// Copyright (C) 2024 Storj Labs, Inc.
// See LICENSE for copying information.

//go:build linux

package debug

import (
	"context"
	"net"
	"reflect"
	"runtime/trace"
	"sync"
	"sync/atomic"
	"syscall"
	"time"

	"github.com/google/gopacket"
	"github.com/google/gopacket/layers"
	"github.com/google/gopacket/pcapgo"
)

func capturePackets(ctx context.Context, stop *atomic.Bool) {
	var wg sync.WaitGroup
	defer wg.Wait()

	type Handle struct {
		eh    *pcapgo.EthernetHandle
		iface net.Interface
	}

	var handles []Handle

	ifaces, err := net.Interfaces() // ignore errors because pcap is supplemental
	trace.Logf(ctx, "trace-debug", "found %d interfaces (err:%v)", len(ifaces), err)
	for _, iface := range ifaces {
		trace.Logf(ctx, "trace-debug", "checking interface %q", iface.Name)

		if iface.Flags&net.FlagLoopback != 0 {
			continue
		}
		if iface.Flags&(net.FlagUp|net.FlagRunning) != net.FlagUp|net.FlagRunning {
			continue
		}
		if len(iface.HardwareAddr) == 0 {
			continue
		}
		addrs, err := iface.Addrs()
		if err != nil {
			continue
		}
		hasIPv4 := false
		for _, addr := range addrs {
			ip, _ := addr.(*net.IPNet)
			hasIPv4 = hasIPv4 || len(ip.IP.To4()) == net.IPv4len
		}
		if !hasIPv4 {
			continue
		}

		eh, err := pcapgo.NewEthernetHandle(iface.Name)
		if err != nil {
			trace.Logf(ctx, "trace-debug", "could not open handle for %q: %v", iface.Name, err)
			continue
		}
		trace.Logf(ctx, "trace-debug", "opened handle for %q", iface.Name)
		defer eh.Close()

		handles = append(handles, Handle{
			eh:    eh,
			iface: iface,
		})
	}

	for _, handle := range handles {
		handle := handle // avoid loop capture bug

		wg.Add(1)
		go func() {
			defer wg.Done()

			src := gopacket.NewPacketSource(handle.eh, layers.LinkTypeEthernet)
			for {
				packet, err := src.NextPacket()
				if err != nil || stop.Load() {
					return
				}

				tcp, _ := packet.Layer(layers.LayerTypeTCP).(*layers.TCP)
				ip, _ := packet.Layer(layers.LayerTypeIPv4).(*layers.IPv4)
				if tcp == nil || ip == nil {
					continue
				}

				trace.Logf(ctx, "tcp-packet",
					"if:%s local:%s:%d remote:%s:%d seq:%d ack:%d flags:%d window:%d payload:%d fragoff:%d ts:%d",
					handle.iface.Name,
					ip.SrcIP,
					tcp.SrcPort,
					ip.DstIP,
					tcp.DstPort,
					tcp.Seq,
					tcp.Ack,
					makeTCPFlags(tcp),
					tcp.Window,
					len(tcp.Payload),
					ip.FragOffset,
					packet.Metadata().Timestamp.UnixNano(),
				)
			}
		}()
	}

	// wait for all of the handles to be done and send a signal when they are
	doneWaiting := make(chan struct{})
	go func() {
		wg.Wait()
		close(doneWaiting)
	}()

	select {
	case <-doneWaiting:
		// all of the handles are not being used anymore, so we can do clean
		// close calls
		for _, handle := range handles {
			handle.eh.Close()
		}

	case <-time.After(10 * time.Second):
		// at least one handle is still being used and so after 10 seconds is
		// almost certainly blocked in the syscall, so we can safely interrupt
		// it with a close call on the fd and be reasonably sure that no reads
		// will happen on a new socket that got the same fd.
		for _, handle := range handles {
			fd := reflect.ValueOf(handle.eh).Elem().FieldByName("fd").Int()
			_ = syscall.Close(int(fd))
		}

		<-doneWaiting
	}
}

func makeTCPFlags(t *layers.TCP) (f uint32) {
	if t.FIN {
		f |= 0x0001
	}
	if t.SYN {
		f |= 0x0002
	}
	if t.RST {
		f |= 0x0004
	}
	if t.PSH {
		f |= 0x0008
	}
	if t.ACK {
		f |= 0x0010
	}
	if t.URG {
		f |= 0x0020
	}
	if t.ECE {
		f |= 0x0040
	}
	if t.CWR {
		f |= 0x0080
	}
	if t.NS {
		f |= 0x0100
	}
	return f
}
