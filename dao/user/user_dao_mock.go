package user

import (
	"fmt"

	"github.com/korora-social/korora/models"
)

type MockDao struct {
	Dao
	Users []*models.User
}

func (m *MockDao) GetByUsername(username string) (*models.User, error) {
	for _, user := range m.Users {
		if user.Username == username {
			return user, nil
		}
	}

	return nil, fmt.Errorf("No user found")
}
