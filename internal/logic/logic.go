package logic

import (
	"context"

	"github.com/tyrm/mcp-dbmem/internal/models"
)

// Error represents a database specific error.
type Error error

type Logic interface {
	Entities
	Observations
	Relations
}

// Entities is the interface that wraps the basic entity operations.
type Entities interface {
	// CreateEntity creates a new entity and observations.
	CreateEntity(ctx context.Context, entity *models.Entity) error
	// DeleteEntity deletes an entity and all its associated observations and relations.
	DeleteEntity(ctx context.Context, entity *models.Entity) error
	// ReadAllEntities reads all entities and their observations.
	ReadAllEntities(ctx context.Context) ([]*models.Entity, error)
	// ReadEntityByName reads an entity by its name and returns the entity and its observations.
	ReadEntityByName(ctx context.Context, name string) (*models.Entity, error)
}

// Observations is the interface that wraps the basic observation operations.
type Observations interface {
	// CreateObservation creates a new observation.
	CreateObservation(ctx context.Context, observation *models.Observation) error
	//DeleteAllObservationsByEntityID(ctx context.Context, entityID int64) error

	// DeleteObservation deletes an observation.
	DeleteObservation(ctx context.Context, observation *models.Observation) error
	// ReadObservationByTextForEntityID reads an observation by its text and entity ID.
	ReadObservationByTextForEntityID(ctx context.Context, entityID int64, text string) (*models.Observation, error)
}

// Relations is the interface that wraps the basic relation operations.
type Relations interface {
	// CreateRelation creates a new relation.
	CreateRelation(ctx context.Context, relation *models.Relation) error
	//DeleteAllRelationsByEntityID(ctx context.Context, entityID int64) error

	// ReadAllRelations reads all relations.
	ReadAllRelations(ctx context.Context) ([]*models.Relation, error)
	// ReadExactRelation reads a relation by its from entity, to entity, and relation type.
	ReadExactRelation(ctx context.Context, fromID, toID int64, relationType string) (*models.Relation, error)
	// DeleteRelation deletes a relation.
	DeleteRelation(ctx context.Context, relation *models.Relation) error
}
