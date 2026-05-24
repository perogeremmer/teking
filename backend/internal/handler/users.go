package handler

import (
	"net/http"
	"strconv"

	"github.com/ayomendaki/ayomendaki-admin/internal/auth"
	"github.com/ayomendaki/ayomendaki-admin/internal/common"
)

type userItem struct {
	ID       int64
	Username string
	Name     string
	Role     string
}

func (h *Handler) UserList(w http.ResponseWriter, r *http.Request) {
	rows, err := h.db.Query("SELECT id, username, name, role FROM operators ORDER BY id")
	users := []userItem{}
	if err == nil {
		defer rows.Close()
		for rows.Next() {
			var u userItem
			rows.Scan(&u.ID, &u.Username, &u.Name, &u.Role)
			users = append(users, u)
		}
	}
	h.renderer.RenderTemplate(w, "users/index", map[string]interface{}{"Users": users})
}

func (h *Handler) UserForm(w http.ResponseWriter, r *http.Request) {
	h.renderer.RenderTemplate(w, "users/form", map[string]interface{}{
		"EditMode": false,
		"User":     nil,
		"Error":    "",
	})
}

func (h *Handler) UserCreate(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	username := r.FormValue("username")
	password := r.FormValue("password")
	name := r.FormValue("name")
	role := r.FormValue("role")

	if username == "" || password == "" || name == "" {
		h.renderer.RenderTemplate(w, "users/form", map[string]interface{}{
			"EditMode": false,
			"User":     nil,
			"Error":    "Username, password, dan nama harus diisi",
		})
		return
	}
	if role != common.RoleSuperadmin && role != common.RoleAdmin && role != common.RoleUser {
		role = common.RoleUser
	}

	hash, err := auth.HashPassword(password)
	if err != nil {
		http.Error(w, "Internal error", http.StatusInternalServerError)
		return
	}

	_, err = h.db.Exec("INSERT INTO operators (username, password_hash, name, role, description) VALUES (?, ?, ?, ?, ?)",
		username, hash, name, role, "")
	if err != nil {
		h.renderer.RenderTemplate(w, "users/form", map[string]interface{}{
			"EditMode": false,
			"User":     nil,
			"Error":    "Gagal membuat user: " + err.Error(),
		})
		return
	}

	http.Redirect(w, r, "/users?flash=User berhasil dibuat&flash_type=success", http.StatusSeeOther)
}

func (h *Handler) UserFormEdit(w http.ResponseWriter, r *http.Request) {
	id, _ := strconv.ParseInt(r.PathValue("id"), 10, 64)

	var u struct {
		ID       int64
		Username string
		Name     string
		Role     string
	}
	err := h.db.QueryRow("SELECT id, username, name, role FROM operators WHERE id = ?", id).Scan(&u.ID, &u.Username, &u.Name, &u.Role)
	if err != nil {
		http.Redirect(w, r, "/users?flash=User tidak ditemukan&flash_type=error", http.StatusSeeOther)
		return
	}

	h.renderer.RenderTemplate(w, "users/form", map[string]interface{}{
		"EditMode": true,
		"User":     u,
		"Error":    "",
	})
}

func (h *Handler) UserUpdate(w http.ResponseWriter, r *http.Request) {
	id, _ := strconv.ParseInt(r.PathValue("id"), 10, 64)
	r.ParseForm()
	username := r.FormValue("username")
	password := r.FormValue("password")
	name := r.FormValue("name")
	role := r.FormValue("role")

	if username == "" || name == "" {
		http.Redirect(w, r, "/users/"+strconv.FormatInt(id, 10)+"/edit?flash=Username dan nama harus diisi&flash_type=error", http.StatusSeeOther)
		return
	}
	if role != common.RoleSuperadmin && role != common.RoleAdmin && role != common.RoleUser {
		role = common.RoleUser
	}

	if password != "" {
		hash, err := auth.HashPassword(password)
		if err == nil {
			h.db.Exec("UPDATE operators SET password_hash = ? WHERE id = ?", hash, id)
		}
	}

	_, err := h.db.Exec("UPDATE operators SET username = ?, name = ?, role = ? WHERE id = ?",
		username, name, role, id)
	if err != nil {
		http.Redirect(w, r, "/users/"+strconv.FormatInt(id, 10)+"/edit?flash=Gagal menyimpan&flash_type=error", http.StatusSeeOther)
		return
	}

	http.Redirect(w, r, "/users?flash=User berhasil diperbarui&flash_type=success", http.StatusSeeOther)
}

func (h *Handler) UserDelete(w http.ResponseWriter, r *http.Request) {
	operatorID := common.GetOperatorID(r)
	id, _ := strconv.ParseInt(r.PathValue("id"), 10, 64)

	if id == operatorID {
		w.Write([]byte(`<div class="bg-red-100 border border-red-400 text-red-700 px-4 py-3 rounded">Tidak bisa menghapus akun sendiri.</div>`))
		return
	}

	h.db.Exec("DELETE FROM operators WHERE id = ?", id)
	http.Redirect(w, r, "/users?flash=User berhasil dihapus&flash_type=success", http.StatusSeeOther)
}
