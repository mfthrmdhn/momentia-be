package services_test

import (
	"errors"
	"testing"

	"momentia-be/model"
	"momentia-be/services"

	"golang.org/x/crypto/bcrypt"
)

// stubRepo is a configurable in-memory stub that satisfies repository.UserRepository.
type stubRepo struct {
	registerFn          func(*model.User) error
	loginFn             func(string) (*model.User, error)
	logoutFn            func(int) (*model.User, error)
	getUserByIDFn       func(int) (*model.User, error)
	getUserByEmailFn    func(string) (*model.User, error)
	getUserByUsernameFn func(string) (*model.User, error)
	getUserByMsisdnFn   func(string) (*model.User, error)
	updateUserFn        func(*model.User) error
}

func (s *stubRepo) Register(u *model.User) error {
	if s.registerFn != nil {
		return s.registerFn(u)
	}
	return nil
}
func (s *stubRepo) Login(email string) (*model.User, error) {
	if s.loginFn != nil {
		return s.loginFn(email)
	}
	return &model.User{}, nil
}
func (s *stubRepo) Logout(id int) (*model.User, error) {
	if s.logoutFn != nil {
		return s.logoutFn(id)
	}
	return &model.User{ID: id}, nil
}
func (s *stubRepo) GetUserByID(id int) (*model.User, error) {
	if s.getUserByIDFn != nil {
		return s.getUserByIDFn(id)
	}
	return &model.User{ID: id}, nil
}
func (s *stubRepo) GetUserByEmail(email string) (*model.User, error) {
	if s.getUserByEmailFn != nil {
		return s.getUserByEmailFn(email)
	}
	return &model.User{}, nil
}
func (s *stubRepo) GetUserByUsername(username string) (*model.User, error) {
	if s.getUserByUsernameFn != nil {
		return s.getUserByUsernameFn(username)
	}
	return &model.User{}, nil
}
func (s *stubRepo) GetUserByMsisdn(msisdn string) (*model.User, error) {
	if s.getUserByMsisdnFn != nil {
		return s.getUserByMsisdnFn(msisdn)
	}
	return &model.User{}, nil
}
func (s *stubRepo) UpdateUser(u *model.User) error {
	if s.updateUserFn != nil {
		return s.updateUserFn(u)
	}
	return nil
}

// --- Register ---

func TestRegister_Success(t *testing.T) {
	svc := services.NewUserService(&stubRepo{})

	user, err := svc.Register(services.RegisterInput{
		Username: "alice",
		Email:    "alice@example.com",
		Msisdn:   "+628123456789",
		Password: "secret123",
	})

	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if user.Username != "alice" {
		t.Errorf("expected username %q, got %q", "alice", user.Username)
	}
	if user.PasswordHash == "" {
		t.Error("expected password hash to be set")
	}
}

func TestRegister_DuplicateUsername(t *testing.T) {
	repo := &stubRepo{
		getUserByUsernameFn: func(string) (*model.User, error) {
			return &model.User{ID: 1, Username: "alice"}, nil
		},
	}
	svc := services.NewUserService(repo)

	_, err := svc.Register(services.RegisterInput{
		Username: "alice",
		Email:    "alice@example.com",
		Msisdn:   "+628123456789",
		Password: "secret123",
	})

	if err == nil {
		t.Fatal("expected error, got nil")
	}
	if err.Error() != "username already taken" {
		t.Errorf("unexpected error: %v", err)
	}
}

func TestRegister_DuplicateEmail(t *testing.T) {
	repo := &stubRepo{
		getUserByEmailFn: func(string) (*model.User, error) {
			return &model.User{ID: 2, Email: "alice@example.com"}, nil
		},
	}
	svc := services.NewUserService(repo)

	_, err := svc.Register(services.RegisterInput{
		Username: "alice",
		Email:    "alice@example.com",
		Msisdn:   "+628123456789",
		Password: "secret123",
	})

	if err == nil {
		t.Fatal("expected error, got nil")
	}
	if err.Error() != "email already registered" {
		t.Errorf("unexpected error: %v", err)
	}
}

func TestRegister_DuplicateMsisdn(t *testing.T) {
	repo := &stubRepo{
		getUserByMsisdnFn: func(string) (*model.User, error) {
			return &model.User{ID: 3, Msisdn: "+628123456789"}, nil
		},
	}
	svc := services.NewUserService(repo)

	_, err := svc.Register(services.RegisterInput{
		Username: "alice",
		Email:    "alice@example.com",
		Msisdn:   "+628123456789",
		Password: "secret123",
	})

	if err == nil {
		t.Fatal("expected error, got nil")
	}
	if err.Error() != "msisdn already registered" {
		t.Errorf("unexpected error: %v", err)
	}
}

func TestRegister_RepoError(t *testing.T) {
	repo := &stubRepo{
		registerFn: func(*model.User) error {
			return errors.New("db connection failed")
		},
	}
	svc := services.NewUserService(repo)

	_, err := svc.Register(services.RegisterInput{
		Username: "alice",
		Email:    "alice@example.com",
		Msisdn:   "+628123456789",
		Password: "secret123",
	})

	if err == nil {
		t.Fatal("expected error, got nil")
	}
}

// --- Login ---

func TestLogin_Success(t *testing.T) {
	t.Setenv("JWT_SECRET", "test-secret-key")

	hash, _ := bcrypt.GenerateFromPassword([]byte("password123"), bcrypt.MinCost)
	repo := &stubRepo{
		loginFn: func(email string) (*model.User, error) {
			return &model.User{ID: 1, Email: email, PasswordHash: string(hash)}, nil
		},
	}
	svc := services.NewUserService(repo)

	token, err := svc.Login(services.LoginInput{
		Email:    "alice@example.com",
		Password: "password123",
	})

	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if token == "" {
		t.Error("expected non-empty token")
	}
}

func TestLogin_WrongPassword(t *testing.T) {
	t.Setenv("JWT_SECRET", "test-secret-key")

	hash, _ := bcrypt.GenerateFromPassword([]byte("correctpassword"), bcrypt.MinCost)
	repo := &stubRepo{
		loginFn: func(email string) (*model.User, error) {
			return &model.User{ID: 1, Email: email, PasswordHash: string(hash)}, nil
		},
	}
	svc := services.NewUserService(repo)

	_, err := svc.Login(services.LoginInput{
		Email:    "alice@example.com",
		Password: "wrongpassword",
	})

	if err == nil {
		t.Fatal("expected error, got nil")
	}
	if err.Error() != "invalid email or password" {
		t.Errorf("unexpected error: %v", err)
	}
}

func TestLogin_UserNotFound(t *testing.T) {
	repo := &stubRepo{
		loginFn: func(string) (*model.User, error) {
			return &model.User{}, nil // ID == 0 means not found
		},
	}
	svc := services.NewUserService(repo)

	_, err := svc.Login(services.LoginInput{
		Email:    "ghost@example.com",
		Password: "anypassword",
	})

	if err == nil {
		t.Fatal("expected error, got nil")
	}
	if err.Error() != "invalid email or password" {
		t.Errorf("unexpected error: %v", err)
	}
}

func TestLogin_MissingJWTSecret(t *testing.T) {
	t.Setenv("JWT_SECRET", "")

	hash, _ := bcrypt.GenerateFromPassword([]byte("password123"), bcrypt.MinCost)
	repo := &stubRepo{
		loginFn: func(email string) (*model.User, error) {
			return &model.User{ID: 1, Email: email, PasswordHash: string(hash)}, nil
		},
	}
	svc := services.NewUserService(repo)

	_, err := svc.Login(services.LoginInput{
		Email:    "alice@example.com",
		Password: "password123",
	})

	if err == nil {
		t.Fatal("expected error, got nil")
	}
}

// --- Logout ---

func TestLogout_Success(t *testing.T) {
	repo := &stubRepo{
		logoutFn: func(id int) (*model.User, error) {
			return &model.User{ID: id, Username: "alice"}, nil
		},
	}
	svc := services.NewUserService(repo)

	user, err := svc.Logout(1)

	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if user.ID != 1 {
		t.Errorf("expected user ID 1, got %d", user.ID)
	}
}

func TestLogout_UserNotFound(t *testing.T) {
	repo := &stubRepo{
		logoutFn: func(int) (*model.User, error) {
			return nil, errors.New("record not found")
		},
	}
	svc := services.NewUserService(repo)

	_, err := svc.Logout(99)

	if err == nil {
		t.Fatal("expected error, got nil")
	}
}
