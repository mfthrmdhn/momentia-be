package repository

import (
	"momentia-be/model"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type PersonDatesRepository interface {
	CreatePersonDates(personID uuid.UUID, personDates *model.PersonDate) (*model.PersonDate, error)
	GetPersonDatesByID(personID uuid.UUID, id uuid.UUID) (*model.PersonDate, error)
	GetAllPersonDates(personID uuid.UUID, page int, pageSize int) ([]*model.PersonDate, int, error)
	UpdatePersonDates(id uuid.UUID, personID uuid.UUID, updates map[string]interface{}) error
	DeletePersonDates(personID uuid.UUID, id uuid.UUID) error
}

type personDates struct {
	db *gorm.DB
}

func NewPersonDatesRepository(db *gorm.DB) PersonDatesRepository {
	return &personDates{db}
}

func (r *personDates) CreatePersonDates(personID uuid.UUID, personDates *model.PersonDate) (*model.PersonDate, error) {
	personDates.PersonID = personID.String()
	if err := r.db.Create(personDates).Error; err != nil {
		return nil, err
	}
	return personDates, nil
}

func (r *personDates) GetAllPersonDates(personID uuid.UUID, page int, pageSize int) ([]*model.PersonDate, int, error) {
	var personDates []*model.PersonDate
	var total int64

	if err := r.db.Model(&model.PersonDate{}).Where("person_id = ?", personID).Count(&total).Error; err != nil {
		return nil, 0, err
	}

	if err := r.db.Where("person_id = ?", personID).Limit(pageSize).Offset((page - 1) * pageSize).Find(&personDates).Error; err != nil {
		return nil, 0, err
	}

	return personDates, int(total), nil
}

func (r *personDates) GetPersonDatesByID(personID uuid.UUID, id uuid.UUID) (*model.PersonDate, error) {
	var personDates model.PersonDate
	if err := r.db.First(&personDates, "id = ? AND person_id = ?", id, personID).Error; err != nil {
		return nil, err
	}
	return &personDates, nil
}

func (r *personDates) UpdatePersonDates(id uuid.UUID, personID uuid.UUID, updates map[string]interface{}) error {
	result := r.db.Model(&model.PersonDate{}).Where("id = ? AND person_id = ?", id, personID).Updates(updates)
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}
	return nil
}

func (r *personDates) DeletePersonDates(personID uuid.UUID, id uuid.UUID) error {
	return r.db.Where("person_id = ? AND id = ?", personID, id).Delete(&model.PersonDate{}).Error
}
