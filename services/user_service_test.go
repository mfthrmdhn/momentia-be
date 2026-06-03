package services_test

import (
	"errors"
	"testing"
	"time"

	"momentia-be/model"
	"momentia-be/services"

	"golang.org/x/crypto/bcrypt"
)

// stubUserRepo satisfies repository.UserRepository.
type stubUserRepo struct {
	registerFn          func(*model.User) error
	loginFn             func(string) (*model.User, error)
	getUserByIDFn       func(int) (*model.User, error)
	getUserByEmailFn    func(string) (*model.User, error)
	getUserByUsernameFn func(string) (*model.User, error)
	getUserByMsisdnFn   func(string) (*model.User, error)
	updateUserFn        func(*model.User) error
}

func (s *stubUserRepo) Register(u *model.User) error {
	if s.registerFn != nil {
		return s.registerFn(u)
	}
	return nil
}
func (s *stubUserRepo) Login(email string) (*model.User, error) {
	if s.loginFn != nil {
		return s.loginFn(email)
	}
	return &model.User{}, nil
}
func (s *stubUserRepo) GetUserByID(id int) (*model.User, error) {
	if s.getUserByIDFn != nil {
		return s.getUserByIDFn(id)
	}
	return &model.User{ID: id}, nil
}
func (s *stubUserRepo) GetUserByEmail(email string) (*model.User, error) {
	if s.getUserByEmailFn != nil {
		return s.getUserByEmailFn(email)
	}
	return nil, nil
}
func (s *stubUserRepo) GetUserByUsername(username string) (*model.User, error) {
	if s.getUserByUsernameFn != nil {
		return s.getUserByUsernameFn(username)
	}
	return nil, nil
}
func (s *stubUserRepo) GetUserByMsisdn(msisdn string) (*model.User, error) {
	if s.getUserByMsisdnFn != nil {
		return s.getUserByMsisdnFn(msisdn)
	}
	return nil, nil
}
func (s *stubUserRepo) UpdateUser(u *model.User) error {
	if s.updateUserFn != nil {
		return s.updateUserFn(u)
	}
	return nil
}

// stubSessionRepo satisfies repository.UserSessionRepository.
type stubSessionRepo struct {
	createFn           func(*model.UserSession) error
	deleteByTokenFn    func(string) error
	findByTokenFn      func(string) (*model.UserSession, error)
	deleteExpiredFn    func() error
}

func (s *stubSessionRepo) Create(sess *model.UserSession) error {
	if s.createFn != nil {
		return s.createFn(sess)
	}
	return nil
}
func (s *stubSessionRepo) DeleteByTokenHash(hash string) error {
	if s.deleteByTokenFn != nil {
		return s.deleteByTokenFn(hash)
	}
	return nil
}
func (s *stubSessionRepo) FindByTokenHash(hash string) (*model.UserSession, error) {
	if s.findByTokenFn != nil {
		return s.findByTokenFn(hash)
	}
	return &model.UserSession{}, nil
}
func (s *stubSessionRepo) DeleteExpired() error {
	if s.deleteExpiredFn != nil {
		return s.deleteExpiredFn()
	}
	return nil
}

// --- Register ---

func TestRegister_Success(t *testing.T) {
	svc := services.NewUserService(&stubUserRepo{}, &stubSessionRepo{})

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
	repo := &stubUserRepo{
		getUserByUsernameFn: func(string) (*model.User, error) {
			return &model.User{ID: 1, Username: "alice"}, nil
		},
	}
	svc := services.NewUserService(repo, &stubSessionRepo{})

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
	repo := &stubUserRepo{
		getUserByEmailFn: func(string) (*model.User, error) {
			return &model.User{ID: 2, Email: "alice@example.com"}, nil
		},
	}
	svc := services.NewUserService(repo, &stubSessionRepo{})

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
	repo := &stubUserRepo{
		getUserByMsisdnFn: func(string) (*model.User, error) {
			return &model.User{ID: 3, Msisdn: "+628123456789"}, nil
		},
	}
	svc := services.NewUserService(repo, &stubSessionRepo{})

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
	repo := &stubUserRepo{
		registerFn: func(*model.User) error {
			return errors.New("db connection failed")
		},
	}
	svc := services.NewUserService(repo, &stubSessionRepo{})

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
	repo := &stubUserRepo{
		loginFn: func(email string) (*model.User, error) {
			return &model.User{ID: 1, Email: email, PasswordHash: string(hash)}, nil
		},
	}
	svc := services.NewUserService(repo, &stubSessionRepo{})

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
	repo := &stubUserRepo{
		loginFn: func(email string) (*model.User, error) {
			return &model.User{ID: 1, Email: email, PasswordHash: string(hash)}, nil
		},
	}
	svc := services.NewUserService(repo, &stubSessionRepo{})

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
	repo := &stubUserRepo{
		loginFn: func(string) (*model.User, error) {
			return &model.User{}, nil // ID == 0 means not found
		},
	}
	svc := services.NewUserService(repo, &stubSessionRepo{})

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
	repo := &stubUserRepo{
		loginFn: func(email string) (*model.User, error) {
			return &model.User{ID: 1, Email: email, PasswordHash: string(hash)}, nil
		},
	}
	svc := services.NewUserService(repo, &stubSessionRepo{})

	_, err := svc.Login(services.LoginInput{
		Email:    "alice@example.com",
		Password: "password123",
	})

	if err == nil {
		t.Fatal("expected error, got nil")
	}
}

func TestLogin_SessionCreateError(t *testing.T) {
	t.Setenv("JWT_SECRET", "test-secret-key")

	hash, _ := bcrypt.GenerateFromPassword([]byte("password123"), bcrypt.MinCost)
	repo := &stubUserRepo{
		loginFn: func(email string) (*model.User, error) {
			return &model.User{ID: 1, Email: email, PasswordHash: string(hash)}, nil
		},
	}
	sessionRepo := &stubSessionRepo{
		createFn: func(*model.UserSession) error {
			return errors.New("db error")
		},
	}
	svc := services.NewUserService(repo, sessionRepo)

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
	deleted := ""
	sessionRepo := &stubSessionRepo{
		deleteByTokenFn: func(hash string) error {
			deleted = hash
			return nil
		},
	}
	svc := services.NewUserService(&stubUserRepo{}, sessionRepo)

	err := svc.Logout("somehash")

	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if deleted != "somehash" {
		t.Errorf("expected token hash %q to be deleted, got %q", "somehash", deleted)
	}
}

func TestLogout_SessionRepoError(t *testing.T) {
	sessionRepo := &stubSessionRepo{
		deleteByTokenFn: func(string) error {
			return errors.New("record not found")
		},
	}
	svc := services.NewUserService(&stubUserRepo{}, sessionRepo)

	err := svc.Logout("badhash")

	if err == nil {
		t.Fatal("expected error, got nil")
	}
}

func TestLogout_SessionExpiry(t *testing.T) {
	t.Setenv("JWT_SECRET", "test-secret-key")

	hash, _ := bcrypt.GenerateFromPassword([]byte("password123"), bcrypt.MinCost)
	repo := &stubUserRepo{
		loginFn: func(email string) (*model.User, error) {
			return &model.User{ID: 1, Email: email, PasswordHash: string(hash)}, nil
		},
	}
	var storedSession *model.UserSession
	sessionRepo := &stubSessionRepo{
		createFn: func(s *model.UserSession) error {
			storedSession = s
			return nil
		},
	}
	svc := services.NewUserService(repo, sessionRepo)

	_, err := svc.Login(services.LoginInput{Email: "alice@example.com", Password: "password123"})
	if err != nil {
		t.Fatalf("login failed: %v", err)
	}

	if storedSession == nil {
		t.Fatal("expected session to be created")
	}
	if storedSession.ExpiresAt.Before(time.Now()) {
		t.Error("expected session expiry to be in the future")
	}
}
