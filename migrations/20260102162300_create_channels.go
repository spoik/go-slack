package migrations

import (
	"context"
	"github.com/uptrace/bun"
	"go-slack/channels"
)

func init() {
	Migrations.MustRegister(func(ctx context.Context, db *bun.DB) error {
		_, err := db.NewCreateTable().
			Model((*channels.Channel)(nil)).
			IfNotExists().
			Exec(ctx)

		return err
	}, func(ctx context.Context, db *bun.DB) error {
		_, err := db.NewDropTable().
			Model((*channels.Channel)(nil)).
			IfExists().
			Exec(ctx)
		return err
	})
}
