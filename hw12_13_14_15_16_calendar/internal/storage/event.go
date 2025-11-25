package storage

import "time"

type Event struct {
	ID        int64     `db:"id" json:"id"`
	Title     string    `db:"title" json:"title"`
	StartDate time.Time `db:"start_date" json:"startDate"`
	EndDate   time.Time `db:"end_date" json:"endDate"`
	UserID    int64     `db:"user_id" json:"userId"`
}
