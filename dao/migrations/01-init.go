package migrations

import (
	"context"

	"github.com/korora-social/korora/dao/user"
	"github.com/uptrace/bun"
)

var _ = AddMigration("2023-12-29 - init", func(tx bun.Tx) error {
	_, err := tx.NewCreateTable().Model(&user.User{}).Exec(context.Background())
	return err
})
