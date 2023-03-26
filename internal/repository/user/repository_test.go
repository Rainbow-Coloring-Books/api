package repo_test

import (
	"context"
	"errors"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"github.com/lib/pq"

	"rainbowcoloringbooks/internal/model/user"
	repo "rainbowcoloringbooks/internal/repository/user"
)

func setupValidCreation(t *testing.T, sqlMock sqlmock.Sqlmock, tc *struct {
	email    string
	password string
}) {
	rows := sqlmock.NewRows([]string{"id", "email", "password"}).AddRow(1, tc.email, tc.password)
	sqlMock.ExpectQuery("INSERT INTO users").WithArgs(tc.email, tc.password).WillReturnRows(rows)
}

func setupUniqueEmailConstraintViolation(t *testing.T, sqlMock sqlmock.Sqlmock, tc *struct {
	email    string
	password string
}) {
	sqlMock.ExpectQuery("INSERT INTO users").
		WithArgs(tc.email, tc.password).
		WillReturnError(&pq.Error{Code: "23505"}) // 23505 is the PostgreSQL error code for unique constraint violation
}

func setupDatabaseConnectionError(t *testing.T, sqlMock sqlmock.Sqlmock, tc *struct {
	email    string
	password string
}) {
	sqlMock.ExpectQuery("INSERT INTO users").
		WithArgs(tc.email, tc.password).
		WillReturnError(errors.New("database connection error"))
}

func TestCreateUser(t *testing.T) {
	testCases := []struct {
		name      string
		email     string
		password  string
		expectErr bool
		setupFunc func(*testing.T, sqlmock.Sqlmock, *struct {
			email    string
			password string
		})
	}{
		{
			name:      "valid creation",
			email:     "test@example.com",
			password:  "TestPass123!",
			expectErr: false,
			setupFunc: setupValidCreation,
		},
		{
			name:      "unique email constraint violation",
			email:     "test@example.com",
			password:  "TestPass123!",
			expectErr: true,
			setupFunc: setupUniqueEmailConstraintViolation,
		},
		{
			name:      "database connection error",
			email:     "test@example.com",
			password:  "TestPass123!",
			expectErr: true,
			setupFunc: setupDatabaseConnectionError,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			mockDB, sqlMock, err := sqlmock.New()
			if err != nil {
				t.Fatalf("failed to create sqlmock: %v", err)
			}
			defer mockDB.Close()

			tc.setupFunc(t, sqlMock, &struct {
				email    string
				password string
			}{email: tc.email, password: tc.password})

			testRepo := repo.NewPostgresUserRepository(mockDB)

			_, err = testRepo.CreateUser(context.Background(), &model.User{Email: tc.email, Password: tc.password, ID: 1})

			if tc.expectErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}