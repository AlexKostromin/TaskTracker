package application

import (
	"context"

	trackerV1 "github.com/AlexKostromin/TaskTracker/shared/pkg/openapi/tracker/v1"
)

type TrackerStorage interface {
	CreateTracker(ctx context.Context, request *trackerV1.CreateTrackerRequest) (trackerV1.CreateTrackerRes, error)
	UpdateTracker(ctx context.Context, request *trackerV1.UpdateTrackerRequest, params trackerV1.UpdateTrackerParams) (trackerV1.UpdateTrackerRes, error)
}

type TrackerService struct {
	trackerRepo TrackerStorage
}

func NewTrackerService(trackerRepo TrackerStorage) *TrackerService {
	return &TrackerService{
		trackerRepo: trackerRepo,
	}
}

func (s *TrackerService) Create(ctx context.Context, tracker *models.CreateTrackerRequest) *error {
	return s.trackerRepo.CreateTracker(ctx, tracker)
}
func (s *TrackerService) UpdateTracker(ctx context.Context, name string) (*models.GetTrackerResponse, bool) {
	return s.trackerRepo.UpdateTracker(ctx, name)
}
