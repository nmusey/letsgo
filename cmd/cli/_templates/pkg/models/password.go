package models

type Password struct {
    ID          int    `json:"id"`
    Password    string `json:"password"`
    UserID      int    `json:"user_id"`
}

func (p Password) Table() string {
    return "passwords"
}

func (p Password) Columns() string {
    return "id, password"
}

func (p Password) ColumnValues() []interface{} {
    return []interface{}{p.ID, p.Password, p.UserID}
}

func (p Password) Populate() []interface{} {
    return []interface{}{&p.ID, &p.Password, &p.UserID}
}
