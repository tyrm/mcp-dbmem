package flag

import "github.com/tyrm/mcp-dbmem/internal/config"

var usage = config.KeyNames{
	LogLevel:        "Log level",
	SoftwareVersion: "Software version",
	DBType:          "Database type [postgres, sqlite]",
	DBAddress:       "Database address",
	DBPort:          "Database port",
	DBUser:          "Database user",
	DBPassword:      "Database password",
	DBDatabase:      "Database name",
	DBTLSMode:       "Database TLS mode",
	DBTLSCACert:     "Database TLS CA certificate",
}
