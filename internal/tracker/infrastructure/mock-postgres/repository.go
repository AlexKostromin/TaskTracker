package mock_postgres

import (
	"context"
	"sync"

	trackerV1 "github.com/AlexKostromin/TaskTracker/shared/pkg/openapi/tracker/v1"
	"github.com/google/uuid"
)

type TrackerStorage struct {
	mu       sync.RWMutex
	trackers map[uuid.UUID]*trackerV1.TrackerResponse
}

func NewTrackerStorage() *TrackerStorage {
	return &TrackerStorage{
		trackers: make(map[uuid.UUID]*trackerV1.TrackerResponse),
	}
}

func (s *TrackerStorage) CreateTracker(_ context.Context, request *trackerV1.CreateTrackerRequest) (trackerV1.CreateTrackerRes, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	return request.Name, nil
}

func (s *TrackerStorage) UpdateTracker(_ context.Context, request *trackerV1.UpdateTrackerRequest, params trackerV1.UpdateTrackerParams) (*trackerV1.TrackerResponse, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	return
}
