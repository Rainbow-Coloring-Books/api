package user_test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"database/sql"
	"log"
	"os"
	_ "github.com/lib/pq"
	
	"github.com/gorilla/mux"
	"github.com/stretchr/testify/assert"
	"rainbowcoloringbooks/internal"
	"rainbowcoloringbooks/internal/db"
	service "rainbowcoloringbooks/internal/service/user"
	userHandler "rainbowcoloringbooks/internal/handler/user"
)

var (
    postgresDB  *db.PostgresDatabase
    userService service.UserService
)

func setupTables(db *sql.DB) error {
	_, err := db.Exec(`CREATE TABLE IF NOT EXISTS users (
		id SERIAL PRIMARY KEY,
		email TEXT NOT NULL UNIQUE,
		password TEXT NOT NULL
	);`)
	return err
}

func dropTables(db *sql.DB) error {
	_, err := db.Exec(`DROP TABLE IF EXISTS users;`)
	return err
}

func TestMain(m *testing.M) {
    _, postgresDB, userService = internal.SetupApp("../../../config_test.yaml")
    defer postgresDB.Close()

	if err := dropTables(postgresDB.DB); err != nil {
        log.Fatalf("Failed to drop test tables: %v", err)
    }

    if err := setupTables(postgresDB.DB); err != nil {
        log.Fatalf("Failed to set up test tables: %v", err)
    }

    exitCode := m.Run()

    if err := dropTables(postgresDB.DB); err != nil {
        log.Fatalf("Failed to drop test tables: %v", err)
    }

    os.Exit(exitCode)
}

func TestRegisterIntegration(t *testing.T) {
    testCases := []struct {
        name               string
        email              string
        password           string
        expectedStatusCode int
    }{
        {
            name:               "successful registration",
            email:              "test@example.com",
            password:           "TestPass123!",
            expectedStatusCode: http.StatusCreated,
        },
        {
            name:               "invalid input",
            email:              "invalid_email",
            password:           "short",
            expectedStatusCode: http.StatusBadRequest,
        },
    }

    userH := userHandler.NewUserHandler(userService)

    r := mux.NewRouter()
    userH.RegisterRoutes(r)

    for _, tc := range testCases {
        t.Run(tc.name, func(t *testing.T) {
            body, _ := json.Marshal(map[string]string{
                "email":    tc.email,
                "password": tc.password,
            })
            req, _ := http.NewRequest(http.MethodPost, "/register", bytes.NewReader(body))
            resp := httptest.NewRecorder()
            r.ServeHTTP(resp, req)

            assert.Equal(t, tc.expectedStatusCode, resp.Code)
        })
    }
}
