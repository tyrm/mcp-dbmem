package v1

import (
	"context"

	mcp "github.com/metoro-io/mcp-golang"
	"github.com/tyrm/mcp-dbmem/internal/models"
	"go.uber.org/zap"
)

// Entity represents an entity in the knowledge graph
type Entity struct {
	Name         string   `json:"name" jsonschema:"required,description=The name of the entity"`
	Type         string   `json:"entityType" jsonschema:"required,description=The type of the entity"`
	Observations []string `json:"observations" jsonschema:"required,description=An array of observation contents associated with the entity"`
}

// CreateEntitiesArgs represents the arguments for creating entities
type CreateEntitiesArgs struct {
	Entities []Entity `json:"entities" jsonschema:"required,description=An array of observation contents associated with the entity"`
}

// CreateEntities creates entities in the knowledge graph
func (l *Logic) CreateEntities(ctx context.Context, args CreateEntitiesArgs) (*mcp.ToolResponse, error) {
	ctx, span := tracer.Start(ctx, "CreateEntities", tracerAttrs...)
	defer span.End()

	response := make([]Entity, 0)
	for _, entity := range args.Entities {
		// Process each entity
		newEntity := &models.Entity{
			Name: entity.Name,
			Type: entity.Type,
		}
		if err := l.DB.CreateEntity(ctx, newEntity); err != nil {
			zap.L().Error("Can't create entity from database", zap.Error(err), zap.Any("entity", newEntity))
			span.RecordError(err)
			return nil, err
		}
		newEntityResponse := Entity{
			Name:         newEntity.Name,
			Type:         newEntity.Type,
			Observations: make([]string, 0),
		}

		for _, observation := range entity.Observations {
			newObservation := &models.Observation{
				EntityID: newEntity.ID,
				Contents: observation,
			}
			if err := l.DB.CreateObservation(ctx, newObservation); err != nil {
				zap.L().Error("Can't create observation in database", zap.Error(err), zap.Any("observation", newObservation))
				span.RecordError(err)
				return nil, err
			}

			newEntityResponse.Observations = append(newEntityResponse.Observations, newObservation.Contents)
		}

		response = append(response, newEntityResponse)
	}

	// convert response to json string
	toolResponse, err := toolJSONResponse(ctx, response)
	if err != nil {
		zap.L().Error("Can't marshal response json", zap.Error(err))
		span.RecordError(err)
		return nil, err
	}

	return toolResponse, nil
}

// DeleteEntitiesArgs represents the arguments for deleting entities
type DeleteEntitiesArgs struct {
	EntityNames []string `json:"entityNames" jsonschema:"required,description=An array of entity names to delete"`
}

// DeleteEntities deletes entities from the knowledge graph
func (l *Logic) DeleteEntities(ctx context.Context, args DeleteEntitiesArgs) (*mcp.ToolResponse, error) {
	ctx, span := tracer.Start(ctx, "DeleteEntities", tracerAttrs...)
	defer span.End()

	for _, entityName := range args.EntityNames {
		// Process each entity
		entity, err := l.DB.ReadEntityByName(ctx, entityName)
		if err != nil {
			zap.L().Error("Can't read entity from database", zap.Error(err), zap.String("entityName", entityName))
			span.RecordError(err)
			return nil, err
		}
		if entity == nil {
			zap.L().Warn("Entity not found in database", zap.String("entityName", entityName))
			continue
		}

		if err := l.DB.DeleteAllObservationsByEntityID(ctx, entity.ID); err != nil {
			zap.L().Error("Can't delete observations", zap.Error(err), zap.String("entityName", entityName))
			span.RecordError(err)
			return nil, err
		}

		if err := l.DB.DeleteEntity(ctx, entity); err != nil {
			zap.L().Error("Can't delete entity", zap.Error(err), zap.String("entityName", entityName))
			span.RecordError(err)
			return nil, err
		}
	}

	return mcp.NewToolResponse(
		mcp.NewTextContent("Entities deleted successfully"),
	), nil
}
