package bun

import (
	"context"
	"errors"

	"github.com/tyrm/mcp-dbmem/internal/db"
	"github.com/tyrm/mcp-dbmem/internal/models"
	"github.com/uptrace/bun"
)

func (c *Client) CreateRelation(ctx context.Context, relation *models.Relation) db.Error {
	ctx, span := tracer.Start(ctx, "CreateRelation", tracerAttrs...)
	defer span.End()

	err := c.create(ctx, relation)
	span.RecordError(err)
	return err
}

func (c *Client) DeleteRelation(ctx context.Context, relation *models.Relation) db.Error {
	ctx, span := tracer.Start(ctx, "DeleteRelation", tracerAttrs...)
	defer span.End()

	err := c.delete(ctx, relation)
	span.RecordError(err)
	return err
}

func (c *Client) ReadAllRelations(ctx context.Context) ([]*models.Relation, db.Error) {
	ctx, span := tracer.Start(ctx, "ReadAllRelations", tracerAttrs...)
	defer span.End()

	var relations []*models.Relation
	query := newRelationsQ(c.db, &relations)

	if err := query.Scan(ctx); err != nil {
		span.RecordError(err)
		return nil, c.ProcessError(err)
	}

	return relations, nil
}

func (c *Client) ReadExactRelation(ctx context.Context, fromID, toID int64, relationType string) (*models.Relation, db.Error) {
	ctx, span := tracer.Start(ctx, "ReadExactRelation", tracerAttrs...)
	defer span.End()

	relation := new(models.Relation)
	query := newRelationQ(c.db, relation).
		Where("from_id = ?", fromID).
		Where("to_id = ?", toID).
		Where("relation_type = ?", relationType)

	if err := query.Scan(ctx); err != nil {
		err := c.ProcessError(err)
		if !errors.Is(err, db.ErrNoEntries) {
			span.RecordError(err)
		}
		return nil, err
	}

	return relation, nil
}

func deleteAllRelationsByEntityID(ctx context.Context, c bun.IDB, entityID int64) db.Error {
	ctx, span := tracer.Start(ctx, "deleteAllRelationsByEntityID", tracerAttrs...)
	defer span.End()

	query := c.
		NewDelete().
		Model((*models.Relation)(nil)).
		Where("from_id = ?", entityID).
		WhereOr("to_id = ?", entityID)

	if _, err := query.Exec(ctx); err != nil {
		span.RecordError(err)
		return err
	}

	return nil
}
func newRelationQ(c bun.IDB, i *models.Relation) *bun.SelectQuery {
	return c.
		NewSelect().
		Model(i).
		Relation("From").
		Relation("To")
}

func newRelationsQ(c bun.IDB, i *[]*models.Relation) *bun.SelectQuery {
	return c.
		NewSelect().
		Model(i).
		Relation("From").
		Relation("To")
}
