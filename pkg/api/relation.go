package api

// Relation represents a relationship between two entities in the knowledge graph.
type Relation struct {
	From string `json:"from"         validate:"required"`
	To   string `json:"to"           validate:"required"`
	Type string `json:"relationType" validate:"required"`
}

type CreateRelationsRequest struct {
}

type CreateRelationsResponse struct {
}

type DeleteRelationsRequest struct {
}

type DeleteRelationsResponse struct {
}
