package db

import (
	"database/sql"
)

func ConnectToPostgres() (*sql.DB, error) {
	connStr := "user=saus password=postgres dbname=rainbow-coloring-books sslmode=disable"

	db, err := sql.Open("postgres", connStr)
	if err != nil {
		return nil, err
	}

	err = db.Ping()
	if err != nil {
		return nil, err
	}

	return db, nil
}
