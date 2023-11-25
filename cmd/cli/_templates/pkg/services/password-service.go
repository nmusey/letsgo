package services

import (
    "database/sql"

    "golang.org/x/crypto/bcrypt"
	"$appRepo/pkg/core"
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

    return core.Write(s.ctx.DB, "INSERT INTO passwords (user_id, password) VALUES ($1, $2)", userId, hashed)
}

func (s PasswordService) CheckPassword(password string, userId int) (bool, error) {
    var hashed string
    err := core.Read(s.ctx.DB, "SELECT password FROM passwords WHERE user_id=$1", func (row *sql.Rows) error {
        if err := row.Scan(&hashed); err != nil {
            return err
        }

        return nil
    }, userId)

    if err != nil {
        return false, err
    }

    return s.checkPasswordHash(password, hashed), nil
}

func (s PasswordService) hashPassword(password string) (string, error) {
    bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
    return string(bytes), err
}

func (s PasswordService) checkPasswordHash(password, hash string) bool {
    err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
    return err == nil
}
