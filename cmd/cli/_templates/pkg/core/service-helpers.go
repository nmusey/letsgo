package core

import (
    "database/sql"
)

func ReadOne(
    db *sql.DB, 
    query string, 
    callback func(*sql.Row) error,
    arguments ...any, 
) error {
    row := db.QueryRow(query, arguments...)
    if row.Err() != nil {
        return row.Err()
    }

    return callback(row)
}

func Read(
    db *sql.DB, 
    query string, 
    callback func(*sql.Rows) error,
    arguments ...any, 
) error {
    rows, err := db.Query(query, arguments...)
	if err != nil {
		return err
	}
    defer rows.Close()

	for rows.Next() {
        if err := callback(rows); err != nil {
            return err
        }
	}

	if err = rows.Err(); err != nil {
		return err
	}

    return nil
}

func Write(
    db *sql.DB, 
    query string, 
    arguments ...any,
) error {
    _, err := db.Exec(query, arguments...)
    if err != nil {
        return err
    }

    return nil
}
