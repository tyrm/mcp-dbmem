package models

import "time"

// Relation represents a relation between two entities in a knowledge graph.
type Relation struct {
	ID        int64     `bun:",pk,autoincrement"`
	CreatedAt time.Time `bun:",nullzero,notnull,default:current_timestamp"`
	UpdatedAt time.Time `bun:",nullzero,notnull,default:current_timestamp"`

	Type string `bun:"type,notnull" json:"type"`

	FromID int64   `bun:"from_id,notnull" json:"from_id"`
	From   *Entity `bun:"rel:belongs-to,join:from_id=id" json:"from"`
	ToID   int64   `bun:"to_id,notnull" json:"to_id"`
	To     *Entity `bun:"rel:belongs-to,join:to_id=id" json:"to"`
}
