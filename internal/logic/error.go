package logic

import (
	"errors"

	"github.com/tyrm/mcp-dbmem/internal/db"
)

var (
	// ErrNotFound is returned when an entity is not found in the database.
	ErrNotFound = errors.New("not found")
)

// ProcessError replaces any known values with our own db.Error types.
func ProcessError(err error) Error {
	switch {
	case err == nil:
		return nil
	case errors.Is(err, db.ErrNoEntries):
		return ErrNotFound
	default:
		return err
	}
}
