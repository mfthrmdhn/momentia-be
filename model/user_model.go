package model

import "time"

type User struct {
	ID           int    `gorm:"primaryKey;autoIncrement" json:"id"`
	Username     string `gorm:"type:varchar(255);not null;unique" json:"username"`
	Email        string `gorm:"type:varchar(255);not null;unique" json:"email"`
	Msisdn       string `gorm:"type:varchar(255);not null;unique" json:"msisdn"`
	PasswordHash string `gorm:"type:varchar(255);not null" json:"-"`
	CreatedAt 	 time.Time `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt 	 time.Time `gorm:"autoUpdateTime" json:"updated_at"`
}