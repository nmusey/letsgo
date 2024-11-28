package users

import (
	"$appRepo/internal/core"
)

type UserService struct {
    Store UserStore
}

func NewUserService(router *core.Router) *UserService {
    return &UserService{Store: newSQLUserStore(router.DB.SqlDB)}
}

type UserStore interface {
	SaveUser(user *User) error

	GetUsers() ([]User, error)
	GetUserByEmail(email string) (*User, error)
	GetUserById(id int) (*User, error)
}
