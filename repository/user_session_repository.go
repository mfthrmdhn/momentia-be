package repository

import (
	"momentia-be/model"
	"time"

	"gorm.io/gorm"
)

type UserSessionRepository interface {
	Create(session *model.UserSession) error
	DeleteByTokenHash(tokenHash string) error
	FindByTokenHash(tokenHash string) (*model.UserSession, error)
	DeleteExpired() error
}

type userSessionRepository struct {
	db *gorm.DB
}

func NewUserSessionRepository(db *gorm.DB) UserSessionRepository {
	return &userSessionRepository{db}
}

func (r *userSessionRepository) Create(session *model.UserSession) error {
	return r.db.Create(session).Error
}

func (r *userSessionRepository) DeleteByTokenHash(tokenHash string) error {
	return r.db.Where("token_hash = ?", tokenHash).Delete(&model.UserSession{}).Error
}

func (r *userSessionRepository) FindByTokenHash(tokenHash string) (*model.UserSession, error) {
	var session model.UserSession
	err := r.db.Where("token_hash = ? AND expires_at > ?", tokenHash, time.Now()).First(&session).Error
	if err != nil {
		return nil, err
	}
	return &session, nil
}

func (r *userSessionRepository) DeleteExpired() error {
	return r.db.Where("expires_at <= ?", time.Now()).Delete(&model.UserSession{}).Error
}
