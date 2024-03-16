package models

type Password struct {
    ID          int    `json:"id" db:"id"`
    Password    string `json:"password" db:"password"`
    UserID      int    `json:"user_id" db:"user_id"`
}
