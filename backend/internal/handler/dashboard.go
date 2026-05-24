package handler

import (
	"database/sql"
	"net/http"

	"github.com/ayomendaki/ayomendaki-admin/internal/common"
)

type dashboardData struct {
	OperatorName   string
	TotalTrips     int
	BookingsMonth  int
	RevenueMonth   int64
	NearFullTrips  []tripQuota
	RecentBookings []recentBooking
	MonthlyStats   []monthlyStat
}

type tripQuota struct {
	ID             int64
	Name           string
	QuotaRemaining int
	QuotaTotal     int
}

type recentBooking struct {
	ID       int64
	Customer string
	TripName string
	Date     string
	Total    int64
	Status   string
}

type monthlyStat struct {
	Month   string
	Bookings int
	Revenue  int64
}

func (h *Handler) Dashboard(w http.ResponseWriter, r *http.Request) {
	operatorID := common.GetOperatorID(r)

	var operatorName string
	h.db.QueryRow("SELECT name FROM operators WHERE id = ?", operatorID).Scan(&operatorName)

	var totalTrips int
	h.db.QueryRow("SELECT COUNT(*) FROM trips WHERE operator_id = ?", operatorID).Scan(&totalTrips)

	var bookingsMonth int
	h.db.QueryRow(`SELECT COUNT(*) FROM bookings b JOIN trips t ON b.trip_id = t.id WHERE t.operator_id = ? AND strftime('%Y-%m', b.created_at) = strftime('%Y-%m', 'now')`, operatorID).Scan(&bookingsMonth)

	var revenueMonth sql.NullInt64
	h.db.QueryRow(`SELECT COALESCE(SUM(b.total), 0) FROM bookings b JOIN trips t ON b.trip_id = t.id WHERE t.operator_id = ? AND b.status IN ('confirmed', 'completed') AND strftime('%Y-%m', b.created_at) = strftime('%Y-%m', 'now')`, operatorID).Scan(&revenueMonth)
	revMonth := revenueMonth.Int64

	rows, err := h.db.Query(`SELECT t.id, t.name, COALESCE(s.quota_remaining, 0), COALESCE(s.quota_total, 0) FROM trips t LEFT JOIN schedules s ON s.trip_id = t.id WHERE t.operator_id = ? AND s.quota_remaining < 5 AND s.quota_remaining > 0 GROUP BY t.id LIMIT 5`, operatorID)
	nearFull := []tripQuota{}
	if err == nil {
		defer rows.Close()
		for rows.Next() {
			var t tripQuota
			rows.Scan(&t.ID, &t.Name, &t.QuotaRemaining, &t.QuotaTotal)
			nearFull = append(nearFull, t)
		}
	}

	rRows, err := h.db.Query(`SELECT b.id, b.lead_name, t.name, b.created_at, b.total, b.status FROM bookings b JOIN trips t ON b.trip_id = t.id WHERE t.operator_id = ? ORDER BY b.created_at DESC LIMIT 5`, operatorID)
	recent := []recentBooking{}
	if err == nil {
		defer rRows.Close()
		for rRows.Next() {
			var rb recentBooking
			rRows.Scan(&rb.ID, &rb.Customer, &rb.TripName, &rb.Date, &rb.Total, &rb.Status)
			recent = append(recent, rb)
		}
	}

	months := []monthlyStat{}
	for i := 5; i >= 0; i-- {
		var bookings int
		var revenue sql.NullInt64
		h.db.QueryRow(`SELECT COUNT(*), COALESCE(SUM(b.total), 0) FROM bookings b JOIN trips t ON b.trip_id = t.id WHERE t.operator_id = ? AND b.status IN ('confirmed', 'completed') AND strftime('%Y-%m', b.created_at) = strftime('%Y-%m', 'now', ? || ' months')`, operatorID, -i).Scan(&bookings, &revenue)
		monthLabel := ""
		h.db.QueryRow(`SELECT strftime('%m', 'now', ? || ' months')`, -i).Scan(&monthLabel)
		months = append(months, monthlyStat{
			Month:    monthLabel,
			Bookings: bookings,
			Revenue:  revenue.Int64,
		})
	}

	data := dashboardData{
		OperatorName:   operatorName,
		TotalTrips:     totalTrips,
		BookingsMonth:  bookingsMonth,
		RevenueMonth:   revMonth,
		NearFullTrips:  nearFull,
		RecentBookings: recent,
		MonthlyStats:   months,
	}

	h.renderer.RenderTemplate(w, "dashboard", data)
}
