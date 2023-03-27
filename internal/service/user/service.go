//go:generate mockgen -destination=mocks.go -package=user rainbowcoloringbooks/internal/service/user UserService

package user

import (
	"context"
	"errors"

	"github.com/go-playground/validator/v10"
	"golang.org/x/crypto/bcrypt"

	"rainbowcoloringbooks/internal/model/user"
	repo "rainbowcoloringbooks/internal/repository/user"
)

type UserService interface {
	Register(context context.Context, email, password string) (model.User, error)
}

type userService struct {
	validate *validator.Validate
	repo     repo.UserRepository
}

var ErrEmailAlreadyInUse = errors.New("email already in use")

func NewUserService(repo repo.UserRepository) UserService {
	return &userService{
		validate: validator.New(),
		repo: repo,
	}
}

func (s *userService) Register(ctx context.Context, email, password string) (model.User, error) {
	existingUser, _ := s.repo.FindUserByEmail(ctx, email)
	if existingUser != nil {
		return model.User{}, ErrEmailAlreadyInUse
	}

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

// TODO: Move this to a separate package
func hashPassword(password string) (string, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(hashedPassword), nil
}
