package auth

import (
	"database/sql"
)

type SQLAuthStore struct {
	db *sql.DB
}

func newAuthStore(db *sql.DB) SQLAuthStore {
	return SQLAuthStore{db: db}
}

func (as SQLAuthStore) savePassword(hashed *Password) error {
    query := "insert into passwords(password, user_id) values ($1, $2) returning id;"

    row := as.db.QueryRow(query, hashed.Password, hashed.UserId)
    if err := row.Err(); err != nil {
        return err
    }

    var password Password
    err := row.Scan(&password.Id)
    if err != nil {
        return err
    }

    return nil
}

func (as SQLAuthStore) getPassword(userId int) (*Password, error) {
    query := "select id, password, user_id from passwords where user_id=$1;"

    row := as.db.QueryRow(query, userId)
    if err := row.Err(); err != nil {
        return nil, err
    }

    var password Password
    err := row.Scan(&password.Id, &password.Password, &password.UserId)
    if err != nil {
        return nil, err
    }

    return &password, nil
}
