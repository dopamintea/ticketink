package migrations

import (
	"log"
	"ticketink/models"

	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

func RunMigrations(db *gorm.DB) {
	err := db.AutoMigrate(
		&models.User{},
		&models.Ticket{},
		&models.Event{},
		&models.Report{},
		&models.TokenBlacklist{},
	)
	if err != nil {
		log.Fatal("Failed to run migrations:", err)
	} else {
		log.Println("Database migrated successfully!")
	}

	var adminCount int64
	db.Model(&models.User{}).Where("role = ?", "admin").Count(&adminCount)

	if adminCount == 0 {

		hashedPassword, err := bcrypt.GenerateFromPassword([]byte("admin"), bcrypt.DefaultCost)
		if err != nil {
			log.Println(err)
		}

		admin := models.User{
			Name:     "Rinan",
			Email:    "Admin@gmail.com",
			Password: string(hashedPassword),
			Role:     "admin",
		}

		if err := db.Create(&admin).Error; err != nil {
			log.Println(err)
		}
	}
	if err != nil {
		log.Fatal("Failed to run migrations:", err)
	} else {
		log.Println("Database migrated successfully!")
	}
}
