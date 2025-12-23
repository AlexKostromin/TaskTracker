package bootstrap

import (
	"gitlab.com/godevs2/micro/internal/payment/application"
	grpc_server "gitlab.com/godevs2/micro/internal/payment/infrastructure/grpc-server"
	"gitlab.com/godevs2/micro/internal/payment/infrastructure/mock"
)

func providePaymentStorage() application.PaymentStorage {
	return mock.NewPaymentStorage()
}
func providePaymentService(s application.PaymentStorage) grpc_server.PaymentProcessor {
	return application.NewPaymentService(s)
}
func providePaymentHandler(port string, s grpc_server.PaymentProcessor) *grpc_server.Server {
	return grpc_server.NewServer(port, s)
}
