package main

import (
	"log"
	"net/http"

	// TODO add actual paths
	"your_project_name/internal/handler/user"
	"your_project_name/internal/service"

	"github.com/gorilla/mux"
	"github.com/go-playground/validator/v10"
)

type User struct {
	ID       int    `json:"id"`
	Email    string    `json:"email"`
	Password string `json:"-"`
}

var validate *validator.Validate

func main() {
	validate = validator.New()
	userService := service.NewUserService(validate)
	userHandler := handler.NewUserHandler(userService)

	r := mux.NewRouter()
	userHandler.RegisterRoutes(r)

	http.Handle("/", r)
	log.Fatal(http.ListenAndServe(":8080", nil))
}