package handler

import (
	"database/sql"
	"net/http"
	"time"

	"github.com/ayomendaki/ayomendaki-admin/internal/auth"
	"github.com/ayomendaki/ayomendaki-admin/internal/common"
)

type loginData struct {
	Error string
}

func (h *Handler) Login(w http.ResponseWriter, r *http.Request) {
	h.renderer.RenderLogin(w, "login", nil)
}

func (h *Handler) LoginPost(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	username := r.FormValue("username")
	password := r.FormValue("password")

	var id int64
	var passwordHash string
	err := h.db.QueryRow("SELECT id, password_hash FROM operators WHERE username = ?", username).Scan(&id, &passwordHash)

	if err == sql.ErrNoRows || !auth.CheckPassword(password, passwordHash) {
		h.renderer.RenderLogin(w, "login", loginData{Error: "Username atau password salah"})
		return
	}
	if err != nil {
		http.Error(w, "Internal error", http.StatusInternalServerError)
		return
	}

	sessionID, err := auth.CreateSession(h.db, id)
	if err != nil {
		http.Error(w, "Internal error", http.StatusInternalServerError)
		return
	}

	http.SetCookie(w, &http.Cookie{
		Name:     common.CookieSessionName,
		Value:    *sessionID,
		Path:     "/",
		HttpOnly: true,
		Expires:  time.Now().Add(24 * time.Hour),
	})

	http.Redirect(w, r, "/", http.StatusSeeOther)
}

func (h *Handler) Logout(w http.ResponseWriter, r *http.Request) {
	cookie, err := r.Cookie(common.CookieSessionName)
	if err == nil {
		auth.DeleteSession(h.db, cookie.Value)
	}

	http.SetCookie(w, &http.Cookie{
		Name:   common.CookieSessionName,
		Value:  "",
		Path:   "/",
		MaxAge: -1,
	})

	http.Redirect(w, r, "/login?flash=Anda telah logout&flash_type=info", http.StatusSeeOther)
}
