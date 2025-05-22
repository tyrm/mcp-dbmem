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
	"github.com/tyrm/mcp-dbmem/cmd/mcp_dbmem/action/server"
	"github.com/tyrm/mcp-dbmem/cmd/mcp_dbmem/flag"
	"github.com/tyrm/mcp-dbmem/internal/config"
	"go.uber.org/zap"
)

const minimumHashLength = 7

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
	var version string
	if len(Commit) < minimumHashLength {
		version = "version" + Version
	} else {
		version = "version" + Version + "-" + Commit[:7]
	}

	viper.Set(config.Keys.SoftwareVersion, version)

	rootCmd := &cobra.Command{
		Use:           config.ApplicationName,
		Short:         "", // TODO
		Version:       version,
		SilenceErrors: true,
		SilenceUsage:  true,
	}
	flag.Global(rootCmd, config.Defaults)

	rootCmd.AddCommand(directCommands())
	rootCmd.AddCommand(migrateCommands())
	rootCmd.AddCommand(serverCommands())

	err = rootCmd.Execute()
	if err != nil {
		zap.L().Fatal("Error executing command", zap.Error(err))
	}
}

func preRun(cmd *cobra.Command) error {
	if err := config.Init(cmd.Flags()); err != nil {
		return fmt.Errorf("error initializing config: %w", err)
	}

	return nil
}

func run(ctx context.Context, action action.Action, args []string) error {
	return action(ctx, args)
}

// directCommands returns the 'direct' subcommand.
func directCommands() *cobra.Command {
	rootCmd := &cobra.Command{
		Use:   "direct",
		Short: "the mcp server will connect directly to the database",
		PreRunE: func(cmd *cobra.Command, args []string) error {
			return preRun(cmd)
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			return run(cmd.Context(), direct.Direct, args)
		},
	}
	flag.Direct(rootCmd, config.Defaults)

	return rootCmd
}

// migrateCommands returns the 'migrate' subcommand.
func migrateCommands() *cobra.Command {
	rootCmd := &cobra.Command{
		Use:   "migrate",
		Short: "run db migrations",
		PreRunE: func(cmd *cobra.Command, _ []string) error {
			return preRun(cmd)
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			return run(cmd.Context(), migrate.Migrate, args)
		},
	}
	flag.Migrate(rootCmd, config.Defaults)

	return rootCmd
}

// serverCommands returns the 'server' subcommand.
func serverCommands() *cobra.Command {
	serverCmd := &cobra.Command{
		Use:   "server",
		Short: "start the server",
		PreRunE: func(cmd *cobra.Command, args []string) error {
			return preRun(cmd)
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			return run(cmd.Context(), server.Server, args)
		},
	}
	flag.Server(serverCmd, config.Defaults)

	return serverCmd
}
