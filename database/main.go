package database

import (
	"database/sql"

	// prevent 'panic: sql: unknown driver "postgres" (forgotten import?)' message
	_ "github.com/lib/pq"
)

var (
	// schemaVersion = 1
	conninfo     = "user=root dbname=main password=pass sslmode=disable"
	driver       = "postgres"
	testPassword = "password"
)

func init() {
	// query := schema()
	// if _, err := Exec(query); err != nil {
	// 	panic(err)
	// }
}

// Exec executes a query without returning any rows.
func Exec(query string, args ...interface{}) (sql.Result, error) {
	db, err := sql.Open(driver, conninfo)
	if err != nil {
		return nil, err
	}
	defer db.Close()

	res, err := db.Exec(query, args...)
	if err != nil {
		return nil, err
	}

	return res, nil
}

// Query executes a query that returns rows, typically a SELECT.
func Query(query string, args ...interface{}) (*sql.Rows, error) {
	db, err := sql.Open(driver, conninfo)
	if err != nil {
		return nil, err
	}
	defer db.Close()

	rows, err := db.Query(query, args...)
	if err != nil {
		return nil, err
	}

	return rows, nil
}

// QueryRow executes a query that is expected to return at most one row.
func QueryRow(query string, args ...interface{}) (*sql.Row, error) {
	db, err := sql.Open(driver, conninfo)
	if err != nil {
		return nil, err
	}
	defer db.Close()

	row := db.QueryRow(query, args...)

	return row, nil
}
