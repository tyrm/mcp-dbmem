package v1

import (
	"context"
	"errors"
	"fmt"

	mcp "github.com/metoro-io/mcp-golang"
	"github.com/tyrm/mcp-dbmem/internal/db"
	"github.com/tyrm/mcp-dbmem/internal/models"
	"go.uber.org/zap"
)

// AddObservationsArgs represents the arguments for creating Observations.
type AddObservationsArgs struct {
	Observations []AddObservation `json:"observations" jsonschema:"required,description=An array of observation contents to add"`
}

// AddObservation represents an observation associated with an entity.
type AddObservation struct {
	EntityName string   `json:"entityName" jsonschema:"required,description=The name of the entity to add the observations to"`
	Contents   []string `json:"contents"   jsonschema:"required,description=An array of observation contents to addAn array of observations"`
}

// AddedObservationsResp represents the response for creating Observations.
type AddedObservationsResp struct {
	EntityName        string   `json:"entityName"        jsonschema:"required,description=The name of the entity containing the observations"`
	AddedObservations []string `json:"addedObservations" jsonschema:"required,description=An array of observations"`
}

// AddObservations creates Observations on entities in the knowledge graph.
func (l *Logic) AddObservations(ctx context.Context, args AddObservationsArgs) (*mcp.ToolResponse, error) {
	ctx, span := tracer.Start(ctx, "AddObservations", tracerAttrs...)
	defer span.End()

	response := make([]AddedObservationsResp, 0, len(args.Observations))
	for _, observation := range args.Observations {
		entity, err := l.DB.ReadEntityByName(ctx, observation.EntityName)
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

			if err := l.DB.CreateObservation(ctx, newObservation); err != nil {
				zap.L().Error("Failed to create observation", zap.Error(err), zap.String("entity_name", observation.EntityName), zap.String("content", content))
				span.RecordError(err)
				return nil, err
			}

			response = append(response, newResponse)
		}
	}

	// convert response to json string
	toolResponse, err := toolJSONResponse(ctx, response)
	if err != nil {
		zap.L().Error("json marshal error", zap.Error(err))
		span.RecordError(err)
		return nil, err
	}

	return toolResponse, nil
}

// DeleteObservationsArgs represents the arguments for deleting Observations.
type DeleteObservationsArgs struct {
	Deletions []DeleteObservation `json:"deletions" jsonschema:"required,description=An array of observations to delete"`
}

// DeleteObservation represents an observation associated with an entity.
type DeleteObservation struct {
	EntityName   string   `json:"entityName"   jsonschema:"required,description=The name of the entity containing the observations"`
	Observations []string `json:"observations" jsonschema:"required,description=An array of observations to delete"`
}

// DeleteObservations deletes Observations on entities in the knowledge graph.
func (l *Logic) DeleteObservations(ctx context.Context, args DeleteObservationsArgs) (*mcp.ToolResponse, error) {
	ctx, span := tracer.Start(ctx, "DeleteObservations", tracerAttrs...)
	defer span.End()

	for _, observation := range args.Deletions {
		entity, err := l.DB.ReadEntityByName(ctx, observation.EntityName)
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
			observationToDelete, err := l.DB.ReadObservationByTextForEntityID(ctx, entity.ID, content)
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
			if err := l.DB.DeleteObservation(ctx, observationToDelete); err != nil {
				zap.L().Error("Failed to create observation", zap.Error(err), zap.Int64("id", observationToDelete.ID), zap.String("content", content))
				span.RecordError(err)
				return nil, err
			}
		}
	}

	return mcp.NewToolResponse(
		mcp.NewTextContent("Observations deleted successfully"),
	), nil
}
