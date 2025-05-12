package bun

import (
	"context"
	"crypto/tls"
	"crypto/x509"
	"database/sql"
	"encoding/pem"
	"errors"
	"fmt"
	"os"
	"runtime"
	"strconv"
	"strings"

	"github.com/go-sql-driver/mysql"
	"github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/stdlib"
	"github.com/tyrm/mcp-dbmem/internal/config"
	"github.com/tyrm/mcp-dbmem/internal/db"
	"github.com/tyrm/mcp-dbmem/internal/db/bun/migrations"
	"github.com/uptrace/bun"
	"github.com/uptrace/bun/dialect"
	"github.com/uptrace/bun/dialect/mysqldialect"
	"github.com/uptrace/bun/dialect/pgdialect"
	"github.com/uptrace/bun/dialect/sqlitedialect"
	"github.com/uptrace/bun/extra/bunotel"
	"github.com/uptrace/bun/migrate"
	"go.uber.org/zap"
	"modernc.org/sqlite"
)

const (
	dbTypeMysql    = "mysql"
	dbTypePostgres = "postgres"
	dbTypeSqlite   = "sqlite"

	dbTLSModeDisable = "disable"
	dbTLSModeEnable  = "enable"
	dbTLSModeRequire = "require"
	dbTLSModeUnset   = ""
)

// Client is a DB interface compatible client for Bun.
type Client struct {
	db      *bun.DB
	errProc func(error) db.Error
}

var _ db.DB = (*Client)(nil)

type ClientConfig struct {
	Type      string
	Address   string
	Port      uint16
	User      string
	Password  string
	Database  string
	TLSMode   string
	TLSCACert string
}

// New creates a new bun database client
func New(ctx context.Context, c ClientConfig) (*Client, error) {
	var newBun *Client
	var err error
	dbType := c.Type

	switch dbType {
	case dbTypeMysql:
		newBun, err = myConn(ctx, c)
		if err != nil {
			return nil, err
		}
	case dbTypePostgres:
		newBun, err = pgConn(ctx, c)
		if err != nil {
			return nil, err
		}
	case dbTypeSqlite:
		newBun, err = sqliteConn(ctx, c)
		if err != nil {
			return nil, err
		}
	default:
		return nil, fmt.Errorf("database type %s not supported for bundb", dbType)
	}

	newBun.db.AddQueryHook(bunotel.NewQueryHook(bunotel.WithDBName(c.Database)))

	// Add a query hook to log all queries (debug)
	// newBun.db.AddQueryHook(bunzap.NewQueryHook(bunzap.QueryHookOptions{
	// 	Logger: zap.L(),
	// 	//SlowDuration: 200 * time.Millisecond, // Omit to log all operations as debug
	// }))

	return newBun, nil
}

// privates

func sqliteConn(ctx context.Context, c ClientConfig) (*Client, error) {
	// validate bun address has actually been set
	dbAddress := c.Address
	if dbAddress == "" {
		return nil, fmt.Errorf("'%s' was not set when attempting to start sqlite", c.Address)
	}

	// Drop anything fancy from DB address
	dbAddress = strings.Split(dbAddress, "?")[0]
	dbAddress = strings.TrimPrefix(dbAddress, "file:")

	// Append our own SQLite preferences
	dbAddress = "file:" + dbAddress + "?cache=shared"

	// Open new DB instance
	sqldb, err := sql.Open("sqlite", dbAddress)
	if err != nil {
		errWithCode := &sqlite.Error{}
		if errors.As(err, &errWithCode) {
			err = errors.New(sqlite.ErrorCodeString[errWithCode.Code()])
		}
		return nil, fmt.Errorf("could not open sqlite bun: %w", err)
	}

	setConnectionValues(sqldb)

	if dbAddress == "file::memory:?cache=shared" {
		zap.L().Warn("sqlite in-memory database should only be used for debugging")
		// don't close connections on disconnect -- otherwise
		// the SQLite database will be deleted when there
		// are no active connections
		sqldb.SetConnMaxLifetime(0)
	}

	conn := getErrConn(bun.NewDB(sqldb, sqlitedialect.New()))

	// ping to check the bun is there and listening
	if err := conn.db.PingContext(ctx); err != nil {
		errWithCode := &sqlite.Error{}
		if errors.As(err, &errWithCode) {
			err = errors.New(sqlite.ErrorCodeString[errWithCode.Code()])
		}
		return nil, fmt.Errorf("sqlite ping: %w", err)
	}

	zap.L().Info("Connected to database", zap.String("db_type", "sqlite"))
	return conn, nil
}

func myConn(ctx context.Context, c ClientConfig) (*Client, error) {
	opts, err := deriveBunDBMyOptions(c)
	if err != nil {
		return nil, fmt.Errorf("could not create bundb mysql options: %w", err)
	}

	sqldb, err := sql.Open("mysql", opts)
	if err != nil {
		return nil, fmt.Errorf("could not open mysql connection: %w", err)
	}

	setConnectionValues(sqldb)

	conn := getErrConn(bun.NewDB(sqldb, mysqldialect.New()))

	// ping to check the bun is there and listening
	if err := conn.db.PingContext(ctx); err != nil {
		return nil, fmt.Errorf("postgres ping: %w", err)
	}

	zap.L().Info("Connected to database", zap.String("db_type", "mysql"))
	return conn, nil
}

func pgConn(ctx context.Context, c ClientConfig) (*Client, error) {
	opts, err := deriveBunDBPGOptions(c)
	if err != nil {
		return nil, fmt.Errorf("could not doCreate bundb postgres options: %w", err)
	}

	sqldb := stdlib.OpenDB(*opts)

	setConnectionValues(sqldb)

	conn := getErrConn(bun.NewDB(sqldb, pgdialect.New()))

	// ping to check the bun is there and listening
	if err := conn.db.PingContext(ctx); err != nil {
		return nil, fmt.Errorf("postgres ping: %w", err)
	}

	zap.L().Info("Connected to database", zap.String("db_type", "postgres"))
	return conn, nil
}

func deriveBunDBMyOptions(c ClientConfig) (string, error) {
	// these are all optional, the bun adapter figures out defaults
	port := c.Port
	address := c.Address
	username := c.User
	password := c.Password

	// validate database
	database := c.Database
	if database == "" {
		return "", errors.New("no database set")
	}

	tlsConfig, err := makeTLSConfig(c)
	if err != nil {
		zap.L().Error("Error creating TLS config", zap.Error(err))
		return "", fmt.Errorf("could not create tls config: %w", err)
	}

	cfg := ""
	if username != "" {
		cfg += username
		if password != "" {
			cfg += ":" + password
		}
		cfg += "@"
	}
	if address != "" {
		cfg += "tcp(" + address
		if port > 0 {
			cfg += ":" + strconv.Itoa(int(port))
		}
		cfg += ")"
	}
	cfg += "/" + database

	// options
	if tlsConfig != nil {
		if err := mysql.RegisterTLSConfig("bun", tlsConfig); err != nil {
			return "", fmt.Errorf("could not register tls config: %w", err)
		}

		cfg += "?tls=bun"
	}

	return cfg, nil
}

func deriveBunDBPGOptions(c ClientConfig) (*pgx.ConnConfig, error) {
	// these are all optional, the bun adapter figures out defaults
	port := c.Port
	address := c.Address
	username := c.User
	password := c.Password

	// validate database
	database := c.Database
	if database == "" {
		return nil, errors.New("no database set")
	}

	tlsConfig, err := makeTLSConfig(c)
	if err != nil {
		zap.L().Error("Error creating TLS config", zap.Error(err))
		return nil, fmt.Errorf("could not create tls config: %w", err)
	}

	cfg, _ := pgx.ParseConfig("")
	if address != "" {
		cfg.Host = address
	}
	if port > 0 {
		cfg.Port = port
	}

	if username != "" {
		cfg.User = username
	}
	if password != "" {
		cfg.Password = password
	}
	if tlsConfig != nil {
		cfg.TLSConfig = tlsConfig
	}
	cfg.Database = database
	cfg.PreferSimpleProtocol = true
	cfg.RuntimeParams["application_name"] = config.ApplicationName

	return cfg, nil
}

// https://bun.uptrace.dev/postgres/running-bun-in-production.html#database-sql
func setConnectionValues(sqldb *sql.DB) {
	maxOpenConns := 4 * runtime.GOMAXPROCS(0)
	sqldb.SetMaxOpenConns(maxOpenConns)
	sqldb.SetMaxIdleConns(maxOpenConns)
}

func makeTLSConfig(c ClientConfig) (*tls.Config, error) {
	var tlsConfig *tls.Config
	tlsMode := c.TLSMode
	switch tlsMode {
	case dbTLSModeDisable, dbTLSModeUnset:
		break // nothing to do
	case dbTLSModeEnable:
		/* #nosec G402 */
		tlsConfig = &tls.Config{
			InsecureSkipVerify: true,
		}
	case dbTLSModeRequire:
		tlsConfig = &tls.Config{
			InsecureSkipVerify: false,
			ServerName:         c.Address,
			MinVersion:         tls.VersionTLS12,
		}
	}

	caCertPath := c.TLSCACert
	if tlsConfig != nil && caCertPath != "" {
		// load the system cert pool first -- we'll append the given CA cert to this
		certPool, err := x509.SystemCertPool()
		if err != nil {
			return nil, fmt.Errorf("error fetching system CA cert pool: %w", err)
		}

		// read the CA cert from the file
		caCert, err := readCertFile(caCertPath)
		if err != nil {
			return nil, fmt.Errorf("error reading CA cert file %s: %w", caCertPath, err)
		}

		// we're happy, add it to the existing pool and then use this pool in our tls config
		certPool.AddCert(caCert)
		tlsConfig.RootCAs = certPool
	}

	return tlsConfig, nil
}

func readCertFile(caCertPath string) (*x509.Certificate, error) {
	// open the file itself and make sure there's something in it
	/* #nosec G304 */
	caCertBytes, err := os.ReadFile(caCertPath)
	if err != nil {
		return nil, fmt.Errorf("error opening CA certificate at %s: %w", caCertPath, err)
	}
	if len(caCertBytes) == 0 {
		return nil, fmt.Errorf("ca cert at %s was empty", caCertPath)
	}

	// make sure we have a PEM block
	caPem, _ := pem.Decode(caCertBytes)
	if caPem == nil {
		return nil, fmt.Errorf("could not parse cert at %s into PEM", caCertPath)
	}

	// parse the PEM block into the certificate
	caCert, err := x509.ParseCertificate(caPem.Bytes)
	if err != nil {
		return nil, fmt.Errorf("could not parse cert at %s into x509 certificate: %w", caCertPath, err)
	}

	return caCert, nil
}

func getErrConn(dbConn *bun.DB) *Client {
	var errProc func(error) db.Error
	switch dbConn.Dialect().Name() {
	case dialect.PG:
		errProc = processPostgresError
	case dialect.SQLite:
		errProc = processSQLiteError
	default:
		panic("unknown dialect name: " + dbConn.Dialect().Name().String())
	}
	return &Client{
		errProc: errProc,
		db:      dbConn,
	}
}

// Close closes the bun db connection.
func (c *Client) Close() db.Error {
	zap.L().Info("Closing db connection", zap.String("db_dialect", c.db.Dialect().Name().String()))

	return c.db.Close()
}

// DoMigration runs schema migrations on the database.
func (c *Client) DoMigration(ctx context.Context) db.Error {
	migrator := migrate.NewMigrator(c.db, migrations.Migrations)

	if err := migrator.Init(ctx); err != nil {
		return err
	}

	group, err := migrator.Migrate(ctx)
	if err != nil {
		if err.Error() == "migrate: there are no any migrations" {
			zap.L().Info("No migrations to run", zap.String("db_dialect", c.db.Dialect().Name().String()))
			return nil
		}

		return err
	}

	if group.ID == 0 {
		zap.L().Info("No migrations to run", zap.String("db_dialect", c.db.Dialect().Name().String()))
		return nil
	}
	zap.L().Info("Migration successful", zap.String("db_dialect", c.db.Dialect().Name().String()), zap.String("group", group.String()))

	return nil
}
