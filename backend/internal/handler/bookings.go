package handler

import (
	"database/sql"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/ayomendaki/ayomendaki-admin/internal/common"
	"github.com/ayomendaki/ayomendaki-admin/internal/model"
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
	PaymentStatus string
}

type bookingDetailData struct {
	BookingID     int64
	TripName      string
	ScheduleDate  string
	LeadName      string
	LeadPhone     string
	LeadEmail     string
	Total         int64
	PaidAmount    int64
	Status        string
	PaymentStatus string
	CreatedAt     string
	Participants  []struct{ Name, KTP string }
	Payments      []model.Payment
}

func (h *Handler) BookingList(w http.ResponseWriter, r *http.Request) {
	operatorID := common.GetOperatorID(r)

	statusFilter := r.URL.Query().Get("status")
	tripFilter := r.URL.Query().Get("trip_id")
	search := r.URL.Query().Get("search")

	query := `
		SELECT b.id, b.lead_name, b.lead_phone, COALESCE(t.name, ''), b.created_at,
			(SELECT COUNT(*) FROM booking_participants bp WHERE bp.booking_id = b.id),
			b.total, b.status, b.payment_status
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
			rows.Scan(&b.ID, &b.Customer, &b.Phone, &b.TripName, &b.TripDate, &b.Participants, &b.Total, &b.Status, &b.PaymentStatus)
			bookings = append(bookings, b)
		}
	}

	tripRows, _ := h.db.Query("SELECT id, name FROM trips WHERE operator_id = ? ORDER BY name", operatorID)
	type tripOpt struct{ ID, Name string }
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

	h.renderer.RenderTemplate(w, "bookings/index", map[string]interface{}{
		"Bookings": bookings, "Trips": trips,
		"StatusFilter": statusFilter, "Search": search,
	})
}

func (h *Handler) BookingDetail(w http.ResponseWriter, r *http.Request) {
	operatorID := common.GetOperatorID(r)
	id, _ := strconv.ParseInt(r.PathValue("id"), 10, 64)

	d := bookingDetailData{BookingID: id}
	err := h.db.QueryRow(`
		SELECT b.lead_name, b.lead_phone, b.lead_email, b.total, b.status, b.payment_status, b.created_at,
			COALESCE(t.name, ''), COALESCE(s.date_start || ' - ' || s.date_end, '')
		FROM bookings b
		JOIN trips t ON t.id = b.trip_id
		LEFT JOIN schedules s ON s.id = b.schedule_id
		WHERE b.id = ? AND t.operator_id = ?
	`, id, operatorID).Scan(&d.LeadName, &d.LeadPhone, &d.LeadEmail, &d.Total, &d.Status, &d.PaymentStatus, &d.CreatedAt, &d.TripName, &d.ScheduleDate)
	if err != nil {
		http.Redirect(w, r, "/bookings?flash=Booking tidak ditemukan&flash_type=error", http.StatusSeeOther)
		return
	}

	// Load participants
	pRows, _ := h.db.Query("SELECT name, ktp FROM booking_participants WHERE booking_id = ?", id)
	if pRows != nil {
		defer pRows.Close()
		for pRows.Next() {
			var p struct{ Name, KTP string }
			pRows.Scan(&p.Name, &p.KTP)
			d.Participants = append(d.Participants, p)
		}
	}

	// Load payments
	payRows, _ := h.db.Query("SELECT id, amount, notes, proof_file, created_at FROM payments WHERE booking_id = ? ORDER BY created_at", id)
	if payRows != nil {
		defer payRows.Close()
		for payRows.Next() {
			var p model.Payment
			payRows.Scan(&p.ID, &p.Amount, &p.Notes, &p.ProofFile, &p.CreatedAt)
			d.Payments = append(d.Payments, p)
		}
	}

	// Calculate paid amount
	d.PaidAmount = 0
	for _, p := range d.Payments {
		d.PaidAmount += p.Amount
	}

	h.renderer.RenderTemplate(w, "bookings/detail", map[string]interface{}{"Booking": d})
}

func (h *Handler) BookingStatus(w http.ResponseWriter, r *http.Request) {
	operatorID := common.GetOperatorID(r)
	id, _ := strconv.ParseInt(r.PathValue("id"), 10, 64)

	r.ParseForm()
	newStatus := r.FormValue("status")

	validStatuses := map[string]bool{"confirmed": true, "completed": true, "cancelled": true}
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

	// Calculate participants count for quota
	var participants int
	h.db.QueryRow("SELECT COUNT(*) FROM booking_participants WHERE booking_id = ?", id).Scan(&participants)
	if participants < 1 {
		participants = 1
	}

	if newStatus == "cancelled" {
		// Restore quota but keep booking row
		h.db.Exec("UPDATE schedules SET quota_remaining = quota_remaining + ? WHERE id = ?", participants, scheduleID)
	} else if newStatus == "confirmed" {
		h.db.Exec("UPDATE schedules SET quota_remaining = quota_remaining - ? WHERE id = ?", participants, scheduleID)
	}

	h.db.Exec("UPDATE bookings SET status = ? WHERE id = ?", newStatus, id)

	if r.Header.Get("HX-Request") == "true" {
		w.Header().Set("HX-Redirect", fmt.Sprintf("/bookings/%d?flash=Status berhasil diperbarui&flash_type=success", id))
		w.WriteHeader(http.StatusOK)
		return
	}
	http.Redirect(w, r, "/bookings/"+strconv.FormatInt(id, 10)+"?flash=Status booking berhasil diperbarui&flash_type=success", http.StatusSeeOther)
}

// ── Manual Booking Creation ──

func (h *Handler) BookingManualCreate(w http.ResponseWriter, r *http.Request) {
	operatorID := common.GetOperatorID(r)
	tripID, _ := strconv.ParseInt(r.PathValue("tripID"), 10, 64)

	// Verify trip belongs to operator
	var tripOperatorID int64
	h.db.QueryRow("SELECT operator_id FROM trips WHERE id = ?", tripID).Scan(&tripOperatorID)
	if tripOperatorID != operatorID {
		http.Error(w, "Forbidden", http.StatusForbidden)
		return
	}

	r.ParseForm()
	custName := strings.TrimSpace(r.FormValue("name"))
	custPhone := strings.TrimSpace(r.FormValue("phone"))
	custNIK := r.FormValue("nik")
	custEmail := r.FormValue("email")
	scheduleIDStr := r.FormValue("schedule_id")
	mpIDStr := r.FormValue("meeting_point_id")
	packageIDStr := r.FormValue("package_id")
	priceStr := r.FormValue("price")
	participantNames := r.Form["participant_name"]
	participantKTPs := r.Form["participant_ktp"]

	if custName == "" || custPhone == "" || scheduleIDStr == "" {
		http.Error(w, "Nama, telepon, dan jadwal harus diisi", http.StatusBadRequest)
		return
	}

	scheduleID, _ := strconv.ParseInt(scheduleIDStr, 10, 64)
	mpID, _ := strconv.ParseInt(mpIDStr, 10, 64)
	packageID, _ := strconv.ParseInt(packageIDStr, 10, 64)
	total, _ := strconv.ParseInt(priceStr, 10, 64)

	// Find or create customer
	var custID int64
	err := h.db.QueryRow("SELECT id FROM customers WHERE phone = ?", custPhone).Scan(&custID)
	if err == sql.ErrNoRows {
		res, e := h.db.Exec("INSERT INTO customers (name, phone, nik, email) VALUES (?, ?, ?, ?)",
			custName, custPhone, custNIK, custEmail)
		if e == nil {
			custID, _ = res.LastInsertId()
		}
	} else if err == nil && custID > 0 {
		h.db.Exec("UPDATE customers SET name = ?, nik = ?, email = ? WHERE id = ?", custName, custNIK, custEmail, custID)
	}

	// Create booking
	res, err := h.db.Exec(`INSERT INTO bookings (trip_id, schedule_id, lead_name, lead_phone, lead_email, total,
		status, payment_status, customer_id, meeting_point_id, package_id) VALUES (?, ?, ?, ?, ?, ?, 'confirmed', 'unpaid', ?, ?, ?)`,
		tripID, scheduleID, custName, custPhone, custEmail, total, custID, mpID, packageID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	bookingID, _ := res.LastInsertId()

	// Add participants
	participantCount := 0
	for i, name := range participantNames {
		name = strings.TrimSpace(name)
		if name == "" {
			continue
		}
		ktp := ""
		if i < len(participantKTPs) {
			ktp = participantKTPs[i]
		}
		h.db.Exec("INSERT INTO booking_participants (booking_id, name, ktp) VALUES (?, ?, ?)", bookingID, name, ktp)
		participantCount++
	}
	if participantCount == 0 {
		// At least the lead participant
		h.db.Exec("INSERT INTO booking_participants (booking_id, name) VALUES (?, ?)", bookingID, custName)
		participantCount = 1
	}

	// Reduce quota
	h.db.Exec("UPDATE schedules SET quota_remaining = quota_remaining - ? WHERE id = ?", participantCount, scheduleID)

	http.Redirect(w, r, fmt.Sprintf("/bookings/%d?flash=Booking berhasil dibuat&flash_type=success", bookingID), http.StatusSeeOther)
}

// ── Payment Creation ──

func (h *Handler) PaymentCreate(w http.ResponseWriter, r *http.Request) {
	operatorID := common.GetOperatorID(r)
	bookingID, _ := strconv.ParseInt(r.PathValue("id"), 10, 64)

	// Verify booking belongs to operator's trip
	var exists int
	h.db.QueryRow(`SELECT COUNT(*) FROM bookings b JOIN trips t ON t.id = b.trip_id WHERE b.id = ? AND t.operator_id = ?`,
		bookingID, operatorID).Scan(&exists)
	if exists == 0 {
		http.Error(w, "Not found", http.StatusNotFound)
		return
	}

	r.ParseMultipartForm(10 << 20)
	amountStr := r.FormValue("amount")
	notes := r.FormValue("notes")
	amount, _ := strconv.ParseInt(amountStr, 10, 64)

	if amount <= 0 {
		http.Error(w, "Jumlah harus > 0", http.StatusBadRequest)
		return
	}

	// Handle file upload
	var proofFile string
	file, header, err := r.FormFile("proof")
	if err == nil {
		defer file.Close()
		dir := "data/proofs"
		os.MkdirAll(dir, 0755)
		filename := fmt.Sprintf("%d_%d_%s", bookingID, time.Now().Unix(), header.Filename)
		dst, err := os.Create(filepath.Join(dir, filename))
		if err == nil {
			defer dst.Close()
			io.Copy(dst, file)
			proofFile = filename
		}
	}

	// Insert payment
	_, err = h.db.Exec("INSERT INTO payments (booking_id, amount, notes, proof_file) VALUES (?, ?, ?, ?)",
		bookingID, amount, notes, proofFile)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Update payment status
	var total int64
	var totalPaid int64
	h.db.QueryRow("SELECT total FROM bookings WHERE id = ?", bookingID).Scan(&total)
	payRows, _ := h.db.Query("SELECT COALESCE(SUM(amount),0) FROM payments WHERE booking_id = ?", bookingID)
	if payRows != nil {
		defer payRows.Close()
		payRows.Next()
		payRows.Scan(&totalPaid)
	}

	newPaymentStatus := "dp_paid"
	if totalPaid >= total {
		newPaymentStatus = "full_paid"
	}
	h.db.Exec("UPDATE bookings SET payment_status = ? WHERE id = ?", newPaymentStatus, bookingID)

	http.Redirect(w, r, fmt.Sprintf("/bookings/%d?flash=Pembayaran berhasil dicatat&flash_type=success", bookingID), http.StatusSeeOther)
}

// ── Serve Payment Proof ──

func (h *Handler) PaymentProof(w http.ResponseWriter, r *http.Request) {
	filename := r.PathValue("file")
	path := filepath.Join("data/proofs", filepath.Base(filename))
	http.ServeFile(w, r, path)
}
