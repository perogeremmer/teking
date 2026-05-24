package model

type Booking struct {
	ID         int64                `json:"id"`
	TripID     int64                `json:"trip_id"`
	ScheduleID int64                `json:"schedule_id"`
	LeadName   string               `json:"lead_name"`
	LeadPhone  string               `json:"lead_phone"`
	LeadEmail  string               `json:"lead_email"`
	Total      int64                `json:"total"`
	Status     string               `json:"status"`
	CreatedAt  string               `json:"created_at"`
	Participants []BookingParticipant `json:"participants,omitempty"`
	Addons     []BookingAddon       `json:"addons,omitempty"`
}

type BookingParticipant struct {
	ID        int64  `json:"id"`
	BookingID int64  `json:"booking_id"`
	Name      string `json:"name"`
	KTP       string `json:"ktp"`
}

type BookingAddon struct {
	ID        int64  `json:"id"`
	BookingID int64  `json:"booking_id"`
	Name      string `json:"name"`
	Price     int64  `json:"price"`
}
