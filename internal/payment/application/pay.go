/*package payment_service

type name struct {
}

func Pay() {
	// domain.validateUser
	// domain.validatePaymentReq
	// domain.toDbStruct
	// infrastructure.SendAtDb
}*/

package application

import (
	"context"

	"gitlab.com/godevs2/micro/internal/payment/domain/model"
)

type PaymentStorage interface {
	Pay(ctx context.Context, order *model.PayOrder) (*model.PayOrderResponse, error)
}

type PaymentService struct {
	paymentRepo PaymentStorage // interface
}

func NewPaymentService(paymentRepo PaymentStorage) *PaymentService {
	return &PaymentService{
		paymentRepo: paymentRepo,
	}
}

func (s *PaymentService) Pay(ctx context.Context, order *model.PayOrder) (*model.PayOrderResponse, error) {
	/*if err := domain.ValidateUser(ctx, order.UserUuid); err != nil {
		return nil, err
	}
	if err := domain.ValidatePaymentRequest(order); err != nil {
		return nil, err
	}*/

	return s.paymentRepo.Pay(ctx, order)
}
