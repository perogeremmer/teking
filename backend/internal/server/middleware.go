package server

import (
	"context"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/ayomendaki/ayomendaki-admin/internal/auth"
	"github.com/ayomendaki/ayomendaki-admin/internal/common"
	"github.com/ayomendaki/ayomendaki-admin/internal/database"
)

func (s *Server) loggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		next.ServeHTTP(w, r)
		log.Printf("%s %s %s", r.Method, r.URL.Path, time.Since(start))
	})
}

func (s *Server) authMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		path := r.URL.Path

		// Allow static files without auth only for /static path
		if len(path) > 8 && path[:8] == "/static/" {
			next.ServeHTTP(w, r)
			return
		}

		// Login/logout don't need auth
		if path == "/login" || path == "/logout" {
			next.ServeHTTP(w, r)
			return
		}

		// Validate session
		cookie, err := r.Cookie(common.CookieSessionName)
		if err != nil {
			http.Redirect(w, r, "/login", http.StatusSeeOther)
			return
		}

		operatorID, err := auth.ValidateSession(database.DB, cookie.Value)
		if err != nil || operatorID == nil {
			http.SetCookie(w, &http.Cookie{
				Name:   common.CookieSessionName,
				Value:  "",
				Path:   "/",
				MaxAge: -1,
			})
			http.Redirect(w, r, "/login", http.StatusSeeOther)
			return
		}

		// Get role from DB
		var role string
		database.DB.QueryRow("SELECT role FROM operators WHERE id = ?", *operatorID).Scan(&role)
		if role == "" {
			role = common.RoleAdmin
		}

		// Role-based access check
		allow := false
		switch role {
		case common.RoleSuperadmin:
			allow = true
		case common.RoleAdmin:
			allow = !strings.HasPrefix(path, "/users")
		case common.RoleUser:
			allow = path == "/profile" ||
				strings.HasPrefix(path, "/reports") ||
				strings.HasPrefix(path, "/static/")
		}
		if !allow {
			http.Error(w, "Forbidden", http.StatusForbidden)
			return
		}

		// Set context and proceed
		ctx := r.Context()
		ctx = context.WithValue(ctx, common.ContextKeyOperatorID, *operatorID)
		ctx = context.WithValue(ctx, common.ContextKeyRole, role)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
