package models

import "github.com/google/uuid"

type User struct {
	ID           uuid.UUID `gorm:"type:uuid;default:gen_random_uuid();primaryKey"`
	Email        string    `gorm:"unique;not null"`
	PasswordHash string    `gorm:"not null"`
	Name         string    `gorm:"not null"`
	Status       string    `gorm:"default:active"`
	CreatedAt    int64     `gorm:"autoCreateTime"`
	UpdatedAt    int64     `gorm:"autoUpdateTime"`
}
