package memorystorage

import (
	"context"
	"sync"
	"time"

	"github.com/romakorinenko/homework-test/hw12_13_14_15_calendar/internal/storage"
)

type Storage struct {
	store map[int64]*storage.Event
	mutex sync.Mutex
}

func NewStorage(store map[int64]*storage.Event) *Storage {
	return &Storage{store: store, mutex: sync.Mutex{}}
}

func (s *Storage) Create(_ context.Context, event *storage.Event) (*storage.Event, error) {
	s.mutex.Lock()
	defer s.mutex.Unlock()
	ID := len(s.store) + 1
	event.ID = int64(ID)

	s.store[event.ID] = event
	return event, nil
}

func (s *Storage) Update(_ context.Context, event *storage.Event) error {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	s.store[event.ID] = event
	return nil
}

func (s *Storage) DeleteByID(_ context.Context, id int64) error {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	delete(s.store, id)
	return nil
}

func (s *Storage) GetByPeriod(_ context.Context, startDate, endDate time.Time) ([]storage.Event, error) {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	res := make([]storage.Event, 0)
	for _, event := range s.store {
		if (event.StartDate.Before(startDate) || event.StartDate == startDate) &&
			(event.EndDate.After(endDate) || event.EndDate == endDate) {
			res = append(res, *event)
		}
	}
	return res, nil
}

func (s *Storage) Close() error {
	return nil
}
