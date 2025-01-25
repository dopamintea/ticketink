package models

type Report struct {
	TotalTicketsSold int     `json:"total_tickets_sold"`
	TotalRevenue     float64 `json:"total_revenue"`
	EventID          uint    `json:"event_id"`
	EventTitle       string  `json:"event_title"`
	TicketsSold      int64   `json:"tickets_sold"`
	RevenueGenerated float64 `json:"revenue_generated"`
}
