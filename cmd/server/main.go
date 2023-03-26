package main

import (
	"context"
	"log"
	"net/http"

	"rainbowcoloringbooks/internal"
	"rainbowcoloringbooks/internal/db"
	userHandler "rainbowcoloringbooks/internal/handler/user"
	userRepo "rainbowcoloringbooks/internal/repository/user"
	userService "rainbowcoloringbooks/internal/service/user"

	"github.com/gorilla/mux"

	_ "github.com/lib/pq"
)

func main() {
	ctx := context.Background()
	config, err := internal.LoadConfig("config.yaml")
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	postgresDB := &db.PostgresDatabase{
		User:     config.DBUser,
		Password: config.DBPassword,
		DBName:   config.DBName,
		SSLMode:  "disable",
	}

	err = postgresDB.Connect(ctx)

	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	defer postgresDB.Close()

	userRepo := userRepo.NewPostgresUserRepository(postgresDB.DB)
	userService := userService.NewUserService(userRepo)
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
