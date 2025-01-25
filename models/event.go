package models

import (
	"time"

	"gorm.io/gorm"
)

type Event struct {
	ID          uint           `gorm:"primaryKey"`
	Title       string         `gorm:"size:200;not null;unique"`
	Description string         `gorm:"type:text"`
	Date        time.Time      `gorm:"not null"`
	Location    string         `gorm:"size:255;not null"`
	Price       float64        `gorm:"not null;check:price >= 0"`
	Capacity    int64          `gorm:"not null;check:capacity >= 0"`
	Status      string         `gorm:"size:20;not null"` // e.g., "active", "ongoing", "completed"
	CreatedAt   time.Time      `gorm:"autoCreateTime"`
	UpdatedAt   time.Time      `gorm:"autoUpdateTime"`
	DeletedAt   gorm.DeletedAt `gorm:"index"`
	Tickets     []Ticket       `gorm:"constraint:OnDelete:CASCADE"`
}
