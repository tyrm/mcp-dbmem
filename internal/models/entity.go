package models

import "time"

// Entity represents an entity in a knowledge graph.
type Entity struct {
	ID        int64     `bun:",pk,autoincrement"`
	CreatedAt time.Time `bun:",nullzero,notnull,default:current_timestamp"`
	UpdatedAt time.Time `bun:",nullzero,notnull,default:current_timestamp"`

	Name         string         `bun:"name,notnull"                   json:"name"`
	Type         string         `bun:"type,notnull"                   json:"type"`
	Observations []*Observation `bun:"rel:has-many,join:id=entity_id" json:"observations"`
}
