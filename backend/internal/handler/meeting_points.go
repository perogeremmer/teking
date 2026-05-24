package handler

import (
	"net/http"
	"strconv"

	"github.com/ayomendaki/ayomendaki-admin/internal/common"
)

func (h *Handler) MeetingPointList(w http.ResponseWriter, r *http.Request) {
	operatorID := common.GetOperatorID(r)

	search := r.URL.Query().Get("search")

	query := `
		SELECT mp.id, mp.type, mp.name, mp.address, mp.lat, mp.lng,
			COALESCE((SELECT COUNT(*) FROM trip_meeting_points WHERE meeting_point_id = mp.id), 0)
		FROM meeting_points mp
		WHERE mp.operator_id = ?
	`
	args := []interface{}{operatorID}

	if search != "" {
		query += " AND (mp.name LIKE ? OR mp.address LIKE ?)"
		s := "%" + search + "%"
		args = append(args, s, s)
	}
	query += " ORDER BY mp.name"

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

	h.renderer.RenderTemplate(w, "meeting_points/index", map[string]interface{}{
		"Items":  items,
		"Search": search,
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
