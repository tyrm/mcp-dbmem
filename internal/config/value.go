package config

// Values contains the type of each value.
type Values struct {
	LogLevel        string
	SoftwareVersion string

	// database
	DBType      string
	DBAddress   string
	DBPort      int
	DBUser      string
	DBPassword  string
	DBDatabase  string
	DBTLSMode   string
	DBTLSCACert string
}

// Defaults contains the default values.
var Defaults = Values{
	LogLevel: "info",

	// database
	DBType:      "postgres",
	DBAddress:   "localhost",
	DBPort:      5432,
	DBUser:      "mcp-pgmem",
	DBPassword:  "mcp-pgmem",
	DBDatabase:  "mcp-pgmem",
	DBTLSMode:   "disable",
	DBTLSCACert: "",
}
