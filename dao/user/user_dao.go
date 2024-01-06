package user

import (
	"database/sql"

	"github.com/korora-social/korora/dao/daoerrors"
	"github.com/korora-social/korora/models"
	sq "github.com/sleepdeprecation/squirrelly"
)

type Dao interface {
	GetByUsername(username string) (*models.User, error)
	Save(*models.User) error
}

type dao struct {
	db *sq.Db
}

func New(db *sq.Db) Dao {
	return &dao{db: db}
}

func (d *dao) GetByUsername(username string) (*models.User, error) {
	u := &models.User{}
	query := sq.Select("*").From("users").Where(sq.Eq{"username": username})
	err := d.db.Get(query, u)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, daoerrors.NotFound
		}
		return nil, err
	}

	return u, err
}

func (d *dao) Save(u *models.User) error {
	query := sq.Insert("users").Struct(u).OnConflict("uri").UpdateColumns("public_key", "private_key")
	_, err := d.db.Exec(query)

	return err
}
