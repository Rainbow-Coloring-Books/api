package main

import (
	"log"
	"net/http"

	"rainbowcoloringbooks/internal"
	"rainbowcoloringbooks/internal/db"
	userHandler "rainbowcoloringbooks/internal/handler/user"
	userRepo "rainbowcoloringbooks/internal/repository/user"
	userService "rainbowcoloringbooks/internal/service/user"

	"github.com/go-playground/validator/v10"
	"github.com/gorilla/mux"
	_ "github.com/lib/pq"
)

type User struct {
	ID       int    `json:"id"`
	Email    string `json:"email"`
	Password string `json:"-"`
}

var validate *validator.Validate

func main() {
	config, err := internal.LoadConfig("config.yaml")
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	postgresDB := &db.PostgresDatabase{
		User:     config.DBUser,
		Password: config.DBPassword,
		DBName:   config.DBName,
		SSLMode: "disable",
	}

	err = postgresDB.Connect()

	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	defer postgresDB.Close()

	validate = validator.New()
	userRepo := userRepo.NewPostgresUserRepository(postgresDB)
	userService := userService.NewUserService(validate, userRepo)
	userHandler := userHandler.NewUserHandler(userService)

	r := mux.NewRouter()
	userHandler.RegisterRoutes(r)

	http.Handle("/", r)

	address := ":8080"
	log.Printf("Server starting at %s", address)
	if err := http.ListenAndServe(address, r); err != nil {
		log.Fatalf("Server error: %v", err)
	}
}
