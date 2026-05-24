package database

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"path/filepath"

	_ "modernc.org/sqlite"
)

var DB *sql.DB

func Init(dbPath string) error {
	dir := filepath.Dir(dbPath)
	if dir != "" {
		if err := os.MkdirAll(dir, 0755); err != nil {
			return fmt.Errorf("create db directory: %w", err)
		}
	}

	var err error
	DB, err = sql.Open("sqlite", dbPath)
	if err != nil {
		return fmt.Errorf("open database: %w", err)
	}

	DB.SetMaxOpenConns(1)

	if err := DB.Ping(); err != nil {
		return fmt.Errorf("ping database: %w", err)
	}

	if err := migrate(); err != nil {
		return fmt.Errorf("migrate: %w", err)
	}

	log.Println("Database initialized:", dbPath)
	return nil
}

func migrate() error {
	schema := `
	CREATE TABLE IF NOT EXISTS operators (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		username TEXT UNIQUE NOT NULL,
		password_hash TEXT NOT NULL,
		name TEXT NOT NULL DEFAULT '',
		role TEXT NOT NULL DEFAULT 'admin',
		logo TEXT NOT NULL DEFAULT '',
		rating REAL NOT NULL DEFAULT 0,
		verified INTEGER NOT NULL DEFAULT 0,
		description TEXT NOT NULL DEFAULT '',
		phone TEXT NOT NULL DEFAULT '',
		whatsapp TEXT NOT NULL DEFAULT '',
		instagram TEXT NOT NULL DEFAULT '',
		created_at DATETIME NOT NULL DEFAULT (datetime('now'))
	);

	CREATE TABLE IF NOT EXISTS sessions (
		id TEXT PRIMARY KEY,
		operator_id INTEGER NOT NULL REFERENCES operators(id) ON DELETE CASCADE,
		created_at DATETIME NOT NULL DEFAULT (datetime('now')),
		expires_at TEXT NOT NULL
	);

	CREATE TABLE IF NOT EXISTS provinces (
		id TEXT PRIMARY KEY,
		name TEXT NOT NULL,
		image TEXT NOT NULL DEFAULT '',
		count INTEGER NOT NULL DEFAULT 0
	);

	CREATE TABLE IF NOT EXISTS mountains (
		id TEXT PRIMARY KEY,
		name TEXT NOT NULL,
		province_id TEXT NOT NULL REFERENCES provinces(id),
		height INTEGER NOT NULL DEFAULT 0,
		difficulty TEXT NOT NULL DEFAULT '',
		image TEXT NOT NULL DEFAULT '',
		description TEXT NOT NULL DEFAULT '',
		trending INTEGER NOT NULL DEFAULT 0,
		lat REAL NOT NULL DEFAULT 0,
		lng REAL NOT NULL DEFAULT 0,
		zoom INTEGER NOT NULL DEFAULT 12
	);

	CREATE TABLE IF NOT EXISTS packages (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		operator_id INTEGER NOT NULL REFERENCES operators(id),
		name TEXT NOT NULL,
		description TEXT NOT NULL DEFAULT '',
		facilities TEXT NOT NULL DEFAULT '[]'
	);

	CREATE TABLE IF NOT EXISTS meeting_points (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		operator_id INTEGER NOT NULL REFERENCES operators(id),
		type TEXT NOT NULL DEFAULT 'titik_jemput',
		name TEXT NOT NULL,
		address TEXT NOT NULL DEFAULT '',
		lat REAL NOT NULL DEFAULT 0,
		lng REAL NOT NULL DEFAULT 0
	);

	CREATE TABLE IF NOT EXISTS trips (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		operator_id INTEGER NOT NULL REFERENCES operators(id),
		mountain_id TEXT NOT NULL REFERENCES mountains(id),
		package_id INTEGER REFERENCES packages(id),
		name TEXT NOT NULL,
		route TEXT NOT NULL DEFAULT '',
		duration TEXT NOT NULL DEFAULT '',
		price INTEGER NOT NULL DEFAULT 0,
		created_at DATETIME NOT NULL DEFAULT (datetime('now'))
	);

	CREATE TABLE IF NOT EXISTS trip_meeting_points (
		trip_id INTEGER NOT NULL REFERENCES trips(id) ON DELETE CASCADE,
		meeting_point_id INTEGER NOT NULL REFERENCES meeting_points(id) ON DELETE CASCADE,
		order_index INTEGER NOT NULL DEFAULT 0,
		estimated_departure TEXT NOT NULL DEFAULT '',
		PRIMARY KEY (trip_id, meeting_point_id)
	);

	CREATE TABLE IF NOT EXISTS trip_package_prices (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		trip_id INTEGER NOT NULL REFERENCES trips(id) ON DELETE CASCADE,
		meeting_point_id INTEGER NOT NULL REFERENCES meeting_points(id),
		package_id INTEGER NOT NULL REFERENCES packages(id),
		price INTEGER NOT NULL DEFAULT 0
	);

	CREATE TABLE IF NOT EXISTS schedules (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		trip_id INTEGER NOT NULL REFERENCES trips(id) ON DELETE CASCADE,
		date_start TEXT NOT NULL,
		date_end TEXT NOT NULL,
		quota_total INTEGER NOT NULL DEFAULT 0,
		quota_remaining INTEGER NOT NULL DEFAULT 0,
		created_at DATETIME NOT NULL DEFAULT (datetime('now'))
	);

	CREATE TABLE IF NOT EXISTS customers (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		name TEXT NOT NULL,
		phone TEXT NOT NULL,
		nik TEXT NOT NULL DEFAULT '',
		email TEXT NOT NULL DEFAULT '',
		password_hash TEXT NOT NULL DEFAULT '',
		created_at DATETIME NOT NULL DEFAULT (datetime('now'))
	);

	CREATE TABLE IF NOT EXISTS bookings (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		trip_id INTEGER NOT NULL REFERENCES trips(id),
		schedule_id INTEGER NOT NULL REFERENCES schedules(id),
		lead_name TEXT NOT NULL,
		lead_phone TEXT NOT NULL,
		lead_email TEXT NOT NULL DEFAULT '',
		total INTEGER NOT NULL DEFAULT 0,
		status TEXT NOT NULL DEFAULT 'pending',
		payment_status TEXT NOT NULL DEFAULT 'unpaid',
		customer_id INTEGER REFERENCES customers(id),
		meeting_point_id INTEGER REFERENCES meeting_points(id),
		package_id INTEGER REFERENCES packages(id),
		created_at DATETIME NOT NULL DEFAULT (datetime('now'))
	);

	CREATE TABLE IF NOT EXISTS booking_participants (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		booking_id INTEGER NOT NULL REFERENCES bookings(id) ON DELETE CASCADE,
		name TEXT NOT NULL,
		ktp TEXT NOT NULL DEFAULT ''
	);

	CREATE TABLE IF NOT EXISTS booking_addons (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		booking_id INTEGER NOT NULL REFERENCES bookings(id) ON DELETE CASCADE,
		name TEXT NOT NULL,
		price INTEGER NOT NULL DEFAULT 0
	);

	CREATE TABLE IF NOT EXISTS payments (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		booking_id INTEGER NOT NULL REFERENCES bookings(id) ON DELETE CASCADE,
		amount INTEGER NOT NULL DEFAULT 0,
		notes TEXT NOT NULL DEFAULT '',
		proof_file TEXT NOT NULL DEFAULT '',
		created_at DATETIME NOT NULL DEFAULT (datetime('now'))
	);

	CREATE TABLE IF NOT EXISTS addon_templates (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		operator_id INTEGER NOT NULL REFERENCES operators(id),
		name TEXT NOT NULL,
		price INTEGER NOT NULL DEFAULT 0,
		icon TEXT NOT NULL DEFAULT 'bx-package'
	);
	`

	_, err := DB.Exec(schema)
	if err != nil {
		return fmt.Errorf("execute schema: %w", err)
	}
	_, err = DB.Exec("CREATE UNIQUE INDEX IF NOT EXISTS idx_mp_name_type ON meeting_points(name, type)")
	if err != nil {
		return fmt.Errorf("create index: %w", err)
	}
	return nil
}

func Close() {
	if DB != nil {
		DB.Close()
	}
}
