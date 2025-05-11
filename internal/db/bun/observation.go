package bun

import (
	"context"
	"errors"

	"github.com/tyrm/mcp-dbmem/internal/db"
	"github.com/tyrm/mcp-dbmem/internal/models"
	"github.com/uptrace/bun"
)

func (c *Client) CreateObservation(ctx context.Context, observation *models.Observation) db.Error {
	ctx, span := tracer.Start(ctx, "CreateObservation", tracerAttrs...)
	defer span.End()

	err := c.create(ctx, observation)
	span.RecordError(err)
	return err
}

func (c *Client) DeleteAllObservationsByEntityID(ctx context.Context, entityID int64) db.Error {
	ctx, span := tracer.Start(ctx, "DeleteAllObservationsByEntityID", tracerAttrs...)
	defer span.End()

	query := c.db.NewDelete().
		Model((*models.Observation)(nil)).
		Where("entity_id = ?", entityID)

	if _, err := query.Exec(ctx); err != nil {
		span.RecordError(err)
		return c.ProcessError(err)
	}

	return nil
}

func (c *Client) DeleteObservation(ctx context.Context, observation *models.Observation) db.Error {
	ctx, span := tracer.Start(ctx, "DeleteObservation", tracerAttrs...)
	defer span.End()

	err := c.delete(ctx, observation)
	span.RecordError(err)
	return err
}

func (c *Client) ReadObservationByTextForEntityID(ctx context.Context, entityID int64, text string) (*models.Observation, db.Error) {
	ctx, span := tracer.Start(ctx, "ReadObservationByTextForEntityID", tracerAttrs...)
	defer span.End()

	observation := new(models.Observation)
	query := newObservationQ(c.db, observation).
		Where("contents = ?", text).
		Where("entity_id = ?", entityID)

	if err := query.Scan(ctx); err != nil {
		err := c.ProcessError(err)
		if !errors.Is(err, db.ErrNoEntries) {
			span.RecordError(err)
		}
		return nil, err
	}

	return observation, nil
}

func newObservationQ(c bun.IDB, i *models.Observation) *bun.SelectQuery {
	return c.
		NewSelect().
		Model(i)
}
