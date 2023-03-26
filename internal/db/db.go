//go:generate mockgen -destination=mocks.go -package=db rainbowcoloringbooks/internal/db Database

package db

import (
	"context"
	"database/sql"
	"fmt"
)

type Database interface {
	Connect(ctx context.Context) error
	Close() error
	Exec(ctx context.Context, query string, args ...interface{}) (sql.Result, error)
	Query(ctx context.Context, query string, args ...interface{}) (*sql.Rows, error)
	QueryRow(ctx context.Context, query string, args ...interface{}) *sql.Row
	QueryRowContext(ctx context.Context, query string, args ...interface{}) *sql.Row
}

type PostgresDatabase struct {
	DB       *sql.DB
	Host     string
	Port     string
	User     string
	Password string
	DBName   string
	SSLMode  string
}

func (p *PostgresDatabase) Connect(ctx context.Context) error {
	connStr := fmt.Sprintf(
		"user=%s password=%s dbname=%s sslmode=%s",
		p.User, p.Password, p.DBName, p.SSLMode,
	)

	db, err := sql.Open("postgres", connStr)
	if err != nil {
		return err
	}

	p.DB = db

	return nil
}

func (p *PostgresDatabase) Close() error {
	return p.DB.Close()
}

func (p *PostgresDatabase) Exec(ctx context.Context, query string, args ...interface{}) (sql.Result, error) {
	return p.DB.Exec(query, args...)
}

func (p *PostgresDatabase) Query(ctx context.Context, query string, args ...interface{}) (*sql.Rows, error) {
	return p.DB.Query(query, args...)
}

func (p *PostgresDatabase) QueryRow(ctx context.Context, query string, args ...interface{}) *sql.Row {
	return p.DB.QueryRow(query, args...)
}

func (p *PostgresDatabase) QueryRowContext(ctx context.Context, query string, args ...interface{}) *sql.Row {
	return p.DB.QueryRowContext(ctx, query, args...)
}
