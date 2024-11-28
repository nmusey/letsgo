package auth

import (
	"errors"
	"net/http"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"

	"$appRepo/internal/core"
	"$appRepo/internal/users"
)

const (
    PasswordHashCost = 14

    PasswordMinLength = 8
)

var (
    PasswordTooShortError = errors.New("Password is too short")
    InvalidEmailError = errors.New("Invalid email")
    UserNotFoundError = errors.New("User not found")
)

type AuthService struct {
	Store       AuthStore
	UserService *users.UserService
}

func NewAuthService(router *core.Router) *AuthService {
	return &AuthService{
		Store:       newAuthStore(router.DB.SqlDB),
		UserService: users.NewUserService(router),
	}
}

type AuthStore interface {
	savePassword(hashed *Password) error
	getPassword(userId int) (*Password, error)
}

func (as AuthService) StorePassword(password string, userId int) error {
	hashed, err := as.hashPassword(password)
	if err != nil {
		return err
	}

	passwordObj := Password{
		Password: hashed,
		UserId:   userId,
	}

	return as.Store.savePassword(&passwordObj)
}

func (as AuthService) CheckPassword(password string, userId int) (bool, error) {
	passwordObj, err := as.Store.getPassword(userId)
	if err != nil {
		return false, err
	}

	return as.checkPasswordHash(password, passwordObj.Password), nil
}

func (as AuthService) GetAuthCookie(user *users.User) (*http.Cookie, error) {
	expiry := time.Now().Add(time.Hour * 24)
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"uid": user.Id,
		"exp": expiry.Unix(),
	})

	tokenString, err := token.SignedString([]byte(os.Getenv("JWT_SECRET")))
	if err != nil {
		return nil, err
	}

	cookie := &http.Cookie{
		Name:    "Authorization",
		Value:   tokenString,
		Expires: expiry,
	}

	return cookie, nil
}

func (as AuthService) hashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), PasswordHashCost)
	return string(bytes), err
}

func (as AuthService) checkPasswordHash(password string, hash string) bool {
    err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
    return err == nil
}

