package auth

import (
    "golang.org/x/crypto/bcrypt"
	"$appRepo/pkg/core"
)

type PasswordService interface {
    SavePassword(string, int) error
    CheckPassword(string, int) (bool, error)
}

type SQLPasswordService struct {
    ctx *core.RouterContext
}

func NewPasswordService(ctx *core.RouterContext) *SQLPasswordService {
    return &SQLPasswordService{
        ctx: ctx,
    }
}

func (s SQLPasswordService) SavePassword(password string, userId int) error {
    hashed, err := s.hashPassword(password)
    if err != nil {
        return err
    }

    passwordObj := Password{
        Password: hashed,
        UserId: userId,
    }

    query := "INSERT INTO passwords(user_id, password) VALUES :password, :userId"
    _, err = s.ctx.DB.NamedExec(query, passwordObj)
    return err
}

func (s SQLPasswordService) CheckPassword(password string, userId int) (bool, error) {
    passwordObj := Password{}
    if err := s.ctx.DB.Get(&passwordObj, "user_id = $1", userId); err != nil {
        return false, err
    }

    return s.checkPasswordHash(password, passwordObj.Password), nil
}

func (s SQLPasswordService) hashPassword(password string) (string, error) {
    bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
    return string(bytes), err
}

func (s SQLPasswordService) checkPasswordHash(password, hash string) bool {
    err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
    return err == nil
}
