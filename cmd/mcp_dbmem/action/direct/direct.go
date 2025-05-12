package direct

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
	"github.com/tyrm/mcp-dbmem/internal/config"
	"github.com/tyrm/mcp-dbmem/internal/db/bun"
	v1 "github.com/tyrm/mcp-dbmem/internal/logic/v1"
	"github.com/uptrace/uptrace-go/uptrace"
	"go.uber.org/zap"
)

// Direct is the action to start the mcp server with a direct connection to the database.
var Direct action.Action = func(ctx context.Context, _ []string) error {
	zap.L().Info("starting pgmcp")

	// Setup tracing
	if viper.GetString(config.Keys.UptraceDSN) != "" {
		uptrace.ConfigureOpentelemetry(
			uptrace.WithServiceName("mcp-dbmem"),
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
		err := dbClient.Close()
		if err != nil {
			zap.L().Error("Error closing bun client", zap.Error(err))
		}
	}()

	// build logic
	logic := v1.Logic{
		DB: dbClient,
	}

	// add tools
	server := mcp.NewServer(stdio.NewStdioServerTransport())
	if err := server.RegisterTool("create_entities", "Create multiple new entities in the knowledge graph", logic.CreateEntities); err != nil {
		return err
	}
	if err := server.RegisterTool("create_relations", "Create multiple new relations between entities in the knowledge graph. Relations should be in active voice", logic.CreateRelations); err != nil {
		return err
	}
	if err := server.RegisterTool("add_observations", "Add new observations to existing entities in the knowledge graph", logic.AddObservations); err != nil {
		return err
	}
	if err := server.RegisterTool("delete_entities", "Delete multiple entities and their associated relations from the knowledge graph", logic.DeleteEntities); err != nil {
		return err
	}
	if err := server.RegisterTool("delete_observations", "Delete specific observations from entities in the knowledge graph", logic.DeleteObservations); err != nil {
		return err
	}
	if err := server.RegisterTool("delete_relations", "Delete multiple relations from the knowledge graph", logic.DeleteRelations); err != nil {
		return err
	}
	if err := server.RegisterTool("read_graph", "Read the entire knowledge graph", logic.ReadGraph); err != nil {
		return err
	}
	if err := server.RegisterTool("search_nodes", "Search for nodes in the knowledge graph based on a query", logic.SearchNodes); err != nil {
		return err
	}
	if err := server.RegisterTool("open_nodes", "Open specific nodes in the knowledge graph by their names", logic.OpenNodes); err != nil {
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
