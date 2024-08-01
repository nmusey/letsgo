package auth

import (
    "golang.org/x/crypto/bcrypt"
	"$appRepo/internal/core"
)

type SQLPasswordService struct {
    router core.Router
}

func NewPasswordService(router core.Router) *SQLPasswordService {
    return &SQLPasswordService{
        router: router,
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
    _, err = s.router.DB.NamedExec(query, passwordObj)
    return err
}

func (s SQLPasswordService) CheckPassword(password string, userId int) (bool, error) {
    passwordObj := Password{}
    if err := s.router.DB.Get(&passwordObj, "user_id = $1", userId); err != nil {
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
