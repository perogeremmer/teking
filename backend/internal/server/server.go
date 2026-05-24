package server

import (
	"embed"
	"fmt"
	"html/template"
	"io/fs"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/ayomendaki/ayomendaki-admin/internal/database"
	"github.com/ayomendaki/ayomendaki-admin/internal/handler"
)

type Server struct {
	http.Handler
	templates map[string]*template.Template
	assets    embed.FS
	dev       bool
}

func New(assets embed.FS, dev bool) *Server {
	s := &Server{
		templates: make(map[string]*template.Template),
		assets:    assets,
		dev:       dev,
	}

	if !dev {
		s.loadTemplates()
	}
	s.setupRoutes()

	return s
}

func (s *Server) templateFuncs() template.FuncMap {
	return template.FuncMap{
		"formatCurrency": func(n interface{}) string {
			var val int64
			switch v := n.(type) {
			case int64:
				val = v
			case int:
				val = int64(v)
			default:
				return "Rp 0"
			}
			if val == 0 {
				return "Gratis"
			}
			sign := ""
			if val < 0 {
				sign = "-"
				val = -val
			}
			ns := fmt.Sprintf("%d", val)
			var parts []string
			for i := len(ns); i > 0; i -= 3 {
				start := i - 3
				if start < 0 {
					start = 0
				}
				parts = append([]string{ns[start:i]}, parts...)
			}
			return sign + "Rp " + strings.Join(parts, ".")
		},
		"formatDate": func(date string) string {
			if len(date) >= 10 {
				parts := strings.Split(date[:10], "-")
				if len(parts) == 3 {
					months := []string{"", "Jan", "Feb", "Mar", "Apr", "Mei", "Jun", "Jul", "Agu", "Sep", "Okt", "Nov", "Des"}
					var monthNum int
					fmt.Sscanf(parts[1], "%d", &monthNum)
					if monthNum >= 1 && monthNum <= 12 {
						return parts[2] + " " + months[monthNum] + " " + parts[0]
					}
				}
			}
			return date
		},
		"formatDateTime": func(date string) string {
			if len(date) >= 16 {
				dp := date[:10]
				tp := date[11:16]
				parts := strings.Split(dp, "-")
				if len(parts) == 3 {
					months := []string{"", "Jan", "Feb", "Mar", "Apr", "Mei", "Jun", "Jul", "Agu", "Sep", "Okt", "Nov", "Des"}
					var monthNum int
					fmt.Sscanf(parts[1], "%d", &monthNum)
					if monthNum >= 1 && monthNum <= 12 {
						return parts[2] + " " + months[monthNum] + " " + parts[0] + " " + tp
					}
				}
			}
			return date
		},
		"statusBadge": func(status string) template.HTML {
			color := "bg-gray-100 text-gray-800"
			switch status {
			case "pending":
				color = "bg-yellow-100 text-yellow-800"
			case "confirmed":
				color = "bg-blue-100 text-blue-800"
			case "completed":
				color = "bg-green-100 text-green-800"
			case "cancelled":
				color = "bg-red-100 text-red-800"
			}
			return template.HTML(fmt.Sprintf(`<span class="inline-flex items-center px-2.5 py-0.5 rounded-full text-xs font-medium %s">%s</span>`, color, status))
		},
		"safeHTML": func(s string) template.HTML {
			return template.HTML(s)
		},
		"seq": func(start, end int) []int {
			var s []int
			for i := start; i <= end; i++ {
				s = append(s, i)
			}
			return s
		},
		"sub": func(a, b int64) int64 {
			return a - b
		},
		"add": func(a, b int) int {
			return a + b
		},
	}
}

func (s *Server) loadTemplates() {
	funcMap := s.templateFuncs()
	layout := "web/templates/layout.html"

	pages := map[string][]string{
		"dashboard":            {layout, "web/templates/dashboard.html"},
		"trips/index":          {layout, "web/templates/trips/index.html"},
		"trips/form":           {layout, "web/templates/trips/form.html"},
		"trips/detail":         {layout, "web/templates/trips/detail.html"},
		"bookings/index":       {layout, "web/templates/bookings/index.html"},
		"bookings/detail":      {layout, "web/templates/bookings/detail.html"},
		"meeting_points/index": {layout, "web/templates/meeting_points/index.html"},
		"meeting_points/form":  {layout, "web/templates/meeting_points/form.html"},
		"packages/index":       {layout, "web/templates/packages/index.html"},
		"packages/form":        {layout, "web/templates/packages/form.html"},
		"reports/revenue":      {layout, "web/templates/reports/revenue.html"},
		"profile/index":        {layout, "web/templates/profile/index.html"},
		"schedules/form":       {layout, "web/templates/schedules/form.html"},
		"users/index":          {layout, "web/templates/users/index.html"},
		"users/form":           {layout, "web/templates/users/form.html"},
		"404":                  {layout, "web/templates/404.html"},
	}

	for name, files := range pages {
		tmpl := template.New(filepath.Base(files[0])).Funcs(funcMap)
		for _, file := range files {
			content, err := s.assets.ReadFile(file)
			if err != nil {
				log.Printf("Warning: could not read template %s: %v", file, err)
				continue
			}
			_, err = tmpl.Parse(string(content))
			if err != nil {
				log.Printf("Warning: could not parse template %s: %v", file, err)
				continue
			}
		}
		s.templates[name] = tmpl
	}

	loginContent, err := s.assets.ReadFile("web/templates/login.html")
	if err == nil {
		tmpl := template.New("login").Funcs(funcMap)
		tmpl.Parse(string(loginContent))
		s.templates["login"] = tmpl
	}

	log.Printf("Loaded %d templates", len(s.templates))
}

func (s *Server) RenderTemplate(w http.ResponseWriter, name string, data interface{}) {
	if s.dev {
		s.renderDev(w, name, data)
		return
	}
	tmpl, ok := s.templates[name]
	if !ok {
		http.Error(w, "Template not found: "+name, http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	err := tmpl.ExecuteTemplate(w, "layout", data)
	if err != nil {
		log.Printf("Template error (%s): %v", name, err)
		http.Error(w, "Template error: "+err.Error(), http.StatusInternalServerError)
	}
}

func (s *Server) RenderLogin(w http.ResponseWriter, name string, data interface{}) {
	if s.dev {
		s.renderLoginDev(w, name, data)
		return
	}
	tmpl, ok := s.templates[name]
	if !ok {
		http.Error(w, "Template not found: "+name, http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	err := tmpl.Execute(w, data)
	if err != nil {
		log.Printf("Template error (%s): %v", name, err)
		http.Error(w, "Template error: "+err.Error(), http.StatusInternalServerError)
	}
}

func (s *Server) renderDev(w http.ResponseWriter, name string, data interface{}) {
	funcMap := s.templateFuncs()
	pagePaths := map[string][]string{
		"dashboard":            {"web/templates/layout.html", "web/templates/dashboard.html"},
		"trips/index":          {"web/templates/layout.html", "web/templates/trips/index.html"},
		"trips/form":           {"web/templates/layout.html", "web/templates/trips/form.html"},
		"trips/detail":         {"web/templates/layout.html", "web/templates/trips/detail.html"},
		"bookings/index":       {"web/templates/layout.html", "web/templates/bookings/index.html"},
		"bookings/detail":      {"web/templates/layout.html", "web/templates/bookings/detail.html"},
		"meeting_points/index": {"web/templates/layout.html", "web/templates/meeting_points/index.html"},
		"meeting_points/form":  {"web/templates/layout.html", "web/templates/meeting_points/form.html"},
		"packages/index":       {"web/templates/layout.html", "web/templates/packages/index.html"},
		"packages/form":        {"web/templates/layout.html", "web/templates/packages/form.html"},
		"reports/revenue":      {"web/templates/layout.html", "web/templates/reports/revenue.html"},
		"profile/index":        {"web/templates/layout.html", "web/templates/profile/index.html"},
		"schedules/form":       {"web/templates/layout.html", "web/templates/schedules/form.html"},
		"users/index":          {"web/templates/layout.html", "web/templates/users/index.html"},
		"users/form":           {"web/templates/layout.html", "web/templates/users/form.html"},
		"404":                  {"web/templates/layout.html", "web/templates/404.html"},
	}

	files, ok := pagePaths[name]
	if !ok {
		files = []string{"web/templates/layout.html", "web/templates/" + name + ".html"}
	}

	tmpl := template.New(filepath.Base(files[0])).Funcs(funcMap)
	for _, f := range files {
		content, err := os.ReadFile(f)
		if err != nil {
			http.Error(w, "Template error: "+err.Error(), http.StatusInternalServerError)
			return
		}
		_, err = tmpl.Parse(string(content))
		if err != nil {
			http.Error(w, "Template parse error: "+err.Error(), http.StatusInternalServerError)
			return
		}
	}
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	err := tmpl.ExecuteTemplate(w, "layout", data)
	if err != nil {
		http.Error(w, "Template error: "+err.Error(), http.StatusInternalServerError)
	}
}

func (s *Server) renderLoginDev(w http.ResponseWriter, name string, data interface{}) {
	content, err := os.ReadFile("web/templates/login.html")
	if err != nil {
		http.Error(w, "Template error: "+err.Error(), http.StatusInternalServerError)
		return
	}
	tmpl := template.New("login").Funcs(s.templateFuncs())
	_, err = tmpl.Parse(string(content))
	if err != nil {
		http.Error(w, "Template parse error: "+err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	err = tmpl.Execute(w, data)
	if err != nil {
		http.Error(w, "Template error: "+err.Error(), http.StatusInternalServerError)
	}
}

func (s *Server) setupRoutes() {
	mux := http.NewServeMux()

	if s.dev {
		staticDisk := http.FileServer(http.Dir("web/static"))
		mux.Handle("GET /static/", http.StripPrefix("/static/", staticDisk))
	} else {
		staticFS, err := fs.Sub(s.assets, "web/static")
		if err != nil {
			log.Fatal("Static FS:", err)
		}
		fileServer := http.FileServer(http.FS(staticFS))
		mux.Handle("GET /static/", http.StripPrefix("/static/", fileServer))
	}

	h := handler.New(database.DB, s, s.readFacilitiesJSON())

	mux.HandleFunc("GET /{$}", h.Dashboard)
	mux.HandleFunc("GET /login", h.Login)
	mux.HandleFunc("POST /login", h.LoginPost)
	mux.HandleFunc("GET /logout", h.Logout)

	mux.HandleFunc("GET /trips", h.TripList)
	mux.HandleFunc("GET /trips/new", h.TripForm)
	mux.HandleFunc("POST /trips", h.TripCreate)
	mux.HandleFunc("GET /trips/{id}", h.TripDetail)
	mux.HandleFunc("GET /trips/{id}/edit", h.TripFormEdit)
	mux.HandleFunc("PUT /trips/{id}", h.TripUpdate)
	mux.HandleFunc("DELETE /trips/{id}", h.TripDelete)

	mux.HandleFunc("GET /trips/{id}/schedules/new", h.ScheduleForm)
	mux.HandleFunc("POST /trips/{tripID}/schedules", h.ScheduleCreate)
	mux.HandleFunc("GET /schedules/{id}/edit", h.ScheduleFormEdit)
	mux.HandleFunc("PUT /schedules/{id}", h.ScheduleUpdate)
	mux.HandleFunc("DELETE /schedules/{id}", h.ScheduleDelete)

	mux.HandleFunc("GET /bookings", h.BookingList)
	mux.HandleFunc("GET /bookings/{id}", h.BookingDetail)
	mux.HandleFunc("PATCH /bookings/{id}/status", h.BookingStatus)

	mux.HandleFunc("GET /meeting-points", h.MeetingPointList)
	mux.HandleFunc("GET /meeting-points/new", h.MeetingPointForm)
	mux.HandleFunc("POST /meeting-points", h.MeetingPointCreate)
	mux.HandleFunc("GET /meeting-points/{id}/edit", h.MeetingPointFormEdit)
	mux.HandleFunc("PUT /meeting-points/{id}", h.MeetingPointUpdate)
	mux.HandleFunc("DELETE /meeting-points/{id}", h.MeetingPointDelete)

	mux.HandleFunc("GET /packages", h.PackageList)
	mux.HandleFunc("GET /packages/new", h.PackageForm)
	mux.HandleFunc("POST /packages", h.PackageCreate)
	mux.HandleFunc("GET /packages/{id}/edit", h.PackageFormEdit)
	mux.HandleFunc("PUT /packages/{id}", h.PackageUpdate)
	mux.HandleFunc("DELETE /packages/{id}", h.PackageDelete)

	mux.HandleFunc("GET /reports/revenue", h.RevenueReport)

	mux.HandleFunc("GET /profile", h.Profile)
	mux.HandleFunc("PUT /profile", h.ProfileUpdate)

	mux.HandleFunc("POST /trips/{tripID}/bookings/manual", h.BookingManualCreate)
	mux.HandleFunc("POST /bookings/{id}/payments", h.PaymentCreate)
	mux.HandleFunc("GET /proofs/{file}", h.PaymentProof)

	mux.HandleFunc("GET /users", h.UserList)
	mux.HandleFunc("GET /users/new", h.UserForm)
	mux.HandleFunc("POST /users", h.UserCreate)
	mux.HandleFunc("GET /users/{id}/edit", h.UserFormEdit)
	mux.HandleFunc("PUT /users/{id}", h.UserUpdate)
	mux.HandleFunc("DELETE /users/{id}", h.UserDelete)

	loggedMux := s.loggingMiddleware(mux)
	methodMux := s.methodOverrideMiddleware(loggedMux)
	authMux := s.authMiddleware(methodMux)

	// Catch-all for 404
	mux.HandleFunc("/{path...}", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotFound)
		s.RenderTemplate(w, "404", nil)
	})

	s.Handler = authMux
}

func (s *Server) methodOverrideMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method == "POST" {
			if err := r.ParseForm(); err == nil {
				if method := r.FormValue("_method"); method != "" {
					r.Method = strings.ToUpper(method)
				}
			}
		}
		next.ServeHTTP(w, r)
	})
}

func (s *Server) readFacilitiesJSON() string {
	var data []byte
	var err error
	if s.dev {
		data, err = os.ReadFile("web/static/data/facilities.json")
	} else {
		data, err = s.assets.ReadFile("web/static/data/facilities.json")
	}
	if err != nil {
		log.Printf("Warning: could not read facilities.json: %v", err)
		return "[]"
	}
	return string(data)
}
