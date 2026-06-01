package handler_test

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"momentia-be/internal/handler"
	"momentia-be/model"
	"momentia-be/services"

	"github.com/gin-gonic/gin"
)

func init() {
	gin.SetMode(gin.TestMode)
}

// stubService implements services.UserService.
type stubService struct {
	registerFn          func(services.RegisterInput) (*model.User, error)
	loginFn             func(services.LoginInput) (string, error)
	logoutFn            func(int) (*model.User, error)
	getUserByIDFn       func(int) (*model.User, error)
	getUserByEmailFn    func(string) (*model.User, error)
	getUserByUsernameFn func(string) (*model.User, error)
	getUserByMsisdnFn   func(string) (*model.User, error)
	updateUserFn        func(*model.User) error
}

func (s *stubService) Register(input services.RegisterInput) (*model.User, error) {
	if s.registerFn != nil {
		return s.registerFn(input)
	}
	return &model.User{ID: 1, Username: input.Username, Email: input.Email}, nil
}
func (s *stubService) Login(input services.LoginInput) (string, error) {
	if s.loginFn != nil {
		return s.loginFn(input)
	}
	return "mock-token", nil
}
func (s *stubService) Logout(id int) (*model.User, error) {
	if s.logoutFn != nil {
		return s.logoutFn(id)
	}
	return &model.User{ID: id}, nil
}
func (s *stubService) GetUserByID(id int) (*model.User, error) {
	if s.getUserByIDFn != nil {
		return s.getUserByIDFn(id)
	}
	return &model.User{ID: id}, nil
}
func (s *stubService) GetUserByEmail(email string) (*model.User, error) {
	if s.getUserByEmailFn != nil {
		return s.getUserByEmailFn(email)
	}
	return &model.User{}, nil
}
func (s *stubService) GetUserByUsername(username string) (*model.User, error) {
	if s.getUserByUsernameFn != nil {
		return s.getUserByUsernameFn(username)
	}
	return &model.User{}, nil
}
func (s *stubService) GetUserByMsisdn(msisdn string) (*model.User, error) {
	if s.getUserByMsisdnFn != nil {
		return s.getUserByMsisdnFn(msisdn)
	}
	return &model.User{}, nil
}
func (s *stubService) UpdateUser(u *model.User) error {
	if s.updateUserFn != nil {
		return s.updateUserFn(u)
	}
	return nil
}

// stubRepo implements repository.UserRepository for the handler's direct repo calls.
type stubRepo struct {
	getUserByIDFn func(int) (*model.User, error)
}

func (r *stubRepo) Register(u *model.User) error                          { return nil }
func (r *stubRepo) Login(email string) (*model.User, error)               { return &model.User{}, nil }
func (r *stubRepo) Logout(id int) (*model.User, error)                    { return &model.User{ID: id}, nil }
func (r *stubRepo) GetUserByEmail(email string) (*model.User, error)      { return &model.User{}, nil }
func (r *stubRepo) GetUserByUsername(u string) (*model.User, error)       { return &model.User{}, nil }
func (r *stubRepo) GetUserByMsisdn(m string) (*model.User, error)         { return &model.User{}, nil }
func (r *stubRepo) UpdateUser(u *model.User) error                        { return nil }
func (r *stubRepo) GetUserByID(id int) (*model.User, error) {
	if r.getUserByIDFn != nil {
		return r.getUserByIDFn(id)
	}
	return &model.User{ID: id, Username: "alice"}, nil
}

func newRouter(h *handler.UserHandler) *gin.Engine {
	r := gin.New()
	r.POST("/register", h.Register)
	r.POST("/login", h.Login)
	r.GET("/profile", h.GetUserByID)
	r.POST("/logout", h.Logout)
	return r
}

// injectUserID simulates what AuthMiddleware does.
func injectUserID(id int) gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Set("userID", id)
		c.Next()
	}
}

func newRouterWithAuth(h *handler.UserHandler, userID int) *gin.Engine {
	r := gin.New()
	r.POST("/register", h.Register)
	r.POST("/login", h.Login)
	r.GET("/profile", injectUserID(userID), h.GetUserByID)
	r.POST("/logout", injectUserID(userID), h.Logout)
	return r
}

// --- POST /register ---

func TestRegisterHandler_Success(t *testing.T) {
	h := handler.NewUserHandler(&stubRepo{}, &stubService{})
	r := newRouter(h)

	body, _ := json.Marshal(map[string]string{
		"username": "alice",
		"email":    "alice@example.com",
		"msisdn":   "+628123456789",
		"password": "secret123",
	})

	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodPost, "/register", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	r.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("expected 200, got %d — body: %s", w.Code, w.Body.String())
	}
}

func TestRegisterHandler_InvalidJSON(t *testing.T) {
	h := handler.NewUserHandler(&stubRepo{}, &stubService{})
	r := newRouter(h)

	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodPost, "/register", bytes.NewReader([]byte(`{bad json}`)))
	req.Header.Set("Content-Type", "application/json")
	r.ServeHTTP(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("expected 400, got %d", w.Code)
	}
}

func TestRegisterHandler_MissingRequiredFields(t *testing.T) {
	h := handler.NewUserHandler(&stubRepo{}, &stubService{})
	r := newRouter(h)

	body, _ := json.Marshal(map[string]string{"username": "alice"}) // missing email, msisdn, password

	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodPost, "/register", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	r.ServeHTTP(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("expected 400, got %d", w.Code)
	}
}

func TestRegisterHandler_ServiceError(t *testing.T) {
	svc := &stubService{
		registerFn: func(services.RegisterInput) (*model.User, error) {
			return nil, errors.New("username already taken")
		},
	}
	h := handler.NewUserHandler(&stubRepo{}, svc)
	r := newRouter(h)

	body, _ := json.Marshal(map[string]string{
		"username": "alice",
		"email":    "alice@example.com",
		"msisdn":   "+628123456789",
		"password": "secret123",
	})

	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodPost, "/register", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	r.ServeHTTP(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("expected 400, got %d", w.Code)
	}
}

// --- POST /login ---

func TestLoginHandler_Success(t *testing.T) {
	svc := &stubService{
		loginFn: func(services.LoginInput) (string, error) {
			return "jwt-token-string", nil
		},
	}
	h := handler.NewUserHandler(&stubRepo{}, svc)
	r := newRouter(h)

	body, _ := json.Marshal(map[string]string{
		"email":    "alice@example.com",
		"password": "secret123",
	})

	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodPost, "/login", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	r.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("expected 200, got %d — body: %s", w.Code, w.Body.String())
	}

	var resp map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &resp)
	data, ok := resp["data"].(map[string]interface{})
	if !ok || data["token"] == "" {
		t.Errorf("expected token in response, got: %s", w.Body.String())
	}
}

func TestLoginHandler_InvalidCredentials(t *testing.T) {
	svc := &stubService{
		loginFn: func(services.LoginInput) (string, error) {
			return "", errors.New("invalid email or password")
		},
	}
	h := handler.NewUserHandler(&stubRepo{}, svc)
	r := newRouter(h)

	body, _ := json.Marshal(map[string]string{
		"email":    "alice@example.com",
		"password": "wrongpassword",
	})

	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodPost, "/login", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	r.ServeHTTP(w, req)

	if w.Code != http.StatusUnauthorized {
		t.Errorf("expected 401, got %d", w.Code)
	}
}

func TestLoginHandler_InvalidJSON(t *testing.T) {
	h := handler.NewUserHandler(&stubRepo{}, &stubService{})
	r := newRouter(h)

	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodPost, "/login", bytes.NewReader([]byte(`{bad}`)))
	req.Header.Set("Content-Type", "application/json")
	r.ServeHTTP(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("expected 400, got %d", w.Code)
	}
}

// --- GET /profile ---

func TestGetProfile_Success(t *testing.T) {
	repo := &stubRepo{
		getUserByIDFn: func(id int) (*model.User, error) {
			return &model.User{ID: id, Username: "alice", Email: "alice@example.com"}, nil
		},
	}
	h := handler.NewUserHandler(repo, &stubService{})
	r := newRouterWithAuth(h, 1)

	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/profile", nil)
	r.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("expected 200, got %d — body: %s", w.Code, w.Body.String())
	}
}

func TestGetProfile_Unauthenticated(t *testing.T) {
	h := handler.NewUserHandler(&stubRepo{}, &stubService{})
	r := newRouter(h) // no injectUserID middleware

	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/profile", nil)
	r.ServeHTTP(w, req)

	if w.Code != http.StatusUnauthorized {
		t.Errorf("expected 401, got %d", w.Code)
	}
}

func TestGetProfile_UserNotFound(t *testing.T) {
	repo := &stubRepo{
		getUserByIDFn: func(int) (*model.User, error) {
			return nil, errors.New("record not found")
		},
	}
	h := handler.NewUserHandler(repo, &stubService{})
	r := newRouterWithAuth(h, 99)

	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/profile", nil)
	r.ServeHTTP(w, req)

	if w.Code != http.StatusNotFound {
		t.Errorf("expected 404, got %d", w.Code)
	}
}

// --- POST /logout ---

func TestLogoutHandler_Success(t *testing.T) {
	svc := &stubService{
		logoutFn: func(id int) (*model.User, error) {
			return &model.User{ID: id, Username: "alice"}, nil
		},
	}
	h := handler.NewUserHandler(&stubRepo{}, svc)
	r := newRouterWithAuth(h, 1)

	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodPost, "/logout", nil)
	r.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("expected 200, got %d — body: %s", w.Code, w.Body.String())
	}
}

func TestLogoutHandler_Unauthenticated(t *testing.T) {
	h := handler.NewUserHandler(&stubRepo{}, &stubService{})
	r := newRouter(h) // no injectUserID middleware

	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodPost, "/logout", nil)
	r.ServeHTTP(w, req)

	if w.Code != http.StatusUnauthorized {
		t.Errorf("expected 401, got %d", w.Code)
	}
}

func TestLogoutHandler_UserNotFound(t *testing.T) {
	svc := &stubService{
		logoutFn: func(int) (*model.User, error) {
			return nil, errors.New("record not found")
		},
	}
	h := handler.NewUserHandler(&stubRepo{}, svc)
	r := newRouterWithAuth(h, 99)

	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodPost, "/logout", nil)
	r.ServeHTTP(w, req)

	if w.Code != http.StatusNotFound {
		t.Errorf("expected 404, got %d", w.Code)
	}
}
