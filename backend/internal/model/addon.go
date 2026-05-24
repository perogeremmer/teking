package model

type AddonTemplate struct {
	ID         int64  `json:"id"`
	OperatorID int64  `json:"operator_id"`
	Name       string `json:"name"`
	Price      int64  `json:"price"`
	Icon       string `json:"icon"`
}
