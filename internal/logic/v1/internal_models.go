package v1

// Entity represents an entity in the knowledge graph
type Entity struct {
	Name         string   `json:"name" jsonschema:"required,description=The name of the entity"`
	Type         string   `json:"entityType" jsonschema:"required,description=The type of the entity"`
	Observations []string `json:"observations" jsonschema:"required,description=An array of observation contents associated with the entity"`
}

// Observation represents an observation associated with an entity
//type Observation struct {
//	EntityName string   `json:"entityName" jsonschema:"required,description=The name of the entity containing the observations"`
//	Contents   []string `json:"contents" jsonschema:"required,description=An array of observations"`
//}

// Relation represents a relationship between two entities in the knowledge graph
//type Relation struct {
//	From string `json:"from" jsonschema:"required,description=The name of the entity where the relation starts"`
//	To   string `json:"to" jsonschema:"required,description=The name of the entity where the relation ends"`
//	Type string `json:"relationType" jsonschema:"required,description=The type of the relation"`
//}

// KnowledgeGraph represents the entire knowledge graph
type KnowledgeGraph struct {
	Entities  []Entity   `json:"entities"`
	Relations []Relation `json:"relations"`
}
