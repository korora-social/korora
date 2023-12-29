package user

import (
	"github.com/korora-social/korora/models"
	"github.com/uptrace/bun"
)

type User struct {
	bun.BaseModel `bun:"table:users"`
	Pk            int    `bun:"pk,pk,autoincrement"`
	Username      string `bun:"username,unique,notnull"`
	Uri           string `bun:"uri,unique,notnull"`
	PublicKey     []byte `bun:"public_key,nullzero"`
	PrivateKey    []byte `bun:"private_key,nullzero"`
}

func (u *User) ToUser() *models.User {
	return &models.User{
		Username:   u.Username,
		Uri:        u.Uri,
		PrivateKey: u.PrivateKey,
		PublicKey:  u.PublicKey,
	}
}
