// Copyright (C) 2024 Storj Labs, Inc.
// See LICENSE for copying information.

package eventkitbq

import (
	"context"

	"go.uber.org/zap"

	"storj.io/eventkit"
	"storj.io/eventkit/bigquery"
)

// BQDestination initializes the BQ destination.
// Context should be cancelled to stop internal goroutines.
func BQDestination(ctx context.Context, log *zap.Logger, destConfig string, eventRegistry *eventkit.Registry, appName string, instanceID string) {
	c, err := bigquery.CreateDestination(ctx, destConfig)
	if err != nil {
		log.Error("Eventkit BQ destination couldn't be initialized", zap.Error(err))
		return
	}

	eventRegistry.AddDestination(c)
	go c.Run(ctx)
}
