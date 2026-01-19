package mock_postgres

import (
	"context"
	"errors"
	"math/rand"
	"sync"
	"time"

	"github.com/AlexKostromin/TaskTracker/internal/tracker/domain"
)

type TrackerStorage struct {
	mu       sync.RWMutex
	trackers map[int]models.Tracker
}

func NewTrackerStorage() *TrackerStorage {
	return &TrackerStorage{
		trackers: make(map[int]models.Tracker),
	}
}

func generateRandomIntID() int {
	// Инициализация генератора случайных чисел
	rand.New(rand.NewSource(time.Now().UnixNano()))

	// Генерируем число от 1000 до 999999 (или другой диапазон)
	return 1000 + rand.Intn(999000)
}

func (s *TrackerStorage) CreateTracker(ctx context.Context, tracker models.CreateTrackerRequest) (models.Tracker, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	id := generateRandomIntID()
	newTracker := models.Tracker{
		ID:          id,
		Name:        tracker.Name,
		Description: tracker.Description,
	}

	s.trackers[id] = newTracker

	return newTracker, nil
}

func (s *TrackerStorage) UpdateTracker(ctx context.Context, request models.UpdateTrackerRequest, params models.UpdateTrackerParams) (models.Tracker, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	trackerID := params.ID

	tracker, ok := s.trackers[trackerID]

	if !ok {
		return models.Tracker{}, errors.New("tracker not found")
	}

	if request.Name != "" {
		tracker.Name = request.Name
	}
	if request.Description != "" {
		tracker.Description = request.Description
	}

	s.trackers[trackerID] = tracker

	return models.Tracker{
		ID:          trackerID,
		Name:        tracker.Name,
		Description: tracker.Description,
	}, nil
}
