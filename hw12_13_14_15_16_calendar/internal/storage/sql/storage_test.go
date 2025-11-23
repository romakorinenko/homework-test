package sqlstorage

import (
	"context"
	"testing"
	"time"

	storage2 "github.com/romakorinenko/homework-test/hw12_13_14_15_calendar/internal/storage"
	test "github.com/romakorinenko/homework-test/hw12_13_14_15_calendar/tests"
	"github.com/stretchr/testify/require"
)

func TestStorage_CRUD(t *testing.T) {
	ctx := context.Background()
	testDB := test.CreateDBForTest(t, "/migrations")

	start := time.Now().Add(-1 * time.Hour)
	end := time.Now()
	storage := Storage{testDB.DBPool}
	event, err := storage.Create(ctx, &storage2.Event{
		Title:     "Test title",
		StartDate: start,
		EndDate:   end,
		UserID:    1,
	})
	require.NoError(t, err)
	require.Equal(t, "Test title", event.Title)

	event.Title = "Test title 2"
	err = storage.Update(ctx, event)
	require.NoError(t, err)

	byPeriod, err := storage.GetByPeriod(ctx, start, end)
	require.NoError(t, err)
	require.Len(t, byPeriod, 1)
	require.Equal(t, byPeriod[0].ID, event.ID)
	require.Equal(t, byPeriod[0].Title, event.Title)

	err = storage.DeleteByID(ctx, event.ID)
	require.NoError(t, err)

	eventsByPeriod, err := storage.GetByPeriod(ctx, start, end)
	require.NoError(t, err)
	require.Len(t, eventsByPeriod, 0)
}
