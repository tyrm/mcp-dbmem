package logic

import "errors"

var (
	// ErrNotFound is returned when an entity is not found in the database.
	ErrNotFound = errors.New("not found")
)
