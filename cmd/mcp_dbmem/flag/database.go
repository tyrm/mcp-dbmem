package flag

import (
	"github.com/spf13/cobra"
	"github.com/tyrm/mcp-dbmem/internal/config"
)

// Database adds flags for database configuration.
func Database(cmd *cobra.Command, values config.Values) {
	cmd.PersistentFlags().String(config.Keys.DBType, values.DBType, usage.DBType)
	cmd.PersistentFlags().String(config.Keys.DBAddress, values.DBAddress, usage.DBAddress)
	cmd.PersistentFlags().Uint16(config.Keys.DBPort, values.DBPort, usage.DBPort)
	cmd.PersistentFlags().String(config.Keys.DBUser, values.DBUser, usage.DBUser)
	cmd.PersistentFlags().String(config.Keys.DBPassword, values.DBPassword, usage.DBPassword)
	cmd.PersistentFlags().String(config.Keys.DBDatabase, values.DBDatabase, usage.DBDatabase)
	cmd.PersistentFlags().String(config.Keys.DBTLSMode, values.DBTLSMode, usage.DBTLSMode)
	cmd.PersistentFlags().String(config.Keys.DBTLSCACert, values.DBTLSCACert, usage.DBTLSCACert)
}
