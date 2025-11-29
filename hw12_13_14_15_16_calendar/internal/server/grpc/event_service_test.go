package grpcinternal

import (
	"context"
	"testing"
	"time"

	"github.com/romakorinenko/homework-test/hw12_13_14_15_calendar/internal/app"
	"github.com/romakorinenko/homework-test/hw12_13_14_15_calendar/internal/logger"
	"github.com/romakorinenko/homework-test/hw12_13_14_15_calendar/internal/server/grpc/pb"
	sqlstorage "github.com/romakorinenko/homework-test/hw12_13_14_15_calendar/internal/storage/sql"
	test "github.com/romakorinenko/homework-test/hw12_13_14_15_calendar/tests"
	"github.com/stretchr/testify/require"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func TestEventService_CRUD(t *testing.T) {
	dbForTest := test.CreateDBForTest(t, "/migrations")

	es := &EventService{
		app: app.NewApp(
			logger.NewLogger(0),
			&sqlstorage.Storage{DBPool: dbForTest.DBPool},
		),
	}

	now := time.Now()

	ctx := context.Background()
	event, err := es.Create(ctx, &pb.Event{
		Id:        1,
		Title:     "test",
		StartDate: timestamppb.New(now),
		EndDate:   timestamppb.New(now),
		UserId:    1,
	})
	require.NoError(t, err)

	event.Title = "test2"
	updated, err := es.Update(ctx, event)
	require.NoError(t, err)
	require.Empty(t, updated)

	byDay, err := es.ListByDay(ctx, event)
	require.NoError(t, err)
	require.Len(t, byDay.Events, 1)

	deleted, err := es.Delete(ctx, event)
	require.NoError(t, err)
	require.Empty(t, deleted)

	byDay, err = es.ListByDay(ctx, event)
	require.NoError(t, err)
	require.Len(t, byDay.Events, 0)
}
