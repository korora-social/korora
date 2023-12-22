package daomocks

import (
	"github.com/korora-social/korora/dao"
	"github.com/korora-social/korora/dao/user"
	"github.com/korora-social/korora/models"
)

type UserDao struct {
	user.Dao
	Users []*models.User
}

func (m *UserDao) GetByUsername(username string) (*models.User, error) {
	for _, user := range m.Users {
		if user.Username == username {
			return user, nil
		}
	}

	return nil, dao.NotFound
}
