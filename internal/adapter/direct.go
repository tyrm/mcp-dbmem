package adapter

import (
	"context"
	"errors"
	"fmt"

	mcp "github.com/metoro-io/mcp-golang"
	"github.com/tyrm/mcp-dbmem/internal/db"
	"github.com/tyrm/mcp-dbmem/internal/logic"
	"github.com/tyrm/mcp-dbmem/internal/models"
	"github.com/tyrm/mcp-dbmem/internal/util"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/trace"
	"go.uber.org/zap"
)

var directTracer = otel.Tracer("internal/adapter.DirectAdapter")
var directTracerAttrs []trace.SpanStartOption

type DirectAdapter struct {
	logic logic.Logic
}

func (d *DirectAdapter) Apply(server *mcp.Server) error {
	return apply(d, server)
}

func NewDirectAdapter(logic logic.Logic) *DirectAdapter {
	return &DirectAdapter{
		logic: logic,
	}
}

func (d *DirectAdapter) CreateEntities(ctx context.Context, args CreateEntitiesArgs) (*mcp.ToolResponse, error) {
	ctx, span := directTracer.Start(ctx, "CreateEntities", directTracerAttrs...)
	defer span.End()

	response := make([]Entity, 0)
	for _, entity := range args.Entities {
		// Process each entity
		newEntity := &models.Entity{
			Name: entity.Name,
			Type: entity.Type,
		}
		if err := d.logic.CreateEntity(ctx, newEntity); err != nil {
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
			if err := d.logic.CreateObservation(ctx, newObservation); err != nil {
				zap.L().Error("Can't create observation in database", zap.Error(err), zap.Any("observation", newObservation))
				span.RecordError(err)
				return nil, err
			}

			newEntityResponse.Observations = append(newEntityResponse.Observations, newObservation.Contents)
		}

		response = append(response, newEntityResponse)
	}

	// convert response to json string
	toolResponse, err := util.ToolJSONResponse(ctx, response)
	if err != nil {
		zap.L().Error("Can't marshal response json", zap.Error(err))
		span.RecordError(err)
		return nil, err
	}

	return toolResponse, nil
}

func (d *DirectAdapter) DeleteEntities(ctx context.Context, args DeleteEntitiesArgs) (*mcp.ToolResponse, error) {
	ctx, span := directTracer.Start(ctx, "DeleteEntities", directTracerAttrs...)
	defer span.End()

	for _, entityName := range args.EntityNames {
		// Process each entity
		entity, err := d.logic.ReadEntityByName(ctx, entityName)
		if err != nil {
			zap.L().Error("Can't read entity from database", zap.Error(err), zap.String("entityName", entityName))
			span.RecordError(err)
			return nil, err
		}
		if entity == nil {
			zap.L().Warn("Entity not found in database", zap.String("entityName", entityName))
			continue
		}

		if err := d.logic.DeleteAllObservationsByEntityID(ctx, entity.ID); err != nil {
			zap.L().Error("Can't delete observations", zap.Error(err), zap.String("entityName", entityName))
			span.RecordError(err)
			return nil, err
		}

		if err := d.logic.DeleteEntity(ctx, entity); err != nil {
			zap.L().Error("Can't delete entity", zap.Error(err), zap.String("entityName", entityName))
			span.RecordError(err)
			return nil, err
		}
	}

	return mcp.NewToolResponse(
		mcp.NewTextContent("Entities deleted successfully"),
	), nil
}

func (d *DirectAdapter) ReadGraph(ctx context.Context, _ ReadGraphArgs) (*mcp.ToolResponse, error) {
	ctx, span := directTracer.Start(ctx, "ReadGraph", directTracerAttrs...)
	defer span.End()

	// Read entities
	zap.L().Debug("Reading all entities from the database")
	entities, err := d.logic.ReadAllEntities(ctx)
	if err != nil && !errors.Is(err, logic.ErrNotFound) {
		zap.L().Error("Can't read entities from the database", zap.Error(err))
		span.RecordError(err)
		return nil, err
	}

	// Convert entities to response format
	zap.L().Debug("Converting entities to response format", zap.Any("entities", entities))
	entitiesResponse := make([]Entity, 0)
	for _, entity := range entities {
		newEntity := Entity{
			Name:         entity.Name,
			Type:         entity.Type,
			Observations: make([]string, 0),
		}
		for _, observation := range entity.Observations {
			newEntity.Observations = append(newEntity.Observations, observation.Contents)
		}
		entitiesResponse = append(entitiesResponse, newEntity)
	}

	// Read relations
	zap.L().Debug("Reading all relations from the database")
	relations, err := d.logic.ReadAllRelations(ctx)
	if err != nil && !errors.Is(err, logic.ErrNotFound) {
		span.RecordError(err)
		return nil, err
	}

	// Convert relations to response format
	zap.L().Debug("Converting relations to response format", zap.Any("relations", relations))
	relationsResponse := make([]Relation, 0)
	for _, relation := range relations {
		newRelation := Relation{
			From: relation.From.Name,
			To:   relation.To.Name,
			Type: relation.Type,
		}
		relationsResponse = append(relationsResponse, newRelation)
	}

	// Create the knowledge graph
	zap.L().Debug("Creating knowledge graph", zap.Int("entities", len(entitiesResponse)), zap.Int("relations", len(relationsResponse)))
	graph := KnowledgeGraph{
		Entities:  entitiesResponse,
		Relations: relationsResponse,
	}

	// Convert response to json string
	zap.L().Debug("Converting knowledge graph to JSON", zap.Any("knowledge_graph", graph))
	jsonResponse, err := util.ToolJSONResponse(ctx, graph)
	if err != nil {
		span.RecordError(err)
		return nil, err
	}
	return jsonResponse, nil
}

func (d *DirectAdapter) OpenNodes(ctx context.Context, args OpenNodesArgs) (*mcp.ToolResponse, error) {
	_, span := directTracer.Start(ctx, "OpenNodes", directTracerAttrs...)
	defer span.End()

	return mcp.NewToolResponse(
		mcp.NewTextContent("Open Nodes not implemented yet"),
	), nil
}

func (d *DirectAdapter) SearchNodes(ctx context.Context, args SearchNodesArgs) (*mcp.ToolResponse, error) {
	_, span := directTracer.Start(ctx, "SearchNodes", directTracerAttrs...)
	defer span.End()

	return mcp.NewToolResponse(
		mcp.NewTextContent("Search Nodes not implemented yet"),
	), nil
}

func (d *DirectAdapter) AddObservations(ctx context.Context, args AddObservationsArgs) (*mcp.ToolResponse, error) {
	ctx, span := directTracer.Start(ctx, "AddObservations", directTracerAttrs...)
	defer span.End()

	response := make([]AddedObservationsResp, 0, len(args.Observations))
	for _, observation := range args.Observations {
		entity, err := d.logic.ReadEntityByName(ctx, observation.EntityName)
		switch {
		case err != nil && !errors.Is(err, db.ErrNoEntries):
			zap.L().Error("Failed to read entity by name", zap.Error(err), zap.String("entity_name", observation.EntityName))
			span.RecordError(err)
			return nil, err
		case errors.Is(err, db.ErrNoEntries):
			return mcp.NewToolResponse(mcp.NewTextContent(fmt.Sprintf("The entity %s was not found", observation.EntityName))), nil
		}
		newResponse := AddedObservationsResp{
			EntityName: observation.EntityName,
		}

		for _, content := range observation.Contents {
			newObservation := &models.Observation{
				EntityID: entity.ID,
				Contents: content,
			}

			if err := d.logic.CreateObservation(ctx, newObservation); err != nil {
				zap.L().Error("Failed to create observation", zap.Error(err), zap.String("entity_name", observation.EntityName), zap.String("content", content))
				span.RecordError(err)
				return nil, err
			}
			newResponse.AddedObservations = append(newResponse.AddedObservations, newObservation.Contents)
		}

		response = append(response, newResponse)
	}

	// convert response to json string
	toolResponse, err := util.ToolJSONResponse(ctx, response)
	if err != nil {
		zap.L().Error("json marshal error", zap.Error(err))
		span.RecordError(err)
		return nil, err
	}

	return toolResponse, nil
}

func (d *DirectAdapter) DeleteObservations(ctx context.Context, args DeleteObservationsArgs) (*mcp.ToolResponse, error) {
	ctx, span := directTracer.Start(ctx, "DeleteObservations", directTracerAttrs...)
	defer span.End()

	for _, observation := range args.Deletions {
		entity, err := d.logic.ReadEntityByName(ctx, observation.EntityName)
		switch {
		case err != nil && !errors.Is(err, db.ErrNoEntries):
			zap.L().Error("Failed to read entity by name", zap.Error(err), zap.String("entity_name", observation.EntityName))
			span.RecordError(err)
			return nil, err
		case errors.Is(err, db.ErrNoEntries):
			return mcp.NewToolResponse(mcp.NewTextContent(fmt.Sprintf("The entity %s was not found", observation.EntityName))), nil
		}

		for _, content := range observation.Observations {
			// Read the observation by text
			observationToDelete, err := d.logic.ReadObservationByTextForEntityID(ctx, entity.ID, content)
			if err != nil {
				if errors.Is(err, db.ErrNoEntries) {
					// Observation not found, continue to the next one
					zap.L().Debug("Observation not found, skipping deletion", zap.String("entity_name", observation.EntityName), zap.String("content", content))
					continue
				}
				zap.L().Error("Failed to read observation by text", zap.Error(err), zap.String("entity_name", observation.EntityName), zap.String("content", content))
				span.RecordError(err)
				return nil, err
			}

			// Delete the observation
			zap.L().Debug("Deleting observation", zap.Int64("id", observationToDelete.ID), zap.String("content", content))
			if err := d.logic.DeleteObservation(ctx, observationToDelete); err != nil {
				zap.L().Error("Failed to delete observation", zap.Error(err), zap.Int64("id", observationToDelete.ID), zap.String("content", content))
				span.RecordError(err)
				return nil, err
			}
		}
	}

	return mcp.NewToolResponse(
		mcp.NewTextContent("Observations deleted successfully"),
	), nil
}

func (d *DirectAdapter) CreateRelations(ctx context.Context, args CreateRelationsArgs) (*mcp.ToolResponse, error) {
	ctx, span := directTracer.Start(ctx, "CreateRelations", directTracerAttrs...)
	defer span.End()

	response := make([]Relation, 0, len(args.Relations))
	for _, relation := range args.Relations {
		entityFrom, err := d.logic.ReadEntityByName(ctx, relation.From)
		if err != nil {
			return nil, err
		}
		zap.L().Debug("got from entity", zap.Any("entity", entityFrom))

		entityTo, err := d.logic.ReadEntityByName(ctx, relation.To)
		if err != nil {
			return nil, err
		}
		zap.L().Debug("got to entity", zap.Any("entity", entityTo))

		newRelation := &models.Relation{
			FromID: entityFrom.ID,
			ToID:   entityTo.ID,
			Type:   relation.Type,
		}
		if err := d.logic.CreateRelation(ctx, newRelation); err != nil {
			return nil, err
		}

		response = append(response, Relation{
			From: entityFrom.Name,
			To:   entityTo.Name,
			Type: newRelation.Type,
		})
	}

	// convert response to json string
	toolResponse, err := util.ToolJSONResponse(ctx, response)
	if err != nil {
		zap.L().Error("json marshal error", zap.Error(err))
		span.RecordError(err)
		return nil, err
	}

	return toolResponse, err

}

func (d *DirectAdapter) DeleteRelations(ctx context.Context, args DeleteRelationsArgs) (*mcp.ToolResponse, error) {
	ctx, span := directTracer.Start(ctx, "DeleteRelations", directTracerAttrs...)
	defer span.End()

	for _, relation := range args.Relations {
		entityFrom, err := d.logic.ReadEntityByName(ctx, relation.From)
		switch {
		case errors.Is(err, db.ErrNoEntries):
			zap.L().Debug("entity not found", zap.String("entity", relation.From), zap.String("position", "from"))
			return mcp.NewToolResponse(
				mcp.NewTextContent(fmt.Sprintf("Entity %s was not found", relation.From)),
			), nil
		case err != nil:
			zap.L().Error("read from entity error", zap.Error(err))
			span.RecordError(err)
			return nil, err
		default:
			zap.L().Debug("got from entity", zap.Any("entity", entityFrom))
		}

		entityTo, err := d.logic.ReadEntityByName(ctx, relation.To)
		switch {
		case errors.Is(err, db.ErrNoEntries):
			zap.L().Debug("entity not found", zap.String("entity", relation.To), zap.String("position", "to"))
			return mcp.NewToolResponse(
				mcp.NewTextContent(fmt.Sprintf("Entity %s was not found", relation.To)),
			), nil
		case err != nil:
			zap.L().Error("read to entity error", zap.Error(err))
			span.RecordError(err)
			return nil, err
		default:
			zap.L().Debug("got to entity", zap.Any("entity", entityFrom))
		}

		// find the relation
		existingRelation, err := d.logic.ReadExactRelation(ctx, entityFrom.ID, entityTo.ID, relation.Type)
		if err != nil {
			return nil, err
		}

		if err := d.logic.DeleteRelation(ctx, existingRelation); err != nil {
			return nil, err
		}
	}

	return mcp.NewToolResponse(
		mcp.NewTextContent("Relations deleted successfully"),
	), nil
}

var _ Adapter = (*DirectAdapter)(nil)
