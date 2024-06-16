package users

import (
	"$appRepo/pkg/core"
)

type SQLUserService struct {
    ctx *core.RouterContext
}

func NewUserService(ctx *core.RouterContext) SQLUserService {
    return SQLUserService{
        ctx: ctx,
    }
}

func (u SQLUserService) SaveUser(user *User) error {
    query := "INSERT INTO users(email) VALUES (:email)"
    _, err := u.ctx.DB.NamedExec(query, user)
    return err
}

func (u SQLUserService) GetUsers() ([]User, error) {
    var users []User
    query := "select * from users";

    err := u.ctx.DB.Select(&users, query)
    return users, err
}

func (u SQLUserService) GetUserById(id int) (*User, error) {
    user := &User{}
    query := "SELECT * FROM users WHERE id = $1"

    err := u.ctx.DB.Get(user, query, id)
    return user, err
}

func (u SQLUserService) GetUserByEmail(email string) (*User, error) {
    user := &User{}
    query := "SELECT * FROM users WHERE email = $1"

    err := u.ctx.DB.Get(user, query, email)
    return user, err
}

