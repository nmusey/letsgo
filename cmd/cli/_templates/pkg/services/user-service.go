package services

import (
	"$appRepo/pkg/core"
    "$appRepo/pkg/models"
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
    return u.ctx.DB.Insert(user)
}

func (u UserService) GetUsers() ([]models.User, error) {
    users, err := u.ctx.DB.Select(users, "1=1")
    return users.([]models.User), err
}

func (u UserService) GetUserByID(id int) (models.User, error) {
    user := models.User{}
    err := u.ctx.DB.SelectOne(&user, "id = $1", id)
    return user, err
}

func (u UserService) GetUserByEmail(email string) (models.User, error) {
    user := models.User{}
    err := u.ctx.DB.SelectOne(&user, "email = $1", email)
    return user, err
}

