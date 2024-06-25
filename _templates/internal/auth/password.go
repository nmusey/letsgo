package auth 

type Password struct {
    Id          int    `json:"id" db:"id"`
    Password    string `json:"password" db:"password"`
    UserId      int    `json:"user_id" db:"user_id"`
}
