package auth

import (
	"errors"

	"golang.org/x/crypto/bcrypt"
)

type LocalAuthStore struct {
	ShouldError bool
}

func (as LocalAuthStore) savePassword(hashed *Password) error {
	return as.error()
}

func (as LocalAuthStore) getPassword(userId int) (*Password, error) {
    bytes, _ := bcrypt.GenerateFromPassword([]byte("test"), PasswordHashCost)
    return &Password{Password: string(bytes)}, as.error()
}

func (as LocalAuthStore) error() error {
	if as.ShouldError {
		return errors.New("Test error")
	}

	return nil
}
