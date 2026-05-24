package handler

import (
	"context"
	"database/sql"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/ayomendaki/ayomendaki-admin/internal/common"
	_ "modernc.org/sqlite"
)

// ── Mock Renderer ──

type testRenderer struct {
	lastTemplate string
	lastData     map[string]interface{}
	loginData    interface{}
}

func (r *testRenderer) RenderTemplate(w http.ResponseWriter, name string, data interface{}) {
	r.lastTemplate = name
	if m, ok := data.(map[string]interface{}); ok {
		r.lastData = m
	}
}

func (r *testRenderer) RenderLogin(w http.ResponseWriter, name string, data interface{}) {
	r.loginData = data
}

// ── Test Helpers ──

func setupTestDB(t *testing.T) *sql.DB {
	t.Helper()
	db, err := sql.Open("sqlite", ":memory:")
	if err != nil {
		t.Fatal(err)
	}
	schema := `
	CREATE TABLE operators (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		username TEXT UNIQUE NOT NULL,
		password_hash TEXT NOT NULL,
		name TEXT NOT NULL DEFAULT '',
		role TEXT NOT NULL DEFAULT 'admin',
		rating REAL NOT NULL DEFAULT 0,
		verified INTEGER NOT NULL DEFAULT 0,
		description TEXT NOT NULL DEFAULT '',
		phone TEXT NOT NULL DEFAULT '',
		whatsapp TEXT NOT NULL DEFAULT '',
		instagram TEXT NOT NULL DEFAULT '',
		created_at DATETIME NOT NULL DEFAULT (datetime('now'))
	);
	CREATE TABLE sessions (
		id TEXT PRIMARY KEY,
		operator_id INTEGER NOT NULL REFERENCES operators(id) ON DELETE CASCADE,
		created_at DATETIME NOT NULL DEFAULT (datetime('now')),
		expires_at TEXT NOT NULL
	);
	CREATE TABLE provinces (
		id TEXT PRIMARY KEY, name TEXT NOT NULL, image TEXT NOT NULL DEFAULT '', count INTEGER NOT NULL DEFAULT 0
	);
	CREATE TABLE mountains (
		id TEXT PRIMARY KEY, name TEXT NOT NULL, province_id TEXT NOT NULL REFERENCES provinces(id),
		height INTEGER NOT NULL DEFAULT 0, difficulty TEXT NOT NULL DEFAULT '', image TEXT NOT NULL DEFAULT '',
		description TEXT NOT NULL DEFAULT '', trending INTEGER NOT NULL DEFAULT 0,
		lat REAL NOT NULL DEFAULT 0, lng REAL NOT NULL DEFAULT 0, zoom INTEGER NOT NULL DEFAULT 12
	);
	CREATE TABLE packages (
		id INTEGER PRIMARY KEY AUTOINCREMENT, operator_id INTEGER NOT NULL REFERENCES operators(id),
		name TEXT NOT NULL, description TEXT NOT NULL DEFAULT '', facilities TEXT NOT NULL DEFAULT '[]'
	);
	CREATE TABLE meeting_points (
		id INTEGER PRIMARY KEY AUTOINCREMENT, operator_id INTEGER NOT NULL REFERENCES operators(id),
		type TEXT NOT NULL DEFAULT 'titik_jemput', name TEXT NOT NULL, address TEXT NOT NULL DEFAULT '',
		lat REAL NOT NULL DEFAULT 0, lng REAL NOT NULL DEFAULT 0
	);
	CREATE TABLE trips (
		id INTEGER PRIMARY KEY AUTOINCREMENT, operator_id INTEGER NOT NULL REFERENCES operators(id),
		mountain_id TEXT NOT NULL REFERENCES mountains(id), package_id INTEGER REFERENCES packages(id),
		name TEXT NOT NULL, route TEXT NOT NULL DEFAULT '', duration TEXT NOT NULL DEFAULT '',
		price INTEGER NOT NULL DEFAULT 0, created_at DATETIME NOT NULL DEFAULT (datetime('now'))
	);
	CREATE TABLE trip_meeting_points (
		trip_id INTEGER NOT NULL REFERENCES trips(id) ON DELETE CASCADE,
		meeting_point_id INTEGER NOT NULL REFERENCES meeting_points(id) ON DELETE CASCADE,
		order_index INTEGER NOT NULL DEFAULT 0, estimated_departure TEXT NOT NULL DEFAULT '',
		PRIMARY KEY (trip_id, meeting_point_id)
	);
	CREATE TABLE trip_package_prices (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		trip_id INTEGER NOT NULL REFERENCES trips(id) ON DELETE CASCADE,
		meeting_point_id INTEGER NOT NULL REFERENCES meeting_points(id),
		package_id INTEGER NOT NULL REFERENCES packages(id), price INTEGER NOT NULL DEFAULT 0
	);
	CREATE TABLE schedules (
		id INTEGER PRIMARY KEY AUTOINCREMENT, trip_id INTEGER NOT NULL REFERENCES trips(id) ON DELETE CASCADE,
		date_start TEXT NOT NULL, date_end TEXT NOT NULL, quota_total INTEGER NOT NULL DEFAULT 0,
		quota_remaining INTEGER NOT NULL DEFAULT 0, created_at DATETIME NOT NULL DEFAULT (datetime('now'))
	);
	CREATE TABLE bookings (
		id INTEGER PRIMARY KEY AUTOINCREMENT, trip_id INTEGER NOT NULL REFERENCES trips(id),
		schedule_id INTEGER NOT NULL REFERENCES schedules(id),
		lead_name TEXT NOT NULL, lead_phone TEXT NOT NULL, lead_email TEXT NOT NULL DEFAULT '',
		total INTEGER NOT NULL DEFAULT 0, status TEXT NOT NULL DEFAULT 'pending',
		created_at DATETIME NOT NULL DEFAULT (datetime('now'))
	);
	CREATE TABLE booking_participants (
		id INTEGER PRIMARY KEY AUTOINCREMENT, booking_id INTEGER NOT NULL REFERENCES bookings(id) ON DELETE CASCADE,
		name TEXT NOT NULL, ktp TEXT NOT NULL DEFAULT ''
	);
	CREATE TABLE booking_addons (
		id INTEGER PRIMARY KEY AUTOINCREMENT, booking_id INTEGER NOT NULL REFERENCES bookings(id) ON DELETE CASCADE,
		name TEXT NOT NULL, price INTEGER NOT NULL DEFAULT 0
	);
	`
	if _, err := db.Exec(schema); err != nil {
		t.Fatal("schema:", err)
	}
	return db
}

func setupHandler(t *testing.T, db *sql.DB) *Handler {
	t.Helper()
	return &Handler{db: db, renderer: &testRenderer{}, FacilitiesJSON: "[]"}
}

func newRequest(method, path string, body string) *http.Request {
	req := httptest.NewRequest(method, path, strings.NewReader(body))
	if body != "" {
		req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	}
	return req
}

func withPathValue(req *http.Request, key, value string) *http.Request {
	req.SetPathValue(key, value)
	return req
}

func authRequest(method, path string, id string, body string) *http.Request {
	return withAuth(withPathValue(newRequest(method, path, body), "id", id), 1)
}

func authRequest2(method, path string, key, value string, body string) *http.Request {
	return withAuth(withPathValue(newRequest(method, path, body), key, value), 1)
}

func withAuth(req *http.Request, opID int64) *http.Request {
	ctx := context.WithValue(req.Context(), common.ContextKeyOperatorID, opID)
	ctx = context.WithValue(ctx, common.ContextKeyRole, common.RoleAdmin)
	return req.WithContext(ctx)
}

func withUserAuth(req *http.Request, opID int64, role string) *http.Request {
	ctx := context.WithValue(req.Context(), common.ContextKeyOperatorID, opID)
	ctx = context.WithValue(ctx, common.ContextKeyRole, role)
	return req.WithContext(ctx)
}

func assertStatus(t *testing.T, w *httptest.ResponseRecorder, code int) {
	t.Helper()
	if w.Code != code {
		t.Errorf("Expected status %d, got %d. Body: %s", code, w.Code, w.Body.String())
	}
}

func assertContains(t *testing.T, body, substr string) {
	t.Helper()
	if !strings.Contains(body, substr) {
		t.Errorf("Expected body to contain %q. Body: %s", substr, body)
	}
}

func assertRedirect(t *testing.T, w *httptest.ResponseRecorder, target string) {
	t.Helper()
	if w.Code != http.StatusSeeOther {
		t.Errorf("Expected 303 redirect, got %d", w.Code)
	}
	if loc := w.Header().Get("Location"); !strings.Contains(loc, target) {
		t.Errorf("Expected Location to contain %q, got %q", target, loc)
	}
}

func seedOperator(t *testing.T, db *sql.DB, username, passwordHash, name, role string) {
	t.Helper()
	_, err := db.Exec("INSERT INTO operators (username, password_hash, name, role) VALUES (?, ?, ?, ?)",
		username, passwordHash, name, role)
	if err != nil {
		t.Fatal("seed operator:", err)
	}
}

func seedMountain(t *testing.T, db *sql.DB, id, name string) {
	t.Helper()
	db.Exec("INSERT OR IGNORE INTO provinces (id, name) VALUES (?, ?)", "test", "Test Province")
	_, err := db.Exec("INSERT OR IGNORE INTO mountains (id, name, province_id, height, difficulty) VALUES (?, ?, 'test', 1000, 'Mudah')",
		id, name)
	if err != nil {
		t.Fatal("seed mountain:", err)
	}
}

func seedTrip(db *sql.DB, opID int64, name, mountainID, route, duration string) int64 {
	res, _ := db.Exec("INSERT INTO trips (operator_id, mountain_id, name, route, duration, price) VALUES (?, ?, ?, ?, ?, 0)",
		opID, mountainID, name, route, duration)
	id, _ := res.LastInsertId()
	return id
}
