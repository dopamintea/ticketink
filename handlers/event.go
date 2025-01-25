package handlers

import (
	"net/http"
	"strconv"
	"strings"
	"time"

	"ticketink/models"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type EventRequest struct {
	Title       string  `json:"title"`
	Description string  `json:"description"`
	Date        string  `json:"date"`
	Location    string  `json:"location"`
	Price       float64 `json:"price"`
	Capacity    int64   `json:"capacity"`
	Status      string  `json:"status"` // Active, Ongoing, Finished
}

func ListEvents(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
		limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))
		offset := (page - 1) * limit

		category := c.Query("category")
		status := c.Query("status")
		search := c.Query("search")

		query := db.Model(&models.Event{})

		if category != "" {
			query = query.Where("category = ?", category)
		}
		if status != "" {
			query = query.Where("status = ?", status)
		}
		if search != "" {
			query = query.Where("title LIKE ? OR description LIKE ?", "%"+search+"%", "%"+search+"%")
		}

		var totalItems int64
		query.Count(&totalItems)

		var events []models.Event
		query.Limit(limit).Offset(offset).Find(&events)

		totalPages := (int(totalItems) + limit - 1) / limit

		c.JSON(http.StatusOK, gin.H{
			"events": events,
			"pagination": gin.H{
				"current_page": page,
				"total_pages":  totalPages,
				"total_items":  totalItems,
			},
		})
	}
}

func CreateEvent(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		var req EventRequest

		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		var existingEvent models.Event
		if err := db.Where("title = ?", req.Title).First(&existingEvent).Error; err == nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Event title must be unique"})
			return
		}

		eventDate, err := time.Parse("2006-01-02", req.Date)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid date format, use YYYY-MM-DD"})
			return
		}

		event := models.Event{
			Title:       req.Title,
			Description: req.Description,
			Date:        eventDate,
			Location:    req.Location,
			Price:       req.Price,
			Capacity:    req.Capacity,
			Status:      "Active",
		}

		if err := db.Create(&event).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create event"})
			return
		}

		c.JSON(http.StatusCreated, event)
	}
}

func UpdateEvent(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		id := c.Param("id")

		var event models.Event
		if err := db.First(&event, id).Error; err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "Event not found"})
			return
		}

		var req EventRequest
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		if strings.EqualFold(event.Status, "finished") || event.Date.Before(time.Now()) {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Cannot update an event that has already started or finished"})
			return
		}

		event.Title = req.Title
		event.Description = req.Description
		event.Location = req.Location
		event.Price = req.Price
		event.Capacity = req.Capacity

		if err := db.Save(&event).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update event"})
			return
		}

		c.JSON(http.StatusOK, event)
	}
}

func UpdateEventStatus(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {

		id := c.Param("id")

		var event models.Event
		if err := db.First(&event, id).Error; err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "Event not found"})
			return
		}

		var req struct {
			Status string `json:"status" binding:"required,oneof=ongoing finished"`
		}

		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		if strings.EqualFold(event.Status, "finished") {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Event is already finished and cannot be updated"})
			return
		}

		event.Status = req.Status

		if err := db.Save(&event).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update event status"})
			return
		}

		c.JSON(http.StatusOK, event)
	}
}

func DeleteEvent(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		id := c.Param("id")

		var event models.Event
		if err := db.First(&event, id).Error; err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "Event not found"})
			return
		}

		var ticketsSold int64
		db.Model(&models.Ticket{}).Where("event_id = ?", id).Count(&ticketsSold)
		if ticketsSold > 0 {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Cannot delete an event with sold tickets"})
			return
		}

		if err := db.Delete(&event).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete event"})
			return
		}

		c.JSON(http.StatusOK, gin.H{"message": "Event deleted successfully"})
	}
}
