package http_server

import (
	"context"
	"errors"

	"github.com/AlexKostromin/TaskTracker/internal/tracker/infrastructure/http-server/converter"
	trackerV1 "github.com/AlexKostromin/TaskTracker/shared/pkg/openapi/tracker/v1"
)

func (s *Server) CreateTracker(ctx context.Context, request *trackerV1.CreateTrackerRequest) (trackerV1.CreateTrackerRes, error) {

	if request == nil {
		return nil, errors.New("request is nil")
	}
	tracker, err := s.storage.CreateTracker(ctx, converter.CreateTrackerRequestToModel(request))
	if err != nil {
		return nil, err
	}

	return &trackerV1.TrackerResponse{
		ID:          trackerV1.NewOptInt(tracker.ID),
		Name:        trackerV1.NewOptString(tracker.Name),
		Description: trackerV1.NewOptString(tracker.Description),
	}, nil
}
