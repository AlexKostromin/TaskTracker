package mock_postgres

import (
	"context"
	"sync"

	"github.com/google/uuid"
	orderV1 "gitlab.com/godevs2/micro/shared/pkg/openapi/order/v1"
)

type OrderStorage struct {
	mu     sync.RWMutex
	orders map[uuid.UUID]*orderV1.GetOrderResponse
}

func NewOrderStorage() *OrderStorage {
	return &OrderStorage{
		orders: make(map[uuid.UUID]*orderV1.GetOrderResponse),
	}
}

func (s *OrderStorage) Create(_ context.Context, order *orderV1.GetOrderResponse) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.orders[order.OrderUUID] = order

	return nil
}

func (s *OrderStorage) Get(_ context.Context, orderUUID uuid.UUID) (*orderV1.GetOrderResponse, bool) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	order, ok := s.orders[orderUUID]
	return order, ok
}

func CancelOrder(_ context.Context, params orderV1.CancelOrderParams) (orderV1.CancelOrderRes, error) {

}
