package handler

import (
	"net/http"

	"github.com/ayomendaki/ayomendaki-admin/internal/auth"
	"github.com/ayomendaki/ayomendaki-admin/internal/common"
)

type profileData struct {
	Name        string
	Description string
	Phone       string
	Whatsapp    string
	Instagram   string
	Rating      float64
	Trips       int
	Bookings    int
	Verified    bool
	Error       string
}

func (h *Handler) Profile(w http.ResponseWriter, r *http.Request) {
	operatorID := common.GetOperatorID(r)

	data := profileData{}
	err := h.db.QueryRow(`
		SELECT name, description, phone, whatsapp, instagram, rating, verified FROM operators WHERE id = ?
	`, operatorID).Scan(&data.Name, &data.Description, &data.Phone, &data.Whatsapp, &data.Instagram, &data.Rating, &data.Verified)

	if err != nil {
		http.Error(w, "Operator not found", http.StatusNotFound)
		return
	}

	h.db.QueryRow("SELECT COUNT(*) FROM trips WHERE operator_id = ?", operatorID).Scan(&data.Trips)
	h.db.QueryRow("SELECT COUNT(*) FROM bookings b JOIN trips t ON b.trip_id = t.id WHERE t.operator_id = ?", operatorID).Scan(&data.Bookings)

	h.renderer.RenderTemplate(w, "profile/index", data)
}

func (h *Handler) ProfileUpdate(w http.ResponseWriter, r *http.Request) {
	operatorID := common.GetOperatorID(r)
	r.ParseForm()

	currentPassword := r.FormValue("current_password")
	newPassword := r.FormValue("new_password")

	if currentPassword != "" && newPassword != "" {
		var hash string
		h.db.QueryRow("SELECT password_hash FROM operators WHERE id = ?", operatorID).Scan(&hash)
		if !auth.CheckPassword(currentPassword, hash) {
			var data profileData
			h.db.QueryRow(`SELECT name, description, phone, whatsapp, instagram, rating, verified FROM operators WHERE id = ?`, operatorID).
				Scan(&data.Name, &data.Description, &data.Phone, &data.Whatsapp, &data.Instagram, &data.Rating, &data.Verified)
			h.db.QueryRow("SELECT COUNT(*) FROM trips WHERE operator_id = ?", operatorID).Scan(&data.Trips)
			h.db.QueryRow("SELECT COUNT(*) FROM bookings b JOIN trips t ON b.trip_id = t.id WHERE t.operator_id = ?", operatorID).Scan(&data.Bookings)
			data.Error = "Password saat ini salah"
			h.renderer.RenderTemplate(w, "profile/index", data)
			return
		}
		newHash, err := auth.HashPassword(newPassword)
		if err == nil {
			h.db.Exec("UPDATE operators SET password_hash = ? WHERE id = ?", newHash, operatorID)
		}
		http.Redirect(w, r, "/profile?flash=Password berhasil diganti&flash_type=success", http.StatusSeeOther)
		return
	}

	name := r.FormValue("name")
	if name == "" {
		http.Redirect(w, r, "/profile", http.StatusSeeOther)
		return
	}

	description := r.FormValue("description")
	phone := r.FormValue("phone")
	whatsapp := r.FormValue("whatsapp")
	instagram := r.FormValue("instagram")

	_, err := h.db.Exec(`
		UPDATE operators SET name = ?, description = ?, phone = ?, whatsapp = ?, instagram = ?
		WHERE id = ?
	`, name, description, phone, whatsapp, instagram, operatorID)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, "/profile?flash=Profil berhasil diperbarui&flash_type=success", http.StatusSeeOther)
}
