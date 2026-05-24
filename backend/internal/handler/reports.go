package handler

import (
	"database/sql"
	"encoding/json"
	"net/http"

	"github.com/ayomendaki/ayomendaki-admin/internal/common"
)

type revenueData struct {
	RevenueYear   int64
	TotalBookings int
	AvgBooking    float64
	TopTrip       string
	MonthlyStats  []monthlyStat
	TripBreakdown []tripRevenue
}

type tripRevenue struct {
	Name       string
	Bookings   int
	Participants int
	Revenue    int64
}

func (h *Handler) RevenueReport(w http.ResponseWriter, r *http.Request) {
	operatorID := common.GetOperatorID(r)

	var revenueYear sql.NullInt64
	h.db.QueryRow(`SELECT COALESCE(SUM(b.total), 0) FROM bookings b JOIN trips t ON b.trip_id = t.id WHERE t.operator_id = ? AND b.status IN ('confirmed', 'completed') AND strftime('%Y', b.created_at) = strftime('%Y', 'now')`, operatorID).Scan(&revenueYear)

	var totalBookings int
	h.db.QueryRow(`SELECT COUNT(*) FROM bookings b JOIN trips t ON b.trip_id = t.id WHERE t.operator_id = ? AND b.status IN ('confirmed', 'completed') AND strftime('%Y', b.created_at) = strftime('%Y', 'now')`, operatorID).Scan(&totalBookings)

	var avgBooking int64
	var rawAvg sql.NullFloat64
	if totalBookings > 0 {
		h.db.QueryRow(`SELECT ROUND(CAST(COALESCE(SUM(b.total), 0) AS REAL) / COUNT(*), 0) FROM bookings b JOIN trips t ON b.trip_id = t.id WHERE t.operator_id = ? AND b.status IN ('confirmed', 'completed') AND strftime('%Y', b.created_at) = strftime('%Y', 'now')`, operatorID).Scan(&rawAvg)
		avgBooking = int64(rawAvg.Float64)
	}

	var topTrip string
	h.db.QueryRow(`SELECT t.name FROM trips t JOIN bookings b ON b.trip_id = t.id WHERE t.operator_id = ? AND b.status IN ('confirmed', 'completed') AND strftime('%Y', b.created_at) = strftime('%Y', 'now') GROUP BY t.id ORDER BY SUM(b.total) DESC LIMIT 1`, operatorID).Scan(&topTrip)

	months := []monthlyStat{}
	for i := 11; i >= 0; i-- {
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

	tRows, err := h.db.Query(`
		SELECT t.name,
			COUNT(b.id),
			COALESCE(SUM((SELECT COUNT(*) FROM booking_participants bp WHERE bp.booking_id = b.id)), 0),
			COALESCE(SUM(b.total), 0)
		FROM trips t
		LEFT JOIN bookings b ON b.trip_id = t.id AND b.status IN ('confirmed', 'completed') AND strftime('%Y', b.created_at) = strftime('%Y', 'now')
		WHERE t.operator_id = ?
		GROUP BY t.id
		ORDER BY COALESCE(SUM(b.total), 0) DESC
	`, operatorID)
	tripBreakdown := []tripRevenue{}
	if err == nil {
		defer tRows.Close()
		for tRows.Next() {
			var tr tripRevenue
			tRows.Scan(&tr.Name, &tr.Bookings, &tr.Participants, &tr.Revenue)
			tripBreakdown = append(tripBreakdown, tr)
		}
	}

	monthsJSON, _ := json.Marshal(months)

	h.renderer.RenderTemplate(w, "reports/revenue", map[string]interface{}{
		"RevenueYear":    revenueYear.Int64,
		"TotalBookings":  totalBookings,
		"AvgBooking":     avgBooking,
		"TopTrip":        topTrip,
		"MonthlyStats":   months,
		"MonthlyStatsJSON": string(monthsJSON),
		"TripBreakdown":  tripBreakdown,
	})
}
