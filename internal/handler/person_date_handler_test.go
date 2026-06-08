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

type stubPersonDatesService struct {
	createFn      func(personID uuid.UUID, input services.CreatePersonDatesInput) (*model.PersonDate, error)
	getByIDFn     func(id uuid.UUID, personID uuid.UUID) (*model.PersonDate, error)
	getAllFn      func(personID uuid.UUID, page int, pageSize int) ([]*model.PersonDate, *pagination.PaginationMeta, error)
	updateFn      func(id uuid.UUID, personID uuid.UUID, input services.UpdatePersonDatesInput) (*model.PersonDate, error)
	deleteFn      func(personID uuid.UUID, id uuid.UUID) error
}

func (s *stubPersonDatesService) CreatePersonDates(personID uuid.UUID, input services.CreatePersonDatesInput) (*model.PersonDate, error) {
	if s.createFn != nil {
		return s.createFn(personID, input)
	}
	return &model.PersonDate{
		ID:       uuid.New(),
		PersonID: personID.String(),
		Date:     input.Date,
		Label:    input.Label,
	}, nil
}

func (s *stubPersonDatesService) GetPersonDatesByID(id uuid.UUID, personID uuid.UUID) (*model.PersonDate, error) {
	if s.getByIDFn != nil {
		return s.getByIDFn(id, personID)
	}
	return &model.PersonDate{ID: id}, nil
}

func (s *stubPersonDatesService) GetAllPersonDates(personID uuid.UUID, page int, pageSize int) ([]*model.PersonDate, *pagination.PaginationMeta, error) {
	if s.getAllFn != nil {
		return s.getAllFn(personID, page, pageSize)
	}
	return []*model.PersonDate{}, &pagination.PaginationMeta{}, nil
}

func (s *stubPersonDatesService) UpdatePersonDates(id uuid.UUID, personID uuid.UUID, input services.UpdatePersonDatesInput) (*model.PersonDate, error) {
	if s.updateFn != nil {
		return s.updateFn(id, personID, input)
	}
	return &model.PersonDate{ID: id}, nil
}

func (s *stubPersonDatesService) DeletePersonDates(personID uuid.UUID, id uuid.UUID) error {
	if s.deleteFn != nil {
		return s.deleteFn(personID, id)
	}
	return nil
}

func newPersonDateRouter(h *handler.PersonDateHandler, userID int) *gin.Engine {
	r := gin.New()
	r.POST("/persons/:id/dates", injectUserID(userID), h.CreatePersonDate)
	r.GET("/persons/:id/dates", injectUserID(userID), h.GetAllPersonDates)
	r.GET("/persons/:id/dates/:dateId", injectUserID(userID), h.GetPersonDatesByID)
	r.PUT("/persons/:id/dates/:dateId", injectUserID(userID), h.UpdatePersonDates)
	r.DELETE("/persons/:id/dates/:dateId", injectUserID(userID), h.DeletePersonDates)
	return r
}

// --- POST /persons/:id/dates ---

func TestCreatePersonDateHandler_Success(t *testing.T) {
	svc := &stubPersonDatesService{}
	h := handler.NewPersonDateHandler(svc)
	r := newPersonDateRouter(h, 1)

	personID := uuid.New()
	body, _ := json.Marshal(map[string]interface{}{
		"date":  "2026-01-01",
		"label": "Birthday",
	})

	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodPost, "/persons/"+personID.String()+"/dates", bytes.NewReader(body))
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
	if data["label"] != "Birthday" {
		t.Errorf("expected label 'Birthday', got: %v", data["label"])
	}
	if data["date"] != "2026-01-01" {
		t.Errorf("expected date '2026-01-01', got: %v", data["date"])
	}
}

func TestCreatePersonDateHandler_MissingLabel(t *testing.T) {
	svc := &stubPersonDatesService{}
	h := handler.NewPersonDateHandler(svc)
	r := newPersonDateRouter(h, 1)

	personID := uuid.New()
	body, _ := json.Marshal(map[string]interface{}{
		"date": "2026-01-01",
	})

	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodPost, "/persons/"+personID.String()+"/dates", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	r.ServeHTTP(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("expected 400, got %d — body: %s", w.Code, w.Body.String())
	}
}

func TestCreatePersonDateHandler_MissingDate(t *testing.T) {
	svc := &stubPersonDatesService{}
	h := handler.NewPersonDateHandler(svc)
	r := newPersonDateRouter(h, 1)

	personID := uuid.New()
	body, _ := json.Marshal(map[string]interface{}{
		"label": "Birthday",
	})

	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodPost, "/persons/"+personID.String()+"/dates", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	r.ServeHTTP(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("expected 400, got %d — body: %s", w.Code, w.Body.String())
	}
}

func TestCreatePersonDateHandler_InvalidPersonID(t *testing.T) {
	svc := &stubPersonDatesService{}
	h := handler.NewPersonDateHandler(svc)
	r := newPersonDateRouter(h, 1)

	body, _ := json.Marshal(map[string]interface{}{
		"date":  "2026-01-01",
		"label": "Birthday",
	})

	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodPost, "/persons/not-a-uuid/dates", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	r.ServeHTTP(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("expected 400, got %d — body: %s", w.Code, w.Body.String())
	}
}

func TestCreatePersonDateHandler_Unauthenticated(t *testing.T) {
	svc := &stubPersonDatesService{}
	h := handler.NewPersonDateHandler(svc)
	r := newPersonDateRouter(h, 0)

	personID := uuid.New()
	body, _ := json.Marshal(map[string]interface{}{
		"date":  "2026-01-01",
		"label": "Birthday",
	})

	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodPost, "/persons/"+personID.String()+"/dates", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	r.ServeHTTP(w, req)

	if w.Code != http.StatusUnauthorized {
		t.Errorf("expected 401, got %d — body: %s", w.Code, w.Body.String())
	}
}

func TestCreatePersonDateHandler_RepoError(t *testing.T) {
	svc := &stubPersonDatesService{
		createFn: func(personID uuid.UUID, input services.CreatePersonDatesInput) (*model.PersonDate, error) {
			return nil, errors.New("db error")
		},
	}
	h := handler.NewPersonDateHandler(svc)
	r := newPersonDateRouter(h, 1)

	personID := uuid.New()
	body, _ := json.Marshal(map[string]interface{}{
		"date":  "2026-01-01",
		"label": "Birthday",
	})

	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodPost, "/persons/"+personID.String()+"/dates", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	r.ServeHTTP(w, req)

	if w.Code != http.StatusInternalServerError {
		t.Errorf("expected 500, got %d — body: %s", w.Code, w.Body.String())
	}
}

// --- GET /persons/:id/dates ---

func TestGetAllPersonDatesHandler_Success(t *testing.T) {
	svc := &stubPersonDatesService{}
	h := handler.NewPersonDateHandler(svc)
	r := newPersonDateRouter(h, 1)

	personID := uuid.New()
	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/persons/"+personID.String()+"/dates", nil)
	r.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("expected 200, got %d — body: %s", w.Code, w.Body.String())
	}
}

func TestGetAllPersonDatesHandler_Unauthenticated(t *testing.T) {
	svc := &stubPersonDatesService{}
	h := handler.NewPersonDateHandler(svc)
	r := newPersonDateRouter(h, 0)

	personID := uuid.New()
	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/persons/"+personID.String()+"/dates", nil)
	r.ServeHTTP(w, req)

	if w.Code != http.StatusUnauthorized {
		t.Errorf("expected 401, got %d — body: %s", w.Code, w.Body.String())
	}
}

func TestGetAllPersonDatesHandler_InvalidPersonID(t *testing.T) {
	svc := &stubPersonDatesService{}
	h := handler.NewPersonDateHandler(svc)
	r := newPersonDateRouter(h, 1)

	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/persons/not-a-uuid/dates", nil)
	r.ServeHTTP(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("expected 400, got %d — body: %s", w.Code, w.Body.String())
	}
}

func TestGetAllPersonDatesHandler_RepoError(t *testing.T) {
	svc := &stubPersonDatesService{
		getAllFn: func(personID uuid.UUID, page int, pageSize int) ([]*model.PersonDate, *pagination.PaginationMeta, error) {
			return nil, nil, errors.New("db error")
		},
	}
	h := handler.NewPersonDateHandler(svc)
	r := newPersonDateRouter(h, 1)

	personID := uuid.New()
	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/persons/"+personID.String()+"/dates", nil)
	r.ServeHTTP(w, req)

	if w.Code != http.StatusInternalServerError {
		t.Errorf("expected 500, got %d — body: %s", w.Code, w.Body.String())
	}
}

// --- GET /persons/:id/dates/:dateId ---

func TestGetPersonDatesByIDHandler_Success(t *testing.T) {
	svc := &stubPersonDatesService{}
	h := handler.NewPersonDateHandler(svc)
	r := newPersonDateRouter(h, 1)

	personID := uuid.New()
	dateID := uuid.New()
	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/persons/"+personID.String()+"/dates/"+dateID.String(), nil)
	r.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("expected 200, got %d — body: %s", w.Code, w.Body.String())
	}
}

func TestGetPersonDatesByIDHandler_InvalidPersonID(t *testing.T) {
	svc := &stubPersonDatesService{}
	h := handler.NewPersonDateHandler(svc)
	r := newPersonDateRouter(h, 1)

	dateID := uuid.New()
	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/persons/not-a-uuid/dates/"+dateID.String(), nil)
	r.ServeHTTP(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("expected 400, got %d — body: %s", w.Code, w.Body.String())
	}
}

func TestGetPersonDatesByIDHandler_InvalidDateID(t *testing.T) {
	svc := &stubPersonDatesService{}
	h := handler.NewPersonDateHandler(svc)
	r := newPersonDateRouter(h, 1)

	personID := uuid.New()
	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/persons/"+personID.String()+"/dates/not-a-uuid", nil)
	r.ServeHTTP(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("expected 400, got %d — body: %s", w.Code, w.Body.String())
	}
}

func TestGetPersonDatesByIDHandler_NotFound(t *testing.T) {
	svc := &stubPersonDatesService{
		getByIDFn: func(id uuid.UUID, personID uuid.UUID) (*model.PersonDate, error) {
			return nil, gorm.ErrRecordNotFound
		},
	}
	h := handler.NewPersonDateHandler(svc)
	r := newPersonDateRouter(h, 1)

	personID := uuid.New()
	dateID := uuid.New()
	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/persons/"+personID.String()+"/dates/"+dateID.String(), nil)
	r.ServeHTTP(w, req)

	if w.Code != http.StatusNotFound {
		t.Errorf("expected 404, got %d — body: %s", w.Code, w.Body.String())
	}
}

func TestGetPersonDatesByIDHandler_RepoError(t *testing.T) {
	svc := &stubPersonDatesService{
		getByIDFn: func(id uuid.UUID, personID uuid.UUID) (*model.PersonDate, error) {
			return nil, errors.New("db error")
		},
	}
	h := handler.NewPersonDateHandler(svc)
	r := newPersonDateRouter(h, 1)

	personID := uuid.New()
	dateID := uuid.New()
	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/persons/"+personID.String()+"/dates/"+dateID.String(), nil)
	r.ServeHTTP(w, req)

	if w.Code != http.StatusInternalServerError {
		t.Errorf("expected 500, got %d — body: %s", w.Code, w.Body.String())
	}
}

// --- PUT /persons/:id/dates/:dateId ---

func TestUpdatePersonDatesHandler_Success(t *testing.T) {
	svc := &stubPersonDatesService{}
	h := handler.NewPersonDateHandler(svc)
	r := newPersonDateRouter(h, 1)

	personID := uuid.New()
	dateID := uuid.New()
	body, _ := json.Marshal(map[string]interface{}{
		"label": "Anniversary",
	})

	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodPut, "/persons/"+personID.String()+"/dates/"+dateID.String(), bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	r.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("expected 200, got %d — body: %s", w.Code, w.Body.String())
	}
}

func TestUpdatePersonDatesHandler_InvalidPersonID(t *testing.T) {
	svc := &stubPersonDatesService{}
	h := handler.NewPersonDateHandler(svc)
	r := newPersonDateRouter(h, 1)

	dateID := uuid.New()
	body, _ := json.Marshal(map[string]interface{}{"label": "Anniversary"})

	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodPut, "/persons/not-a-uuid/dates/"+dateID.String(), bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	r.ServeHTTP(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("expected 400, got %d — body: %s", w.Code, w.Body.String())
	}
}

func TestUpdatePersonDatesHandler_InvalidDateID(t *testing.T) {
	svc := &stubPersonDatesService{}
	h := handler.NewPersonDateHandler(svc)
	r := newPersonDateRouter(h, 1)

	personID := uuid.New()
	body, _ := json.Marshal(map[string]interface{}{"label": "Anniversary"})

	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodPut, "/persons/"+personID.String()+"/dates/not-a-uuid", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	r.ServeHTTP(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("expected 400, got %d — body: %s", w.Code, w.Body.String())
	}
}

func TestUpdatePersonDatesHandler_Unauthenticated(t *testing.T) {
	svc := &stubPersonDatesService{}
	h := handler.NewPersonDateHandler(svc)
	r := newPersonDateRouter(h, 0)

	personID := uuid.New()
	dateID := uuid.New()
	body, _ := json.Marshal(map[string]interface{}{"label": "Anniversary"})

	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodPut, "/persons/"+personID.String()+"/dates/"+dateID.String(), bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	r.ServeHTTP(w, req)

	if w.Code != http.StatusUnauthorized {
		t.Errorf("expected 401, got %d — body: %s", w.Code, w.Body.String())
	}
}

func TestUpdatePersonDatesHandler_NotFound(t *testing.T) {
	svc := &stubPersonDatesService{
		updateFn: func(id uuid.UUID, personID uuid.UUID, input services.UpdatePersonDatesInput) (*model.PersonDate, error) {
			return nil, gorm.ErrRecordNotFound
		},
	}
	h := handler.NewPersonDateHandler(svc)
	r := newPersonDateRouter(h, 1)

	personID := uuid.New()
	dateID := uuid.New()
	body, _ := json.Marshal(map[string]interface{}{"label": "Anniversary"})

	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodPut, "/persons/"+personID.String()+"/dates/"+dateID.String(), bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	r.ServeHTTP(w, req)

	if w.Code != http.StatusNotFound {
		t.Errorf("expected 404, got %d — body: %s", w.Code, w.Body.String())
	}
}

func TestUpdatePersonDatesHandler_RepoError(t *testing.T) {
	svc := &stubPersonDatesService{
		updateFn: func(id uuid.UUID, personID uuid.UUID, input services.UpdatePersonDatesInput) (*model.PersonDate, error) {
			return nil, errors.New("db error")
		},
	}
	h := handler.NewPersonDateHandler(svc)
	r := newPersonDateRouter(h, 1)

	personID := uuid.New()
	dateID := uuid.New()
	body, _ := json.Marshal(map[string]interface{}{"label": "Anniversary"})

	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodPut, "/persons/"+personID.String()+"/dates/"+dateID.String(), bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	r.ServeHTTP(w, req)

	if w.Code != http.StatusInternalServerError {
		t.Errorf("expected 500, got %d — body: %s", w.Code, w.Body.String())
	}
}

func TestUpdatePersonDatesHandler_InvalidJSON(t *testing.T) {
	svc := &stubPersonDatesService{}
	h := handler.NewPersonDateHandler(svc)
	r := newPersonDateRouter(h, 1)

	personID := uuid.New()
	dateID := uuid.New()

	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodPut, "/persons/"+personID.String()+"/dates/"+dateID.String(), bytes.NewReader([]byte("{invalid")))
	req.Header.Set("Content-Type", "application/json")
	r.ServeHTTP(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("expected 400, got %d — body: %s", w.Code, w.Body.String())
	}
}

// --- DELETE /persons/:id/dates/:dateId ---

func TestDeletePersonDatesHandler_Success(t *testing.T) {
	svc := &stubPersonDatesService{}
	h := handler.NewPersonDateHandler(svc)
	r := newPersonDateRouter(h, 1)

	personID := uuid.New()
	dateID := uuid.New()
	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodDelete, "/persons/"+personID.String()+"/dates/"+dateID.String(), nil)
	r.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("expected 200, got %d — body: %s", w.Code, w.Body.String())
	}
}

func TestDeletePersonDatesHandler_InvalidDateID(t *testing.T) {
	svc := &stubPersonDatesService{}
	h := handler.NewPersonDateHandler(svc)
	r := newPersonDateRouter(h, 1)

	personID := uuid.New()
	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodDelete, "/persons/"+personID.String()+"/dates/not-a-uuid", nil)
	r.ServeHTTP(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("expected 400, got %d — body: %s", w.Code, w.Body.String())
	}
}

func TestDeletePersonDatesHandler_InvalidPersonID(t *testing.T) {
	svc := &stubPersonDatesService{}
	h := handler.NewPersonDateHandler(svc)
	r := newPersonDateRouter(h, 1)

	dateID := uuid.New()
	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodDelete, "/persons/not-a-uuid/dates/"+dateID.String(), nil)
	r.ServeHTTP(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("expected 400, got %d — body: %s", w.Code, w.Body.String())
	}
}

func TestDeletePersonDatesHandler_NotFound(t *testing.T) {
	svc := &stubPersonDatesService{
		deleteFn: func(personID uuid.UUID, id uuid.UUID) error {
			return gorm.ErrRecordNotFound
		},
	}
	h := handler.NewPersonDateHandler(svc)
	r := newPersonDateRouter(h, 1)

	personID := uuid.New()
	dateID := uuid.New()
	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodDelete, "/persons/"+personID.String()+"/dates/"+dateID.String(), nil)
	r.ServeHTTP(w, req)

	if w.Code != http.StatusNotFound {
		t.Errorf("expected 404, got %d — body: %s", w.Code, w.Body.String())
	}
}

func TestDeletePersonDatesHandler_RepoError(t *testing.T) {
	svc := &stubPersonDatesService{
		deleteFn: func(personID uuid.UUID, id uuid.UUID) error {
			return errors.New("db error")
		},
	}
	h := handler.NewPersonDateHandler(svc)
	r := newPersonDateRouter(h, 1)

	personID := uuid.New()
	dateID := uuid.New()
	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodDelete, "/persons/"+personID.String()+"/dates/"+dateID.String(), nil)
	r.ServeHTTP(w, req)

	if w.Code != http.StatusInternalServerError {
		t.Errorf("expected 500, got %d — body: %s", w.Code, w.Body.String())
	}
}
