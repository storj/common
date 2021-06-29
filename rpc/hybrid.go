// Copyright (C) 2021 Storj Labs, Inc.
// See LICENSE for copying information.

package rpc

import (
	"context"
	"crypto/tls"
	"fmt"
	"sync"

	"github.com/zeebo/errs"

	"storj.io/common/memory"
)

// HybridConnector implements a dialer that creates a connection using any of
// (potentially) multiple connector candidates. The fastest one is kept, and
// all others are closed and discarded.
type HybridConnector struct {
	connectors []candidateConnector
}

// candidateConnector encapsulates a Connector paired with a name and a
// priority value. The priority value is used to break ties when multiple
// connectors succeed in creating a connection at or near the same time.
type candidateConnector struct {
	name      string
	connector Connector
	priority  int
}

// candidateConnectorType encapsulates a connector type (by way of what might
// be called a "factory function") paired with a name and a priority value.
// The name and priority value will be inherited by the corresponding
// candidateConnector instance owned by a HybridConnector.
type candidateConnectorType struct {
	name          string
	connectorType func() Connector
	priority      int
}

const tcpConnectorPriority = 10

var (
	// connectorRegistryLock must be held when changing or accessing
	// connectorRegistry.
	connectorRegistryLock sync.Mutex

	// connectorRegistry contains a list of connector types that all
	// new HybridConnector instances should have by default. Other
	// packages may add to this list using RegisterCandidateConnectorType.
	connectorRegistry = []candidateConnectorType{
		{
			name:          "tcp",
			connectorType: func() Connector { return NewDefaultTCPConnector(nil) },
			priority:      tcpConnectorPriority,
		},
	}
)

// RegisterCandidateConnectorType registers a type of connector for use with
// all HybridConnector instances created in the future. If the new connector
// type has the same name as one that is already registered, it will replace
// the preexisting entry.
func RegisterCandidateConnectorType(name string, connectorType func() Connector, priority int) {
	connectorRegistryLock.Lock()
	defer connectorRegistryLock.Unlock()
	newConnDefinition := candidateConnectorType{
		name:          name,
		connectorType: connectorType,
		priority:      priority,
	}
	for i, connDefinition := range connectorRegistry {
		if connDefinition.name == name {
			connectorRegistry[i] = newConnDefinition
			return
		}
	}
	connectorRegistry = append(connectorRegistry, newConnDefinition)
}

// NewHybridConnector instantiates a new instance of HybridConnector with
// all registered connector types.
func NewHybridConnector() HybridConnector {
	connectorRegistryLock.Lock()
	candidates := make([]candidateConnectorType, len(connectorRegistry))
	copy(candidates, connectorRegistry)
	connectorRegistryLock.Unlock()

	connectors := make([]candidateConnector, len(candidates))
	for i, connDefinition := range candidates {
		connectors[i].name = connDefinition.name
		connectors[i].connector = connDefinition.connectorType()
		connectors[i].priority = connDefinition.priority
	}
	return HybridConnector{connectors: connectors}
}

type candidateConnection struct {
	conn     ConnectorConn
	name     string
	priority int
}

// AddCandidateConnector adds a candidate connector to this HybridConnector
// instance. (Other HybridConnector instances, both current and future, will
// not be affected by this call.
//
// It is recommended that this be used before any connections are attempted
// with the HybridConnector, because no concurrency protection is built in to
// accesses to c.connectors.
func (c *HybridConnector) AddCandidateConnector(name string, connector Connector, priority int) {
	c.connectors = append(c.connectors, candidateConnector{
		name:      name,
		connector: connector,
		priority:  priority,
	})
}

// RemoveCandidateConnector removes a candidate connector from this
// HybridConnector instance, if there is one with the given name.
//
// It is recommended that this be used before any connections are attempted
// with the HybridConnector, because no concurrency protection is built in to
// accesses to c.connectors.
func (c *HybridConnector) RemoveCandidateConnector(name string) {
	for i, candidate := range c.connectors {
		if candidate.name == name {
			c.connectors = append(c.connectors[:i], c.connectors[i+1:]...)
			return
		}
	}
}

// DialContext creates an encrypted connection using one of the candidate
// connectors. All connectors are started at the same time, and the first one
// to finish will have its connection returned (and other connectors will be
// canceled). If multiple connectors finish before they are canceled, the
// connection with the highest priority value is kept.
func (c HybridConnector) DialContext(ctx context.Context, tlsConfig *tls.Config, address string) (_ ConnectorConn, err error) {
	defer mon.Task()(&ctx)(&err)

	if tlsConfig == nil {
		return nil, Error.New("tls config is not set")
	}

	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	var chosen candidateConnection
	errChan := make(chan error)
	readyChan := make(chan candidateConnection)

	for _, entry := range c.connectors {
		entry := entry
		go func() {
			conn, err := entry.connector.DialContext(ctx, tlsConfig.Clone(), address)
			if err != nil {
				errChan <- fmt.Errorf("%s connector failed: %w", entry.name, err)
				return
			}
			readyChan <- candidateConnection{
				conn:     conn,
				name:     entry.name,
				priority: entry.priority,
			}
		}()
	}

	var errors []error
	var numFinished int
	// make sure all dials are finished either with an established connection or
	// an error. This allows us to appropriately close extra connection if multiple
	// connections are ready around the same time
	for numFinished < len(c.connectors) {
		select {
		case candidate := <-readyChan:
			numFinished++
			// cancel all other dials (they might be ready too, or they might
			// become ready before they receive this cancellation message; that
			// is ok). if cancel() was already called, that's fine too; this
			// will do nothing.
			cancel()

			if candidate.priority > chosen.priority {
				// if multiple connectors finish successfully, we will keep
				// the one with the highest priority value
				if chosen.conn != nil {
					// discard a previous choice
					_ = chosen.conn.Close()
				}
				chosen = candidate
			} else {
				// discard the new candidate
				_ = candidate.conn.Close()
			}

		case err := <-errChan:
			numFinished++
			errors = append(errors, err)
		}
	}

	if chosen.priority == 0 {
		// no connectors succeeded!
		mon.Event("hybrid_connector_established_no_connection")
		return nil, errs.Combine(errors...)
	}
	mon.Event(fmt.Sprintf("hybrid_connector_established_%s_connection", chosen.name))
	return chosen.conn, nil
}

// SetTransferRate calls SetTransferRate with the given transfer rate on all of
// its candidate connectors (if they have a SetTransferRate method).
func (c *HybridConnector) SetTransferRate(rate memory.Size) {
	for _, entry := range c.connectors {
		if entry, ok := entry.connector.(interface{ SetTransferRate(rate memory.Size) }); ok {
			entry.SetTransferRate(rate)
		}
	}
}
