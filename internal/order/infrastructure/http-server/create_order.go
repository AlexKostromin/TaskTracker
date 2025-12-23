package http_server

import (
	"context"
	"net/http"

	"github.com/google/uuid"
	orderV1 "gitlab.com/godevs2/micro/shared/pkg/openapi/order/v1"
)

func (s *Server) CreateOrder(ctx context.Context, orderRes *orderV1.CreateOrderRequest) (orderV1.CreateOrderRes, error) {
	newUuid := uuid.New()

	totalPrice := float64(len(orderRes.PartUuids)) * 100

	order := &orderV1.GetOrderResponse{
		OrderUUID:       newUuid,
		UserUUID:        orderRes.UserUUID,
		PartUuids:       orderRes.PartUuids,
		TotalPrice:      totalPrice,
		TransactionUUID: orderV1.OptNilUUID{},
		PaymentMethod:   orderV1.OptNilString{},
		Status: orderV1.OrderStatus{
			OrderStatus: orderV1.NewOptOrderStatusOrderStatus(orderV1.OrderStatusOrderStatusPENDINGPAYMENT),
		},
	}

	err := s.storage.Create(ctx, order)
	if err != nil {
		// log

		return nil, http.ErrServerClosed
	}

	return &orderV1.CreateOrderResponse{
		OrderUUID:  newUuid,
		TotalPrice: totalPrice,
	}, nil
}
