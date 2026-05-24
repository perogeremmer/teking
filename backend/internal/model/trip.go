package model

type Trip struct {
	ID          int64   `json:"id"`
	OperatorID  int64   `json:"operator_id"`
	MountainID  string  `json:"mountain_id"`
	PackageID   *int64  `json:"package_id"`
	Name        string  `json:"name"`
	Route       string  `json:"route"`
	Duration    string  `json:"duration"`
	Price       int64   `json:"price"`
	CreatedAt   string  `json:"created_at"`
}

type Schedule struct {
	ID             int64  `json:"id"`
	TripID         int64  `json:"trip_id"`
	DateStart      string `json:"date_start"`
	DateEnd        string `json:"date_end"`
	QuotaTotal     int    `json:"quota_total"`
	QuotaRemaining int    `json:"quota_remaining"`
	CreatedAt      string `json:"created_at"`
}
