package http_server

import (
	"context"

	models "gitlab.com/godevs2/micro/internal/order/domain/model"
	orderV1 "gitlab.com/godevs2/micro/shared/pkg/openapi/order/v1"
)

func (s *Server) GetOrder(ctx context.Context, params models.GetOrderParams) (orderV1.GetOrderRes, error) {

	order, ok := s.storage.Get(ctx, params.OrderUuid)
	if !ok {
		return &orderV1.NotFoundError{
			Error:   orderV1.NewOptString("NOT_FOUND"),
			Code:    orderV1.NewOptInt(404),
			Message: orderV1.NewOptString("Order not found"),
		}, nil
	}

	return order, nil
}
