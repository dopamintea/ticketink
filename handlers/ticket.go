package handlers

import (
	"net/http"
	"strconv"
	"strings"
	"ticketink/models"
	"time"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func GetTickets(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {

		page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
		limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))
		offset := (page - 1) * limit

		userID := c.Query("user_id")
		eventID := c.Query("event_id")
		status := c.Query("status")

		query := db.Model(&models.Ticket{}).Preload("Event").Preload("User")
		if userID != "" {
			query = query.Where("user_id = ?", userID)
		}
		if eventID != "" {
			query = query.Where("event_id = ?", eventID)
		}
		if status != "" {
			query = query.Where("status = ?", status)
		}

		var tickets []models.Ticket
		var totalItems int64
		query.Count(&totalItems).Limit(limit).Offset(offset).Find(&tickets)

		totalPages := int((totalItems + int64(limit) - 1) / int64(limit))

		c.JSON(http.StatusOK, gin.H{
			"tickets": tickets,
			"pagination": gin.H{
				"current_page": page,
				"total_pages":  totalPages,
				"total_items":  totalItems,
			},
		})
	}
}

func GetTicketByID(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		id := c.Param("id")

		var ticket models.Ticket
		if err := db.Preload("Event").Preload("User").First(&ticket, id).Error; err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "Ticket not found"})
			return
		}

		c.JSON(http.StatusOK, ticket)
	}
}

func PurchaseTicket(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {

		var req struct {
			Email   string `json:"email" binding:"required"`
			EventID uint   `json:"event_id" binding:"required"`
		}

		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request: " + err.Error()})
			return
		}

		var user models.User
		if err := db.Where("email = ?", req.Email).First(&user).Error; err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
			return
		}

		var event models.Event
		if err := db.First(&event, req.EventID).Error; err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "Event not found"})
			return
		}

		if strings.EqualFold(event.Status, "finished") || event.Date.Before(time.Now()) {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Cannot purchase a ticket for event that has already finished"})
			return
		}

		var ticketsSold int64
		if err := db.Model(&models.Ticket{}).Where("event_id = ? AND status = ?", req.EventID, "purchased").Count(&ticketsSold).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to check ticket availability"})
			return
		}
		if ticketsSold >= event.Capacity {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Event is sold out"})
			return
		}

		ticket := models.Ticket{
			UserID:  user.ID,
			EventID: req.EventID,
			Status:  "purchased",
			Price:   event.Price,
		}

		if err := db.Create(&ticket).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to purchase ticket"})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"message": "Ticket purchased successfully",
			"ticket":  ticket,
		})
	}
}

func UpdateTicket(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		id := c.Param("id")

		var req struct {
			Status string `json:"status" binding:"required"`
		}

		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
			return
		}

		validStatuses := map[string]bool{
			"purchased": true,
			"cancelled": true,
		}
		if !validStatuses[req.Status] {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid status"})
			return
		}

		var ticket models.Ticket
		if err := db.Preload("Event").First(&ticket, id).Error; err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "Ticket not found"})
			return
		}

		if ticket.Event.Date.Before(time.Now()) {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Cannot update tickets for past events"})
			return
		}

		ticket.Status = req.Status
		if err := db.Save(&ticket).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update ticket"})
			return
		}

		c.JSON(http.StatusOK, gin.H{"message": "Ticket updated successfully", "ticket": ticket})
	}
}
