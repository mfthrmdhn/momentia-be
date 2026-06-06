package model

import (
	"time"

	"github.com/google/uuid"
)

type Person struct {
	ID           	uuid.UUID `gorm:"type:uuid;default:uuid_generate_v7()" json:"id"`
	CreatorUserID   int    `gorm:"not null" json:"creator_user_id"`
	Name 			string `gorm:"type:varchar(255);not null" json:"name"`
	Relationship 	string `gorm:"type:varchar(255);not null" json:"relationship"`
	IsPinned 		bool `gorm:"default:false" json:"is_pinned"`
	CreatedAt 	 	time.Time `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt 	 	time.Time `gorm:"autoUpdateTime" json:"updated_at"`
}

func (Person) TableName() string {
	return "person"
}