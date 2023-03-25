package user

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/go-playground/validator/v10"
	"github.com/gorilla/mux"

	userService "rainbowcoloringbooks/internal/service/user"
)

type UserHandler struct {
	userService userService.UserService
}

func NewUserHandler(userService userService.UserService) *UserHandler {
	return &UserHandler{userService: userService}
}

func (h *UserHandler) Register(w http.ResponseWriter, r *http.Request) {
	var req struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		http.Error(w, "Invalid input, please provide valid email and password", http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	user, err := h.userService.Register(r.Context(), req.Email, req.Password)

	if err != nil {
		if errors.Is(err, validator.ValidationErrors{}) {
			http.Error(w, "Invalid input, please provide valid email and password", http.StatusBadRequest)
		} else {
			http.Error(w, "Failed to register user", http.StatusInternalServerError)
		}
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"message": "User Created Successfully",
		"user_id": user.ID,
	})

	if err != nil {
		http.Error(w, "Failed to encode response", http.StatusInternalServerError)
		return
	}
}

func (h *UserHandler) RegisterRoutes(r *mux.Router) {
	r.HandleFunc("/register", h.Register).Methods(http.MethodPost)
}