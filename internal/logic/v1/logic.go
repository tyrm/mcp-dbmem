package v1

import (
	"context"

	"github.com/tyrm/mcp-dbmem/internal/db"
	"github.com/tyrm/mcp-dbmem/internal/logic"
	"github.com/tyrm/mcp-dbmem/internal/models"
)

// Logic implements the program logic.
type Logic struct {
	db db.DB
}

var _ logic.Logic = (*Logic)(nil)

// LogicConfig contains the configuration for the Logic instance.
type LogicConfig struct {
	DB db.DB
}

// NewLogic creates a new Logic instance.
func NewLogic(cfg LogicConfig) *Logic {
	return &Logic{
		db: cfg.DB,
	}
}

func (l *Logic) CreateEntity(ctx context.Context, entity *models.Entity) error {
	ctx, span := tracer.Start(ctx, "CreateEntity", tracerAttrs...)
	defer span.End()

	return logic.ProcessError(l.db.CreateEntity(ctx, entity))
}

func (l *Logic) DeleteEntity(ctx context.Context, entity *models.Entity) error {
	ctx, span := tracer.Start(ctx, "DeleteEntity", tracerAttrs...)
	defer span.End()

	return logic.ProcessError(l.db.DeleteEntity(ctx, entity))
}

func (l *Logic) ReadAllEntities(ctx context.Context) ([]*models.Entity, error) {
	ctx, span := tracer.Start(ctx, "ReadAllEntities", tracerAttrs...)
	defer span.End()

	entities, err := l.db.ReadAllEntities(ctx)
	if err != nil {
		return nil, logic.ProcessError(err)
	}
	return entities, nil
}

func (l *Logic) ReadEntityByName(ctx context.Context, name string) (*models.Entity, error) {
	ctx, span := tracer.Start(ctx, "ReadEntityByName", tracerAttrs...)
	defer span.End()

	entities, err := l.db.ReadEntityByName(ctx, name)
	if err != nil {
		return nil, logic.ProcessError(err)
	}
	return entities, nil
}

func (l *Logic) CreateObservation(ctx context.Context, observation *models.Observation) error {
	ctx, span := tracer.Start(ctx, "CreateObservation", tracerAttrs...)
	defer span.End()

	return logic.ProcessError(l.db.CreateObservation(ctx, observation))
}

func (l *Logic) DeleteAllObservationsByEntityID(ctx context.Context, entityID int64) error {
	ctx, span := tracer.Start(ctx, "DeleteAllObservationsByEntityID", tracerAttrs...)
	defer span.End()

	return logic.ProcessError(l.db.DeleteAllObservationsByEntityID(ctx, entityID))
}

func (l *Logic) DeleteObservation(ctx context.Context, observation *models.Observation) error {
	ctx, span := tracer.Start(ctx, "DeleteObservation", tracerAttrs...)
	defer span.End()

	return logic.ProcessError(l.db.DeleteObservation(ctx, observation))
}

func (l *Logic) ReadObservationByTextForEntityID(ctx context.Context, entityID int64, text string) (*models.Observation, error) {
	ctx, span := tracer.Start(ctx, "ReadObservationByTextForEntityID", tracerAttrs...)
	defer span.End()

	observations, err := l.db.ReadObservationByTextForEntityID(ctx, entityID, text)
	if err != nil {
		return nil, logic.ProcessError(err)
	}
	return observations, nil
}

func (l *Logic) CreateRelation(ctx context.Context, relation *models.Relation) error {
	ctx, span := tracer.Start(ctx, "CreateRelation", tracerAttrs...)
	defer span.End()

	return logic.ProcessError(l.db.CreateRelation(ctx, relation))
}

func (l *Logic) DeleteAllRelationsByEntityID(ctx context.Context, entityID int64) error {
	ctx, span := tracer.Start(ctx, "DeleteAllRelationsByEntityID", tracerAttrs...)
	defer span.End()

	return logic.ProcessError(l.db.DeleteAllRelationsByEntityID(ctx, entityID))
}

func (l *Logic) ReadAllRelations(ctx context.Context) ([]*models.Relation, error) {
	ctx, span := tracer.Start(ctx, "ReadAllRelations", tracerAttrs...)
	defer span.End()

	relations, err := l.db.ReadAllRelations(ctx)
	if err != nil {
		return nil, logic.ProcessError(err)
	}
	return relations, nil
}

func (l *Logic) ReadExactRelation(ctx context.Context, fromID, toID int64, relationType string) (*models.Relation, error) {
	ctx, span := tracer.Start(ctx, "ReadExactRelation", tracerAttrs...)
	defer span.End()

	relation, err := l.db.ReadExactRelation(ctx, fromID, toID, relationType)
	if err != nil {
		return nil, logic.ProcessError(err)
	}
	return relation, nil
}

func (l *Logic) DeleteRelation(ctx context.Context, relation *models.Relation) error {
	ctx, span := tracer.Start(ctx, "DeleteRelation", tracerAttrs...)
	defer span.End()

	return logic.ProcessError(l.db.DeleteRelation(ctx, relation))
}
