package user_test

import (
	"context"
	"testing"
	"errors"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"

	model "rainbowcoloringbooks/internal/model/user"
	repo "rainbowcoloringbooks/internal/repository/user"
	user "rainbowcoloringbooks/internal/service/user"
)

func TestRegister(t *testing.T) {
	testCases := []struct {
		name      string
		email     string
		password  string
		expectErr bool
	}{
		{
			name:      "valid registration",
			email:     "test@example.com",
			password:  "TestPass123!",
			expectErr: false,
		},
		{
			name:      "invalid email",
			email:     "badEmail",
			password:  "TestPass123!",
			expectErr: true,
		},
		{
			name:      "no email",
			email:     "",
			password:  "TestPass123!",
			expectErr: true,
		},
		{
			name:      "no password",
			email:     "badEmail",
			password:  "",
			expectErr: true,
		},
		{
			name:      "password too short",
			email:     "test@example.com",
			password:  "short",
			expectErr: true,
		},
		{
			name:      "create user error",
			email:     "test@example.com",
			password:  "TestPass123",
			expectErr: true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			mockUserRepository := repo.NewMockUserRepository(ctrl)

			service := user.NewUserService(mockUserRepository)

			if !tc.expectErr {
				mockUserRepository.EXPECT().CreateUser(gomock.Any(), gomock.Any()).Return(&model.User{Email: tc.email, Password: "hashed_password"}, nil)
			} else {
				mockUserRepository.EXPECT().CreateUser(gomock.Any(), gomock.Any()).Return(nil, errors.New("CreateUser error")).AnyTimes()
			}

			registeredUser, err := service.Register(context.Background(), tc.email, tc.password)

			if tc.expectErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tc.email, registeredUser.Email)
			}
		})
	}
}
