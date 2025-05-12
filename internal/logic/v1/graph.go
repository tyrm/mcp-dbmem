package v1

import (
	"context"
	"errors"

	mcp "github.com/metoro-io/mcp-golang"
	"github.com/tyrm/mcp-dbmem/internal/db"
	"go.uber.org/zap"
)

// KnowledgeGraph represents the entire knowledge graph.
type KnowledgeGraph struct {
	Entities  []Entity   `json:"entities"`
	Relations []Relation `json:"relations"`
}

// ReadGraphArgs represents the arguments for reading the knowledge graph.
type ReadGraphArgs struct {
}

// ReadGraph reads the entire knowledge graph.
func (l *Logic) ReadGraph(ctx context.Context, _ ReadGraphArgs) (*mcp.ToolResponse, error) {
	ctx, span := tracer.Start(ctx, "ReadGraph", tracerAttrs...)
	defer span.End()

	// Read entities
	zap.L().Debug("Reading all entities from the database")
	entities, err := l.DB.ReadAllEntities(ctx)
	if err != nil && !errors.Is(err, db.ErrNoEntries) {
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
	relations, err := l.DB.ReadAllRelations(ctx)
	if err != nil && !errors.Is(err, db.ErrNoEntries) {
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
	jsonResponse, err := toolJSONResponse(ctx, graph)
	if err != nil {
		span.RecordError(err)
		return nil, err
	}
	return jsonResponse, nil
}
