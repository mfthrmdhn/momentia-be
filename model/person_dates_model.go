package model

import (
	"time"

	"github.com/google/uuid"
)

type PersonDate struct {
	ID 		 uuid.UUID `gorm:"type:uuid;default:uuid_generate_v7();primaryKey" json:"id"`
	PersonID string `gorm:"type:uuid;not null" json:"person_id"`
	Label    string `gorm:"type:varchar(100);not null" json:"label"`
	Date     string `gorm:"type:date;not null" json:"date"`
	Note     string `gorm:"type:varchar(255)" json:"note"`
	CreatedAt time.Time `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt time.Time `gorm:"autoUpdateTime" json:"updated_at"`
}
