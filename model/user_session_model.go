package model

import "time"

type UserSession struct {
	ID        int       `gorm:"primaryKey;autoIncrement" json:"id"`
	UserID    int       `gorm:"not null;index" json:"user_id"`
	TokenHash string    `gorm:"type:varchar(64);not null;uniqueIndex" json:"token_hash"`
	ExpiresAt time.Time `gorm:"not null" json:"expires_at"`
	CreatedAt time.Time `gorm:"autoCreateTime" json:"created_at"`
}

func (UserSession) TableName() string {
	return "user_sessions"
}
