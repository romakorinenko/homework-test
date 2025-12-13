package grpcinternal

import (
	"context"
	"time"

	"github.com/romakorinenko/homework-test/hw12_13_14_15_calendar/internal/app"
	"github.com/romakorinenko/homework-test/hw12_13_14_15_calendar/internal/server/grpc/pb"
	"github.com/romakorinenko/homework-test/hw12_13_14_15_calendar/internal/storage"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type EventService struct {
	app *app.App
	pb.UnimplementedEventServiceServer
}

func (h *EventService) Create(ctx context.Context, ev *pb.Event) (*pb.Event, error) {
	event := &storage.Event{
		Title:     ev.Title,
		StartDate: ev.StartDate.AsTime().Local(),
		EndDate:   ev.EndDate.AsTime().Local(),
		UserID:    ev.UserId,
	}
	if event.Title == "" || event.UserID == 0 {
		return nil, status.Error(codes.InvalidArgument, "failed to create event")
	}

	createdEvent, err := h.app.CreateEvent(ctx, event)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &pb.Event{
		Id:        createdEvent.ID,
		Title:     createdEvent.Title,
		StartDate: timestamppb.New(createdEvent.StartDate),
		EndDate:   timestamppb.New(createdEvent.EndDate),
		UserId:    createdEvent.UserID,
	}, nil
}

func (h *EventService) Delete(ctx context.Context, ev *pb.Event) (*pb.EmptyEventResponse, error) {
	if ev.Id == 0 {
		return nil, status.Error(codes.InvalidArgument, "incorrect event with id 0")
	}

	err := h.app.DeleteEvent(ctx, ev.Id)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &pb.EmptyEventResponse{}, nil
}

func (h *EventService) Update(ctx context.Context, ev *pb.Event) (*pb.EmptyEventResponse, error) {
	err := h.app.UpdateEvent(ctx, &storage.Event{
		Title:     ev.Title,
		StartDate: ev.StartDate.AsTime(),
		EndDate:   ev.EndDate.AsTime(),
		UserID:    ev.UserId,
	})
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &pb.EmptyEventResponse{}, nil
}

func (h *EventService) ListByDay(ctx context.Context, ev *pb.Event) (*pb.ListEventResponse, error) {
	startDate := ev.StartDate.AsTime()
	if startDate.IsZero() {
		return nil, status.Error(codes.InvalidArgument, "incorrect start date")
	}

	events, err := h.app.GetByPeriod(ctx, startDate, startDate.Add(24*time.Hour))
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	res := make([]*pb.Event, 0)
	for _, event := range events {
		res = append(res, &pb.Event{
			Id:        event.ID,
			Title:     event.Title,
			StartDate: timestamppb.New(event.StartDate),
			EndDate:   timestamppb.New(event.EndDate),
			UserId:    event.UserID,
		})
	}
	return &pb.ListEventResponse{Events: res}, nil
}

func (h *EventService) ListByWeek(ctx context.Context, ev *pb.Event) (*pb.ListEventResponse, error) {
	startDate := ev.StartDate.AsTime()
	if startDate.IsZero() {
		return nil, status.Error(codes.InvalidArgument, "incorrect start date")
	}

	events, err := h.app.GetByPeriod(ctx, startDate, startDate.Add(7*24*time.Hour))
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	res := make([]*pb.Event, 0)
	for _, event := range events {
		res = append(res, &pb.Event{
			Id:        event.ID,
			Title:     event.Title,
			StartDate: timestamppb.New(event.StartDate),
			EndDate:   timestamppb.New(event.EndDate),
			UserId:    event.UserID,
		})
	}
	return &pb.ListEventResponse{Events: res}, nil
}

func (h *EventService) ListByMonth(ctx context.Context, ev *pb.Event) (*pb.ListEventResponse, error) {
	startDate := ev.StartDate.AsTime()
	if startDate.IsZero() {
		return nil, status.Error(codes.InvalidArgument, "incorrect start date")
	}

	events, err := h.app.GetByPeriod(ctx, startDate, startDate.AddDate(0, 1, 0))
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	res := make([]*pb.Event, 0)
	for _, event := range events {
		res = append(res, &pb.Event{
			Id:        event.ID,
			Title:     event.Title,
			StartDate: timestamppb.New(event.StartDate),
			EndDate:   timestamppb.New(event.EndDate),
			UserId:    event.UserID,
		})
	}
	return &pb.ListEventResponse{Events: res}, nil
}
