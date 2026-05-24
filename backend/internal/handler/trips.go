package handler

import (
	"encoding/json"
	"fmt"
	"html/template"
	"net/http"
	"strconv"

	"github.com/ayomendaki/ayomendaki-admin/internal/common"
)

type tripListItem struct {
	ID        int64
	Name      string
	Mountain  string
	Schedules int
	Bookings  int
}

type simpleOption struct {
	ID   string
	Name string
}

type mpOption struct {
	ID     int64
	Name   string
	MPType string
}

type tripMPDetail struct {
	ID                 int64  `json:"id"`
	Name               string `json:"name"`
	MPType             string `json:"type"`
	OrderIndex         int    `json:"order_index"`
	EstimatedDeparture string `json:"estimated_departure"`
	Prices             map[string]int64 `json:"prices"`
}

type tripFormData struct {
	Trip              *tripDetailData
	Mountains         []simpleOption
	Packages          []simpleOption
	MeetingPoints     []mpOption
	TripMeetingPoints []tripMPDetail
	Error             string
	EditMode          bool
	SelectedPackageIDs []int64
	SelectedPackageIDsJSON string
	MPPricesJSON       string
}

type tripDetailData struct {
	ID              int64
	Name            string
	MountainID      string
	Mountain        string
	Route           string
	Duration        string
	Facilities      []string
	MeetingPoints   []tripMPDetail
	Schedules       []scheduleItem
	Packages        []pkgInfo
	CreatedAt       string
}

type pkgInfo struct {
	ID   int64  `json:"id"`
	Name string `json:"name"`
}

type scheduleItem struct {
	ID             int64
	DateStart      string
	DateEnd        string
	QuotaTotal     int
	QuotaRemaining int
}

type tripSaveData struct {
	Packages      []int64              `json:"packages"`
	MeetingPoints []tripSaveMP         `json:"meetingPoints"`
}

type tripSaveMP struct {
	ID                 int64            `json:"id"`
	EstimatedDeparture string           `json:"estimated_departure"`
	Prices             map[string]int64 `json:"prices"`
}

// ── Handler ──────────────────────────────────────

func (h *Handler) TripList(w http.ResponseWriter, r *http.Request) {
	operatorID := common.GetOperatorID(r)
	rows, err := h.db.Query(`
		SELECT t.id, t.name, COALESCE(m.name, ''),
			(SELECT COUNT(*) FROM schedules s WHERE s.trip_id = t.id),
			(SELECT COUNT(*) FROM bookings b WHERE b.trip_id = t.id)
		FROM trips t
		LEFT JOIN mountains m ON m.id = t.mountain_id
		WHERE t.operator_id = ?
		ORDER BY t.created_at DESC
	`, operatorID)
	trips := []tripListItem{}
	if err == nil {
		defer rows.Close()
		for rows.Next() {
			var t tripListItem
			rows.Scan(&t.ID, &t.Name, &t.Mountain, &t.Schedules, &t.Bookings)
			trips = append(trips, t)
		}
	}
	h.renderer.RenderTemplate(w, "trips/index", map[string]interface{}{"Trips": trips})
}

func (h *Handler) TripForm(w http.ResponseWriter, r *http.Request) {
	operatorID := common.GetOperatorID(r)
	data := h.buildTripFormData(operatorID, 0)
	h.renderer.RenderTemplate(w, "trips/form", data)
}

func (h *Handler) TripFormEdit(w http.ResponseWriter, r *http.Request) {
	operatorID := common.GetOperatorID(r)
	id, _ := strconv.ParseInt(r.PathValue("id"), 10, 64)
	data := h.buildTripFormData(operatorID, id)
	if data.Trip == nil {
		http.Redirect(w, r, "/trips", http.StatusSeeOther)
		return
	}
	data.EditMode = true
	h.renderer.RenderTemplate(w, "trips/form", data)
}

func (h *Handler) TripCreate(w http.ResponseWriter, r *http.Request) {
	operatorID := common.GetOperatorID(r)
	r.ParseForm()
	name := r.FormValue("name")
	mountainID := r.FormValue("mountain_id")
	route := r.FormValue("route")
	duration := r.FormValue("duration")

	if name == "" || mountainID == "" {
		data := h.buildTripFormData(operatorID, 0)
		data.Error = "Nama trip dan gunung harus diisi"
		h.renderer.RenderTemplate(w, "trips/form", data)
		return
	}

	result, err := h.db.Exec(
		"INSERT INTO trips (operator_id, mountain_id, name, route, duration, price) VALUES (?, ?, ?, ?, ?, 0)",
		operatorID, mountainID, name, route, duration,
	)
	if err != nil {
		data := h.buildTripFormData(operatorID, 0)
		data.Error = "Gagal menyimpan trip: " + err.Error()
		h.renderer.RenderTemplate(w, "trips/form", data)
		return
	}
	tripID, _ := result.LastInsertId()
	h.saveTripData(tripID, r.FormValue("trip_json"))
	http.Redirect(w, r, fmt.Sprintf("/trips/%d?flash=Trip berhasil dibuat&flash_type=success", tripID), http.StatusSeeOther)
}

func (h *Handler) TripUpdate(w http.ResponseWriter, r *http.Request) {
	operatorID := common.GetOperatorID(r)
	id, _ := strconv.ParseInt(r.PathValue("id"), 10, 64)
	r.ParseForm()
	name := r.FormValue("name")
	mountainID := r.FormValue("mountain_id")
	route := r.FormValue("route")
	duration := r.FormValue("duration")

	_, err := h.db.Exec(
		"UPDATE trips SET name = ?, mountain_id = ?, route = ?, duration = ? WHERE id = ? AND operator_id = ?",
		name, mountainID, route, duration, id, operatorID,
	)
	if err != nil {
		data := h.buildTripFormData(operatorID, id)
		data.Error = "Gagal menyimpan: " + err.Error()
		h.renderer.RenderTemplate(w, "trips/form", data)
		return
	}

	h.db.Exec("DELETE FROM trip_meeting_points WHERE trip_id = ?", id)
	h.db.Exec("DELETE FROM trip_package_prices WHERE trip_id = ?", id)
	h.saveTripData(id, r.FormValue("trip_json"))
	http.Redirect(w, r, fmt.Sprintf("/trips/%d?flash=Trip berhasil diperbarui&flash_type=success", id), http.StatusSeeOther)
}

func (h *Handler) TripDelete(w http.ResponseWriter, r *http.Request) {
	operatorID := common.GetOperatorID(r)
	id, _ := strconv.ParseInt(r.PathValue("id"), 10, 64)
	var bookingCount int
	h.db.QueryRow("SELECT COUNT(*) FROM bookings WHERE trip_id = ?", id).Scan(&bookingCount)
	if bookingCount > 0 {
		w.Header().Set("Content-Type", "text/html")
		w.Write([]byte(`<div class="bg-red-100 border border-red-400 text-red-700 px-4 py-3 rounded">Tidak bisa dihapus — trip ini memiliki booking aktif.</div>`))
		return
	}
	h.db.Exec("DELETE FROM trips WHERE id = ? AND operator_id = ?", id, operatorID)
	w.Header().Set("HX-Redirect", "/trips")
	w.WriteHeader(http.StatusOK)
}

func (h *Handler) TripDetail(w http.ResponseWriter, r *http.Request) {
	operatorID := common.GetOperatorID(r)
	id, _ := strconv.ParseInt(r.PathValue("id"), 10, 64)

	trip := tripDetailData{ID: id}
	err := h.db.QueryRow(`
		SELECT t.name, t.mountain_id, COALESCE(m.name, ''), t.route, t.duration, t.created_at
		FROM trips t
		LEFT JOIN mountains m ON m.id = t.mountain_id
		WHERE t.id = ? AND t.operator_id = ?
	`, id, operatorID).Scan(&trip.Name, &trip.MountainID, &trip.Mountain, &trip.Route, &trip.Duration, &trip.CreatedAt)
	if err != nil {
		http.Redirect(w, r, "/trips", http.StatusSeeOther)
		return
	}

	// Load meeting points
	mpRows, _ := h.db.Query(`
		SELECT mp.id, mp.name, mp.type, tmp.order_index, tmp.estimated_departure
		FROM meeting_points mp
		JOIN trip_meeting_points tmp ON tmp.meeting_point_id = mp.id
		WHERE tmp.trip_id = ?
		ORDER BY tmp.order_index
	`, id)
	if mpRows != nil {
		defer mpRows.Close()
		for mpRows.Next() {
			var mp tripMPDetail
			mpRows.Scan(&mp.ID, &mp.Name, &mp.MPType, &mp.OrderIndex, &mp.EstimatedDeparture)
			mp.Prices = map[string]int64{}
			trip.MeetingPoints = append(trip.MeetingPoints, mp)
		}
	}

	// Load package prices
	priceRows, _ := h.db.Query(
		"SELECT meeting_point_id, package_id, price FROM trip_package_prices WHERE trip_id = ?", id)
	if priceRows != nil {
		defer priceRows.Close()
		for priceRows.Next() {
			var mpID, pkgID int64
			var price int64
			priceRows.Scan(&mpID, &pkgID, &price)
			for i := range trip.MeetingPoints {
				if trip.MeetingPoints[i].ID == mpID {
					if trip.MeetingPoints[i].Prices == nil {
						trip.MeetingPoints[i].Prices = map[string]int64{}
					}
					trip.MeetingPoints[i].Prices[fmt.Sprintf("%d", pkgID)] = price
				}
			}
		}
	}

	// Load packages
	pkgRows, _ := h.db.Query(`
		SELECT DISTINCT p.id, p.name FROM packages p
		JOIN trip_package_prices tpp ON tpp.package_id = p.id
		WHERE tpp.trip_id = ?
		GROUP BY p.id
	`, id)
	if pkgRows != nil {
		defer pkgRows.Close()
		for pkgRows.Next() {
			var p pkgInfo
			pkgRows.Scan(&p.ID, &p.Name)
			trip.Packages = append(trip.Packages, p)
		}
	}

	// Schedules
	sRows, _ := h.db.Query(`SELECT id, date_start, date_end, quota_total, quota_remaining FROM schedules WHERE trip_id = ? ORDER BY date_start`, id)
	if sRows != nil {
		defer sRows.Close()
		for sRows.Next() {
			var s scheduleItem
			sRows.Scan(&s.ID, &s.DateStart, &s.DateEnd, &s.QuotaTotal, &s.QuotaRemaining)
			trip.Schedules = append(trip.Schedules, s)
		}
	}
	schedulesJSON, _ := json.Marshal(trip.Schedules)

	h.renderer.RenderTemplate(w, "trips/detail", map[string]interface{}{
		"Trip":          trip,
		"SchedulesJSON": template.JS(schedulesJSON),
	})
}

// ── Build trip form data ──────────────────────────

func (h *Handler) buildTripFormData(operatorID int64, tripID int64) tripFormData {
	var data tripFormData

	mRows, _ := h.db.Query("SELECT id, name FROM mountains ORDER BY name")
	if mRows != nil {
		defer mRows.Close()
		for mRows.Next() {
			var o simpleOption
			mRows.Scan(&o.ID, &o.Name)
			data.Mountains = append(data.Mountains, o)
		}
	}

	pRows, _ := h.db.Query("SELECT id, name FROM packages WHERE operator_id = ?", operatorID)
	if pRows != nil {
		defer pRows.Close()
		for pRows.Next() {
			var o simpleOption
			var id int64
			pRows.Scan(&id, &o.Name)
			o.ID = fmt.Sprintf("%d", id)
			data.Packages = append(data.Packages, o)
		}
	}

	mpRows, _ := h.db.Query(`SELECT mp.id, mp.name, mp.type FROM meeting_points mp WHERE mp.operator_id = ? ORDER BY mp.type, mp.name`, operatorID)
	if mpRows != nil {
		defer mpRows.Close()
		for mpRows.Next() {
			var mp mpOption
			mpRows.Scan(&mp.ID, &mp.Name, &mp.MPType)
			data.MeetingPoints = append(data.MeetingPoints, mp)
		}
	}

		if tripID > 0 {
			trip := &tripDetailData{}
			err := h.db.QueryRow(`
				SELECT t.id, t.name, t.mountain_id, COALESCE(m.name, ''), t.route, t.duration, t.created_at
				FROM trips t
				LEFT JOIN mountains m ON m.id = t.mountain_id
				WHERE t.id = ? AND t.operator_id = ?
			`, tripID, operatorID).Scan(&trip.ID, &trip.Name, &trip.MountainID, &trip.Mountain, &trip.Route, &trip.Duration, &trip.CreatedAt)
			if err == nil {
				data.Trip = trip

				// Load meeting points with prices
				selRows, _ := h.db.Query(`
					SELECT mp.id, mp.name, mp.type, tmp.order_index, tmp.estimated_departure
					FROM meeting_points mp
					JOIN trip_meeting_points tmp ON tmp.meeting_point_id = mp.id
					WHERE tmp.trip_id = ?
					ORDER BY tmp.order_index
				`, tripID)
				if selRows != nil {
					defer selRows.Close()
					for selRows.Next() {
						var mp tripMPDetail
						selRows.Scan(&mp.ID, &mp.Name, &mp.MPType, &mp.OrderIndex, &mp.EstimatedDeparture)
						mp.Prices = map[string]int64{}
						data.TripMeetingPoints = append(data.TripMeetingPoints, mp)
					}
				}

				// Load package prices
				selectedPkgMap := map[int64]bool{}
				priceRows, _ := h.db.Query(
					"SELECT meeting_point_id, package_id, price FROM trip_package_prices WHERE trip_id = ?", tripID)
				if priceRows != nil {
					defer priceRows.Close()
					for priceRows.Next() {
						var mpID, pkgID int64
						var price int64
						priceRows.Scan(&mpID, &pkgID, &price)
						for i := range data.TripMeetingPoints {
							if data.TripMeetingPoints[i].ID == mpID {
								data.TripMeetingPoints[i].Prices[fmt.Sprintf("%d", pkgID)] = price
							}
						}
						selectedPkgMap[pkgID] = true
					}
				}

				// Build JSON data for pre-selection in the template
				for id := range selectedPkgMap {
					data.SelectedPackageIDs = append(data.SelectedPackageIDs, id)
				}
				selPkgJSON, _ := json.Marshal(data.SelectedPackageIDs)
				data.SelectedPackageIDsJSON = string(selPkgJSON)

				mpPricesMap := map[string]map[string]int64{}
				for _, mp := range data.TripMeetingPoints {
					mpPricesMap[fmt.Sprintf("%d", mp.ID)] = mp.Prices
				}
				mpPricesJSON, _ := json.Marshal(mpPricesMap)
				data.MPPricesJSON = string(mpPricesJSON)
			}
		}
	return data
}

// ── Save trip data (packages + meeting points + prices) ──

func (h *Handler) saveTripData(tripID int64, jsonStr string) {
	if jsonStr == "" {
		return
	}
	var data tripSaveData
	if err := json.Unmarshal([]byte(jsonStr), &data); err != nil {
		return
	}

	// Save meeting points
	if len(data.MeetingPoints) > 0 {
		for i, mp := range data.MeetingPoints {
			h.db.Exec(`INSERT INTO trip_meeting_points (trip_id, meeting_point_id, order_index, estimated_departure) VALUES (?, ?, ?, ?)`,
				tripID, mp.ID, i+1, mp.EstimatedDeparture)
		}
	}

	// Save package prices
	if len(data.Packages) > 0 && len(data.MeetingPoints) > 0 {
		for _, pkgID := range data.Packages {
			for _, mp := range data.MeetingPoints {
				price := mp.Prices[fmt.Sprintf("%d", pkgID)]
				h.db.Exec(`INSERT INTO trip_package_prices (trip_id, meeting_point_id, package_id, price) VALUES (?, ?, ?, ?)`,
					tripID, mp.ID, pkgID, price)
			}
		}
	}
}
