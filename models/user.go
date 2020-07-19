package models

import (
	"database/sql"
	"errors"
	"gloom/storage"
)

type User struct {
	Id         int64
	Name       string
	InsertDate int64
}

func LoadUser(db *sql.DB, username string, password string) (*User, error) {
	query := `
		SELECT
			u.id,
			u.name,
			u.insert_date
		FROM
			user u
		WHERE
			u.deleted_date IS NULL
			AND u.name = ?
			AND u.password = ?
	`

	row := storage.PreparedQueryRow(
		db, query,
		username, password,
	)

	var user User
	err := row.Scan(&user.Id, &user.Name, &user.InsertDate)

	if err != nil {
		return nil, errors.New("User doesn't exist")
	}

	return &user, nil
}
