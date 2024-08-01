package users

import (
	"$appRepo/internal/core"
)

type SQLUserService struct {
    router core.Router
}

func NewUserService(router core.Router) SQLUserService {
    return SQLUserService{
        router: router,
    }
}

func (u SQLUserService) SaveUser(user *User) error {
    query := "INSERT INTO users(email) VALUES (:email)"
    _, err := u.router.DB.NamedExec(query, user)
    return err
}

func (u SQLUserService) GetUsers() ([]User, error) {
    var users []User
    query := "select * from users";

    err := u.router.DB.Select(&users, query)
    return users, err
}

func (u SQLUserService) GetUserById(id int) (*User, error) {
    user := &User{}
    query := "SELECT * FROM users WHERE id = $1"

    err := u.router.DB.Get(user, query, id)
    return user, err
}

func (u SQLUserService) GetUserByEmail(email string) (*User, error) {
    user := &User{}
    query := "SELECT * FROM users WHERE email = $1"

    err := u.router.DB.Get(user, query, email)
    return user, err
}

