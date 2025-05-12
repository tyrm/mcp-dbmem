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

type Entities interface {
	CreateEntity(ctx context.Context, entity *models.Entity) Error
	DeleteEntity(ctx context.Context, entity *models.Entity) Error
	ReadAllEntities(ctx context.Context) ([]*models.Entity, Error)
	ReadEntityByName(ctx context.Context, name string) (*models.Entity, Error)
}

type Observations interface {
	CreateObservation(ctx context.Context, observation *models.Observation) Error
	DeleteAllObservationsByEntityID(ctx context.Context, entityID int64) Error
	DeleteObservation(ctx context.Context, observation *models.Observation) Error
	ReadObservationByTextForEntityID(ctx context.Context, entityID int64, text string) (*models.Observation, Error)
}

type Relations interface {
	CreateRelation(ctx context.Context, relation *models.Relation) Error
	DeleteAllRelationsByEntityID(ctx context.Context, entityID int64) Error
	ReadAllRelations(ctx context.Context) ([]*models.Relation, Error)
	ReadExactRelation(ctx context.Context, fromID, toID int64, relationType string) (*models.Relation, Error)
	DeleteRelation(ctx context.Context, relation *models.Relation) Error
}
