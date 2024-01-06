// DAOs, or Data Access Objects, are the mechanism used to talk to backing data stores.
package dao

import (
	"github.com/korora-social/korora/dao/user"
	sq "github.com/sleepdeprecation/squirrelly"
)

type Dao interface {
	User() user.Dao
}

type dao struct {
	db      *sq.Db
	userDao user.Dao
}

func New(db *sq.Db) Dao {
	return &dao{
		db:      db,
		userDao: user.New(db),
	}
}

func (d *dao) User() user.Dao {
	return d.userDao
}
