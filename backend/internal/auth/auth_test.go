package auth

import (
	"database/sql"
	"os"
	"testing"
	"time"

	_ "modernc.org/sqlite"
)

func TestHashPassword(t *testing.T) {
	hash, err := HashPassword("secret123")
	if err != nil {
		t.Fatal("HashPassword failed:", err)
	}
	if hash == "" {
		t.Fatal("Hash is empty")
	}
	if hash == "secret123" {
		t.Fatal("Hash should not equal plaintext")
	}
}

func TestCheckPasswordValid(t *testing.T) {
	hash, _ := HashPassword("secret123")
	if !CheckPassword("secret123", hash) {
		t.Fatal("Expected password to match")
	}
}

func TestCheckPasswordInvalid(t *testing.T) {
	hash, _ := HashPassword("secret123")
	if CheckPassword("wrongpass", hash) {
		t.Fatal("Expected password NOT to match")
	}
}

func TestGenerateSessionID(t *testing.T) {
	id1, _ := GenerateSessionID()
	id2, _ := GenerateSessionID()
	if len(id1) != 64 {
		t.Fatalf("Expected 64 hex chars, got %d", len(id1))
	}
	if id1 == id2 {
		t.Fatal("Session IDs should be unique")
	}
}

func setupSessionDB(t *testing.T) *sql.DB {
	db, err := sql.Open("sqlite", ":memory:")
	if err != nil {
		t.Fatal(err)
	}
	_, err = db.Exec(`CREATE TABLE sessions (
		id TEXT PRIMARY KEY,
		operator_id INTEGER NOT NULL,
		expires_at TEXT NOT NULL
	)`)
	if err != nil {
		t.Fatal(err)
	}
	return db
}

func TestCreateSession(t *testing.T) {
	db := setupSessionDB(t)
	defer db.Close()

	sid, err := CreateSession(db, 1)
	if err != nil {
		t.Fatal("CreateSession failed:", err)
	}
	if sid == nil || *sid == "" {
		t.Fatal("Session ID is nil or empty")
	}

	var count int
	db.QueryRow("SELECT COUNT(*) FROM sessions WHERE id = ?", *sid).Scan(&count)
	if count != 1 {
		t.Fatal("Session not found in database")
	}
}

func TestValidateSessionValid(t *testing.T) {
	db := setupSessionDB(t)
	defer db.Close()

	expires := time.Now().Add(1 * time.Hour)
	db.Exec("INSERT INTO sessions (id, operator_id, expires_at) VALUES (?, ?, ?)",
		"test-session-1", 1, expires.Format("2006-01-02 15:04:05"))

	oid, err := ValidateSession(db, "test-session-1")
	if err != nil {
		t.Fatal("ValidateSession failed:", err)
	}
	if oid == nil || *oid != 1 {
		t.Fatal("Expected operator_id = 1")
	}
}

func TestValidateSessionNotFound(t *testing.T) {
	db := setupSessionDB(t)
	defer db.Close()

	oid, err := ValidateSession(db, "nonexistent")
	if err != nil {
		t.Fatal("ValidateSession failed:", err)
	}
	if oid != nil {
		t.Fatal("Expected nil for non-existent session")
	}
}

func TestValidateSessionExpired(t *testing.T) {
	db := setupSessionDB(t)
	defer db.Close()

	expires := time.Now().UTC().Add(-1 * time.Hour)
	db.Exec("INSERT INTO sessions (id, operator_id, expires_at) VALUES (?, ?, ?)",
		"expired-session", 1, expires.Format("2006-01-02 15:04:05"))

	oid, err := ValidateSession(db, "expired-session")
	if err != nil {
		t.Fatal("ValidateSession failed:", err)
	}
	if oid != nil {
		t.Fatal("Expected nil for expired session — should be auto-deleted")
	}

	var count int
	db.QueryRow("SELECT COUNT(*) FROM sessions WHERE id = 'expired-session'").Scan(&count)
	if count != 0 {
		t.Fatal("Expired session should have been deleted")
	}
}

func TestDeleteSession(t *testing.T) {
	db := setupSessionDB(t)
	defer db.Close()

	db.Exec("INSERT INTO sessions (id, operator_id, expires_at) VALUES (?, ?, ?)",
		"delete-test", 1, "2026-12-31 23:59:59")

	err := DeleteSession(db, "delete-test")
	if err != nil {
		t.Fatal("DeleteSession failed:", err)
	}

	var count int
	db.QueryRow("SELECT COUNT(*) FROM sessions WHERE id = 'delete-test'").Scan(&count)
	if count != 0 {
		t.Fatal("Session should have been deleted")
	}
}

func TestMain(m *testing.M) {
	os.Exit(m.Run())
}
