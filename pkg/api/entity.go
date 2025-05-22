package api

// Entity represents an entity in the knowledge graph.
type Entity struct {
	Name         string   `json:"name"         validate:"required"`
	Type         string   `json:"entityType"   validate:"required"`
	Observations []string `json:"observations"`
}

type CreateEntitiesRequest struct {
	Entities []Entity `json:"entities" validate:"required,min=1"`
}

type CreateEntitiesResponse struct {
	Entities []Entity `json:"entities"`
}

type DeleteEntitiesRequest struct {
	EntityNames []string `json:"entityNames" validate:"required,min=1"`
}

type DeleteEntitiesResponse struct {
	Success bool `json:"success"`
}
