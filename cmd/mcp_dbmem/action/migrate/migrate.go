package migrate

import (
	"context"

	"github.com/spf13/viper"
	"github.com/tyrm/mcp-dbmem/cmd/mcp_dbmem/action"
	"github.com/tyrm/mcp-dbmem/internal/config"
	"github.com/tyrm/mcp-dbmem/internal/db/bun"
	"go.uber.org/zap"
)

// Migrate runs database migrations.
var Migrate action.Action = func(ctx context.Context, _ []string) error {
	zap.L().Info("running database migration")

	// create database client
	// create database client
	dbClient, err := bun.New(ctx, bun.ClientConfig{
		Type:      viper.GetString(config.Keys.DBType),
		Address:   viper.GetString(config.Keys.DBAddress),
		Port:      viper.GetUint16(config.Keys.DBPort),
		User:      viper.GetString(config.Keys.DBUser),
		Password:  viper.GetString(config.Keys.DBPassword),
		Database:  viper.GetString(config.Keys.DBDatabase),
		TLSMode:   viper.GetString(config.Keys.DBTLSMode),
		TLSCACert: viper.GetString(config.Keys.DBTLSCACert),
	})
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
