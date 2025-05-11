package main

import (
	"context"
	"fmt"
	"log"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/tyrm/mcp-dbmem/cmd/mcp_pgmem/action"
	"github.com/tyrm/mcp-dbmem/cmd/mcp_pgmem/action/migrate"
	"github.com/tyrm/mcp-dbmem/cmd/mcp_pgmem/action/start"
	"github.com/tyrm/mcp-dbmem/cmd/mcp_pgmem/flag"
	"github.com/tyrm/mcp-dbmem/internal/config"
	"go.uber.org/zap"
)

// Version is the software version.
var Version string

// Commit is the git commit.
var Commit string

func main() {
	// init logger
	zapLogger, err := zap.NewDevelopment()
	if err != nil {
		log.Fatal(err)
		return
	}
	defer func() {
		_ = zapLogger.Sync()
	}()
	zap.ReplaceGlobals(zapLogger)

	// set software version
	var v string
	if len(Commit) < 7 {
		v = "v" + Version
	} else {
		v = "v" + Version + "-" + Commit[:7]
	}

	viper.Set(config.Keys.SoftwareVersion, v)

	rootCmd := &cobra.Command{
		Use:           "mcp-pgmem",
		Short:         "", //TODO
		Version:       v,
		SilenceErrors: true,
		SilenceUsage:  true,
	}

	flag.Global(rootCmd, config.Defaults)

	serverStartCmd := &cobra.Command{
		Use:   "start",
		Short: "start the mcp server",
		PreRunE: func(cmd *cobra.Command, args []string) error {
			return preRun(cmd)
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			return run(cmd.Context(), start.Start, args)
		},
	}
	rootCmd.AddCommand(serverStartCmd)

	databaseMigrateCmd := &cobra.Command{
		Use:   "migrate",
		Short: "run db migrations",
		PreRunE: func(cmd *cobra.Command, args []string) error {
			return preRun(cmd)
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			return run(cmd.Context(), migrate.Migrate, args)
		},
	}
	rootCmd.AddCommand(databaseMigrateCmd)

	err = rootCmd.Execute()
	if err != nil {
		zap.L().Fatal("Error executing command", zap.Error(err))
	}
}

func preRun(cmd *cobra.Command) error {
	if err := config.Init(cmd.Flags()); err != nil {
		return fmt.Errorf("error initializing config: %s", err)
	}

	return nil
}

func run(ctx context.Context, action action.Action, args []string) error {
	return action(ctx, args)
}
