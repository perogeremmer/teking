package handler

import (
	"net/http"
	"strconv"

	"github.com/ayomendaki/ayomendaki-admin/internal/common"
)

type bookingListItem struct {
	ID           int64
	Customer     string
	Phone        string
	TripName     string
	TripDate     string
	Participants int
	Total        int64
	Status       string
}

type bookingDetailData struct {
	BookingID    int64
	TripName     string
	LeadName     string
	LeadPhone    string
	LeadEmail    string
	Total        int64
	Status       string
	CreatedAt    string
	ScheduleDate string
	Participants []struct {
		Name string
		KTP  string
	}
	Addons []struct {
		Name  string
		Price int64
	}
}

func (h *Handler) BookingList(w http.ResponseWriter, r *http.Request) {
	operatorID := common.GetOperatorID(r)

	statusFilter := r.URL.Query().Get("status")
	tripFilter := r.URL.Query().Get("trip_id")
	search := r.URL.Query().Get("search")

	query := `
		SELECT b.id, b.lead_name, b.lead_phone, COALESCE(t.name, ''), b.created_at,
			(SELECT COUNT(*) FROM booking_participants bp WHERE bp.booking_id = b.id),
			b.total, b.status
		FROM bookings b
		JOIN trips t ON t.id = b.trip_id
		WHERE t.operator_id = ?
	`
	args := []interface{}{operatorID}

	if statusFilter != "" {
		query += " AND b.status = ?"
		args = append(args, statusFilter)
	}
	if tripFilter != "" {
		query += " AND b.trip_id = ?"
		tripFilterInt, _ := strconv.ParseInt(tripFilter, 10, 64)
		args = append(args, tripFilterInt)
	}
	if search != "" {
		query += " AND (b.lead_name LIKE ? OR b.lead_phone LIKE ?)"
		s := "%" + search + "%"
		args = append(args, s, s)
	}
	query += " ORDER BY b.created_at DESC"

	rows, err := h.db.Query(query, args...)
	bookings := []bookingListItem{}
	if err == nil {
		defer rows.Close()
		for rows.Next() {
			var b bookingListItem
			rows.Scan(&b.ID, &b.Customer, &b.Phone, &b.TripName, &b.TripDate, &b.Participants, &b.Total, &b.Status)
			bookings = append(bookings, b)
		}
	}

	tripRows, _ := h.db.Query("SELECT id, name FROM trips WHERE operator_id = ? ORDER BY name", operatorID)
	type tripOpt struct {
		ID   string
		Name string
	}
	trips := []tripOpt{}
	if tripRows != nil {
		defer tripRows.Close()
		for tripRows.Next() {
			var t tripOpt
			var id int64
			tripRows.Scan(&id, &t.Name)
			t.ID = strconv.FormatInt(id, 10)
			trips = append(trips, t)
		}
	}

	isHTMX := r.Header.Get("HX-Request") == "true"
	if isHTMX {
		h.renderer.RenderTemplate(w, "bookings/index", map[string]interface{}{
			"Bookings":     bookings,
			"Trips":        trips,
			"StatusFilter": statusFilter,
			"Search":       search,
		})
		return
	}

	h.renderer.RenderTemplate(w, "bookings/index", map[string]interface{}{
		"Bookings":     bookings,
		"Trips":        trips,
		"StatusFilter": statusFilter,
		"Search":       search,
	})
}

func (h *Handler) BookingDetail(w http.ResponseWriter, r *http.Request) {
	operatorID := common.GetOperatorID(r)
	id, _ := strconv.ParseInt(r.PathValue("id"), 10, 64)

	d := bookingDetailData{BookingID: id}

	err := h.db.QueryRow(`
		SELECT b.lead_name, b.lead_phone, b.lead_email, b.total, b.status, b.created_at,
			COALESCE(t.name, ''), COALESCE(s.date_start || ' - ' || s.date_end, '')
		FROM bookings b
		JOIN trips t ON t.id = b.trip_id
		LEFT JOIN schedules s ON s.id = b.schedule_id
		WHERE b.id = ? AND t.operator_id = ?
	`, id, operatorID).Scan(&d.LeadName, &d.LeadPhone, &d.LeadEmail, &d.Total, &d.Status, &d.CreatedAt, &d.TripName, &d.ScheduleDate)

	if err != nil {
		http.Redirect(w, r, "/bookings?flash=Booking tidak ditemukan&flash_type=error", http.StatusSeeOther)
		return
	}

	pRows, _ := h.db.Query("SELECT name, ktp FROM booking_participants WHERE booking_id = ?", id)
	if pRows != nil {
		defer pRows.Close()
		for pRows.Next() {
			var p struct {
				Name string
				KTP  string
			}
			pRows.Scan(&p.Name, &p.KTP)
			d.Participants = append(d.Participants, p)
		}
	}

	aRows, _ := h.db.Query("SELECT name, price FROM booking_addons WHERE booking_id = ?", id)
	if aRows != nil {
		defer aRows.Close()
		for aRows.Next() {
			var a struct {
				Name  string
				Price int64
			}
			aRows.Scan(&a.Name, &a.Price)
			d.Addons = append(d.Addons, a)
		}
	}

	h.renderer.RenderTemplate(w, "bookings/detail", map[string]interface{}{
		"Booking": d,
	})
}

func (h *Handler) BookingStatus(w http.ResponseWriter, r *http.Request) {
	operatorID := common.GetOperatorID(r)
	id, _ := strconv.ParseInt(r.PathValue("id"), 10, 64)

	r.ParseForm()
	newStatus := r.FormValue("status")

	validStatuses := map[string]bool{
		"confirmed": true,
		"completed": true,
		"cancelled": true,
	}
	if !validStatuses[newStatus] {
		http.Error(w, "Invalid status", http.StatusBadRequest)
		return
	}

	var currentStatus string
	var scheduleID int64
	err := h.db.QueryRow(`
		SELECT b.status, b.schedule_id FROM bookings b
		JOIN trips t ON t.id = b.trip_id
		WHERE b.id = ? AND t.operator_id = ?
	`, id, operatorID).Scan(&currentStatus, &scheduleID)

	if err != nil {
		http.Error(w, "Not found", http.StatusNotFound)
		return
	}

	if currentStatus == "cancelled" || currentStatus == "completed" {
		http.Error(w, "Cannot change status", http.StatusBadRequest)
		return
	}

	if newStatus == "confirmed" && currentStatus != "pending" {
		http.Error(w, "Can only confirm pending bookings", http.StatusBadRequest)
		return
	}
	if newStatus == "completed" && currentStatus != "confirmed" {
		http.Error(w, "Can only complete confirmed bookings", http.StatusBadRequest)
		return
	}

	tx, _ := h.db.Begin()

	_, err = tx.Exec("UPDATE bookings SET status = ? WHERE id = ?", newStatus, id)
	if err != nil {
		tx.Rollback()
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if newStatus == "cancelled" {
		tx.Exec("UPDATE schedules SET quota_remaining = quota_remaining + 1 WHERE id = ?", scheduleID)
	} else if newStatus == "confirmed" {
		tx.Exec("UPDATE schedules SET quota_remaining = quota_remaining - 1 WHERE id = ?", scheduleID)
	}

	tx.Commit()

	if r.Header.Get("HX-Request") == "true" {
		var total int64
		var status, tripName, leadName string
		h.db.QueryRow("SELECT total, status FROM bookings WHERE id = ?", id).Scan(&total, &status)
		h.db.QueryRow(`SELECT t.name, b.lead_name FROM bookings b JOIN trips t ON t.id = b.trip_id WHERE b.id = ?`, id).Scan(&tripName, &leadName)

		var toastMsg string
		switch newStatus {
		case "confirmed":
			toastMsg = "Booking berhasil dikonfirmasi"
		case "completed":
			toastMsg = "Booking berhasil diselesaikan"
		case "cancelled":
			toastMsg = "Booking berhasil dibatalkan"
		}

		w.Header().Set("HX-Trigger", `{"showToast":"`+toastMsg+`"}`)

		d := bookingDetailData{
			BookingID: id,
			Total:     total,
			Status:    status,
			TripName:  tripName,
			LeadName:  leadName,
		}

		h.renderer.RenderTemplate(w, "bookings/detail", map[string]interface{}{
			"Booking": d,
		})
		return
	}

	http.Redirect(w, r, "/bookings/"+strconv.FormatInt(id, 10)+"?flash=Status booking berhasil diperbarui&flash_type=success", http.StatusSeeOther)
}
