package api

// Entity represents an entity in the knowledge graph.
type Entity struct {
	Name         string   `json:"name"         validate:"required"`
	Type         string   `json:"entityType"   validate:"required"`
	Observations []string `json:"observations"`
}

type CreateEntitiesRequest struct {
}

type CreateEntitiesResponse struct {
}

type DeleteEntitiesRequest struct {
}

type DeleteEntitiesResponse struct {
}
