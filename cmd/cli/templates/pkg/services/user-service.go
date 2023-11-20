package services

import (
	"database/sql"

	"$appRepo/pkg/core"
    "$appRepo/pkg/models"
)

type UserService struct {
    DB *sql.DB
}

func NewService(ctx *core.RouterContext) *UserService {
    return &UserService{
        DB: ctx.DB,
    }
}

func (u UserService) SaveUser(user models.User) error {
    return core.Write(u.DB, "INSERT INTO users (username) VALUES ($1)", user.Username)
}

func (u UserService) GetUsers() ([]models.User, error) {
    var users []models.User
    return users, core.Read(u.DB, "SELECT * FROM users", func(rows *sql.Rows) error { 
        var user models.User
        if err := rows.Scan(&user.ID, &user.Username); err != nil {
            return err
        }

        users = append(users, user)
        return nil
    })
}

func (u UserService) GetUserByID(id int) (models.User, error) {
    var user models.User
    return user, core.ReadOne(u.DB, "SELECT * FROM users WHERE id = $1", func(row *sql.Row) error {
        if err := row.Scan(&user.ID, &user.Username); err != nil {
            return err
        }

        return nil
    }, id)
}
