package repo

import (
	"context"
	"database/sql"

	"rainbowcoloringbooks/internal/db"
	"rainbowcoloringbooks/internal/model/user"
)

type UserRepository interface {
	CreateUser(ctx context.Context, user *model.User) (*model.User, error)
	FindUserByEmail(ctx context.Context, email string) (*model.User, error)
}

type postgresUserRepository struct {
	db db.Database
}

func NewPostgresUserRepository(db db.Database) UserRepository {
	return &postgresUserRepository{db: db}
}

func (r *postgresUserRepository) CreateUser(ctx context.Context, user *model.User) (*model.User, error) {
	query := "INSERT INTO users (email, password) VALUES ($1, $2) RETURNING id"
	err := r.db.QueryRowContext(ctx, query, user.Email, user.Password).Scan(&user.ID)
	if err != nil {
		return nil, err
	}

	return user, nil
}

func (r *postgresUserRepository) FindUserByEmail(ctx context.Context, email string) (*model.User, error) {
	query := "SELECT id, email, password FROM users WHERE email = $1"
	row := r.db.QueryRowContext(ctx, query, email)

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
