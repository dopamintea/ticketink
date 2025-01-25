package routes

import (
	"ticketink/handlers"
	"ticketink/middleware"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func RegisterRoutes(r *gin.Engine, db *gorm.DB) {

	r.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{"message": "pong"})
	})

	r.POST("/login", handlers.Login(db))
	r.POST("/register", handlers.Register(db))

	api := r.Group("/api")
	api.Use(middleware.AuthMiddleware(db))
	{
		api.GET("/events", handlers.ListEvents(db))

		api.GET("/tickets", handlers.GetTickets(db))
		api.POST("/tickets", handlers.PurchaseTicket(db))
		api.GET("/tickets/:id", handlers.GetTicketByID(db))
		api.PATCH("/tickets/:id", handlers.UpdateTicket(db))

		api.POST("/logout", handlers.Logout(db))
	}

	admin := api.Group("/admin")
	admin.Use(middleware.AdminOnly())
	{
		admin.GET("/reports/summary", handlers.GetRevenueSummary(db))
		admin.GET("/reports/event/:id", handlers.GetEventReport(db))
		admin.POST("/events", handlers.CreateEvent(db))
		admin.PUT("/events/:id", handlers.UpdateEvent(db))
		admin.PATCH("/events/:id", handlers.UpdateEventStatus(db))
		admin.DELETE("/events/:id", handlers.DeleteEvent(db))
	}
}
