package flag

import (
	"github.com/spf13/cobra"
	"github.com/tyrm/mcp-dbmem/internal/config"
)

// Server adds flags for the server command.
func Server(cmd *cobra.Command, values config.Values) {
	Database(cmd, values)
}
