package http_server

import (
	"context"
	orderV1 "gitlab.com/godevs2/micro/shared/pkg/openapi/order/v1"
)

func (s *Server) CancelOrder(ctx context.Context, params orderV1.CancelOrderParams) (orderV1.CancelOrderRes, error) {

	order, ok := s.storage.GetOrder(ctx, params.OrderUUID)
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
			Message: orderV1.NewOptString("Cannot cancel paid order"),
		}, nil
	}

	if order.Status.OrderStatus.Value == orderV1.OrderStatusOrderStatusCANCELLED {
		return &orderV1.ConflictError{
			Error:   orderV1.NewOptString("CONFLICT"),
			Code:    orderV1.NewOptInt(409),
			Message: orderV1.NewOptString("Order already cancelled"),
		}, nil
	}

	updatedOrder := &orderV1.GetOrderResponse{
		OrderUUID:       order.OrderUUID,
		UserUUID:        order.UserUUID,
		PartUuids:       order.PartUuids,
		TotalPrice:      order.TotalPrice,
		TransactionUUID: order.TransactionUUID,
		PaymentMethod:   order.PaymentMethod,
		Status: orderV1.OrderStatus{
			OrderStatus: orderV1.NewOptOrderStatusOrderStatus(orderV1.OrderStatusOrderStatusCANCELLED),
		},
	}

	err := s.storage.CreateOrder(ctx, updatedOrder)
	if err != nil {
		return &orderV1.InternalServerError{
			Error:   orderV1.NewOptString("INTERNAL_SERVER_ERROR"),
			Code:    orderV1.NewOptInt(500),
			Message: orderV1.NewOptString("Failed to cancel order"),
		}, nil
	}

	return &orderV1.CancelOrderNoContent{}, nil
}
