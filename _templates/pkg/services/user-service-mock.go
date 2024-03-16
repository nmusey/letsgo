package services

import (
    "fmt"

    "$appRepo/pkg/models"
)

type MockUserService struct {
    ShouldError bool
}

func (service MockUserService) SaveUser(user *models.User) error {
    return service.getError()
}

func (service MockUserService) GetUsers() ([]models.User, error) {
    users := []models.User{}
    return users, service.getError()
}

func (service MockUserService) GetUserByID(id int) (*models.User, error) {
    return &models.User{ID: id}, service.getError()
}

func (service MockUserService) GetUserByEmail(email string) (*models.User, error) {
    return &models.User{Email: email}, service.getError()
}

func (service MockUserService) getError() error {
    if service.ShouldError {
        return fmt.Errorf("Test error from MockUserService")
    }

    return nil
}
