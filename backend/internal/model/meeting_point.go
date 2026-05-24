package model

type MeetingPoint struct {
	ID         int64   `json:"id"`
	OperatorID int64   `json:"operator_id"`
	Type       string  `json:"type"`
	Name       string  `json:"name"`
	Address    string  `json:"address"`
	Lat        float64 `json:"lat"`
	Lng        float64 `json:"lng"`
}

type TripMeetingPoint struct {
	TripID            int64  `json:"trip_id"`
	MeetingPointID    int64  `json:"meeting_point_id"`
	OrderIndex        int    `json:"order_index"`
	EstimatedDeparture string `json:"estimated_departure"`
}

type TripPackagePrice struct {
	ID              int64 `json:"id"`
	TripID          int64 `json:"trip_id"`
	MeetingPointID  int64 `json:"meeting_point_id"`
	PackageID       int64 `json:"package_id"`
	Price           int64 `json:"price"`
}
