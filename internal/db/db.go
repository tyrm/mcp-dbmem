package db

import (
	"context"

	"github.com/tyrm/mcp-dbmem/internal/models"
)

// DB is the interface that wraps the basic database operations.
type DB interface {
	Entities
	Observations
	Relations
}

// Entities is the interface that wraps the basic entity operations.
type Entities interface {
	// CreateEntity creates a new entity and observations in the database.
	CreateEntity(ctx context.Context, entity *models.Entity) Error
	// DeleteEntity deletes an entity and all its associated observations and relations.
	DeleteEntity(ctx context.Context, entity *models.Entity) Error
	// ReadAllEntities reads all entities and their observations from the database.
	ReadAllEntities(ctx context.Context) ([]*models.Entity, Error)
	// ReadEntityByName reads an entity by its name and returns the entity and its observations.
	ReadEntityByName(ctx context.Context, name string) (*models.Entity, Error)
}

// Observations is the interface that wraps the basic observation operations.
type Observations interface {
	// CreateObservation creates a new observation in the database.
	CreateObservation(ctx context.Context, observation *models.Observation) Error
	//DeleteAllObservationsByEntityID(ctx context.Context, entityID int64) Error

	// DeleteObservation deletes an observation from the database.
	DeleteObservation(ctx context.Context, observation *models.Observation) Error
	// ReadObservationByTextForEntityID reads an observation by its text and entity ID.
	ReadObservationByTextForEntityID(ctx context.Context, entityID int64, text string) (*models.Observation, Error)
}

// Relations is the interface that wraps the basic relation operations.
type Relations interface {
	// CreateRelation creates a new relation in the database.
	CreateRelation(ctx context.Context, relation *models.Relation) Error
	//DeleteAllRelationsByEntityID(ctx context.Context, entityID int64) Error

	// ReadAllRelations reads all relations from the database.
	ReadAllRelations(ctx context.Context) ([]*models.Relation, Error)
	// ReadExactRelation reads a relation by its from entity, to entity, and relation type.
	ReadExactRelation(ctx context.Context, fromID, toID int64, relationType string) (*models.Relation, Error)
	// DeleteRelation deletes a relation from the database.
	DeleteRelation(ctx context.Context, relation *models.Relation) Error
}
