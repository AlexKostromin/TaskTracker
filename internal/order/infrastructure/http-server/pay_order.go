package http_server

import (
	"context"
	"strconv"

	"github.com/google/uuid"
	orderV1 "gitlab.com/godevs2/micro/shared/pkg/openapi/order/v1"
)

func (s *Server) PayOrder(ctx context.Context, req *orderV1.PayOrderRequest, params orderV1.PayOrderParams) (orderV1.PayOrderRes, error) {
	order, ok := s.storage.Get(ctx, params.OrderUUID)
	if !ok {
		return &orderV1.NotFoundError{
			Error:   orderV1.NewOptString("NOT_FOUND"),
			Code:    orderV1.NewOptInt(404),
			Message: orderV1.NewOptString("Order not found"),
		}, nil
	}

	if order.Status.OrderStatus.Value == orderV1.OrderStatusOrderStatusPAID {
		return &orderV1.ConflictError{
			Error:   orderV1.NewOptString("CONFLICT"),
			Code:    orderV1.NewOptInt(409),
			Message: orderV1.NewOptString("Order already paid"),
		}, nil
	}

	newTransactionUuid := uuid.New()

	updatedOrder := &orderV1.GetOrderResponse{
		OrderUUID:  order.OrderUUID,
		UserUUID:   order.UserUUID,
		PartUuids:  order.PartUuids,
		TotalPrice: order.TotalPrice,
		TransactionUUID: orderV1.OptNilUUID{
			Value: newTransactionUuid,
			Set:   true,
		},
		PaymentMethod: orderV1.OptNilString{
			Value: strconv.Itoa(int(req.PaymentMethod.PaymentMethod.Value)), // Исправлено!
			Set:   true,
		},
		Status: orderV1.OrderStatus{
			OrderStatus: orderV1.NewOptOrderStatusOrderStatus(orderV1.OrderStatusOrderStatusPAID),
		},
	}

	err := s.storage.Create(ctx, updatedOrder)
	if err != nil {
		return &orderV1.InternalServerError{
			Error:   orderV1.NewOptString("INTERNAL_SERVER_ERROR"),
			Code:    orderV1.NewOptInt(500),
			Message: orderV1.NewOptString("Failed to save order"),
		}, nil
	}

	return &orderV1.PayOrderResponse{
		TransactionUUID: newTransactionUuid,
	}, nil
}
