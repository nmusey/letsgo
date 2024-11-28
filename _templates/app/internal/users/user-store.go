package users

import (
	"database/sql"
)

type SQLUserStore struct {
	db *sql.DB
}

func newSQLUserStore(db *sql.DB) SQLUserStore {
	return SQLUserStore{db: db}
}

func (us SQLUserStore) SaveUser(user *User) error {
	query := "insert into users(email) values ($1) returning id"
	row := us.db.QueryRow(query, user.Email)
	if row.Err() != nil {
		return row.Err()
	}

	row.Scan(user.Id)
	return nil
}

func (us SQLUserStore) GetUsers() ([]User, error) {
	query := "select id, email from users"

	rows, err := us.db.Query(query)
	if err != nil {
		return nil, err
	}

	var users []User
	for rows.Next() {
		var user User
		if err := rows.Scan(user.Id, user.Email); err != nil {
			return nil, err
		}

		users = append(users, user)
	}

	return users, nil
}

func (us SQLUserStore) GetUserById(id int) (*User, error) {
	user := &User{}
	query := "select id, email from users where id = $1"

	row := us.db.QueryRow(query, id)
	if row.Err() != nil {
		return nil, row.Err()
	}

	if err := row.Scan(user.Id, user.Email); err != nil {
		return nil, err
	}

	return user, nil
}

func (us SQLUserStore) GetUserByEmail(email string) (*User, error) {
	user := &User{}
	query := "select id, email from users where email = $1"

	row := us.db.QueryRow(query, email)
	if row.Err() != nil {
		return nil, row.Err()
	}

	if err := row.Scan(user.Id, user.Email); err != nil {
		return nil, err
	}

	return user, nil
}
