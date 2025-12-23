package grpc_server

import (
	"context"

	"gitlab.com/godevs2/micro/internal/payment/infrastructure/grpc-server/converter"
	paymentV1 "gitlab.com/godevs2/micro/shared/pkg/proto/payment/v1"
)

// PayOrder обрабатывает запрос на оплату заказа
func (s *Server) PayOrder(ctx context.Context, req *paymentV1.PayOrderRequest) (*paymentV1.PayOrderResponse, error) {
	payOrder, err := s.paymentService.Pay(ctx, converter.PayOrderRequestToModel(req))
	if err != nil {
		return nil, err
	}
	return &paymentV1.PayOrderResponse{
		TransactionUuid: payOrder.TransactionUuid,
	}, nil
}
