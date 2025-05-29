package server

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/spf13/viper"
	"github.com/tyrm/mcp-dbmem/cmd/mcp_dbmem/action"
	"github.com/tyrm/mcp-dbmem/internal/config"
	"github.com/tyrm/mcp-dbmem/internal/db/bun"
	"github.com/tyrm/mcp-dbmem/internal/http"
	"github.com/tyrm/mcp-dbmem/internal/http/apiv1"
	v1 "github.com/tyrm/mcp-dbmem/internal/logic/v1"
	"github.com/uptrace/uptrace-go/uptrace"
	"go.uber.org/zap"
)

// Server is the action to start the mcp server with a direct connection to the database.
var Server action.Action = func(ctx context.Context, _ []string) error {
	zap.L().Info(fmt.Sprintf("starting %s server", config.ApplicationName))
	ctx, cancel := context.WithCancel(ctx)

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
		cancel()

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

	httpServer, err := http.NewServer(ctx, http.ServerConfig{
		Logic:           logic,
		ApplicationName: config.ApplicationName,
		HttpBind:        ":4200",
	})
	if err != nil {
		zap.L().Error("can't start http server", zap.Error(err))
		cancel()
		return err
	}
	// create web modules
	var webModules = make([]http.Module, 0)

	zap.L().Info("loading apiv1 module")
	apiV1 := apiv1.New(logic)
	webModules = append(webModules, apiV1)

	// add modules to server
	for _, mod := range webModules {
		mod.Route(httpServer)
	}
	// ** start application **
	errChan := make(chan error)

	// Wait for SIGINT and SIGTERM (HIT CTRL-C)
	stopSigChan := make(chan os.Signal, 1)
	signal.Notify(stopSigChan, syscall.SIGINT, syscall.SIGTERM)

	// start webserver
	go func(s *http.Server, errChan chan error) {
		zap.L().Info("starting http server")
		err := s.Start()
		if err != nil {
			errChan <- fmt.Errorf("http server: %s", err.Error())
		}
	}(httpServer, errChan)

	// wait for event
	select {
	case sig := <-stopSigChan:
		zap.L().Info("got signal", zap.String("signal", sig.String()))
	case err := <-errChan:
		zap.L().Fatal("fatal error", zap.Error(err))
	}

	cancel()
	zap.L().Info("done")
	return nil
}
