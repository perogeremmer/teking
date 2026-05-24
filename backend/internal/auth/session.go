package auth

import (
	"crypto/rand"
	"database/sql"
	"encoding/hex"
	"time"
)

func GenerateSessionID() (string, error) {
	b := make([]byte, 32)
	_, err := rand.Read(b)
	if err != nil {
		return "", err
	}
	return hex.EncodeToString(b), nil
}

func CreateSession(db *sql.DB, operatorID int64) (*string, error) {
	id, err := GenerateSessionID()
	if err != nil {
		return nil, err
	}

	expiresAt := time.Now().Add(24 * time.Hour)

	_, err = db.Exec(
		"INSERT INTO sessions (id, operator_id, expires_at) VALUES (?, ?, ?)",
		id, operatorID, expiresAt.Format("2006-01-02 15:04:05"),
	)
	if err != nil {
		return nil, err
	}

	return &id, nil
}

func ValidateSession(db *sql.DB, sessionID string) (*int64, error) {
	var operatorID int64
	var expiresAt string

	err := db.QueryRow(
		"SELECT operator_id, expires_at FROM sessions WHERE id = ?", sessionID,
	).Scan(&operatorID, &expiresAt)

	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}

	expiresTime, err := time.Parse("2006-01-02 15:04:05", expiresAt)
	if err != nil {
		expiresTime, err = time.Parse(time.RFC3339, expiresAt)
		if err != nil {
			return nil, err
		}
	}

	if time.Now().After(expiresTime) {
		_, _ = db.Exec("DELETE FROM sessions WHERE id = ?", sessionID)
		return nil, nil
	}

	return &operatorID, nil
}

func DeleteSession(db *sql.DB, sessionID string) error {
	_, err := db.Exec("DELETE FROM sessions WHERE id = ?", sessionID)
	return err
}
