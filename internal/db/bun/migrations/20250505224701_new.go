package migrations

import (
	"context"

	models "github.com/tyrm/mcp-dbmem/internal/db/bun/migrations/20250505224701_new"
	"github.com/uptrace/bun"
	"tyr.codes/libs/libmigration"
)

func init() {
	addTables := libmigration.TableList{
		{
			Model: &models.Entity{},
		},
		{
			Model: &models.Observation{},
			ForeignKeys: []string{
				"(entity_id) REFERENCES entities (id) ON DELETE CASCADE",
			},
		},
		{
			Model: &models.Relation{},
			ForeignKeys: []string{
				"(from_id) REFERENCES entities (id) ON DELETE CASCADE",
				"(to_id) REFERENCES entities (id) ON DELETE CASCADE",
			},
		},
	}

	addIndexes := libmigration.IndexList{}

	up := func(ctx context.Context, db *bun.DB) error {
		return db.RunInTx(ctx, nil, func(ctx context.Context, tx bun.Tx) error {
			if err := libmigration.AddTablesUp(ctx, tx, addTables); err != nil {
				return err
			}

			if err := libmigration.AddIndexesUp(ctx, tx, addIndexes); err != nil {
				return err
			}

			return nil
		})
	}

	down := func(ctx context.Context, db *bun.DB) error {
		return db.RunInTx(ctx, nil, func(ctx context.Context, tx bun.Tx) error {
			if err := libmigration.AddIndexesDown(ctx, tx, addIndexes); err != nil {
				return err
			}

			if err := libmigration.AddTablesDown(ctx, tx, addTables); err != nil {
				return err
			}

			return nil
		})
	}

	if err := Migrations.Register(up, down); err != nil {
		panic(err)
	}
}
