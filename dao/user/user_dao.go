package user

import "github.com/korora-social/korora/models"

type Dao interface {
	GetByUsername(username string) (*models.User, error)
}

type dao struct {
}
