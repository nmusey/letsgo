package auth

import (
    "fmt"
)

type MockPasswordService struct {
    ShouldError bool
}

func (service MockPasswordService) SavePassword(password string, userId int) error {
    return service.getError()
}


func (service MockPasswordService) CheckPassword(password string, userId int) (bool, error) {
    return !service.ShouldError, service.getError()
}

func (service MockPasswordService) getError() error {
    if service.ShouldError {
        return fmt.Errorf("Test error from MockPasswordService")
    }

    return nil
}
