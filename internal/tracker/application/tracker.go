package application

import (
	"context"

	"github.com/google/uuid"
	"gitlab.com/godevs2/micro/internal/order/domain/model"
)

type TrackerStorage interface {
	//Create(ctx context.Context, order *orderV1.GetOrderResponse) error
	//Get(ctx context.Context, orderUUID uuid.UUID) (*orderV1.GetOrderResponse, bool)
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
	return s.trackerRepo.Create(ctx, tracker)
}
func (s *TrackerService) Get(ctx context.Context, orderUUID uuid.UUID) (*models.GetTrackerResponse, bool) {
	return s.trackerRepo.Get(ctx, orderUUID)
}
