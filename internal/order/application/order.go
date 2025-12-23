package application

import (
	"context"

	"github.com/google/uuid"
	"gitlab.com/godevs2/micro/internal/order/domain/model"
	orderV1 "gitlab.com/godevs2/micro/shared/pkg/openapi/order/v1"
)

type OrderStorage interface {
	Create(ctx context.Context, order *orderV1.GetOrderResponse) error
	Get(ctx context.Context, orderUUID uuid.UUID) (*orderV1.GetOrderResponse, bool)
}

type OrderService struct {
	orderRepo OrderStorage
}

func NewOrderService(orderRepo OrderStorage) *OrderService {
	return &OrderService{
		orderRepo: orderRepo,
	}
}

func (s *OrderService) Create(ctx context.Context, order *models.CreateOrderRequest) *error {
	return s.orderRepo.Create(ctx, order)
}
func (s *OrderService) Get(ctx context.Context, orderUUID uuid.UUID) (*models.GetOrderResponse, bool) {
	return s.orderRepo.Get(ctx, orderUUID)
}
