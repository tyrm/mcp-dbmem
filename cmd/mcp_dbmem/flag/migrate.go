package flag

import (
	"github.com/spf13/cobra"
	"github.com/tyrm/mcp-dbmem/internal/config"
)

// Migrate adds flags for the migrate command.
func Migrate(cmd *cobra.Command, values config.Values) {
	Database(cmd, values)
}
