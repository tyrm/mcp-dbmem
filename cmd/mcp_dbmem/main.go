package main

import (
	"context"
	"fmt"
	"log"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/tyrm/mcp-dbmem/cmd/mcp_dbmem/action"
	"github.com/tyrm/mcp-dbmem/cmd/mcp_dbmem/action/direct"
	"github.com/tyrm/mcp-dbmem/cmd/mcp_dbmem/action/migrate"
	"github.com/tyrm/mcp-dbmem/cmd/mcp_dbmem/flag"
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
		Use:           "mcp-dbmem",
		Short:         "", //TODO
		Version:       v,
		SilenceErrors: true,
		SilenceUsage:  true,
	}
	flag.Global(rootCmd, config.Defaults)

	directCmd := &cobra.Command{
		Use:   "direct",
		Short: "the mcp server will connect directly to the database",
		PreRunE: func(cmd *cobra.Command, args []string) error {
			return preRun(cmd)
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			return run(cmd.Context(), direct.Direct, args)
		},
	}
	flag.Direct(directCmd, config.Defaults)
	rootCmd.AddCommand(directCmd)

	migrateCmd := &cobra.Command{
		Use:   "migrate",
		Short: "run db migrations",
		PreRunE: func(cmd *cobra.Command, args []string) error {
			return preRun(cmd)
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			return run(cmd.Context(), migrate.Migrate, args)
		},
	}
	flag.Migrate(migrateCmd, config.Defaults)
	rootCmd.AddCommand(migrateCmd)

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
