package handler

import (
	"math"
	"net/http"
	"strconv"

	"github.com/ayomendaki/ayomendaki-admin/internal/common"
)

const mpPerPage = 20

type mpPaginationData struct {
	Page       int
	PerPage    int
	Total      int
	TotalPages int
	HasPrev    bool
	HasNext    bool
}

func (h *Handler) MeetingPointList(w http.ResponseWriter, r *http.Request) {
	operatorID := common.GetOperatorID(r)

	search := r.URL.Query().Get("search")
	filterType := r.URL.Query().Get("type")
	pageStr := r.URL.Query().Get("page")

	page := 1
	if p, err := strconv.Atoi(pageStr); err == nil && p > 0 {
		page = p
	}

	// Count totals per type for filter pills
	typeCounts := map[string]int{}
	typeRows, _ := h.db.Query("SELECT type, COUNT(*) FROM meeting_points WHERE operator_id = ? GROUP BY type", operatorID)
	if typeRows != nil {
		defer typeRows.Close()
		for typeRows.Next() {
			var t string
			var c int
			typeRows.Scan(&t, &c)
			typeCounts[t] = c
		}
	}
	totalAll := 0
	for _, c := range typeCounts {
		totalAll += c
	}

	// Count matching rows with filters
	countQuery := "SELECT COUNT(*) FROM meeting_points mp WHERE mp.operator_id = ?"
	countArgs := []interface{}{operatorID}

	if filterType != "" {
		countQuery += " AND mp.type = ?"
		countArgs = append(countArgs, filterType)
	}
	if search != "" {
		countQuery += " AND (mp.name LIKE ? OR mp.address LIKE ?)"
		s := "%" + search + "%"
		countArgs = append(countArgs, s, s)
	}

	var total int
	h.db.QueryRow(countQuery, countArgs...).Scan(&total)

	totalPages := int(math.Ceil(float64(total) / float64(mpPerPage)))
	if totalPages < 1 {
		totalPages = 1
	}
	if page > totalPages {
		page = totalPages
	}

	// Fetch page
	query := `
		SELECT mp.id, mp.type, mp.name, mp.address, mp.lat, mp.lng,
			COALESCE((SELECT COUNT(*) FROM trip_meeting_points WHERE meeting_point_id = mp.id), 0)
		FROM meeting_points mp
		WHERE mp.operator_id = ?
	`
	args := []interface{}{operatorID}

	if filterType != "" {
		query += " AND mp.type = ?"
		args = append(args, filterType)
	}
	if search != "" {
		query += " AND (mp.name LIKE ? OR mp.address LIKE ?)"
		s := "%" + search + "%"
		args = append(args, s, s)
	}
	query += " ORDER BY mp.type, mp.name"
	query += " LIMIT ? OFFSET ?"
	args = append(args, mpPerPage, (page-1)*mpPerPage)

	rows, err := h.db.Query(query, args...)
	type mpItem struct {
		ID         int64
		Type       string
		Name       string
		Address    string
		Lat, Lng   float64
		UsedInTrip bool
	}
	items := []mpItem{}
	if err == nil {
		defer rows.Close()
		for rows.Next() {
			var it mpItem
			var usedInTrip int
			rows.Scan(&it.ID, &it.Type, &it.Name, &it.Address, &it.Lat, &it.Lng, &usedInTrip)
			it.UsedInTrip = usedInTrip > 0
			items = append(items, it)
		}
	}

	pagination := mpPaginationData{
		Page:       page,
		PerPage:    mpPerPage,
		Total:      total,
		TotalPages: totalPages,
		HasPrev:    page > 1,
		HasNext:    page < totalPages,
	}

	h.renderer.RenderTemplate(w, "meeting_points/index", map[string]interface{}{
		"Items":           items,
		"Search":          search,
		"FilterType":      filterType,
		"Pagination":      pagination,
		"TotalAll":        totalAll,
		"CountBasecamp":   typeCounts["basecamp"],
		"CountTitikJemput": typeCounts["titik_jemput"],
	})
}

func (h *Handler) MeetingPointForm(w http.ResponseWriter, r *http.Request) {
	h.renderer.RenderTemplate(w, "meeting_points/form", map[string]interface{}{
		"EditMode": false,
		"MP":       nil,
		"Error":    "",
	})
}

func (h *Handler) MeetingPointCreate(w http.ResponseWriter, r *http.Request) {
	operatorID := common.GetOperatorID(r)
	r.ParseForm()

	mpType := r.FormValue("type")
	name := r.FormValue("name")
	address := r.FormValue("address")
	latStr := r.FormValue("lat")
	lngStr := r.FormValue("lng")

	if name == "" {
		http.Error(w, "Nama meeting point harus diisi", http.StatusBadRequest)
		return
	}
	if mpType != "titik_jemput" && mpType != "basecamp" {
		mpType = "titik_jemput"
	}

	lat, _ := strconv.ParseFloat(latStr, 64)
	lng, _ := strconv.ParseFloat(lngStr, 64)

	_, err := h.db.Exec(
		"INSERT INTO meeting_points (operator_id, type, name, address, lat, lng) VALUES (?, ?, ?, ?, ?, ?)",
		operatorID, mpType, name, address, lat, lng,
	)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, "/meeting-points?flash=Meeting point berhasil dibuat&flash_type=success", http.StatusSeeOther)
}

func (h *Handler) MeetingPointFormEdit(w http.ResponseWriter, r *http.Request) {
	operatorID := common.GetOperatorID(r)
	id, _ := strconv.ParseInt(r.PathValue("id"), 10, 64)

	var mp struct {
		ID      int64
		Type    string
		Name    string
		Address string
		Lat, Lng float64
	}
	err := h.db.QueryRow(`SELECT id, type, name, address, lat, lng FROM meeting_points WHERE id = ? AND operator_id = ?`,
		id, operatorID).Scan(&mp.ID, &mp.Type, &mp.Name, &mp.Address, &mp.Lat, &mp.Lng)

	if err != nil {
		http.Redirect(w, r, "/meeting-points?flash=Meeting point tidak ditemukan&flash_type=error", http.StatusSeeOther)
		return
	}

	h.renderer.RenderTemplate(w, "meeting_points/form", map[string]interface{}{
		"EditMode": true,
		"MP":       mp,
		"Error":    "",
	})
}

func (h *Handler) MeetingPointUpdate(w http.ResponseWriter, r *http.Request) {
	operatorID := common.GetOperatorID(r)
	id, _ := strconv.ParseInt(r.PathValue("id"), 10, 64)

	r.ParseForm()
	mpType := r.FormValue("type")
	name := r.FormValue("name")
	address := r.FormValue("address")
	latStr := r.FormValue("lat")
	lngStr := r.FormValue("lng")

	if mpType != "titik_jemput" && mpType != "basecamp" {
		mpType = "titik_jemput"
	}

	lat, _ := strconv.ParseFloat(latStr, 64)
	lng, _ := strconv.ParseFloat(lngStr, 64)

	_, err := h.db.Exec(
		"UPDATE meeting_points SET type = ?, name = ?, address = ?, lat = ?, lng = ? WHERE id = ? AND operator_id = ?",
		mpType, name, address, lat, lng, id, operatorID,
	)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, "/meeting-points?flash=Meeting point berhasil diperbarui&flash_type=success", http.StatusSeeOther)
}

func (h *Handler) MeetingPointDelete(w http.ResponseWriter, r *http.Request) {
	operatorID := common.GetOperatorID(r)
	id, _ := strconv.ParseInt(r.PathValue("id"), 10, 64)

	var count int
	h.db.QueryRow("SELECT COUNT(*) FROM trip_meeting_points WHERE meeting_point_id = ?", id).Scan(&count)
	if count > 0 {
		w.Header().Set("Content-Type", "text/html")
		w.Write([]byte(`<div class="bg-red-100 border border-red-400 text-red-700 px-4 py-3 rounded">Tidak bisa dihapus — meeting point ini digunakan di trip.</div>`))
		return
	}

	h.db.Exec("DELETE FROM meeting_points WHERE id = ? AND operator_id = ?", id, operatorID)
	http.Redirect(w, r, "/meeting-points?flash=Meeting point berhasil dihapus&flash_type=success", http.StatusSeeOther)
}
