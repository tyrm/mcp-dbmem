package bun

import (
	"context"
	"database/sql"
	"errors"

	"github.com/tyrm/mcp-dbmem/internal/db"
	"github.com/tyrm/mcp-dbmem/internal/models"
	"github.com/uptrace/bun"
)

func (c *Client) CreateEntity(ctx context.Context, entity *models.Entity) db.Error {
	ctx, span := tracer.Start(ctx, "CreateEntity", tracerAttrs...)
	defer span.End()

	err := c.create(ctx, entity)
	span.RecordError(err)
	return err
}

func (c *Client) DeleteEntity(ctx context.Context, entity *models.Entity) db.Error {
	ctx, span := tracer.Start(ctx, "DeleteEntity", tracerAttrs...)
	defer span.End()

	tx, err := c.db.BeginTx(ctx, &sql.TxOptions{})
	if err != nil {
		return c.ProcessError(err)
	}

	if err := deleteAllRelationsByEntityID(ctx, tx, entity.ID); err != nil {
		span.RecordError(err)
		if err := tx.Rollback(); err != nil {
			span.RecordError(err)
			return c.ProcessError(err)
		}
		return c.ProcessError(err)
	}

	if err := deleteAllObservationsByEntityID(ctx, tx, entity.ID); err != nil {
		span.RecordError(err)
		if err := tx.Rollback(); err != nil {
			span.RecordError(err)
			return c.ProcessError(err)
		}
		return c.ProcessError(err)
	}

	tx.NewDelete().
		Model(&models.Relation{}).
		WherePK()

	err = c.delete(ctx, entity)
	span.RecordError(err)
	return err
}

func (c *Client) ReadAllEntities(ctx context.Context) ([]*models.Entity, db.Error) {
	ctx, span := tracer.Start(ctx, "ReadAllEntities", tracerAttrs...)
	defer span.End()

	var entities []*models.Entity
	query := newEntitiesQ(c.db, &entities)

	if err := query.Scan(ctx); err != nil {
		err := c.ProcessError(err)
		if !errors.Is(err, db.ErrNoEntries) {
			span.RecordError(err)
		}
		return nil, err
	}

	return entities, nil
}

func (c *Client) ReadEntityByName(ctx context.Context, name string) (*models.Entity, db.Error) {
	ctx, span := tracer.Start(ctx, "ReadEntityByName", tracerAttrs...)
	defer span.End()

	entity := new(models.Entity)
	query := newEntityQ(c.db, entity).
		Where("name = ?", name)

	if err := query.Scan(ctx); err != nil {
		err := c.ProcessError(err)
		if !errors.Is(err, db.ErrNoEntries) {
			span.RecordError(err)
		}
		return nil, err
	}

	return entity, nil
}

func newEntityQ(c bun.IDB, i *models.Entity) *bun.SelectQuery {
	return c.
		NewSelect().
		Model(i).
		Relation("Observations")
}

func newEntitiesQ(c bun.IDB, i *[]*models.Entity) *bun.SelectQuery {
	return c.
		NewSelect().
		Model(i).
		Relation("Observations")
}
