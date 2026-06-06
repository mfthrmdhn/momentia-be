package repository

import (
	"momentia-be/model"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type PersonMedicalInfoRepository struct {
	// This struct can hold a database connection or any other dependencies needed for the repository
}

type PersonLikesRepository struct {
	// This struct can hold a database connection or any other dependencies needed for the repository
}

type PersonDislikesRepository struct {
	// This struct can hold a database connection or any other dependencies needed for the repository
}

type PersonDatesRepository struct {
	// This struct can hold a database connection or any other dependencies needed for the repository
}

type PersonRepository interface {
	CreatePerson(userID int, person *model.Person) (*model.Person, error)
	GetPersonByID(id uuid.UUID, userID int) (*model.Person, error)
	GetAllPersons(userID int, page int, pageSize int) ([]*model.Person, int, error)
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

type PersonMedicalInfoRepositoryInterface interface {
	Create(personMedicalInfo *model.PersonMedicalInfo) error
	FindByID(id uuid.UUID) (*model.PersonMedicalInfo, error)
	FindAll() ([]*model.PersonMedicalInfo, error)
	Update(personMedicalInfo *model.PersonMedicalInfo) error
	Delete(id uuid.UUID) error
}

type PersonLikesRepositoryInterface interface {
	Create(personLikes *model.PersonLikes) error
	FindByID(id uuid.UUID) (*model.PersonLikes, error)
	FindAll() ([]*model.PersonLikes, error)
	Update(personLikes *model.PersonLikes) error
	Delete(id uuid.UUID) error
}

type PersonDislikesRepositoryInterface interface {
	Create(personDislikes *model.PersonDislikes) error
	FindByID(id uuid.UUID) (*model.PersonDislikes, error)
	FindAll() ([]*model.PersonDislikes, error)
	Update(personDislikes *model.PersonDislikes) error
	Delete(id uuid.UUID) error
}

type PersonDatesRepositoryInterface interface {
	Create(personDate *model.PersonDates) error
	FindByID(id uuid.UUID) (*model.PersonDates, error)
	FindAll() ([]*model.PersonDates, error)
	Update(personDate *model.PersonDates) error
	Delete(id uuid.UUID) error
}
