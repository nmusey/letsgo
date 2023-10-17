package data

import "time"

type IModel interface {
    GetID() int64
}

type Model struct {
    ID int64
    UUID string

    Created time.Time
    Updated time.Time
    Deleted time.Time
}

func (m Model) GetID() int64 {
    return m.ID
}
