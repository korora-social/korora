package korora

import (
	"github.com/korora-social/korora/dao"
	sq "github.com/sleepdeprecation/squirrelly"
)

type Korora struct {
	dao dao.Dao
}

func New(db *sq.Db) *Korora {
	return &Korora{
		dao: dao.New(db),
	}
}
