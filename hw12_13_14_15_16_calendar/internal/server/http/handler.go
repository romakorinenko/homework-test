package internalhttp

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/romakorinenko/homework-test/hw12_13_14_15_calendar/internal/app"
	"github.com/romakorinenko/homework-test/hw12_13_14_15_calendar/internal/storage"
)

type EventHandler struct {
	app *app.App
}

func NewEventHandler(app *app.App) *EventHandler {
	return &EventHandler{app: app}
}

func (h *EventHandler) Create(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	defer func() { _ = r.Body.Close() }()

	data, err := io.ReadAll(r.Body)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		_, _ = fmt.Fprintf(w, `{"error": "%s"}`, err.Error())

		return
	}

	var event storage.Event
	err = json.Unmarshal(data, &event)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		_, _ = fmt.Fprintf(w, `{"error": "%s"}`, err.Error())

		return
	}

	if event.Title == "" || event.UserID == 0 {
		w.WriteHeader(http.StatusBadRequest)
		_, _ = fmt.Fprintf(w, `{"error": "%s"}`, "incomplete event")

		return
	}

	now := time.Now()
	marshal, err := json.Marshal(now.String())
	if err != nil {
		fmt.Println("11111")
	}
	fmt.Println(string(marshal))

	createdEvent, err := h.app.CreateEvent(r.Context(), &event)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		_, _ = fmt.Fprintf(w, `{"error": "%s"}`, err.Error())

		return
	}

	w.WriteHeader(http.StatusCreated)
	err = json.NewEncoder(w).Encode(createdEvent)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		_, _ = fmt.Fprintf(w, `{"error": "%s"}`, err.Error())

		return
	}
}

func (h *EventHandler) Update(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	defer func() { _ = r.Body.Close() }()

	data, err := io.ReadAll(r.Body)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		_, _ = fmt.Fprintf(w, `{"error": "%s"}`, err.Error())

		return
	}

	var event storage.Event
	err = json.Unmarshal(data, &event)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		_, _ = fmt.Fprintf(w, `{"error": "%s"}`, err.Error())

		return
	}

	err = h.app.UpdateEvent(r.Context(), &event)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		_, _ = fmt.Fprintf(w, `{"error": "%s"}`, err.Error())

		return
	}

	w.WriteHeader(http.StatusOK)
}

func (h *EventHandler) Delete(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	defer func() { _ = r.Body.Close() }()

	data, err := io.ReadAll(r.Body)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		_, _ = fmt.Fprintf(w, `{"error": "%s"}`, err.Error())

		return
	}

	var event storage.Event
	if err := json.Unmarshal(data, &event); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		_, _ = fmt.Fprintf(w, `{"error": "%s"}`, err.Error())

		return
	}

	if event.ID == 0 {
		w.WriteHeader(http.StatusBadRequest)
		_, _ = fmt.Fprintf(w, `{"error": "%s"}`, "id is 0")

		return
	}

	if err = h.app.DeleteEvent(r.Context(), event.ID); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		_, _ = fmt.Fprintf(w, `{"error": "%s"}`, err.Error())

		return
	}

	w.WriteHeader(http.StatusOK)
}

func (h *EventHandler) ListByDay(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	defer func() { _ = r.Body.Close() }()

	data, err := io.ReadAll(r.Body)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		_, _ = fmt.Fprintf(w, `{"error": "%s"}`, err.Error())

		return
	}

	var event storage.Event
	if err := json.Unmarshal(data, &event); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		_, _ = fmt.Fprintf(w, `{"error": "%s"}`, err.Error())

		return
	}

	if event.StartDate.IsZero() {
		w.WriteHeader(http.StatusBadRequest)
		_, _ = fmt.Fprintf(w, `{"error": "%s"}`, "start_date is required")

		return
	}

	endDate := event.StartDate.Add(24 * time.Hour)
	list, err := h.app.GetByPeriod(r.Context(), event.StartDate, endDate)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		_, _ = fmt.Fprintf(w, `{"error": "%s"}`, err.Error())
	}

	if err := json.NewEncoder(w).Encode(list); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		_, _ = fmt.Fprintf(w, `{"error": "%s"}`, err.Error())

		return
	}

	w.WriteHeader(http.StatusOK)
}

func (h *EventHandler) ListByWeek(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	defer func() { _ = r.Body.Close() }()

	data, err := io.ReadAll(r.Body)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		_, _ = fmt.Fprintf(w, `{"error": "%s"}`, err.Error())

		return
	}

	var event storage.Event
	if err := json.Unmarshal(data, &event); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		_, _ = fmt.Fprintf(w, `{"error": "%s"}`, err.Error())

		return
	}

	if event.StartDate.IsZero() {
		w.WriteHeader(http.StatusBadRequest)
		_, _ = fmt.Fprintf(w, `{"error": "%s"}`, "start_date is required")

		return
	}

	endDate := event.StartDate.Add(7 * 24 * time.Hour)
	list, err := h.app.GetByPeriod(r.Context(), event.StartDate, endDate)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		_, _ = fmt.Fprintf(w, `{"error": "%s"}`, err.Error())
	}

	if err := json.NewEncoder(w).Encode(list); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		_, _ = fmt.Fprintf(w, `{"error": "%s"}`, err.Error())

		return
	}

	w.WriteHeader(http.StatusOK)
}

func (h *EventHandler) ListByMonth(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	defer func() { _ = r.Body.Close() }()

	data, err := io.ReadAll(r.Body)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		_, _ = fmt.Fprintf(w, `{"error": "%s"}`, err.Error())

		return
	}

	var event storage.Event
	if err := json.Unmarshal(data, &event); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		_, _ = fmt.Fprintf(w, `{"error": "%s"}`, err.Error())

		return
	}

	if event.StartDate.IsZero() {
		w.WriteHeader(http.StatusBadRequest)
		_, _ = fmt.Fprintf(w, `{"error": "%s"}`, "start_date is required")

		return
	}

	endDate := event.StartDate.AddDate(0, 1, 0)
	list, err := h.app.GetByPeriod(r.Context(), event.StartDate, endDate)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		_, _ = fmt.Fprintf(w, `{"error": "%s"}`, err.Error())
	}

	if err := json.NewEncoder(w).Encode(list); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		_, _ = fmt.Fprintf(w, `{"error": "%s"}`, err.Error())

		return
	}

	w.WriteHeader(http.StatusOK)
}
