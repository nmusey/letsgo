package data

import (
    "database/sql"
)

type IRepository[M IModel] interface {
    ApplyQuery(query string, args ...string) ([]M, error)

    Create(model M) (M, error)   

    Find() ([]M, error)
    FindByID(id int64) (M, error)
    FindByUUID(uuid string) (M, error)

    Update(model M) (M, error)
    Delete(model M) (M, error)
}

type Repository[M Model] struct {
    db *sql.DB
    Table string
}

func (r Repository[M]) ApplyQuery(query string, args ...string) ([]M, error) {
    rows, err := r.db.Query(query, args)
    results := []M{}
    if err != nil {
        return results, err
    }

    rows.Scan(&results)
    return results, nil 
}

func (r Repository[M]) Create(model M) (M, error) {
    rows, err := r.db.Query("INSERT INTO ? VALUES (?)", r.Table , model)
    result := M{}
    if err != nil {
        return result, err
    }

    rows.Scan(&result)
    return result, err
}

func (r Repository[M]) Find() ([]M, error) {
    rows, err := r.db.Query("SELECT * FROM ?", r.Table)
    results := []M{}
    if err != nil {
        return results, err
    }

    rows.Scan(&results)
    return results, err
}

func (r Repository[M]) FindByID(id int64) (M, error) {
    rows, err := r.db.Query("SELECT * FROM ? WHERE id = ?", r.Table, id)
    result := M{}
    if err != nil {
        return result, err
    }

    rows.Scan(&result)
    return result, err
}

func (r Repository[M]) FindByUUID(uuid string) (M, error) {
    rows, err := r.db.Query("SELECT * FROM ? WHERE uuid = ?", r.Table, uuid)
    result := M{}
    if err != nil {
        return result, err
    }

    rows.Scan(&result)
    return result, err
}

func (r Repository[M]) Update(model IModel) (M, error) {
    rows, err := r.db.Query("UPDATE ? SET ? WHERE id = ?", r.Table, model.GetID())
    result := M{}
    if err != nil {
        return result, err
    }

    rows.Scan(&result)
    return result, err
}

func (r Repository[M]) Delete(model IModel) (M, error) {
    rows, err := r.db.Query("DELETE FROM ? WHERE id = ?", r.Table, model.GetID())
    result := M{}
    if err != nil {
        return result, err
    }

    rows.Scan(&result)
    return result, err
}
