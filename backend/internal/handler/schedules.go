package handler

import (
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/ayomendaki/ayomendaki-admin/internal/common"
)

func (h *Handler) ScheduleForm(w http.ResponseWriter, r *http.Request) {
	operatorID := common.GetOperatorID(r)
	tripID, _ := strconv.ParseInt(r.PathValue("id"), 10, 64)

	var tripName string
	err := h.db.QueryRow("SELECT name FROM trips WHERE id = ? AND operator_id = ?", tripID, operatorID).Scan(&tripName)
	if err != nil {
		http.Redirect(w, r, "/trips", http.StatusSeeOther)
		return
	}

	h.renderer.RenderTemplate(w, "schedules/form", map[string]interface{}{
		"TripID":   tripID,
		"TripName": tripName,
		"EditMode": false,
		"Schedule": nil,
	})
}

func (h *Handler) ScheduleCreate(w http.ResponseWriter, r *http.Request) {
	operatorID := common.GetOperatorID(r)
	tripID, _ := strconv.ParseInt(r.PathValue("tripID"), 10, 64)

	var tripOperatorID int64
	h.db.QueryRow("SELECT operator_id FROM trips WHERE id = ?", tripID).Scan(&tripOperatorID)
	if tripOperatorID != operatorID {
		http.Error(w, "Forbidden", http.StatusForbidden)
		return
	}

	r.ParseForm()
	dateStart := r.FormValue("date_start")
	dateEnd := r.FormValue("date_end")
	quotaStr := r.FormValue("quota_total")

	quota, _ := strconv.Atoi(quotaStr)
	if quota < 1 {
		quota = 1
	}

	_, err := h.db.Exec(
		"INSERT INTO schedules (trip_id, date_start, date_end, quota_total, quota_remaining) VALUES (?, ?, ?, ?, ?)",
		tripID, dateStart, dateEnd, quota, quota,
	)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	redirectURL := fmt.Sprintf("/trips/%d?flash=Jadwal berhasil ditambahkan&flash_type=success", tripID)
	if r.Header.Get("HX-Request") == "true" {
		w.Header().Set("HX-Redirect", redirectURL)
		w.WriteHeader(http.StatusOK)
	} else {
		http.Redirect(w, r, redirectURL, http.StatusSeeOther)
	}
}

func (h *Handler) ScheduleFormEdit(w http.ResponseWriter, r *http.Request) {
	operatorID := common.GetOperatorID(r)
	scheduleID, _ := strconv.ParseInt(r.PathValue("id"), 10, 64)

	var sID, tripID int64
	var dateStart, dateEnd string
	var quotaTotal, quotaRemaining int

	err := h.db.QueryRow(`
		SELECT s.id, s.trip_id, s.date_start, s.date_end, s.quota_total, s.quota_remaining
		FROM schedules s
		JOIN trips t ON t.id = s.trip_id
		WHERE s.id = ? AND t.operator_id = ?
	`, scheduleID, operatorID).Scan(&sID, &tripID, &dateStart, &dateEnd, &quotaTotal, &quotaRemaining)

	if err != nil {
		http.Redirect(w, r, "/trips", http.StatusSeeOther)
		return
	}

	var tripName string
	h.db.QueryRow("SELECT name FROM trips WHERE id = ?", tripID).Scan(&tripName)

	h.renderer.RenderTemplate(w, "schedules/form", map[string]interface{}{
		"TripID":   tripID,
		"TripName": tripName,
		"EditMode": true,
		"Schedule": map[string]interface{}{
			"ID":              sID,
			"DateStart":       dateStart,
			"DateEnd":         dateEnd,
			"QuotaTotal":      quotaTotal,
			"QuotaRemaining":  quotaRemaining,
		},
	})
}

func (h *Handler) ScheduleUpdate(w http.ResponseWriter, r *http.Request) {
	operatorID := common.GetOperatorID(r)
	scheduleID, _ := strconv.ParseInt(r.PathValue("id"), 10, 64)

	r.ParseForm()

	var tripID int64
	err := h.db.QueryRow(`
		SELECT s.trip_id FROM schedules s
		JOIN trips t ON t.id = s.trip_id
		WHERE s.id = ? AND t.operator_id = ?
	`, scheduleID, operatorID).Scan(&tripID)
	if err != nil {
		http.Error(w, "Not found", http.StatusNotFound)
		return
	}

	dateStart := r.FormValue("date_start")
	dateEnd := r.FormValue("date_end")
	quotaStr := r.FormValue("quota_total")
	quotaRemainingStr := r.FormValue("quota_remaining")

	quota, _ := strconv.Atoi(quotaStr)
	quotaRemaining, _ := strconv.Atoi(quotaRemainingStr)

	if quota < 1 {
		quota = 1
	}
	if quotaRemaining > quota {
		quotaRemaining = quota
	}
	if quotaRemaining < 0 {
		quotaRemaining = 0
	}

	h.db.Exec(
		"UPDATE schedules SET date_start = ?, date_end = ?, quota_total = ?, quota_remaining = ? WHERE id = ?",
		dateStart, dateEnd, quota, quotaRemaining, scheduleID,
	)

	redirectURL := fmt.Sprintf("/trips/%d?flash=Jadwal berhasil diperbarui&flash_type=success", tripID)
	if r.Header.Get("HX-Request") == "true" {
		w.Header().Set("HX-Redirect", redirectURL)
		w.WriteHeader(http.StatusOK)
	} else {
		http.Redirect(w, r, redirectURL, http.StatusSeeOther)
	}
}

func (h *Handler) ScheduleDelete(w http.ResponseWriter, r *http.Request) {
	operatorID := common.GetOperatorID(r)
	scheduleID, _ := strconv.ParseInt(r.PathValue("id"), 10, 64)

	var bookingCount int
	h.db.QueryRow("SELECT COUNT(*) FROM bookings WHERE schedule_id = ?", scheduleID).Scan(&bookingCount)
	if bookingCount > 0 {
		w.Header().Set("Content-Type", "text/html")
		w.Write([]byte(`<div class="bg-red-100 border border-red-400 text-red-700 px-4 py-3 rounded">Tidak bisa dihapus — ada booking terkait.</div>`))
		return
	}

	var tripID int64
	h.db.QueryRow(`
		SELECT s.trip_id FROM schedules s
		JOIN trips t ON t.id = s.trip_id
		WHERE s.id = ? AND t.operator_id = ?
	`, scheduleID, operatorID).Scan(&tripID)

	h.db.Exec("DELETE FROM schedules WHERE id = ?", scheduleID)

	redirectURL := fmt.Sprintf("/trips/%d?flash=Jadwal berhasil dihapus&flash_type=success", tripID)
	if r.Header.Get("HX-Request") == "true" {
		w.Header().Set("HX-Redirect", redirectURL)
		w.WriteHeader(http.StatusOK)
	} else {
		http.Redirect(w, r, redirectURL, http.StatusSeeOther)
	}
}

func init() {
	// Ensure timezone is loaded
	time.Local = time.UTC
}
