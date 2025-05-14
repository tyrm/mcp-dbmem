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

type Entities interface {
	CreateEntity(ctx context.Context, entity *models.Entity) error
	DeleteEntity(ctx context.Context, entity *models.Entity) error
	ReadAllEntities(ctx context.Context) ([]*models.Entity, error)
	ReadEntityByName(ctx context.Context, name string) (*models.Entity, error)
}

type Observations interface {
	CreateObservation(ctx context.Context, observation *models.Observation) error
	DeleteAllObservationsByEntityID(ctx context.Context, entityID int64) error
	DeleteObservation(ctx context.Context, observation *models.Observation) error
	ReadObservationByTextForEntityID(ctx context.Context, entityID int64, text string) (*models.Observation, error)
}

type Relations interface {
	CreateRelation(ctx context.Context, relation *models.Relation) error
	DeleteAllRelationsByEntityID(ctx context.Context, entityID int64) error
	ReadAllRelations(ctx context.Context) ([]*models.Relation, error)
	ReadExactRelation(ctx context.Context, fromID, toID int64, relationType string) (*models.Relation, error)
	DeleteRelation(ctx context.Context, relation *models.Relation) error
}
