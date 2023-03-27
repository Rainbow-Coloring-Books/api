package main

import (
	"log"
	"net/http"

	"rainbowcoloringbooks/internal"
	userHandler "rainbowcoloringbooks/internal/handler/user"

	"github.com/gorilla/mux"

	_ "github.com/lib/pq"
)

func main() {
	_, postgresDB, userService := internal.SetupApp("config.yaml")
	defer postgresDB.Close()

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
