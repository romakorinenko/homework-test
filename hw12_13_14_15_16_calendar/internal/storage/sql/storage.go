package sqlstorage

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/huandu/go-sqlbuilder"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/romakorinenko/homework-test/hw12_13_14_15_calendar/internal/storage"
)

const (
	EventsSequenceName = "events_sequence"
	EventsTableName    = "events"
)

var Struct = sqlbuilder.NewStruct(new(storage.Event))

type Storage struct {
	DBPool *pgxpool.Pool
}

func NewStorage(ctx context.Context, dbString string) *Storage {
	dbConfig, err := pgxpool.ParseConfig(dbString)
	if err != nil {
		panic(err)
	}
	location, err := time.LoadLocation("Local")
	if err != nil {
		panic(err)
	}
	dbConfig.AfterConnect = func(_ context.Context, conn *pgx.Conn) error {
		conn.TypeMap().RegisterType(&pgtype.Type{
			Name:  "timestamp",
			OID:   pgtype.TimestampOID,
			Codec: &pgtype.TimestampCodec{ScanLocation: location},
		})
		return nil
	}

	dbPool, err := pgxpool.NewWithConfig(ctx, dbConfig)
	if err != nil {
		panic(err)
	}
	return &Storage{
		DBPool: dbPool,
	}
}

func (s *Storage) Create(ctx context.Context, event *storage.Event) (*storage.Event, error) {
	newID, err := s.generateNextEventID(ctx)
	if err != nil {
		return nil, err
	}
	event.ID = newID

	sqlQuery, args := Struct.
		InsertInto(EventsTableName, event).
		BuildWithFlavor(sqlbuilder.PostgreSQL)
	_, err = s.DBPool.Exec(ctx, sqlQuery, args...)
	if err != nil {
		return nil, err
	}

	return event, nil
}

func (s *Storage) Update(ctx context.Context, event *storage.Event) error {
	ub := sqlbuilder.Update(EventsTableName)
	sql, args := ub.Where(ub.Equal("id", event.ID)).
		Set(
			ub.Assign("title", event.Title),
			ub.Assign("start_date", event.StartDate),
			ub.Assign("end_date", event.EndDate),
			ub.Assign("user_id", event.UserID),
		).
		BuildWithFlavor(sqlbuilder.PostgreSQL)

	_, err := s.DBPool.Exec(ctx, sql, args...)
	return err
}

func (s *Storage) DeleteByID(ctx context.Context, id int64) error {
	deleteBuilder := Struct.DeleteFrom(EventsTableName)
	sql, args := deleteBuilder.
		Where(deleteBuilder.Equal("id", id)).
		BuildWithFlavor(sqlbuilder.PostgreSQL)

	_, err := s.DBPool.Exec(ctx, sql, args...)
	return err
}

func (s *Storage) GetByPeriod(ctx context.Context, startData, endData time.Time) ([]storage.Event, error) {
	selectBuilder := Struct.SelectFrom(EventsTableName)
	sql, args := selectBuilder.
		Where(
			selectBuilder.And(
				selectBuilder.GTE("start_date", startData),
				selectBuilder.LTE("end_date", endData),
			),
		).
		BuildWithFlavor(sqlbuilder.PostgreSQL)

	rows, err := s.DBPool.Query(ctx, sql, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	res := make([]storage.Event, 0)
	for rows.Next() {
		var event storage.Event
		err = rows.Scan(Struct.Addr(&event)...)
		if err != nil {
			return nil, err
		}
		res = append(res, event)
	}
	return res, nil
}

func (s *Storage) generateNextEventID(ctx context.Context) (int64, error) {
	rows, err := s.DBPool.Query(ctx, fmt.Sprintf("SELECT nextval('%s')", EventsSequenceName))
	if err != nil {
		return 0, err
	}
	defer rows.Close()

	if rows.Next() {
		var id int64
		err = rows.Scan(&id)
		if err != nil {
			return 0, err
		}
		return id, nil
	}
	return 0, errors.New("something was wrong. there is no next event id")
}

func (s *Storage) Close() error {
	s.DBPool.Close()
	return nil
}
