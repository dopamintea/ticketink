package models

import (
	"time"

	"gorm.io/gorm"
)

type Ticket struct {
	ID        uint           `gorm:"primaryKey"`
	UserID    uint           `gorm:"not null"`
	User      User           `gorm:"foreignKey:UserID"`
	EventID   uint           `gorm:"not null"`
	Event     Event          `gorm:"foreignKey:EventID"`
	Status    string         `gorm:"size:20;not null"` // e.g., "purchased", "cancelled"
	Price     float64        `gorm:"not null;check:price >= 0"`
	CreatedAt time.Time      `gorm:"autoCreateTime"`
	UpdatedAt time.Time      `gorm:"autoUpdateTime"`
	DeletedAt gorm.DeletedAt `gorm:"index"`
}
