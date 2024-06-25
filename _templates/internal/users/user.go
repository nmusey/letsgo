package users 

type User struct {
    Id          int     `json:"id" db:"id"`
    Email       string  `json:"email" db:"email"`
}
