package data

import (
    "database/sql"
)

type IRepository[M Model] interface {
    ApplyQuery(query string, ...args string) ([]M, error)

    Create(model M) (M, error)   

    Find() ([]M, error)
    FindByID(id int64) (M, error)
    FindByUUID(uuid string) (M, error)

    Update(model M) (M, error)
    Delete(model M) (M, error)
}

type Repository[M Model] struct {
    db *sql.DB
}

func (r Repository[M]) ApplyQuery(query string, ...args string) ([]M, error) {
    db.Query(query, args)
}

func (r Repository[M]) Create(model M) (M, error) {
    db.Query("INSERT INTO ? VALUES (?)", moel.Table, model)
}

func (r Repository[M]) Find() ([]M, error) {
    db.Query("SELECT * FROM ?", model.Table)
}

func (r Repository[M]) FindByID(id int64) (M, error) {
    db.Query("SELECT * FROM ? WHERE id = ?", model.Table, id)
}

func (r Repository[M]) FindByUUID(uuid string) (M, error) {
    db.Query("SELECT * FROM ? WHERE uuid = ?", model.Table, uuid)
}

func (r Repository[M]) Update(model M) (M, error) {
    db.Query("UPDATE ? SET ? WHERE id = ?", model.Table, model, model.ID)
}

func (r Repository[M]) Delete(model M) (M, error) {
    db.Query("DELETE FROM ? WHERE id = ?", model.Table, model.ID)
}
