package user

import (
	"github.com/go-playground/validator/v10"
	"golang.org/x/crypto/bcrypt"
)

type User struct {
	ID       int
	Email    string
	Password string
}

type UserService interface {
	Register(email, password string) (User, error)
}

type userService struct {
	validate *validator.Validate
}

func NewUserService(validate *validator.Validate) UserService {
	return &userService{validate: validate}
}

func (s *userService) Register(email, password string) (User, error) {
	req := struct {
		Email    string `validate:"required,email"`
		Password string `validate:"required,min=8"`
	}{
		Email:    email,
		Password: password,
	}

	if err := s.validate.Struct(req); err != nil {
		return User{}, err
	}

	hashedPassword, err := hashPassword(password)
	
	if err != nil {
		return User{}, err
	}

	// create dummy user for now
	user := User{
		ID: 1,
		Email: email,
		Password: hashedPassword,
	}

	return user, nil
}

func hashPassword(password string) (string, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(hashedPassword), nil
}
