package storage

import "time"

type Event struct {
	ID        int64     `db:"id"`
	Title     string    `db:"title"`
	StartDate time.Time `db:"start_date"`
	EndDate   time.Time `db:"end_date"`
	UserID    int64     `db:"user_id"`
}
