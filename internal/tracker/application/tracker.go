package application

import (
	"context"

	models "github.com/AlexKostromin/TaskTracker/internal/tracker/domain"
)

type TrackerStorage interface {
	CreateTracker(ctx context.Context, request models.CreateTrackerRequest) (models.Tracker, error)
	UpdateTracker(ctx context.Context, request models.UpdateTrackerRequest, params models.UpdateTrackerParams) (models.Tracker, error)
}

type TrackerService struct {
	trackerRepo TrackerStorage
}

func NewTrackerService(trackerRepo TrackerStorage) *TrackerService {
	return &TrackerService{
		trackerRepo: trackerRepo,
	}
}

func (s *TrackerService) CreateTracker(ctx context.Context, tracker models.CreateTrackerRequest) (models.Tracker, error) {
	return s.trackerRepo.CreateTracker(ctx, tracker)
}

func (s *TrackerService) UpdateTracker(ctx context.Context, request models.UpdateTrackerRequest, params models.UpdateTrackerParams) (models.Tracker, error) {
	return s.trackerRepo.UpdateTracker(ctx, request, params)
}
