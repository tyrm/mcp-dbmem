package db

import (
	"errors"
)

// Error represents a database specific error.
type Error error

var (
	// ErrGenID is returned when creating a new ID can't be generated for a new model.
	ErrGenID Error = errors.New("can't generate id")
	// ErrNoEntries is returned when a caller expected an entry for a query, but none was found.
	ErrNoEntries Error = errors.New("no entries")
	// ErrMultipleEntries is returned when a caller expected ONE entry for a query, but multiples were found.
	ErrMultipleEntries Error = errors.New("multiple entries")
	// ErrUnknown denotes an unknown database error.
	ErrUnknown Error = errors.New("unknown error")
	// ErrInvalidSort is returned when a sort type is requested the model can't do.
	ErrInvalidSort Error = errors.New("invalid sort")
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
