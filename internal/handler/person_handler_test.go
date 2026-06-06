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
	"momentia-be/pkg/pagination"
	"momentia-be/services"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type stubPersonService struct {
	createPersonFn  func(userID int, input services.CreatePersonInput) (*model.Person, error)
	getPersonByIDFn func(id uuid.UUID, userID int) (*model.Person, error)
	getAllPersonsFn  func(userID int, page int, pageSize int) ([]*model.Person, *pagination.PaginationMeta, error)
	updatePersonFn  func(id uuid.UUID, userID int, input services.UpdatePersonInput) (*model.Person, error)
	deletePersonFn  func(id uuid.UUID) error
}

func (s *stubPersonService) CreatePerson(userID int, input services.CreatePersonInput) (*model.Person, error) {
	if s.createPersonFn != nil {
		return s.createPersonFn(userID, input)
	}
	return &model.Person{
		ID:            uuid.New(),
		Name:          input.Name,
		Relationship:  input.Relationship,
		IsPinned:      input.IsPinned,
		CreatorUserID: userID,
	}, nil
}

func (s *stubPersonService) GetPersonByID(id uuid.UUID, userID int) (*model.Person, error) {
	if s.getPersonByIDFn != nil {
		return s.getPersonByIDFn(id, userID)
	}
	return &model.Person{ID: id}, nil
}

func (s *stubPersonService) GetAllPersons(userID int, page int, pageSize int) ([]*model.Person, *pagination.PaginationMeta, error) {
	if s.getAllPersonsFn != nil {
		return s.getAllPersonsFn(userID, page, pageSize)
	}
	return []*model.Person{}, &pagination.PaginationMeta{}, nil
}

func (s *stubPersonService) UpdatePerson(id uuid.UUID, userID int, input services.UpdatePersonInput) (*model.Person, error) {
	if s.updatePersonFn != nil {
		return s.updatePersonFn(id, userID, input)
	}
	return &model.Person{ID: id}, nil
}

func (s *stubPersonService) DeletePerson(id uuid.UUID) error {
	if s.deletePersonFn != nil {
		return s.deletePersonFn(id)
	}
	return nil
}

func newPersonRouter(h *handler.PersonHandler, userID int) *gin.Engine {
	r := gin.New()
	r.GET("/persons", injectUserID(userID), h.GetPersons)
	r.GET("/persons/:id", injectUserID(userID), h.GetPersonByID)
	r.POST("/persons", injectUserID(userID), h.CreatePerson)
	r.PUT("/persons/:id", injectUserID(userID), h.UpdatePerson)
	r.DELETE("/persons/:id", injectUserID(userID), h.DeletePerson)
	return r
}

// --- POST /persons ---

func TestCreatePersonHandler_Success(t *testing.T) {
	svc := &stubPersonService{}
	h := handler.NewPersonHandler(svc)
	r := newPersonRouter(h, 1)

	body, _ := json.Marshal(map[string]interface{}{
		"name":         "John Doe",
		"relationship": "friend",
		"is_pinned":    false,
	})

	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodPost, "/persons", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	r.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("expected 200, got %d — body: %s", w.Code, w.Body.String())
	}

	var resp map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &resp)
	data, ok := resp["data"].(map[string]interface{})
	if !ok {
		t.Fatalf("expected data object in response, got: %s", w.Body.String())
	}
	if data["name"] != "John Doe" {
		t.Errorf("expected name 'John Doe', got: %v", data["name"])
	}
	if data["relationship"] != "friend" {
		t.Errorf("expected relationship 'friend', got: %v", data["relationship"])
	}
}

func TestCreatePersonHandler_MissingName(t *testing.T) {
	svc := &stubPersonService{}
	h := handler.NewPersonHandler(svc)
	r := newPersonRouter(h, 1)

	body, _ := json.Marshal(map[string]interface{}{
		"relationship": "friend",
	})

	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodPost, "/persons", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	r.ServeHTTP(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("expected 400, got %d — body: %s", w.Code, w.Body.String())
	}
}

func TestCreatePersonHandler_MissingRelationship(t *testing.T) {
	svc := &stubPersonService{}
	h := handler.NewPersonHandler(svc)
	r := newPersonRouter(h, 1)

	body, _ := json.Marshal(map[string]interface{}{
		"name": "John Doe",
	})

	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodPost, "/persons", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	r.ServeHTTP(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("expected 400, got %d — body: %s", w.Code, w.Body.String())
	}
}

func TestCreatePersonHandler_InvalidJSON(t *testing.T) {
	svc := &stubPersonService{}
	h := handler.NewPersonHandler(svc)
	r := newPersonRouter(h, 1)

	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodPost, "/persons", bytes.NewReader([]byte(`{bad json}`)))
	req.Header.Set("Content-Type", "application/json")
	r.ServeHTTP(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("expected 400, got %d", w.Code)
	}
}

func TestCreatePersonHandler_RepoError(t *testing.T) {
	svc := &stubPersonService{
		createPersonFn: func(userID int, input services.CreatePersonInput) (*model.Person, error) {
			return nil, errors.New("db error")
		},
	}
	h := handler.NewPersonHandler(svc)
	r := newPersonRouter(h, 1)

	body, _ := json.Marshal(map[string]interface{}{
		"name":         "John Doe",
		"relationship": "friend",
	})

	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodPost, "/persons", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	r.ServeHTTP(w, req)

	if w.Code != http.StatusInternalServerError {
		t.Errorf("expected 500, got %d — body: %s", w.Code, w.Body.String())
	}
}

func TestCreatePersonHandler_IsPinnedTrue(t *testing.T) {
	svc := &stubPersonService{}
	h := handler.NewPersonHandler(svc)
	r := newPersonRouter(h, 1)

	body, _ := json.Marshal(map[string]interface{}{
		"name":         "Jane Doe",
		"relationship": "partner",
		"is_pinned":    true,
	})

	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodPost, "/persons", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	r.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("expected 200, got %d — body: %s", w.Code, w.Body.String())
	}

	var resp map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &resp)
	data, ok := resp["data"].(map[string]interface{})
	if !ok {
		t.Fatalf("expected data object in response")
	}
	if data["is_pinned"] != true {
		t.Errorf("expected is_pinned true, got: %v", data["is_pinned"])
	}
}

// --- PUT /persons/:id ---

func TestUpdatePersonHandler_Success(t *testing.T) {
	svc := &stubPersonService{}
	h := handler.NewPersonHandler(svc)
	r := newPersonRouter(h, 1)

	id := uuid.New()
	body, _ := json.Marshal(map[string]interface{}{
		"name":         "Updated Name",
		"relationship": "colleague",
	})

	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodPut, "/persons/"+id.String(), bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	r.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("expected 200, got %d — body: %s", w.Code, w.Body.String())
	}
}

func TestUpdatePersonHandler_InvalidID(t *testing.T) {
	svc := &stubPersonService{}
	h := handler.NewPersonHandler(svc)
	r := newPersonRouter(h, 1)

	body, _ := json.Marshal(map[string]interface{}{
		"name":         "Updated Name",
		"relationship": "colleague",
	})

	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodPut, "/persons/not-a-uuid", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	r.ServeHTTP(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("expected 400, got %d — body: %s", w.Code, w.Body.String())
	}
}

func TestUpdatePersonHandler_PartialUpdate(t *testing.T) {
	svc := &stubPersonService{}
	h := handler.NewPersonHandler(svc)
	r := newPersonRouter(h, 1)

	id := uuid.New()
	body, _ := json.Marshal(map[string]interface{}{
		"name": "Only Name",
	})

	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodPut, "/persons/"+id.String(), bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	r.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("expected 200, got %d — body: %s", w.Code, w.Body.String())
	}
}

func TestUpdatePersonHandler_NotFound(t *testing.T) {
	svc := &stubPersonService{
		updatePersonFn: func(id uuid.UUID, userID int, input services.UpdatePersonInput) (*model.Person, error) {
			return nil, gorm.ErrRecordNotFound
		},
	}
	h := handler.NewPersonHandler(svc)
	r := newPersonRouter(h, 1)

	id := uuid.New()
	body, _ := json.Marshal(map[string]interface{}{
		"name":         "Updated Name",
		"relationship": "colleague",
	})

	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodPut, "/persons/"+id.String(), bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	r.ServeHTTP(w, req)

	if w.Code != http.StatusNotFound {
		t.Errorf("expected 404, got %d — body: %s", w.Code, w.Body.String())
	}
}

func TestUpdatePersonHandler_Unauthenticated(t *testing.T) {
	svc := &stubPersonService{}
	h := handler.NewPersonHandler(svc)
	r := newPersonRouter(h, 0)

	id := uuid.New()
	body, _ := json.Marshal(map[string]interface{}{
		"name":         "Updated Name",
		"relationship": "colleague",
	})

	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodPut, "/persons/"+id.String(), bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	r.ServeHTTP(w, req)

	if w.Code != http.StatusUnauthorized {
		t.Errorf("expected 401, got %d — body: %s", w.Code, w.Body.String())
	}
}

// --- GET /persons/:id ---

func TestGetPersonByIDHandler_Success(t *testing.T) {
	id := uuid.New()
	svc := &stubPersonService{
		getPersonByIDFn: func(gotID uuid.UUID, userID int) (*model.Person, error) {
			return &model.Person{ID: gotID, Name: "Alice"}, nil
		},
	}
	h := handler.NewPersonHandler(svc)
	r := newPersonRouter(h, 1)

	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/persons/"+id.String(), nil)
	r.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("expected 200, got %d — body: %s", w.Code, w.Body.String())
	}
}

func TestGetPersonByIDHandler_InvalidID(t *testing.T) {
	svc := &stubPersonService{}
	h := handler.NewPersonHandler(svc)
	r := newPersonRouter(h, 1)

	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/persons/not-a-uuid", nil)
	r.ServeHTTP(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("expected 400, got %d — body: %s", w.Code, w.Body.String())
	}
}

func TestGetPersonByIDHandler_NotFound(t *testing.T) {
	svc := &stubPersonService{
		getPersonByIDFn: func(id uuid.UUID, userID int) (*model.Person, error) {
			return nil, errors.New("record not found")
		},
	}
	h := handler.NewPersonHandler(svc)
	r := newPersonRouter(h, 1)

	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/persons/"+uuid.New().String(), nil)
	r.ServeHTTP(w, req)

	if w.Code != http.StatusInternalServerError {
		t.Errorf("expected 500, got %d — body: %s", w.Code, w.Body.String())
	}
}

// --- GET /persons ---

func TestGetPersonsHandler_Success(t *testing.T) {
	svc := &stubPersonService{
		getAllPersonsFn: func(userID int, page int, pageSize int) ([]*model.Person, *pagination.PaginationMeta, error) {
			return []*model.Person{{ID: uuid.New(), Name: "Alice"}}, &pagination.PaginationMeta{Page: 1, PerPage: 20, TotalItems: 1, TotalPages: 1}, nil
		},
	}
	h := handler.NewPersonHandler(svc)
	r := newPersonRouter(h, 1)

	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/persons", nil)
	r.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("expected 200, got %d — body: %s", w.Code, w.Body.String())
	}
}

func TestGetPersonsHandler_Unauthenticated(t *testing.T) {
	svc := &stubPersonService{}
	h := handler.NewPersonHandler(svc)
	r := newPersonRouter(h, 0)

	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/persons", nil)
	r.ServeHTTP(w, req)

	if w.Code != http.StatusUnauthorized {
		t.Errorf("expected 401, got %d — body: %s", w.Code, w.Body.String())
	}
}

func TestGetPersonsHandler_RepoError(t *testing.T) {
	svc := &stubPersonService{
		getAllPersonsFn: func(userID int, page int, pageSize int) ([]*model.Person, *pagination.PaginationMeta, error) {
			return nil, nil, errors.New("db error")
		},
	}
	h := handler.NewPersonHandler(svc)
	r := newPersonRouter(h, 1)

	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/persons", nil)
	r.ServeHTTP(w, req)

	if w.Code != http.StatusInternalServerError {
		t.Errorf("expected 500, got %d — body: %s", w.Code, w.Body.String())
	}
}

// --- DELETE /persons/:id ---

func TestDeletePersonHandler_Success(t *testing.T) {
	svc := &stubPersonService{}
	h := handler.NewPersonHandler(svc)
	r := newPersonRouter(h, 1)

	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodDelete, "/persons/"+uuid.New().String(), nil)
	r.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("expected 200, got %d — body: %s", w.Code, w.Body.String())
	}
}

func TestDeletePersonHandler_InvalidID(t *testing.T) {
	svc := &stubPersonService{}
	h := handler.NewPersonHandler(svc)
	r := newPersonRouter(h, 1)

	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodDelete, "/persons/not-a-uuid", nil)
	r.ServeHTTP(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("expected 400, got %d — body: %s", w.Code, w.Body.String())
	}
}

func TestDeletePersonHandler_RepoError(t *testing.T) {
	svc := &stubPersonService{
		deletePersonFn: func(id uuid.UUID) error {
			return errors.New("db error")
		},
	}
	h := handler.NewPersonHandler(svc)
	r := newPersonRouter(h, 1)

	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodDelete, "/persons/"+uuid.New().String(), nil)
	r.ServeHTTP(w, req)

	if w.Code != http.StatusInternalServerError {
		t.Errorf("expected 500, got %d — body: %s", w.Code, w.Body.String())
	}
}
