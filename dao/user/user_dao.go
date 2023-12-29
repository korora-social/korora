package user

import (
	"context"

	"github.com/korora-social/korora/models"
	"github.com/uptrace/bun"
)

type Dao interface {
	GetByUsername(username string) (*models.User, error)
}

type dao struct {
	db *bun.DB
}

func New(db *bun.DB) Dao {
	return &dao{db: db}
}

func (d *dao) GetByUsername(username string) (*models.User, error) {
	u := &User{}
	err := d.db.NewSelect().Model(u).Where("? = ?", bun.Ident("username"), username).Scan(context.Background(), u)
	if err != nil {
		return nil, err
	}

	return u.ToUser(), nil
}
