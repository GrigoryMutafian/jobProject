package model

import "time"

type Subscription struct {
	ID        int     `json:"id,omitempty"`
	Service   string  `json:"service"`
	Price     int     `json:"price"`
	UserID    string  `json:"user_id"`
	StartDate string  `json:"start_date"`
	EndDate   *string `json:"end_date,omitempty"`
}

type SubscriptionDB struct {
	ID        int `json:"id"`
	Service   string
	Price     int
	UserID    string
	StartDate time.Time
	EndDate   *time.Time
}
