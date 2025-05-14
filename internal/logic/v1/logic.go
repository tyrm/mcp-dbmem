package v1

import (
	"context"

	"github.com/tyrm/mcp-dbmem/internal/db"
	"github.com/tyrm/mcp-dbmem/internal/logic"
	"github.com/tyrm/mcp-dbmem/internal/models"
)

// Logic implements the program logic.
type Logic struct {
	DB db.DB
}

func (l *Logic) CreateEntity(ctx context.Context, entity *models.Entity) error {
	//TODO implement me
	panic("implement me")
}

func (l *Logic) DeleteEntity(ctx context.Context, entity *models.Entity) error {
	//TODO implement me
	panic("implement me")
}

func (l *Logic) ReadAllEntities(ctx context.Context) ([]*models.Entity, error) {
	//TODO implement me
	panic("implement me")
}

func (l *Logic) ReadEntityByName(ctx context.Context, name string) (*models.Entity, error) {
	//TODO implement me
	panic("implement me")
}

func (l *Logic) CreateObservation(ctx context.Context, observation *models.Observation) error {
	//TODO implement me
	panic("implement me")
}

func (l *Logic) DeleteAllObservationsByEntityID(ctx context.Context, entityID int64) error {
	//TODO implement me
	panic("implement me")
}

func (l *Logic) DeleteObservation(ctx context.Context, observation *models.Observation) error {
	//TODO implement me
	panic("implement me")
}

func (l *Logic) ReadObservationByTextForEntityID(ctx context.Context, entityID int64, text string) (*models.Observation, error) {
	//TODO implement me
	panic("implement me")
}

func (l *Logic) CreateRelation(ctx context.Context, relation *models.Relation) error {
	//TODO implement me
	panic("implement me")
}

func (l *Logic) DeleteAllRelationsByEntityID(ctx context.Context, entityID int64) error {
	//TODO implement me
	panic("implement me")
}

func (l *Logic) ReadAllRelations(ctx context.Context) ([]*models.Relation, error) {
	//TODO implement me
	panic("implement me")
}

func (l *Logic) ReadExactRelation(ctx context.Context, fromID, toID int64, relationType string) (*models.Relation, error) {
	//TODO implement me
	panic("implement me")
}

func (l *Logic) DeleteRelation(ctx context.Context, relation *models.Relation) error {
	//TODO implement me
	panic("implement me")
}

var _ logic.Logic = (*Logic)(nil)

//func toolJSONResponse(ctx context.Context, response any) (*mcp.ToolResponse, error) {
//	_, span := tracer.Start(ctx, "toolJSONResponse", tracerAttrs...)
//	defer span.End()
//
//	// convert response to json string
//	jsonResponse, err := json.MarshalIndent(response, "", "  ")
//	if err != nil {
//		span.RecordError(err)
//		return nil, err
//	}
//
//	return mcp.NewToolResponse(
//		mcp.NewTextContent(string(jsonResponse)),
//	), nil
//}
