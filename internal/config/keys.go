package config

// KeyNames is a struct that contains the names of keys.
type KeyNames struct {
	LogLevel        string
	SoftwareVersion string

	// database
	DBType      string
	DBAddress   string
	DBPort      string
	DBUser      string
	DBPassword  string
	DBDatabase  string
	DBTLSMode   string
	DBTLSCACert string
}

// Keys contains the names of config keys.
var Keys = KeyNames{
	LogLevel:        "log-level",
	SoftwareVersion: "software-version", // Set at build

	// database
	DBType:      "db-type",
	DBAddress:   "db-address",
	DBPort:      "db-port",
	DBUser:      "db-user",
	DBPassword:  "db-password",
	DBDatabase:  "db-database",
	DBTLSMode:   "db-tls-mode",
	DBTLSCACert: "db-tls-ca-cert",
}
