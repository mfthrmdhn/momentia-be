package model

type User struct {
	ID           int    `gorm:"primaryKey;autoIncrement" json:"id"`
	Username     string `gorm:"type:varchar(255);not null;unique" json:"username"`
	Email        string `gorm:"type:varchar(255);not null;unique" json:"email"`
	Msisdn       string `gorm:"type:varchar(255);not null;unique" json:"msisdn"`
	PasswordHash string `gorm:"type:varchar(255);not null" json:"-"`
}