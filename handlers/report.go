package handlers

import (
	"net/http"
	"ticketink/models"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func GetEventReport(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		id := c.Param("id")

		var event models.Event
		if err := db.First(&event, id).Error; err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "Event not found"})
			return
		}

		var ticketsSold int64
		var totalRevenue float64
		db.Model(&models.Ticket{}).
			Where("event_id = ? AND status = ?", id, "purchased").
			Count(&ticketsSold).
			Select("SUM(price)").
			Row().
			Scan(&totalRevenue)

		report := models.Report{
			EventID:          event.ID,
			EventTitle:       event.Title,
			TicketsSold:      int64(ticketsSold),
			RevenueGenerated: totalRevenue,
		}

		c.JSON(http.StatusOK, report)
	}
}

func GetRevenueSummary(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {

		startDate := c.Query("start_date")
		endDate := c.Query("end_date")

		query := db.Model(&models.Ticket{}).Where("status = ?", "purchased")

		if startDate != "" && endDate != "" {
			query = query.Where("created_at BETWEEN ? AND ?", startDate, endDate)
		}

		var totalTicketsSold int64
		var totalRevenue float64
		query.Count(&totalTicketsSold).
			Select("SUM(price)").
			Row().
			Scan(&totalRevenue)

		c.JSON(http.StatusOK, gin.H{
			"total_tickets_sold": totalTicketsSold,
			"total_revenue":      totalRevenue,
		})
	}
}
