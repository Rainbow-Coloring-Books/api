package db

import (
	"database/sql"
)

func ConnectToPostgres(user, password, dbname string) (*sql.DB, error) {
	connStr := "user=" + user + " password=" + password + " dbname=" + dbname + " sslmode=disable"

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
