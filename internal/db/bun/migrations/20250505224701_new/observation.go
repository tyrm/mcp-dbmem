package models

import "time"

type Observation struct {
	ID        int64     `bun:",pk,autoincrement"`
	CreatedAt time.Time `bun:",nullzero,notnull,default:current_timestamp"`
	UpdatedAt time.Time `bun:",nullzero,notnull,default:current_timestamp"`

	Contents string `bun:"contents,notnull" json:"contents"`

	EntityID int64   `bun:"entity_id,notnull"                json:"entity_id"`
	Entity   *Entity `bun:"rel:belongs-to,join:entity_id=id" json:"entity"`
}
