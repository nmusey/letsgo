package services

import (
	"$appRepo/pkg/core"
	"$appRepo/pkg/models"
	"fmt"
)

type UserService struct {
    ctx *core.RouterContext
}

func NewUserService(ctx *core.RouterContext) UserService {
    return UserService{
        ctx: ctx,
    }
}

func (u UserService) SaveUser(user models.User) error {
    query := fmt.Sprintf("INSERT INTO %s VALUES (%s)", user.Table(), user.AllColumns())
    _, err := u.ctx.DB.NamedExec(query, user)
    return err
}

func (u UserService) GetUsers() ([]models.User, error) {
    user := models.User{}
    query := fmt.Sprintf("SELECT %s FROM %s", user.AllColumns(), user.Table())

    var users []models.User
    err := u.ctx.DB.Select(&users, query)
    return users, err
}

func (u UserService) GetUserByID(id int) (models.User, error) {
    user := models.User{}
    query := fmt.Sprintf("SELECT %s FROM %s WHERE id = $1", user.AllColumns(), user.Table())

    err := u.ctx.DB.SelectOne(&user, query, id)
    return user, err
}

func (u UserService) GetUserByEmail(email string) (models.User, error) {
    user := models.User{}
    query := fmt.Sprintf("SELECT %s FROM %s WHERE email = $1", user.AllColumns(), user.Table())

    err := u.ctx.DB.SelectOne(&user, query, email)
    return user, err
}

