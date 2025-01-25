package config

import (
	"log"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func ConnectDatabase() (*gorm.DB, error) {

	dsn := "user:pass@tcp(127.0.0.1:3306)/ticketink?charset=utf8mb4&parseTime=True&loc=Local"

	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Println("Error connecting to database:", err)
		return nil, err
	}

	log.Println("Database connection successful!")
	return db, nil
}
