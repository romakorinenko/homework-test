package dto

import "time"

type Notification struct {
	ID        int64     `json:"id"`
	Title     string    `json:"title"`
	StartDate time.Time `json:"startDate"`
	UserID    int64     `json:"userId"`
}
