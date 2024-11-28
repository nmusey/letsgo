package users

import "time"

type User struct {
	Id    int    `json:"id" db:"id"`
    CreatedAt time.Time `db:"created_at" json:"-"`
    UpdatedAt time.Time `db:"updated_at" json:"-"`

	Email string `json:"email" db:"email"`
}
