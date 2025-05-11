package bun

import (
	"database/sql"
	"errors"

	"github.com/jackc/pgconn"
	"github.com/tyrm/mcp-dbmem/internal/db"
	"go.uber.org/zap"
	"modernc.org/sqlite"
	sqlite3 "modernc.org/sqlite/lib"
)

// AlreadyExistsError is returned when a caller tries to insert a database entry that already exists in the db.
type AlreadyExistsError struct {
	message string
}

// Error returns the error message as a string.
func (e *AlreadyExistsError) Error() string {
	return e.message
}

// NewErrAlreadyExists wraps a message in an AlreadyExistsError object.
func NewErrAlreadyExists(msg string) error {
	return &AlreadyExistsError{message: msg}
}

// ProcessError replaces any known values with our own db.Error types
func (c *Client) ProcessError(err error) db.Error {
	switch {
	case err == nil:
		return nil
	case errors.Is(err, sql.ErrNoRows):
		return db.ErrNoEntries
	default:
		return c.errProc(err)
	}
}

// processPostgresError processes an error, replacing any postgres specific errors with our own error type
func processPostgresError(err error) db.Error {
	// Attempt to cast as postgres
	var pgErr *pgconn.PgError
	ok := errors.As(err, &pgErr)
	if !ok {
		return err
	}

	zap.L().Debug("Postgres error", zap.String("code", pgErr.Code), zap.Error(pgErr))

	// Handle supplied error code:
	// (https://www.postgresql.org/docs/10/errcodes-appendix.html)
	switch pgErr.Code {
	case "23505" /* unique_violation */ :
		return NewErrAlreadyExists(pgErr.Message)
	default:
		return err
	}
}

// processSQLiteError processes an error, replacing any sqlite specific errors with our own error type

func processSQLiteError(err error) db.Error {
	// Attempt to cast as sqlite
	var sqliteErr *sqlite.Error
	ok := errors.As(err, &sqliteErr)
	if !ok {
		return err
	}

	zap.L().Debug("Postgres error", zap.Int("code", sqliteErr.Code()), zap.Error(sqliteErr))
	// Handle supplied error code:
	switch sqliteErr.Code() {
	case sqlite3.SQLITE_CONSTRAINT_UNIQUE, sqlite3.SQLITE_CONSTRAINT_PRIMARYKEY:
		return NewErrAlreadyExists(err.Error())
	default:
		return err
	}
}
