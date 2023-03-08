// Copyright (C) 2021 Storj Labs, Inc.
// See LICENSE for copying information.

package rpcpool_test

import (
	"context"
	"crypto/tls"
	"fmt"
	"log"
	"time"

	"storj.io/common/identity"
	"storj.io/common/peertls/tlsopts"
	"storj.io/common/rpc"
	"storj.io/common/rpc/rpcpool"
	"storj.io/common/rpc/rpctest"
	"storj.io/drpc"
)

// Example shows how the wrapper can be used to wrap connection multiple times.
func Example() {
	ctx := context.Background()

	id, err := identity.NewFullIdentity(ctx, identity.NewCAOptions{
		Difficulty:  0,
		Concurrency: 1,
	})
	if err != nil {
		log.Printf("%+v\n", err)
	}

	tlsOptions, err := tlsopts.NewOptions(id, tlsopts.Config{}, nil)
	if err != nil {
		log.Printf("%+v\n", err)
	}

	d := rpc.NewDefaultDialer(tlsOptions)

	cr := rpctest.NewCallRecorder()

	ctx = rpcpool.WithDialerWrapper(ctx, func(ctx context.Context, dialer rpcpool.Dialer) rpcpool.Dialer {

		return func(context.Context) (conn rpcpool.RawConn, state *tls.ConnectionState, err error) {

			// this is only for testing when connection is mocked
			stub := rpctest.NewStubConnection()
			state = &tls.ConnectionState{}
			conn = &stub

			// for real world example start with delegating the call to the original dialer
			// conn, state, err = dialer(ctx)

			// adding first wrapper (call recorder)
			conn = cr.Attach(conn)

			// add additional latency wrapper
			conn = rpctest.ConnectionWithLatency(conn, 10*time.Millisecond)

			return conn, state, err
		}
	})

	conn, err := d.DialAddressInsecure(ctx, "localhost:1234")
	if err != nil {
		log.Printf("%+v\n", err)
	}

	in := message{}
	out := message{}
	err = conn.Invoke(ctx, "rpccall", messageEncoding{}, &in, &out)

	if err != nil {
		fmt.Println(cr.History())
	}
	// Output: [rpccall]
}

type message struct {
	content string
}

type messageEncoding struct {
}

func (messageEncoding) Marshal(msg drpc.Message) ([]byte, error) {
	return []byte(msg.(*message).content), nil
}

func (messageEncoding) Unmarshal(buf []byte, msg drpc.Message) error {
	msg.(*message).content = string(buf)
	return nil
}
