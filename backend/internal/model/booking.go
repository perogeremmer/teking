package model

type Booking struct {
	ID              int64                `json:"id"`
	TripID          int64                `json:"trip_id"`
	ScheduleID      int64                `json:"schedule_id"`
	LeadName        string               `json:"lead_name"`
	LeadPhone       string               `json:"lead_phone"`
	LeadEmail       string               `json:"lead_email"`
	Total           int64                `json:"total"`
	Status          string               `json:"status"`
	PaymentStatus   string               `json:"payment_status"`
	CustomerID      int64                `json:"customer_id"`
	MeetingPointID  int64                `json:"meeting_point_id"`
	PackageID       int64                `json:"package_id"`
	CreatedAt       string               `json:"created_at"`
	Participants    []BookingParticipant `json:"participants,omitempty"`
	Payments        []Payment            `json:"payments,omitempty"`
}

type BookingParticipant struct {
	ID        int64  `json:"id"`
	BookingID int64  `json:"booking_id"`
	Name      string `json:"name"`
	KTP       string `json:"ktp"`
}

type Payment struct {
	ID        int64  `json:"id"`
	BookingID int64  `json:"booking_id"`
	Amount    int64  `json:"amount"`
	Notes     string `json:"notes"`
	ProofFile string `json:"proof_file"`
	CreatedAt string `json:"created_at"`
}
