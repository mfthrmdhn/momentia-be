package services

import (
	"math"
	"momentia-be/model"
	"momentia-be/pkg/pagination"
	"momentia-be/repository"

	"github.com/google/uuid"
)

type CreatePersonInput struct {
	Name         string `json:"name" binding:"required"`
	Relationship string `json:"relationship" binding:"required"`
	IsPinned     bool   `json:"is_pinned"`
}

type UpdatePersonInput struct {
	Name         *string `json:"name"`
	Relationship *string `json:"relationship"`
	IsPinned     *bool   `json:"is_pinned"`
}

type PersonService interface {
	CreatePerson(userID int, input CreatePersonInput) (*model.Person, error)
	GetPersonByID(id uuid.UUID, userID int) (*model.Person, error)
	GetAllPersons(userID int, page int, pageSize int) ([]*model.Person, *pagination.PaginationMeta, error)
	UpdatePerson(id uuid.UUID, userID int, input UpdatePersonInput) (*model.Person, error)
	DeletePerson(id uuid.UUID) error
}

type personService struct {
	repo repository.PersonRepository
}

func NewPersonService(repo repository.PersonRepository) PersonService {
	return &personService{repo}
}

func (s *personService) CreatePerson(userID int, input CreatePersonInput) (*model.Person, error) {
	person := &model.Person{
		Name:         input.Name,
		Relationship: input.Relationship,
		IsPinned:     input.IsPinned,
	}

	person, err := s.repo.CreatePerson(userID, person)
	if err != nil {
		return nil, err
	}
	return person, nil
}

func (s *personService) GetPersonByID(id uuid.UUID, userID int) (*model.Person, error) {
	return s.repo.GetPersonByID(id, userID)
}

func (s *personService) GetAllPersons(userID int, page int, pageSize int) ([]*model.Person, *pagination.PaginationMeta, error) {
	if page < 1 {
		page = 1
	}
	if pageSize < 1 || pageSize > 100 {
		pageSize = 20
	}
	person, total, err := s.repo.GetAllPersons(userID, page, pageSize)
	if err != nil {
		return nil, nil, err
	}
	meta := &pagination.PaginationMeta{
		Page:       page,
		PerPage:    pageSize,
		TotalItems: int64(total),
		TotalPages: int(math.Ceil(float64(total) / float64(pageSize))),
	}
	return person, meta, nil
}

func (s *personService) UpdatePerson(id uuid.UUID, userID int, input UpdatePersonInput) (*model.Person, error) {
	updates := map[string]interface{}{}
	if input.Name != nil {
		updates["name"] = *input.Name
	}
	if input.Relationship != nil {
		updates["relationship"] = *input.Relationship
	}
	if input.IsPinned != nil {
		updates["is_pinned"] = *input.IsPinned
	}
	if err := s.repo.UpdatePerson(id, userID, updates); err != nil {
		return nil, err
	}
	return s.repo.GetPersonByID(id, userID)
}

func (s *personService) DeletePerson(id uuid.UUID) error {
	return s.repo.DeletePerson(id)
}