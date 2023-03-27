package user_test

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/go-playground/validator/v10"
	"github.com/golang/mock/gomock"
	"github.com/gorilla/mux"
	"github.com/stretchr/testify/assert"

	userHandler "rainbowcoloringbooks/internal/handler/user"
	model "rainbowcoloringbooks/internal/model/user"
	userServiceMock "rainbowcoloringbooks/internal/service/user"
)

func TestRegister(t *testing.T) {
	testCases := []struct {
		email              string
		name               string
		password           string
		serviceError       error
		serviceUser        model.User
		expectedStatusCode int
	}{
		{
			name:               "successful registration",
			email:              "test@example.com",
			password:           "TestPass123!",
			serviceError:       nil,
			serviceUser:        model.User{ID: 1, Email: "test@example.com", Password: "TestPass123!"},
			expectedStatusCode: http.StatusCreated,
		},
		{
			name:               "invalid input",
			email:              "invalid_email",
			password:           "short",
			serviceError:       validator.ValidationErrors{},
			serviceUser:        model.User{},
			expectedStatusCode: http.StatusBadRequest,
		},
		{
			name:               "registration error",
			email:              "test@example.com",
			password:           "TestPass123!",
			serviceError:       errors.New("registration error"),
			serviceUser:        model.User{},
			expectedStatusCode: http.StatusInternalServerError,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			mockUserService := userServiceMock.NewMockUserService(ctrl)
			handler := userHandler.NewUserHandler(mockUserService)

			mockUserService.EXPECT().
				Register(gomock.Any(), tc.email, tc.password).
				Return(tc.serviceUser, tc.serviceError)

			r := mux.NewRouter()
			handler.RegisterRoutes(r)

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
