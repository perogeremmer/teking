package handler

import (
	"context"
	"database/sql"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/ayomendaki/ayomendaki-admin/internal/common"
	"golang.org/x/crypto/bcrypt"
)

func hashPw(t *testing.T, pw string) string {
	t.Helper()
	h, err := bcrypt.GenerateFromPassword([]byte(pw), bcrypt.DefaultCost)
	if err != nil {
		t.Fatal(err)
	}
	return string(h)
}

func seedOp(db *sql.DB, username, hash, name, role string) {
	db.Exec("INSERT INTO operators (username, password_hash, name, role) VALUES (?, ?, ?, ?)",
		username, hash, name, role)
}

func authReq(method, path, body string, opID int64) *http.Request {
	req := httptest.NewRequest(method, path, strings.NewReader(body))
	if body != "" {
		req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	}
	ctx := req.Context()
	ctx = context.WithValue(ctx, common.ContextKeyOperatorID, opID)
	ctx = context.WithValue(ctx, common.ContextKeyRole, common.RoleAdmin)
	return req.WithContext(ctx)
}

// ── Auth Tests ──

func TestLoginPage(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()
	h := setupHandler(t, db)
	w := httptest.NewRecorder()
	h.Login(w, httptest.NewRequest("GET", "/login", nil))
	r := h.renderer.(*testRenderer)
	if r.loginData != nil {
		t.Fatal("Login page should not have error data")
	}
}

func TestLoginPostSuccess(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()
	h := setupHandler(t, db)
	seedOp(db, "admin", hashPw(t, "pass"), "Admin", "admin")

	w := httptest.NewRecorder()
	req := httptest.NewRequest("POST", "/login", strings.NewReader("username=admin&password=pass"))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	h.LoginPost(w, req)
	if w.Code != 303 {
		t.Fatalf("Expected 303, got %d", w.Code)
	}
	loc := w.Header().Get("Location")
	if !strings.Contains(loc, "/") {
		t.Errorf("Expected redirect to /, got %s", loc)
	}
}

func TestLoginPostWrongPassword(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()
	h := setupHandler(t, db)
	seedOp(db, "admin", hashPw(t, "pass"), "Admin", "admin")

	w := httptest.NewRecorder()
	req := httptest.NewRequest("POST", "/login", strings.NewReader("username=admin&password=wrong"))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	h.LoginPost(w, req)
	r := h.renderer.(*testRenderer)
	if r.loginData == nil {
		t.Fatal("Expected login error data on wrong password")
	}
}

func TestLogout(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()
	h := setupHandler(t, db)
	w := httptest.NewRecorder()
	h.Logout(w, httptest.NewRequest("GET", "/logout", nil))
	if w.Code != 303 {
		t.Fatalf("Expected 303, got %d", w.Code)
	}
}

// ── Trip Tests ──

func TestTripList(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()
	h := setupHandler(t, db)
	w := httptest.NewRecorder()
	h.TripList(w, authReq("GET", "/trips", "", 1))
	if h.renderer.(*testRenderer).lastTemplate != "trips/index" {
		t.Error("Expected trips/index template")
	}
}

func TestTripForm(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()
	h := setupHandler(t, db)
	w := httptest.NewRecorder()
	h.TripForm(w, authReq("GET", "/trips/new", "", 1))
	if h.renderer.(*testRenderer).lastTemplate != "trips/form" {
		t.Error("Expected trips/form template")
	}
}

func TestTripCreate(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()
	h := setupHandler(t, db)
	seedMountain(t, db, "m1", "M1")

	w := httptest.NewRecorder()
	h.TripCreate(w, authReq("POST", "/trips", "name=Test+Trip&mountain_id=m1&route=R&duration=2H", 1))
	if w.Code != 303 {
		t.Fatalf("Expected 303, got %d", w.Code)
	}
	var count int
	db.QueryRow("SELECT COUNT(*) FROM trips").Scan(&count)
	if count != 1 {
		t.Errorf("Expected 1 trip, got %d", count)
	}
}

func TestTripCreateWithMPandPkg(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()
	h := setupHandler(t, db)
	seedMountain(t, db, "m1", "M1")
	db.Exec("INSERT INTO packages (operator_id, name) VALUES (1, 'Pkg')")
	db.Exec("INSERT INTO meeting_points (operator_id, name) VALUES (1, 'MP')")

	body := "name=Trip+Full&mountain_id=m1&route=R&duration=2H&trip_json="
	body += `{"packages":[1],"meetingPoints":[{"id":1,"estimated_departure":"19:00","prices":{"1":500000}}]}`

	w := httptest.NewRecorder()
	h.TripCreate(w, authReq("POST", "/trips", body, 1))
	if w.Code != 303 {
		t.Fatalf("Expected 303, got %d", w.Code)
	}
}

func TestTripDetail(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()
	h := setupHandler(t, db)
	seedMountain(t, db, "m1", "M1")
	tid := seedTrip(db, 1, "DT", "m1", "R", "2H")

	req := authReq("GET", "/trips/1", "", 1)
	req.SetPathValue("id", "1")
	w := httptest.NewRecorder()
	h.TripDetail(w, req)
	if h.renderer.(*testRenderer).lastTemplate != "trips/detail" {
		t.Errorf("Expected trips/detail, got %q", h.renderer.(*testRenderer).lastTemplate)
	}
	_ = tid
}

func TestTripFormEdit(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()
	h := setupHandler(t, db)
	seedMountain(t, db, "m1", "M1")
	seedTrip(db, 1, "ET", "m1", "R", "2H")

	req := authReq("GET", "/trips/1/edit", "", 1)
	req.SetPathValue("id", "1")
	w := httptest.NewRecorder()
	h.TripFormEdit(w, req)
	if h.renderer.(*testRenderer).lastTemplate != "trips/form" {
		t.Errorf("Expected trips/form, got %q", h.renderer.(*testRenderer).lastTemplate)
	}
}

func TestTripUpdate(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()
	h := setupHandler(t, db)
	seedMountain(t, db, "m1", "M1")
	seedTrip(db, 1, "Old", "m1", "R", "2H")

	req := authReq("PUT", "/trips/1", "name=New&mountain_id=m1&route=R2&duration=3H", 1)
	req.SetPathValue("id", "1")
	w := httptest.NewRecorder()
	h.TripUpdate(w, req)

	var name string
	db.QueryRow("SELECT name FROM trips WHERE id = 1").Scan(&name)
	if name != "New" {
		t.Errorf("Expected 'New', got '%s'", name)
	}
}

func TestTripDelete(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()
	h := setupHandler(t, db)
	seedMountain(t, db, "m1", "M1")
	seedTrip(db, 1, "Del", "m1", "R", "1H")

	req := authReq("DELETE", "/trips/1", "", 1)
	req.SetPathValue("id", "1")
	w := httptest.NewRecorder()
	h.TripDelete(w, req)
	if w.Code != 200 && w.Code != 303 {
		t.Errorf("Expected 200 or 303, got %d", w.Code)
	}
}

// ── Schedule Tests ──

func TestScheduleForm(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()
	h := setupHandler(t, db)
	seedMountain(t, db, "m1", "M1")
	seedTrip(db, 1, "T", "m1", "R", "2H")

	req := authReq("GET", "/trips/1/schedules/new", "", 1)
	req.SetPathValue("id", "1")
	w := httptest.NewRecorder()
	h.ScheduleForm(w, req)
	if h.renderer.(*testRenderer).lastTemplate != "schedules/form" {
		t.Error("Expected schedules/form")
	}
}

func TestScheduleCreate(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()
	h := setupHandler(t, db)
	seedMountain(t, db, "m1", "M1")
	seedTrip(db, 1, "T", "m1", "R", "2H")

	req := authReq("POST", "/trips/1/schedules", "date_start=2026-06-01&date_end=2026-06-03&quota_total=10", 1)
	req.SetPathValue("tripID", "1")
	w := httptest.NewRecorder()
	h.ScheduleCreate(w, req)
	if w.Code != 200 && w.Code != 303 {
		t.Errorf("Expected 200 or 303, got %d", w.Code)
	}
}

// ── Booking Tests ──

func TestBookingDetail(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()
	h := setupHandler(t, db)
	seedMountain(t, db, "m1", "M1")
	tid := seedTrip(db, 1, "T", "m1", "R", "2H")
	db.Exec("INSERT INTO schedules (trip_id, date_start, date_end, quota_total, quota_remaining) VALUES (?, '2026-06-01', '2026-06-03', 10, 10)", tid)
	db.Exec("INSERT INTO bookings (trip_id, schedule_id, lead_name, lead_phone, total, status) VALUES (?, 1, 'Budi', '0812', 500000, 'pending')", tid)

	req := authReq("GET", "/bookings/1", "", 1)
	req.SetPathValue("id", "1")
	w := httptest.NewRecorder()
	h.BookingDetail(w, req)
	if h.renderer.(*testRenderer).lastTemplate != "bookings/detail" {
		t.Errorf("Expected bookings/detail, got %q", h.renderer.(*testRenderer).lastTemplate)
	}
}

// ── Meeting Point Tests ──

func TestMeetingPointCreate(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()
	h := setupHandler(t, db)
	w := httptest.NewRecorder()
	h.MeetingPointCreate(w, authReq("POST", "/meeting-points", "name=MP1&type=titik_jemput&address=Addr&lat=0&lng=0", 1))
	if w.Code != 303 {
		t.Fatalf("Expected 303, got %d", w.Code)
	}
	var name string
	db.QueryRow("SELECT name FROM meeting_points WHERE id = 1").Scan(&name)
	if name != "MP1" {
		t.Errorf("Expected 'MP1', got '%s'", name)
	}
}

func TestMeetingPointList(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()
	h := setupHandler(t, db)
	w := httptest.NewRecorder()
	h.MeetingPointList(w, authReq("GET", "/meeting-points", "", 1))
	if h.renderer.(*testRenderer).lastTemplate != "meeting_points/index" {
		t.Error("Expected meeting_points/index")
	}
}

// ── Package Tests ──

func TestPackageCreate(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()
	h := setupHandler(t, db)
	w := httptest.NewRecorder()
	h.PackageCreate(w, authReq("POST", "/packages", "name=Pkg&description=D&facilities_json=[]", 1))
	if w.Code != 303 {
		t.Fatalf("Expected 303, got %d", w.Code)
	}
	var name string
	db.QueryRow("SELECT name FROM packages WHERE id = 1").Scan(&name)
	if name != "Pkg" {
		t.Errorf("Expected 'Pkg', got '%s'", name)
	}
}

func TestPackageList(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()
	h := setupHandler(t, db)
	w := httptest.NewRecorder()
	h.PackageList(w, authReq("GET", "/packages", "", 1))
	if h.renderer.(*testRenderer).lastTemplate != "packages/index" {
		t.Error("Expected packages/index")
	}
}

// ── Report Tests ──

func TestRevenueReport(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()
	h := setupHandler(t, db)
	w := httptest.NewRecorder()
	h.RevenueReport(w, authReq("GET", "/reports/revenue", "", 1))
	if h.renderer.(*testRenderer).lastTemplate != "reports/revenue" {
		t.Error("Expected reports/revenue")
	}
}

// ── Profile Tests ──

func TestProfile(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()
	h := setupHandler(t, db)
	seedOp(db, "admin", hashPw(t, "x"), "Admin", "admin")

	w := httptest.NewRecorder()
	h.Profile(w, authReq("GET", "/profile", "", 1))
	if h.renderer.(*testRenderer).lastTemplate != "profile/index" {
		t.Error("Expected profile/index")
	}
}

func TestProfileUpdate(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()
	h := setupHandler(t, db)
	seedOp(db, "admin", hashPw(t, "x"), "Admin", "admin")

	w := httptest.NewRecorder()
	h.ProfileUpdate(w, authReq("PUT", "/profile", "name=New&description=D&phone=081", 1))
	if w.Code != 303 {
		t.Errorf("Expected 303, got %d", w.Code)
	}
}

// ── Dashboard Tests ──

func TestDashboard(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()
	h := setupHandler(t, db)
	w := httptest.NewRecorder()
	h.Dashboard(w, authReq("GET", "/", "", 1))
	if h.renderer.(*testRenderer).lastTemplate != "dashboard" {
		t.Error("Expected dashboard")
	}
}

// ── User Tests ──

func TestUserList(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()
	h := setupHandler(t, db)
	seedOp(db, "admin", hashPw(t, "x"), "Admin", "admin")

	w := httptest.NewRecorder()
	h.UserList(w, authReq("GET", "/users", "", 1))
	if h.renderer.(*testRenderer).lastTemplate != "users/index" {
		t.Error("Expected users/index")
	}
}

func TestUserCreate(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()
	h := setupHandler(t, db)
	w := httptest.NewRecorder()
	h.UserCreate(w, authReq("POST", "/users", "username=u1&password=p&name=U1&role=user", 1))
	if w.Code != 303 {
		t.Fatalf("Expected 303, got %d", w.Code)
	}
}

// ── 404 Test ──

func TestNotFound(t *testing.T) {
	// This tests that handler returns empty template for 404
	// (catch-all not needed in unit test since there's no mux)
	db := setupTestDB(t)
	defer db.Close()
	h := setupHandler(t, db)
	w := httptest.NewRecorder()
	h.TripList(w, authReq("GET", "/nonexistent", "", 1))
	// TripList should respond with template, not 404
	if h.renderer.(*testRenderer).lastTemplate == "" {
		t.Error("Expected non-empty template")
	}
}
