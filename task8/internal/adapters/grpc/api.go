package grpc

import (
	"context"
	"github.com/golang/protobuf/ptypes"
	"github.com/golang/protobuf/ptypes/empty"
	"github.com/golang/protobuf/ptypes/timestamp"
	"github.com/rendau/my-otus/task8/internal/domain/entities"
	"github.com/rendau/my-otus/task8/internal/domain/errors"
	"github.com/rendau/my-otus/task8/internal/domain/usecases"
	"github.com/rendau/my-otus/task8/internal/interfaces"
	"google.golang.org/grpc"
	"net"
	"time"
)

// API - is type for grpc API adapter
type API struct {
	log    interfaces.Logger
	lAddr  string
	server *grpc.Server
	ucs    *usecases.Usecases
}

// CreateAPI - creates new instance
func CreateAPI(log interfaces.Logger, lAddr string, ucs *usecases.Usecases) *API {
	return &API{
		log:   log,
		lAddr: lAddr,
		ucs:   ucs,
	}
}

// Start - starts api-server
func (a *API) Start() {
	a.server = grpc.NewServer()

	RegisterCalendarServiceServer(a.server, a)

	go func() {
		l, err := net.Listen("tcp", a.lAddr)
		if err != nil {
			a.log.Fatalw("Fail to start grpc server", "error", err)
		}
		err = a.server.Serve(l)
		if err != nil {
			a.log.Fatalw("grpc server stopped", "error", err)
		}
	}()
}

// Shutdown - stops api-server
func (a *API) Shutdown() {
	a.server.GracefulStop()
}

// CreateEvent - method for implement proto-spec, creates event
func (a *API) CreateEvent(ctx context.Context, r *CreateEventRequest) (*CreateEventResponse, error) {
	startTime, err := ptypes.Timestamp(r.StartTime)
	if err != nil {
		return nil, err
	}
	endTime, err := ptypes.Timestamp(r.EndTime)
	if err != nil {
		return nil, err
	}

	event, err := a.ucs.Event.Create(ctx, r.Owner, r.Title, r.Text, startTime, endTime)
	if err != nil {
		switch err.(type) {
		case errors.EventError:
			return &CreateEventResponse{Result: &CreateEventResponse_Error{Error: err.Error()}}, nil
		default:
			return nil, err
		}
	}

	return &CreateEventResponse{Result: &CreateEventResponse_Id{Id: event.ID}}, nil
}

// UpdateEvent - method for implement proto-spec, updates event
func (a *API) UpdateEvent(ctx context.Context, r *UpdateEventRequest) (*UpdateEventResponse, error) {
	startTime, err := ptypes.Timestamp(r.StartTime)
	if err != nil {
		return nil, err
	}
	endTime, err := ptypes.Timestamp(r.EndTime)
	if err != nil {
		return nil, err
	}

	err = a.ucs.Event.Update(
		ctx,
		r.Id,
		r.Owner,
		r.Title,
		r.Text,
		startTime,
		endTime,
	)
	if err != nil {
		switch err.(type) {
		case errors.EventError:
			return &UpdateEventResponse{Error: err.Error()}, nil
		default:
			return nil, err
		}
	}

	return &UpdateEventResponse{}, nil
}

// DeleteEvent - method for implement proto-spec, deletes event
func (a *API) DeleteEvent(ctx context.Context, r *DeleteEventRequest) (*empty.Empty, error) {
	err := a.ucs.Event.Delete(ctx, r.Id)
	if err != nil {
		return nil, err
	}

	return &empty.Empty{}, nil
}

// ListEventForDay - method for implement proto-spec, get list of events for day
func (a *API) ListEventForDay(ctx context.Context, r *empty.Empty) (*ListEventResponse, error) {
	startTime := time.Now().Truncate(24 * time.Hour)
	endTime := startTime.Add(24 * time.Hour)

	return a.listEventForFilter(ctx, &entities.EventListFilter{
		StartTimeGt: &startTime,
		StartTimeLt: &endTime,
	})
}

// ListEventForWeek - method for implement proto-spec, get list of events for week
func (a *API) ListEventForWeek(ctx context.Context, r *empty.Empty) (*ListEventResponse, error) {
	startTime := time.Now().Truncate(24 * time.Hour)
	endTime := startTime.Add(7 * 24 * time.Hour)

	return a.listEventForFilter(ctx, &entities.EventListFilter{
		StartTimeGt: &startTime,
		StartTimeLt: &endTime,
	})
}

// ListEventForMonth - method for implement proto-spec, get list of events for month
func (a *API) ListEventForMonth(ctx context.Context, r *empty.Empty) (*ListEventResponse, error) {
	startTime := time.Now().Truncate(24 * time.Hour)
	endTime := startTime.Add(30 * 24 * time.Hour)

	return a.listEventForFilter(ctx, &entities.EventListFilter{
		StartTimeGt: &startTime,
		StartTimeLt: &endTime,
	})
}

func (a *API) listEventForFilter(ctx context.Context, filter *entities.EventListFilter) (*ListEventResponse, error) {
	events, err := a.ucs.Event.List(ctx, filter)
	if err != nil {
		return nil, err
	}

	var pbEvents []*Event
	var pbStartTime *timestamp.Timestamp
	var pbEndTime *timestamp.Timestamp
	for _, e := range events {
		pbStartTime, err = ptypes.TimestampProto(e.StartTime)
		if err != nil {
			return nil, err
		}
		pbEndTime, err = ptypes.TimestampProto(e.EndTime)
		if err != nil {
			return nil, err
		}
		pbEvents = append(pbEvents, &Event{
			Id:        e.ID,
			Owner:     e.Owner,
			Title:     e.Title,
			Text:      e.Text,
			StartTime: pbStartTime,
			EndTime:   pbEndTime,
		})
	}

	return &ListEventResponse{Events: pbEvents}, nil
}
