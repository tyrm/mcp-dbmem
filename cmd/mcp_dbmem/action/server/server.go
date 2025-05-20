package server

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	mcp "github.com/metoro-io/mcp-golang"
	"github.com/metoro-io/mcp-golang/transport/stdio"
	"github.com/spf13/viper"
	"github.com/tyrm/mcp-dbmem/cmd/mcp_dbmem/action"
	"github.com/tyrm/mcp-dbmem/internal/adapter"
	"github.com/tyrm/mcp-dbmem/internal/config"
	"github.com/tyrm/mcp-dbmem/internal/db/bun"
	v1 "github.com/tyrm/mcp-dbmem/internal/logic/v1"
	"github.com/uptrace/uptrace-go/uptrace"
	"go.uber.org/zap"
)

// Server is the action to start the mcp server with a direct connection to the database.
var Server action.Action = func(ctx context.Context, _ []string) error {
	zap.L().Info(fmt.Sprintf("starting %s server", config.ApplicationName))

	// Setup tracing
	if viper.GetString(config.Keys.UptraceDSN) != "" {
		uptrace.ConfigureOpentelemetry(
			uptrace.WithServiceName(config.ApplicationName),
			uptrace.WithServiceVersion(viper.GetString(config.Keys.SoftwareVersion)),
			uptrace.WithDSN(viper.GetString(config.Keys.UptraceDSN)),
		)
		// Send buffered spans and free resources.
		defer func() {
			if err := uptrace.Shutdown(context.Background()); err != nil {
				zap.L().Error("Error shutting down uptrace", zap.Error(err))
			}
		}()
	}

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
		if err := dbClient.Close(); err != nil {
			zap.L().Error("Error closing bun client", zap.Error(err))
		}
	}()

	// build logic
	logic := v1.NewLogic(v1.LogicConfig{
		DB: dbClient,
	})

	direct := adapter.NewServerAdapter(logic)

	// add tools
	server := mcp.NewServer(stdio.NewStdioServerTransport())
	if err := direct.Apply(server); err != nil {
		zap.L().Error("Error applying direct adapter", zap.Error(err))

		return err
	}

	// ** start application **
	errChan := make(chan error)

	// Wait for SIGINT and SIGTERM (HIT CTRL-C)
	stopSigChan := make(chan os.Signal, 1)
	signal.Notify(stopSigChan, syscall.SIGINT, syscall.SIGTERM)

	// start mcp server
	go func(s *mcp.Server, errChan chan error) {
		zap.L().Info("starting mcp server")
		err := s.Serve()
		if err != nil {
			errChan <- fmt.Errorf("mcp server: %s", err.Error())
		}
	}(server, errChan)

	// wait for event
	select {
	case sig := <-stopSigChan:
		zap.L().Info("got signal", zap.String("signal", sig.String()))
	case err := <-errChan:
		zap.L().Fatal("fatal error", zap.Error(err))
	}

	zap.L().Info("done")
	return nil
}
