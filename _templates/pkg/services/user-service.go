package services

import (
	"$appRepo/pkg/core"
	"$appRepo/pkg/models"
)

type UserService interface {
    SaveUser(*models.User) error
    GetUsers() ([]models.User, error)
    GetUserByID(int) (*models.User, error)
    GetUserByEmail(string) (*models.User, error)
}

type SQLUserService struct {
    ctx *core.RouterContext
}

func NewUserService(ctx *core.RouterContext) SQLUserService {
    return SQLUserService{
        ctx: ctx,
    }
}

func (u SQLUserService) SaveUser(user *models.User) error {
    query := "INSERT INTO users(email) VALUES (:email)"
    _, err := u.ctx.DB.NamedExec(query, user)
    return err
}

func (u SQLUserService) GetUsers() ([]models.User, error) {
    var users []models.User
    query := "select * from users";

    err := u.ctx.DB.Select(&users, query)
    return users, err
}

func (u SQLUserService) GetUserByID(id int) (*models.User, error) {
    user := &models.User{}
    query := "SELECT * FROM users WHERE id = $1"

    err := u.ctx.DB.Get(user, query, id)
    return user, err
}

func (u SQLUserService) GetUserByEmail(email string) (*models.User, error) {
    user := &models.User{}
    query := "SELECT * FROM users WHERE email = $1"

    err := u.ctx.DB.Get(user, query, email)
    return user, err
}

