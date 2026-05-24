package handler

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/ayomendaki/ayomendaki-admin/internal/common"
	"github.com/ayomendaki/ayomendaki-admin/internal/model"
)

type pkgListItem struct {
	ID          int64
	Name        string
	Description string
	Facilities  []model.Facility
	UsedInTrip  bool
	Include     []string
	Exclude     []string
}

func (h *Handler) PackageList(w http.ResponseWriter, r *http.Request) {
	operatorID := common.GetOperatorID(r)

	rows, err := h.db.Query(`
		SELECT p.id, p.name, p.description, p.facilities,
			COALESCE((SELECT COUNT(*) FROM trips WHERE package_id = p.id), 0)
		FROM packages p
		WHERE p.operator_id = ?
		ORDER BY p.name
	`, operatorID)
	items := []pkgListItem{}
	if err == nil {
		defer rows.Close()
		for rows.Next() {
			var p pkgListItem
			var facJSON string
			var usedInTrip int
			rows.Scan(&p.ID, &p.Name, &p.Description, &facJSON, &usedInTrip)
			p.UsedInTrip = usedInTrip > 0

			json.Unmarshal([]byte(facJSON), &p.Facilities)
			for _, f := range p.Facilities {
				if f.Type == "include" {
					p.Include = append(p.Include, f.Name)
				} else {
					p.Exclude = append(p.Exclude, f.Name)
				}
			}
			items = append(items, p)
		}
	}

	h.renderer.RenderTemplate(w, "packages/index", map[string]interface{}{
		"Items": items,
	})
}

func (h *Handler) PackageForm(w http.ResponseWriter, r *http.Request) {
	h.renderer.RenderTemplate(w, "packages/form", map[string]interface{}{
		"EditMode":   false,
		"Package":    nil,
		"Selected":   []model.Facility{},
		"SelectedJSON": "[]",
		"Error":      "",
	})
}

func (h *Handler) PackageCreate(w http.ResponseWriter, r *http.Request) {
	operatorID := common.GetOperatorID(r)
	r.ParseForm()

	name := r.FormValue("name")
	description := r.FormValue("description")
	facilitiesJSON := r.FormValue("facilities_json")

	if name == "" {
		h.renderer.RenderTemplate(w, "packages/form", map[string]interface{}{
			"EditMode":   false,
			"Package":    nil,
			"Selected":   []model.Facility{},
			"SelectedJSON": "[]",
			"Error":      "Nama paket harus diisi",
		})
		return
	}

	if facilitiesJSON == "" {
		facilitiesJSON = "[]"
	}

	_, err := h.db.Exec("INSERT INTO packages (operator_id, name, description, facilities) VALUES (?, ?, ?, ?)",
		operatorID, name, description, facilitiesJSON)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, "/packages?flash=Paket berhasil dibuat&flash_type=success", http.StatusSeeOther)
}

func (h *Handler) PackageFormEdit(w http.ResponseWriter, r *http.Request) {
	operatorID := common.GetOperatorID(r)
	id, _ := strconv.ParseInt(r.PathValue("id"), 10, 64)

	var pkgName, pkgDesc, facJSON string
	err := h.db.QueryRow("SELECT name, description, facilities FROM packages WHERE id = ? AND operator_id = ?",
		id, operatorID).Scan(&pkgName, &pkgDesc, &facJSON)
	if err != nil {
		http.Redirect(w, r, "/packages?flash=Paket tidak ditemukan&flash_type=error", http.StatusSeeOther)
		return
	}

	var selected []model.Facility
	json.Unmarshal([]byte(facJSON), &selected)
	selJSON, _ := json.Marshal(selected)

	h.renderer.RenderTemplate(w, "packages/form", map[string]interface{}{
		"EditMode": true,
		"Package": map[string]interface{}{
			"ID":          id,
			"Name":        pkgName,
			"Description": pkgDesc,
		},
		"Selected":     selected,
		"SelectedJSON": string(selJSON),
		"Error":        "",
	})
}

func (h *Handler) PackageUpdate(w http.ResponseWriter, r *http.Request) {
	operatorID := common.GetOperatorID(r)
	id, _ := strconv.ParseInt(r.PathValue("id"), 10, 64)
	r.ParseForm()

	name := r.FormValue("name")
	description := r.FormValue("description")
	facilitiesJSON := r.FormValue("facilities_json")

	if facilitiesJSON == "" {
		facilitiesJSON = "[]"
	}

	_, err := h.db.Exec("UPDATE packages SET name = ?, description = ?, facilities = ? WHERE id = ? AND operator_id = ?",
		name, description, facilitiesJSON, id, operatorID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, "/packages?flash=Paket berhasil diperbarui&flash_type=success", http.StatusSeeOther)
}

func (h *Handler) PackageDelete(w http.ResponseWriter, r *http.Request) {
	operatorID := common.GetOperatorID(r)
	id, _ := strconv.ParseInt(r.PathValue("id"), 10, 64)

	var count int
	h.db.QueryRow("SELECT COUNT(*) FROM trips WHERE package_id = ?", id).Scan(&count)
	if count > 0 {
		w.Header().Set("Content-Type", "text/html")
		w.Write([]byte(`<div class="bg-red-100 border border-red-400 text-red-700 px-4 py-3 rounded">Tidak bisa dihapus — paket ini digunakan di trip.</div>`))
		return
	}

	h.db.Exec("DELETE FROM packages WHERE id = ? AND operator_id = ?", id, operatorID)
	http.Redirect(w, r, "/packages?flash=Paket berhasil dihapus&flash_type=success", http.StatusSeeOther)
}
