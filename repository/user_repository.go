package repository

import (
	"momentia-be/model"

	"gorm.io/gorm"
)

type UserRepository interface {
	Register(user *model.User) error
	Login(email string) (*model.User, error)
	Logout(id int) (*model.User, error)
	GetUserByID(id int) (*model.User, error)
	GetUserByEmail(email string) (*model.User, error)
	GetUserByUsername(username string) (*model.User, error)
	GetUserByMsisdn(msisdn string) (*model.User, error)
	UpdateUser(user *model.User) error
}

type userRepository struct {
	db *gorm.DB
}

func NewUserRepository(db *gorm.DB) UserRepository {
	return &userRepository{db}
}

func (r *userRepository) Register(user *model.User) error {
	return r.db.Create(user).Error
}

func (r *userRepository) Login(email string) (*model.User, error) {
	var user model.User
	if err := r.db.Where("email = ?", email).Limit(1).Find(&user).Error; err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *userRepository) Logout(id int) (*model.User, error) {
	var user model.User
	if err := r.db.First(&user, "id = ?", id).Error; err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *userRepository) GetUserByID(id int) (*model.User, error) {
	var user model.User
	if err := r.db.First(&user, "id = ?", id).Error; err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *userRepository) GetUserByEmail(email string) (*model.User, error) {
	var user model.User
	if err := r.db.Where("email = ?", email).Limit(1).Find(&user).Error; err != nil {
		return nil, err
	}
	if user.ID == 0 {
		return nil, nil
	}
	return &user, nil
}

func (r *userRepository) GetUserByUsername(username string) (*model.User, error) {
	var user model.User
	if err := r.db.Where("username = ?", username).Limit(1).Find(&user).Error; err != nil {
		return nil, err
	}
	if user.ID == 0 {
		return nil, nil
	}
	return &user, nil
}

func (r *userRepository) GetUserByMsisdn(msisdn string) (*model.User, error) {
	var user model.User
	if err := r.db.Where("msisdn = ?", msisdn).Limit(1).Find(&user).Error; err != nil {
		return nil, err
	}
	if user.ID == 0 {
		return nil, nil
	}
	return &user, nil
}

func (r *userRepository) UpdateUser(user *model.User) error {
	return r.db.Save(user).Error
}

