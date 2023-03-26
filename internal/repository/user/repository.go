//go:generate mockgen -destination=mocks.go -package=repo rainbowcoloringbooks/internal/repository/user UserRepository

package repo

import (
	"context"
	"database/sql"

	"rainbowcoloringbooks/internal/model/user"
)

type UserRepository interface {
	CreateUser(ctx context.Context, user *model.User) (*model.User, error)
	FindUserByEmail(ctx context.Context, email string) (*model.User, error)
}

type postgresUserRepository struct {
	db *sql.DB
}

const (
	createUserQuery      = "INSERT INTO users (email, password) VALUES ($1, $2) RETURNING id, email, password"
	findUserByEmailQuery = "SELECT id, email, password FROM users WHERE email = $1"
)

func NewPostgresUserRepository(db *sql.DB) UserRepository {
	return &postgresUserRepository{db: db}
}

func (r *postgresUserRepository) CreateUser(ctx context.Context, user *model.User) (*model.User, error) {
	err := r.db.QueryRowContext(ctx, createUserQuery, user.Email, user.Password).Scan(&user.ID, &user.Email, &user.Password)
	if err != nil {
		return nil, err
	}

	return user, nil
}

func (r *postgresUserRepository) FindUserByEmail(ctx context.Context, email string) (*model.User, error) {
	row := r.db.QueryRowContext(ctx, findUserByEmailQuery, email)

	var user model.User

	err := row.Scan(&user.ID, &user.Email, &user.Password)

	if err == sql.ErrNoRows {
		return nil, nil
	}

	if err != nil {
		return nil, err
	}

	return &user, nil
}
