package main

import (
	"log"
	"net/http"

	// TODO add actual paths
	userHandler "rainbowcoloringbooks/internal/handler/user"
	userService "rainbowcoloringbooks/internal/service/user"

	"github.com/go-playground/validator/v10"
	"github.com/gorilla/mux"
)

type User struct {
	ID       int    `json:"id"`
	Email    string    `json:"email"`
	Password string `json:"-"`
}

var validate *validator.Validate

func main() {
	validate = validator.New()
	userService := userService.NewUserService(validate)
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