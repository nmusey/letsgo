package models

type User struct {
    ID          int     `json:"id"`
    Email       string  `json:"email"`
    Username    string  `json:"username"`
}

func (u User) Table() string {
    return "users"
}

func (u User) Columns() string {
    return "id, email, username"
}

func (u User) ColumnValues() []interface{} {
    return []interface{}{u.ID, u.Email, u.Username}
}

func (u User) Populate() []interface{} {
    return []interface{}{&u.ID, &u.Email, &u.Username}
}
