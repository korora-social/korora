package user

import "github.com/korora-social/korora/models"

type Dao interface {
	GetByWebfinger(user, domain string) (*models.User, error)
}

type dao struct {
}
