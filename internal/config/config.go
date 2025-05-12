package config

import (
	"strings"

	"github.com/spf13/pflag"
	"github.com/spf13/viper"
)

// ApplicationName contains the application identifier.
const ApplicationName = "mcp-dbmem"

// Init starts config collection.
func Init(flags *pflag.FlagSet) error {
	viper.SetEnvKeyReplacer(strings.NewReplacer("-", "_"))
	viper.AutomaticEnv()

	if err := viper.BindPFlags(flags); err != nil {
		return err
	}

	return nil
}
