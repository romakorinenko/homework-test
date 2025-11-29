package internalhttp

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/romakorinenko/homework-test/hw12_13_14_15_calendar/internal/app"
	"github.com/romakorinenko/homework-test/hw12_13_14_15_calendar/internal/logger"
	"github.com/romakorinenko/homework-test/hw12_13_14_15_calendar/internal/storage"
	sqlstorage "github.com/romakorinenko/homework-test/hw12_13_14_15_calendar/internal/storage/sql"
	test "github.com/romakorinenko/homework-test/hw12_13_14_15_calendar/tests"
	"github.com/stretchr/testify/require"
)

func TestEventHandler_Create(t *testing.T) {
	ctx := context.Background()
	testDB := test.CreateDBForTest(t, "/migrations")

	application := app.New(logger.NewLogger(0), &sqlstorage.Storage{DBPool: testDB.DBPool})
	eventHandler := NewEventHandler(application)
	server := httptest.NewServer(http.HandlerFunc(eventHandler.Create))
	defer server.Close()

	event := &storage.Event{
		Title:     "Test title",
		StartDate: time.Now(),
		EndDate:   time.Now(),
		UserID:    1,
	}
	marshal, err := json.Marshal(event)
	require.Nil(t, err)

	req, err := http.NewRequestWithContext(ctx, "POST", server.URL+"/event/create", bytes.NewBuffer(marshal))
	require.Nil(t, err)

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(eventHandler.Create)

	handler.ServeHTTP(rr, req)

	require.Equal(t, http.StatusCreated, rr.Code)
}

func TestEventHandler_Update(t *testing.T) {
	ctx := context.Background()
	testDB := test.CreateDBForTest(t, "/migrations")

	application := app.New(logger.NewLogger(0), &sqlstorage.Storage{DBPool: testDB.DBPool})
	eventHandler := NewEventHandler(application)
	server := httptest.NewServer(http.HandlerFunc(eventHandler.Update))
	defer server.Close()

	event := &storage.Event{
		ID:        1,
		Title:     "Test title",
		StartDate: time.Now(),
		EndDate:   time.Now(),
		UserID:    1,
	}
	marshal, err := json.Marshal(event)
	require.Nil(t, err)

	req, err := http.NewRequestWithContext(ctx, "POST", server.URL+"/event/update", bytes.NewBuffer(marshal))
	require.Nil(t, err)

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(eventHandler.Delete)

	handler.ServeHTTP(rr, req)

	require.Equal(t, http.StatusOK, rr.Code)
}

func TestEventHandler_Delete(t *testing.T) {
	ctx := context.Background()
	testDB := test.CreateDBForTest(t, "/migrations")

	application := app.New(logger.NewLogger(0), &sqlstorage.Storage{DBPool: testDB.DBPool})
	eventHandler := NewEventHandler(application)
	server := httptest.NewServer(http.HandlerFunc(eventHandler.Delete))
	defer server.Close()

	event := &storage.Event{
		ID: 1,
	}
	marshal, err := json.Marshal(event)
	require.Nil(t, err)

	req, err := http.NewRequestWithContext(ctx, "DELETE", "/event/delete", bytes.NewBuffer(marshal))
	require.Nil(t, err)

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(eventHandler.Delete)

	handler.ServeHTTP(rr, req)

	require.Equal(t, http.StatusOK, rr.Code)
}

// TestEventHandler_GetByDay вообще метод GET не должен содержать тело, но из-за упрощения, было допущено такое.
func TestEventHandler_GetByDay(t *testing.T) {
	ctx := context.Background()
	testDB := test.CreateDBForTest(t, "/migrations")

	application := app.New(logger.NewLogger(0), &sqlstorage.Storage{DBPool: testDB.DBPool})
	eventHandler := NewEventHandler(application)
	server := httptest.NewServer(http.HandlerFunc(eventHandler.ListByDay))
	defer server.Close()

	event := &storage.Event{
		StartDate: time.Now(),
	}
	marshal, err := json.Marshal(event)
	require.Nil(t, err)

	req, err := http.NewRequestWithContext(ctx, "GET", server.URL+"/event/list/day", bytes.NewBuffer(marshal))
	require.Nil(t, err)

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(eventHandler.ListByDay)

	handler.ServeHTTP(rr, req)

	require.Equal(t, http.StatusOK, rr.Code)
}
