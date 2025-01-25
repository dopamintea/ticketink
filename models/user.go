package models

import (
	"time"

	"gorm.io/gorm"
)

type User struct {
	ID        uint           `gorm:"primaryKey"`
	Name      string         `gorm:"size:100;not null"`
	Email     string         `gorm:"size:100;unique;not null"`
	Password  string         `gorm:"not null"`
	Role      string         `gorm:"size:20;not null"` // e.g., "user" or "admin"
	CreatedAt time.Time      `gorm:"autoCreateTime"`
	UpdatedAt time.Time      `gorm:"autoUpdateTime"`
	DeletedAt gorm.DeletedAt `gorm:"index"`
}
