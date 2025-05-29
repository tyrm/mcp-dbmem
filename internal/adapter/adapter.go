package adapter

import (
	"context"

	mcp "github.com/metoro-io/mcp-golang"
)

var (
	RespEntityDeleted      = mcp.NewToolResponse(mcp.NewTextContent("Entities deleted successfully"))
	RespObservationDeleted = mcp.NewToolResponse(mcp.NewTextContent("Observations deleted successfully"))
	RespRelationDeleted    = mcp.NewToolResponse(mcp.NewTextContent("Relations deleted successfully"))
)

type Adapter interface {
	CreateEntities(ctx context.Context, args CreateEntitiesArgs) (*mcp.ToolResponse, error)
	DeleteEntities(ctx context.Context, args DeleteEntitiesArgs) (*mcp.ToolResponse, error)
	ReadGraph(ctx context.Context, args ReadGraphArgs) (*mcp.ToolResponse, error)
	OpenNodes(ctx context.Context, args OpenNodesArgs) (*mcp.ToolResponse, error)
	SearchNodes(ctx context.Context, args SearchNodesArgs) (*mcp.ToolResponse, error)
	AddObservations(ctx context.Context, args AddObservationsArgs) (*mcp.ToolResponse, error)
	DeleteObservations(ctx context.Context, args DeleteObservationsArgs) (*mcp.ToolResponse, error)
	CreateRelations(ctx context.Context, args CreateRelationsArgs) (*mcp.ToolResponse, error)
	DeleteRelations(ctx context.Context, args DeleteRelationsArgs) (*mcp.ToolResponse, error)
	Apply(server *mcp.Server) error
}

func apply[A Adapter](adapter A, server *mcp.Server) error {
	if err := server.RegisterTool("create_entities", "Create multiple new entities in the knowledge graph", adapter.CreateEntities); err != nil {
		return err
	}
	if err := server.RegisterTool("create_relations", "Create multiple new relations between entities in the knowledge graph. Relations should be in active voice", adapter.CreateRelations); err != nil {
		return err
	}
	if err := server.RegisterTool("add_observations", "Add new observations to existing entities in the knowledge graph", adapter.AddObservations); err != nil {
		return err
	}
	if err := server.RegisterTool("delete_entities", "Delete multiple entities and their associated relations from the knowledge graph", adapter.DeleteEntities); err != nil {
		return err
	}
	if err := server.RegisterTool("delete_observations", "Delete specific observations from entities in the knowledge graph", adapter.DeleteObservations); err != nil {
		return err
	}
	if err := server.RegisterTool("delete_relations", "Delete multiple relations from the knowledge graph", adapter.DeleteRelations); err != nil {
		return err
	}
	if err := server.RegisterTool("read_graph", "Read the entire knowledge graph", adapter.ReadGraph); err != nil {
		return err
	}
	if err := server.RegisterTool("search_nodes", "Search for nodes in the knowledge graph based on adapter query", adapter.SearchNodes); err != nil {
		return err
	}
	if err := server.RegisterTool("open_nodes", "Open specific nodes in the knowledge graph by their names", adapter.OpenNodes); err != nil {
		return err
	}

	return nil
}

// Models

// Entity represents an entity in the knowledge graph.
type Entity struct {
	Name         string   `json:"name"         jsonschema:"required,description=The name of the entity"`
	Type         string   `json:"entityType"   jsonschema:"required,description=The type of the entity"`
	Observations []string `json:"observations" jsonschema:"required,description=An array of observation contents associated with the entity"`
}

// KnowledgeGraph represents the entire knowledge graph.
type KnowledgeGraph struct {
	Entities  []Entity   `json:"entities"`
	Relations []Relation `json:"relations"`
}

// Relation represents a relationship between two entities in the knowledge graph.
type Relation struct {
	From string `json:"from"         jsonschema:"required,description=The name of the entity where the relation starts"`
	To   string `json:"to"           jsonschema:"required,description=The name of the entity where the relation ends"`
	Type string `json:"relationType" jsonschema:"required,description=The type of the relation"`
}

// Request / Response

// CreateEntitiesArgs represents the arguments for creating entities.
type CreateEntitiesArgs struct {
	Entities []Entity `json:"entities" jsonschema:"required,description=An array of observation contents associated with the entity"`
}

// DeleteEntitiesArgs represents the arguments for deleting entities.
type DeleteEntitiesArgs struct {
	EntityNames []string `json:"entityNames" jsonschema:"required,description=An array of entity names to delete"`
}

// ReadGraphArgs represents the arguments for reading the knowledge graph.
type ReadGraphArgs struct {
}

// OpenNodesArgs represents the arguments for opening nodes.
type OpenNodesArgs struct {
	Names []string `json:"names" jsonschema:"required,description=An array of entity names to retrieve"`
}

// SearchNodesArgs represents the arguments for searching nodes.
type SearchNodesArgs struct {
	Query string `json:"query" jsonschema:"required,description=The search query to match against entity names, types, and observation content"`
}

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

// DeleteObservationsArgs represents the arguments for deleting Observations.
type DeleteObservationsArgs struct {
	Deletions []DeleteObservation `json:"deletions" jsonschema:"required,description=An array of observations to delete"`
}

// DeleteObservation represents an observation associated with an entity.
type DeleteObservation struct {
	EntityName   string   `json:"entityName"   jsonschema:"required,description=The name of the entity containing the observations"`
	Observations []string `json:"observations" jsonschema:"required,description=An array of observations to delete"`
}

// CreateRelationsArgs represents the arguments for creating Relationships.
type CreateRelationsArgs struct {
	Relations []Relation `json:"relations" jsonschema:"required,description=Create multiple new relations between entities in the knowledge graph. Relations should be in active voice"`
}

// DeleteRelationsArgs represents the arguments for deleting Relationships.
type DeleteRelationsArgs struct {
	Relations []Relation `json:"relations" jsonschema:"required,description=Delete multiple relations from the knowledge graph"`
}
