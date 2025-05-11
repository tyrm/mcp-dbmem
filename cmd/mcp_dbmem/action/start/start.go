package start

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	mcp "github.com/metoro-io/mcp-golang"
	"github.com/metoro-io/mcp-golang/transport/stdio"
	"github.com/tyrm/mcp-dbmem/cmd/mcp_dbmem/action"
	"github.com/tyrm/mcp-dbmem/internal/db/bun"
	v1 "github.com/tyrm/mcp-dbmem/internal/logic/v1"
	"go.uber.org/zap"
)

var Start action.Action = func(ctx context.Context, _ []string) error {
	zap.L().Info("starting pgmcp")

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

	ctx, cancel := context.WithCancel(ctx)

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
	cancel()
	return nil
}
