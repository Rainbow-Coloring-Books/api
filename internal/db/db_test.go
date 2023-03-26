package db_test

import (
	"context"
	"testing"

	"rainbowcoloringbooks/internal/db"
	_ "github.com/lib/pq"
)

// TODO - these tests actually connect to the database. We should mock the database connection.

func TestConnect(t *testing.T) {
	testCases := []struct {
		name    string
		config  db.PostgresDatabase
		wantErr bool
	}{
		{
			name: "valid connection",
			config: db.PostgresDatabase{
				User:     "saus",
				Password: "postgres",
				DBName:   "rainbow-coloring-books",
				SSLMode:  "disable",
			},
			wantErr: false,
		},
		{
			name: "invalid credentials",
			config: db.PostgresDatabase{
				User:     "invalid_user",
				Password: "invalid_password",
				DBName:   "invalid_db_name",
				SSLMode:  "disable",
			},
			wantErr: true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			ctx := context.Background()
			err := tc.config.Connect(ctx)
			if tc.wantErr {
				if err == nil {
					t.Errorf("expected an error but got nil")
				}
			} else {
				if err != nil {
					t.Errorf("unexpected error: %v", err)
				}
			}
			if err == nil {
				_ = tc.config.Close()
			}
		})
	}
}
