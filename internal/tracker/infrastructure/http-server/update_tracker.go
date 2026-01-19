package http_server

import (
	"context"
	"log"
	"net/http"

	"github.com/AlexKostromin/TaskTracker/internal/tracker/infrastructure/http-server/converter"
	trackerV1 "github.com/AlexKostromin/TaskTracker/shared/pkg/openapi/tracker/v1"
)

func (s *Server) UpdateTracker(ctx context.Context, request *trackerV1.UpdateTrackerRequest, params trackerV1.UpdateTrackerParams) (trackerV1.UpdateTrackerRes, error) {
	log.Println("UpdateTracker handler called") // Добавьте эту строку
	if request == nil {
		return &trackerV1.BadRequestError{
			Code:    http.StatusBadRequest,
			Message: "request is nil",
		}, nil
	}

	tracker, err := s.storage.UpdateTracker(ctx, converter.UpdateTrackerRequestToModel(request), converter.UpdateTrackerParamsToModel(params))
	if err != nil {
		return &trackerV1.InternalServerError{
			Code:    http.StatusInternalServerError,
			Message: err.Error(),
		}, nil
	}

	return &trackerV1.TrackerResponse{
		ID:          trackerV1.NewOptInt(tracker.ID),
		Name:        trackerV1.NewOptString(tracker.Name),
		Description: trackerV1.NewOptString(tracker.Description),
	}, nil
}

func (s *Server) NewError(ctx context.Context, err error) *trackerV1.GenericErrorStatusCode {
	return &trackerV1.GenericErrorStatusCode{
		StatusCode: http.StatusInternalServerError,
		Response: trackerV1.GenericError{
			Code:    trackerV1.NewOptInt(http.StatusInternalServerError),
			Message: trackerV1.NewOptString(err.Error()),
		},
	}
}
