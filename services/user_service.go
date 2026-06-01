package services

import (
	"fmt"
	"momentia-be/internal/middleware"
	"momentia-be/model"
	"os"

	"momentia-be/repository"

	"golang.org/x/crypto/bcrypt"
)

type RegisterInput struct {
	Username string `json:"username" binding:"required"`
	Email    string `json:"email" binding:"required,email"`
	Msisdn   string `json:"msisdn" binding:"required"`
	Password string `json:"password" binding:"required,min=6"`
}

type LoginInput struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
}

type UserService interface {
	Register(input RegisterInput) (*model.User, error)
	Login(input LoginInput) (string, error)
	GetUserByID(id int) (*model.User, error)
	GetUserByEmail(email string) (*model.User, error)
	GetUserByUsername(username string) (*model.User, error)
	GetUserByMsisdn(msisdn string) (*model.User, error)
	Logout(id int) (*model.User, error)
	UpdateUser(user *model.User) error
}

type userService struct {
	repo repository.UserRepository
}

func NewUserService(repo repository.UserRepository) UserService {
	return &userService{repo}
}

func (s *userService) Register(input RegisterInput) (*model.User, error) {
	if existing, err := s.repo.GetUserByUsername(input.Username); err != nil {
		return nil, err
	} else if existing != nil {
		return nil, fmt.Errorf("username already taken")
	}

	if existing, err := s.repo.GetUserByEmail(input.Email); err != nil {
		return nil, err
	} else if existing != nil {
		return nil, fmt.Errorf("email already registered")
	}

	if existing, err := s.repo.GetUserByMsisdn(input.Msisdn); err != nil {
		return nil, err
	} else if existing != nil {
		return nil, fmt.Errorf("msisdn already registered")
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(input.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, fmt.Errorf("failed to hash password: %w", err)
	}

	user := &model.User{
		Username:     input.Username,
		Email:        input.Email,
		Msisdn:       input.Msisdn,
		PasswordHash: string(hash),
	}

	if err := s.repo.Register(user); err != nil {
		return nil, fmt.Errorf("failed to create user: %w", err)
	}
	return user, nil
}

func (s *userService) Login(input LoginInput) (string, error) {
	if input.Email == "" || input.Password == "" {
		return "", fmt.Errorf("email and password are required")
	}

	user, err := s.repo.Login(input.Email)
	if err != nil {
		return "", fmt.Errorf("failed to fetch user: %w", err)
	}
	if user.ID == 0 {
		return "", fmt.Errorf("invalid email or password")
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(input.Password)); err != nil {
		return "", fmt.Errorf("invalid email or password")
	}

	secret := os.Getenv("JWT_SECRET")
	if secret == "" {
		return "", fmt.Errorf("server misconfiguration: JWT secret not set")
	}

	token, err := middleware.GenerateToken(user.ID, secret, 24)
	if err != nil {
		return "", fmt.Errorf("failed to generate token: %w", err)
	}

	return token, nil
}

func (s *userService) Logout(id int) (*model.User, error) {
	return s.repo.Logout(id)
}

func (s *userService) GetUserByID(id int) (*model.User, error) {
	return s.repo.GetUserByID(id)
}

func (s *userService) GetUserByEmail(email string) (*model.User, error) {
	return s.repo.GetUserByEmail(email)
}

func (s *userService) GetUserByUsername(username string) (*model.User, error) {
	return s.repo.GetUserByUsername(username)
}

func (s *userService) GetUserByMsisdn(msisdn string) (*model.User, error) {
	return s.repo.GetUserByMsisdn(msisdn)
}

func (s *userService) UpdateUser(user *model.User) error {
	return s.repo.UpdateUser(user)
}
