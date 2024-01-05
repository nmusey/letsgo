package services

import (
    "fmt"

    "golang.org/x/crypto/bcrypt"
	"$appRepo/pkg/core"
    "$appRepo/pkg/models"
)

type PasswordService struct {
    ctx *core.RouterContext
}

func NewPasswordService(ctx *core.RouterContext) *PasswordService {
    return &PasswordService{
        ctx: ctx,
    }
}

func (s PasswordService) SavePassword(password string, userId int) error {
    hashed, err := s.hashPassword(password)
    if err != nil {
        return err
    }

    passwordObj := models.Password{
        Password: hashed,
        UserID: userId,
    }

    query := fmt.Sprintf("INSERT INTO %s (%s) VALUES $1", passwordObj.Table(), passwordObj.AllColumns())
    _, err = s.ctx.DB.NamedExec(query, passwordObj)
    return err
}

func (s PasswordService) CheckPassword(password string, userId int) (bool, error) {
    passwordObj := models.Password{}
    if err := s.ctx.DB.SelectOne(&passwordObj, "user_id = $1", userId); err != nil {
        return false, err
    }

    return s.checkPasswordHash(password, passwordObj.Password), nil
}

func (s PasswordService) hashPassword(password string) (string, error) {
    bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
    return string(bytes), err
}

func (s PasswordService) checkPasswordHash(password, hash string) bool {
    err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
    return err == nil
}
