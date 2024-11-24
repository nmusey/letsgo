package users

import (
	"errors"
)

type LocalUserStore struct {
    ShouldError bool
}

func (us LocalUserStore) SaveUser(user *User) error {
    return us.error()
}

func (us LocalUserStore) GetUsers() ([]User, error) {
    return []User{}, us.error()
}

func (us LocalUserStore) GetUserByEmail(email string)(*User, error) {
    return &User{}, us.error()
}

func (us LocalUserStore) GetUserById(id int) (*User, error) {
    return &User{}, us.error()
}

func (us LocalUserStore) error() error {
    if us.ShouldError {
        errors.New("Test error")
    }

    return nil
}
