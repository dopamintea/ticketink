package models

import (
	"time"

	"gorm.io/gorm"
)

type TokenBlacklist struct {
	ID        uint      `gorm:"primaryKey"`
	Token     string    `gorm:"type:text;not null"`
	Email     string    `gorm:"not null"`
	ExpiresAt time.Time `gorm:"not null"`
	CreatedAt time.Time `gorm:"autoCreateTime"`
}

func CleanupBlacklist(db *gorm.DB) {
	db.Where("expires_at < ?", time.Now()).Delete(&TokenBlacklist{})
}
