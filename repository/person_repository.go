package repository

import (
	"momentia-be/model"

	"github.com/google/uuid"
	"gorm.io/gorm"
)


type PersonRepository interface {
	CreatePerson(userID int, person *model.Person) (*model.Person, error)
	GetPersonByID(id uuid.UUID, userID int) (*model.Person, error)
	GetAllPersons(userID int, page int, pageSize int) ([]*model.Person, int, error)
	UpdatePerson(id uuid.UUID, userID int, updates map[string]interface{}) error
	DeletePerson(id uuid.UUID) error
}

type personRepository struct {
	db *gorm.DB
}

func newPersonRepository(db *gorm.DB) PersonRepository {
	return &personRepository{db}
}

func NewPersonRepository(db *gorm.DB) PersonRepository {
	return newPersonRepository(db)
}

func (r *personRepository) CreatePerson(userID int, person *model.Person) (*model.Person, error) {
	person.CreatorUserID = userID
	if err := r.db.Create(person).Error; err != nil {
		return nil, err
	}
	return person, nil
}

func (r *personRepository) GetPersonByID(id uuid.UUID, userID int) (*model.Person, error) {
	var person model.Person
	if err := r.db.First(&person, "id = ? AND creator_user_id = ?", id, userID).Error; err != nil {
		return nil, err
	}
	return &person, nil
}

func (r *personRepository) GetAllPersons(userID int, page int, pageSize int) ([]*model.Person, int, error) {
	var persons []*model.Person
	var total int64

	if err := r.db.Model(&model.Person{}).Where("creator_user_id = ?", userID).Count(&total).Error; err != nil {
		return nil, 0, err
	}

	if err := r.db.Where("creator_user_id = ?", userID).Limit(pageSize).Offset((page - 1) * pageSize).Find(&persons).Error; err != nil {
		return nil, 0, err
	}

	return persons, int(total), nil
}

func (r *personRepository) UpdatePerson(id uuid.UUID, userID int, updates map[string]interface{}) error {
	result := r.db.Model(&model.Person{}).Where("id = ? AND creator_user_id = ?", id, userID).Updates(updates)
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}
	return nil
}

func (r *personRepository) DeletePerson(id uuid.UUID) error {
	result := r.db.Where("id = ?", id).Delete(&model.Person{})
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}
	return nil
}

