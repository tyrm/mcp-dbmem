package api

// KnowledgeGraph represents the entire knowledge graph.
type KnowledgeGraph struct {
	Entities  []Entity   `json:"entities"`
	Relations []Relation `json:"relations"`
}

type ReadGraphRequest struct {
}

type ReadGraphResponse struct {
}
