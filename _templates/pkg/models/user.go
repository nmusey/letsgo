package models

type User struct {
    ID          int     `json:"id" db:"id"`
    Email       string  `json:"email" db:"email"`
    Username    string  `json:"username" db:"username"`
}

func (u User) Table() string {
    return "users"
}

func (u User) AllColumns() string {
    return ":id, :email, :username"
}
