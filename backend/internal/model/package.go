package model

type Package struct {
	ID          int64   `json:"id"`
	OperatorID  int64   `json:"operator_id"`
	Name        string  `json:"name"`
	Description string  `json:"description"`
	FacilitiesJSON string `json:"facilities"`
}

type Facility struct {
	Name     string `json:"name"`
	Detail   string `json:"detail"`
	Category string `json:"category,omitempty"`
	Type     string `json:"type"`
}
