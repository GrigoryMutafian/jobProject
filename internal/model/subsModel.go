package model

type Subscription struct {
	Service   string  `json:"service"`
	Price     int     `json:"price"`
	UserID    string  `json:"user_id"`
	StartDate string  `json:"start_date"`
	EndDate   *string `json:"end_date,omitempty"`
}
