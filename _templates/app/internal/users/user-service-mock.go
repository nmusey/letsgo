package users

import (
    "fmt"
)

type MockUserService struct {
    ShouldError bool
}

func (service MockUserService) SaveUser(user *User) error {
    return service.getError()
}

func (service MockUserService) GetUsers() ([]User, error) {
    users := []User{}
    return users, service.getError()
}

func (service MockUserService) GetUserById(id int) (*User, error) {
    return &User{Id: id}, service.getError()
}

func (service MockUserService) GetUserByEmail(email string) (*User, error) {
    return &User{Email: email}, service.getError()
}

func (service MockUserService) getError() error {
    if service.ShouldError {
        return fmt.Errorf("Test error from MockUserService")
    }

    return nil
}
