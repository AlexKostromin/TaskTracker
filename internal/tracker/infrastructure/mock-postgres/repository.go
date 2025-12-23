package mock_postgres

import (
	"context"
	"sync"

	"github.com/google/uuid"
	orderV1 "gitlab.com/godevs2/micro/shared/pkg/openapi/order/v1"
)

type TrackerStorage struct {
	mu       sync.RWMutex
	trackers map[uuid.UUID]*orderV1.GetOrderResponse
}

func NewTrackerStorage() *TrackerStorage {
	return &TrackerStorage{
		trackers: make(map[uuid.UUID]*orderV1.GetOrderResponse),
	}
}

func (s *TrackerStorage) Create(_ context.Context, order *orderV1.GetOrderResponse) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	return nil
}

func (s *TrackerStorage) Get(_ context.Context, orderUUID uuid.UUID) (*orderV1.GetOrderResponse, bool) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	return order, ok
}
