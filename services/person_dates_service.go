package services

import (
	"math"
	"momentia-be/model"
	"momentia-be/pkg/pagination"
	"momentia-be/repository"

	"github.com/google/uuid"
)

type CreatePersonDatesInput struct {
	Date  string `json:"date" binding:"required"`
	Label string `json:"label" binding:"required"`
}

type UpdatePersonDatesInput struct {
	Date  *string `json:"date"`
	Label *string `json:"label"`
}

type PersonDatesService interface {
	CreatePersonDates(personID uuid.UUID, input CreatePersonDatesInput) (*model.PersonDate, error)
	GetPersonDatesByID(id uuid.UUID, personID uuid.UUID) (*model.PersonDate, error)
	GetAllPersonDates(personID uuid.UUID, page int, pageSize int) ([]*model.PersonDate, *pagination.PaginationMeta, error)
	UpdatePersonDates(id uuid.UUID, personID uuid.UUID, input UpdatePersonDatesInput) (*model.PersonDate, error)
	DeletePersonDates(personID uuid.UUID, id uuid.UUID) error
}

type personDatesService struct {
	repo repository.PersonDatesRepository
}

func NewPersonDatesService(repo repository.PersonDatesRepository) PersonDatesService {
	return &personDatesService{repo}
}

func (s *personDatesService) CreatePersonDates(personID uuid.UUID, input CreatePersonDatesInput) (*model.PersonDate, error) {
	personDates := &model.PersonDate{
		Date:  input.Date,
		Label: input.Label,
	}

	personDates, err := s.repo.CreatePersonDates(personID, personDates)
	if err != nil {
		return nil, err
	}
	return personDates, nil
}

func (s *personDatesService) GetAllPersonDates(personID uuid.UUID, page int, pageSize int) ([]*model.PersonDate, *pagination.PaginationMeta, error) {
	if page < 1 {
		page = 1
	}
	if pageSize < 1 || pageSize > 100 {
		pageSize = 20
	}

	personDates, total, err := s.repo.GetAllPersonDates(personID, page, pageSize)
	if err != nil {
		return nil, nil, err
	}

	meta := &pagination.PaginationMeta{
		Page:       page,
		PerPage:    pageSize,
		TotalItems: int64(total),
		TotalPages: int(math.Ceil(float64(total) / float64(pageSize))),
	}

	return personDates, meta, nil
}

func (s *personDatesService) GetPersonDatesByID(personID uuid.UUID, id uuid.UUID) (*model.PersonDate, error) {
	return s.repo.GetPersonDatesByID(personID, id)
}

func (s *personDatesService) UpdatePersonDates(id uuid.UUID, personID uuid.UUID, input UpdatePersonDatesInput) (*model.PersonDate, error) {
	updates := map[string]interface{}{}
	if input.Date != nil {
		updates["date"] = *input.Date
	}
	if input.Label != nil {
		updates["label"] = *input.Label
	}
	if err := s.repo.UpdatePersonDates(id, personID, updates); err != nil {
		return nil, err
	}

	return s.repo.GetPersonDatesByID(personID, id)
}

func (s *personDatesService) DeletePersonDates(personID uuid.UUID, id uuid.UUID) error {
	return s.repo.DeletePersonDates(personID, id)
}
