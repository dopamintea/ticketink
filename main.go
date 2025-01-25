package main

import (
	"log"
	"ticketink/config"
	"ticketink/migrations"
	"ticketink/routes"

	"github.com/gin-gonic/gin"
)

func main() {
	db, err := config.ConnectDatabase()
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}

	migrations.RunMigrations(db)

	r := gin.Default()

	routes.RegisterRoutes(r, db)

	log.Println("Server is running on http://localhost:8080")
	if err := r.Run(":8080"); err != nil {
		log.Fatal("Failed to start server:", err)
	}
}
