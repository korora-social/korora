// DAOs, or Data Access Objects, are the mechanism used to talk to backing data stores.
package dao

import (
	"errors"

	"github.com/korora-social/korora/dao/user"
	"github.com/uptrace/bun"
)

var (
	NotFound error = errors.New("Record not found")
)

type Dao interface {
	User() user.Dao
}

type dao struct {
	db      *bun.DB
	userDao user.Dao
}

func New(db *bun.DB) Dao {
	return &dao{
		db:      db,
		userDao: user.New(db),
	}
}

func (d *dao) User() user.Dao {
	return d.userDao
}
