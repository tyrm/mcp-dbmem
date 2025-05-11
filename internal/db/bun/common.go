package bun

import (
	"context"

	"github.com/tyrm/mcp-dbmem/internal/db"
)

func (c *Client) create(ctx context.Context, o any) db.Error {
	ctx, span := tracer.Start(ctx, "create", tracerAttrs...)
	defer span.End()

	query := c.db.NewInsert().
		Model(o).
		ExcludeColumn("created_at", "updated_at")

	if _, err := query.Exec(ctx); err != nil {
		span.RecordError(err)
		return c.ProcessError(err)
	}

	return nil
}

func (c *Client) delete(ctx context.Context, o any) db.Error {
	ctx, span := tracer.Start(ctx, "delete", tracerAttrs...)
	defer span.End()

	query := c.db.
		NewDelete().
		Model(o).
		WherePK()

	if _, err := query.Exec(ctx); err != nil {
		span.RecordError(err)
		return c.ProcessError(err)
	}

	return nil
}
