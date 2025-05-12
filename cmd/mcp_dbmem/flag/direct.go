package flag

import (
	"github.com/spf13/cobra"
	"github.com/tyrm/mcp-dbmem/internal/config"
)

func Direct(cmd *cobra.Command, values config.Values) {
	Database(cmd, values)
}
