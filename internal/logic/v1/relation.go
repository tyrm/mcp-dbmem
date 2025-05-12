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

// Relation represents a relationship between two entities in the knowledge graph.
type Relation struct {
	From string `json:"from" jsonschema:"required,description=The name of the entity where the relation starts"`
	To   string `json:"to" jsonschema:"required,description=The name of the entity where the relation ends"`
	Type string `json:"relationType" jsonschema:"required,description=The type of the relation"`
}

// CreateRelationsArgs represents the arguments for creating Relationships.
type CreateRelationsArgs struct {
	Relations []Relation `json:"relations" jsonschema:"required,description=Create multiple new relations between entities in the knowledge graph. Relations should be in active voice"`
}

// CreateRelations creates relationships in the knowledge graph.
func (l *Logic) CreateRelations(ctx context.Context, args CreateRelationsArgs) (*mcp.ToolResponse, error) {
	ctx, span := tracer.Start(ctx, "CreateRelations", tracerAttrs...)
	defer span.End()

	response := make([]Relation, 0, len(args.Relations))
	for _, relation := range args.Relations {
		entityFrom, err := l.DB.ReadEntityByName(ctx, relation.From)
		if err != nil {
			return nil, err
		}
		zap.L().Debug("got from entity", zap.Any("entity", entityFrom))

		entityTo, err := l.DB.ReadEntityByName(ctx, relation.To)
		if err != nil {
			return nil, err
		}
		zap.L().Debug("got to entity", zap.Any("entity", entityTo))

		newRelation := &models.Relation{
			FromID: entityFrom.ID,
			ToID:   entityTo.ID,
			Type:   relation.Type,
		}
		if err := l.DB.CreateRelation(ctx, newRelation); err != nil {
			return nil, err
		}

		response = append(response, Relation{
			From: entityFrom.Name,
			To:   entityTo.Name,
			Type: newRelation.Type,
		})
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

// DeleteRelationsArgs represents the arguments for deleting Relationships.
type DeleteRelationsArgs struct {
	Relations []Relation `json:"relations" jsonschema:"required,description=Delete multiple relations from the knowledge graph"`
}

// DeleteRelations deletes relationships from the knowledge graph.
func (l *Logic) DeleteRelations(ctx context.Context, args DeleteRelationsArgs) (*mcp.ToolResponse, error) {
	ctx, span := tracer.Start(ctx, "DeleteRelations", tracerAttrs...)
	defer span.End()

	for _, relation := range args.Relations {
		entityFrom, err := l.DB.ReadEntityByName(ctx, relation.From)
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

		entityTo, err := l.DB.ReadEntityByName(ctx, relation.To)
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
		existingRelation, err := l.DB.ReadExactRelation(ctx, entityFrom.ID, entityTo.ID, relation.Type)
		if err != nil {
			return nil, err
		}

		if err := l.DB.DeleteRelation(ctx, existingRelation); err != nil {
			return nil, err
		}
	}

	return mcp.NewToolResponse(
		mcp.NewTextContent("Relations deleted successfully"),
	), nil
}
