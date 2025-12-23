package mock

import (
	"context"
	"log"
	"sync"

	"github.com/google/uuid"
	"gitlab.com/godevs2/micro/internal/payment/domain/model"
)

type PaymentStorage struct {
	mu           sync.RWMutex
	transactions map[uuid.UUID]*model.PayOrderResponse
}

func NewPaymentStorage() *PaymentStorage {
	return &PaymentStorage{
		transactions: make(map[uuid.UUID]*model.PayOrderResponse),
	}
}

func (s *PaymentStorage) Pay(_ context.Context, req *model.PayOrder) (*model.PayOrderResponse, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	transactionUUID := uuid.New().String()
	log.Printf("💳 Обработка платежа: TransactionUUID=%s, OrderID=%s",
		transactionUUID, req.OrderUuid)

	return &model.PayOrderResponse{
		TransactionUuid: transactionUUID,
	}, nil

}
