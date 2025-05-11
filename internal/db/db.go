package db

import (
	"context"
	"fmt"

	"github.com/tyrm/mcp-pgmem/internal/models"
)

// DB is the interface that wraps the basic database operations.
type DB interface {
	CreateEntity(ctx context.Context, entity *models.Entity) Error
	DeleteEntity(ctx context.Context, entity *models.Entity) Error
	ReadAllEntities(ctx context.Context) ([]*models.Entity, Error)
	ReadEntityByName(ctx context.Context, name string) (*models.Entity, Error)

	CreateObservation(ctx context.Context, observation *models.Observation) Error
	DeleteAllObservationsByEntityID(ctx context.Context, entityID int64) Error
	DeleteObservation(ctx context.Context, observation *models.Observation) Error
	ReadObservationByTextForEntityID(ctx context.Context, entityID int64, text string) (*models.Observation, Error)

	CreateRelation(ctx context.Context, relation *models.Relation) Error
	DeleteAllRelationsByEntityID(ctx context.Context, entityID int64) Error
	ReadAllRelations(ctx context.Context) ([]*models.Relation, Error)
	ReadExactRelation(ctx context.Context, fromID, toID int64, relationType string) (*models.Relation, Error)
	DeleteRelation(ctx context.Context, relation *models.Relation) Error
}

// Error represents a database specific error.
type Error error

var (
	// ErrGenID is returned when creating a new ID can't be generated for a new model.
	ErrGenID Error = fmt.Errorf("can't generate id")
	// ErrNoEntries is returned when a caller expected an entry for a query, but none was found.
	ErrNoEntries Error = fmt.Errorf("no entries")
	// ErrMultipleEntries is returned when a caller expected ONE entry for a query, but multiples were found.
	ErrMultipleEntries Error = fmt.Errorf("multiple entries")
	// ErrUnknown denotes an unknown database error.
	ErrUnknown Error = fmt.Errorf("unknown error")
	// ErrInvalidSort is returned when a sort type is requested the model can't do.
	ErrInvalidSort Error = fmt.Errorf("invalid sort")
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
