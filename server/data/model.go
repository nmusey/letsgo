package data

import "time"

type Model struct {
    Table string

    ID int64
    UUID string

    Created time.Time
    Updated time.Time
    Deleted time.Time
}
