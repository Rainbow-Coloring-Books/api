package user

import (
	"context"

	"github.com/go-playground/validator/v10"
	"golang.org/x/crypto/bcrypt"

	"rainbowcoloringbooks/internal/model"
	repo "rainbowcoloringbooks/internal/repository/user"
)

type UserService interface {
	Register(email, password string) (model.User, error)
}

type userService struct {
	validate *validator.Validate
	repo     repo.UserRepository
}

func NewUserService(validate *validator.Validate, repo repo.UserRepository) UserService {
	return &userService{validate: validate, repo: repo}
}

func (s *userService) Register(email, password string) (model.User, error) {
	ctx := context.Background()

	req := struct {
		Email    string `validate:"required,email"`
		Password string `validate:"required,min=8"`
	}{
		Email:    email,
		Password: password,
	}

	if err := s.validate.Struct(req); err != nil {
		return model.User{}, err
	}

	hashedPassword, err := hashPassword(password)

	if err != nil {
		return model.User{}, err
	}

	user := model.User{
		Email:    email,
		Password: hashedPassword,
	}

	createdUser, err := s.repo.CreateUser(ctx, &user)

	if err != nil {
		return model.User{}, err
	}

	return *createdUser, nil
}

func hashPassword(password string) (string, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(hashedPassword), nil
}
