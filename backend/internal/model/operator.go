package model

import "time"

type Operator struct {
	ID          int64  `json:"id"`
	Username    string `json:"username"`
	PasswordHash string `json:"-"`
	Name        string `json:"name"`
	Role        string `json:"role"`
	Logo        string `json:"logo"`
	Rating      float64 `json:"rating"`
	Verified    bool   `json:"verified"`
	Description string `json:"description"`
	Phone       string `json:"phone"`
	Whatsapp    string `json:"whatsapp"`
	Instagram   string `json:"instagram"`
	CreatedAt   string `json:"created_at"`
}

type Session struct {
	ID        string    `json:"id"`
	OperatorID int64    `json:"operator_id"`
	CreatedAt time.Time `json:"created_at"`
	ExpiresAt time.Time `json:"expires_at"`
}

type Province struct {
	ID    string `json:"id"`
	Name  string `json:"name"`
	Image string `json:"image"`
	Count int    `json:"count"`
}

type Mountain struct {
	ID          string  `json:"id"`
	Name        string  `json:"name"`
	ProvinceID  string  `json:"province_id"`
	Height      int     `json:"height"`
	Difficulty  string  `json:"difficulty"`
	Image       string  `json:"image"`
	Description string  `json:"description"`
	Trending    bool    `json:"trending"`
	Lat         float64 `json:"lat"`
	Lng         float64 `json:"lng"`
	Zoom        int     `json:"zoom"`
}
