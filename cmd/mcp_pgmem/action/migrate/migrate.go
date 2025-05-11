package migrate

import (
	"context"

	"github.com/tyrm/mcp-pgmem/cmd/mcp_pgmem/action"
	"github.com/tyrm/mcp-pgmem/internal/db/bun"
	"go.uber.org/zap"
)

// Migrate runs database migrations.
var Migrate action.Action = func(ctx context.Context, _ []string) error {
	zap.L().Info("running database migration")

	// create database client
	dbClient, err := bun.New(ctx)
	if err != nil {
		zap.L().Error("Error creating bun client", zap.Error(err))

		return err
	}
	defer func() {
		err := dbClient.Close()
		if err != nil {
			zap.L().Error("Error closing bun client", zap.Error(err))
		}
	}()

	err = dbClient.DoMigration(ctx)
	if err != nil {
		zap.L().Error("Error running migration", zap.Error(err))

		return err
	}

	return nil
}
