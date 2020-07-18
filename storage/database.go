package storage

import (
	"database/sql"
	"os"

	_ "github.com/mattn/go-sqlite3"
)

var Db *sql.DB

type Queryer interface {
	Query(string, ...interface{}) (*sql.Rows, error)
	QueryRow(string, ...interface{}) *sql.Row
	Prepare(string) (*sql.Stmt, error)
	Exec(string, ...interface{}) (sql.Result, error)
}

func InitDB(filepath string) *sql.DB {
	if _, err := os.Stat(filepath); err == nil {
		os.Remove(filepath)
	}
	db, err := sql.Open("sqlite3", filepath)

	if err != nil {
		panic(err)
	}

	return db
}

func CreateSchema(db *sql.DB) {
	_, err := db.Exec(schema)

	if err != nil {
		panic(err)
	}
}

func PreparedExec(db Queryer, query string, args ...interface{}) (int64, error) {
	stmt, err := db.Prepare(query)
	if err != nil {
		return 0, err
	}

	res, err := stmt.Exec(args...)
	if err != nil {
		return 0, err
	}

	insertId, err := res.LastInsertId()
	if err == nil {
		return insertId, nil
	}

	affectedRows, err := res.RowsAffected()
	if err == nil {
		return affectedRows, nil
	}

	return 0, nil
}

func PreparedQuery(db Queryer, query string, args ...interface{}) *sql.Rows {
	stmt, err := db.Prepare(query)
	if err != nil {
		panic(err)
	}

	rows, err := stmt.Query(args...)
	if err != nil {
		panic(err)
	}

	return rows
}

func PreparedQueryRow(db Queryer, query string, args ...interface{}) *sql.Row {
	stmt, err := db.Prepare(query)
	if err != nil {
		panic(err)
	}

	return stmt.QueryRow(args...)
}
