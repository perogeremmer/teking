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
		if r.URL.Path == "/login" || r.URL.Path == "/logout" {
			next.ServeHTTP(w, r)
			return
		}

		if len(r.URL.Path) > 8 && r.URL.Path[:8] == "/static/" {
			next.ServeHTTP(w, r)
			return
		}

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

		var role string
		database.DB.QueryRow("SELECT role FROM operators WHERE id = ?", *operatorID).Scan(&role)
		if role == "" {
			role = common.RoleAdmin
		}

		ctx := r.Context()
		ctx = context.WithValue(ctx, common.ContextKeyOperatorID, *operatorID)
		ctx = context.WithValue(ctx, common.ContextKeyRole, role)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func (s *Server) roleMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		role := common.GetRole(r)
		path := r.URL.Path

		// Superadmin can access everything
		if role == common.RoleSuperadmin {
			next.ServeHTTP(w, r)
			return
		}

		if role == common.RoleAdmin {
			// Block user management routes
			if strings.HasPrefix(path, "/users") {
				http.Error(w, "Forbidden", http.StatusForbidden)
				return
			}
			next.ServeHTTP(w, r)
			return
		}

		// User role — only reports, profile, and login
		if role == common.RoleUser {
			allowed := path == "/login" || path == "/logout" ||
				path == "/profile" || strings.HasPrefix(path, "/reports") ||
				strings.HasPrefix(path, "/static/")
			if !allowed {
				http.Error(w, "Forbidden", http.StatusForbidden)
				return
			}
			next.ServeHTTP(w, r)
			return
		}

		// Unknown role — block
		http.Error(w, "Forbidden", http.StatusForbidden)
	})
}
