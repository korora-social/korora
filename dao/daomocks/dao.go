package daomocks

import (
	"github.com/korora-social/korora/dao"
	"github.com/korora-social/korora/dao/user"
)

type MockDao struct {
	dao.Dao
	UserDao user.Dao
}

func (m *MockDao) User() user.Dao {
	return m.UserDao
}
