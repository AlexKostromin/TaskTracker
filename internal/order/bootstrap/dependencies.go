package bootstrap

import (
	"gitlab.com/godevs2/micro/internal/order/application"
	http_server "gitlab.com/godevs2/micro/internal/order/infrastructure/http-server"
	mock_postgres "gitlab.com/godevs2/micro/internal/order/infrastructure/mock-postgres"
)

/*storage := mock-postgres.NewOrderStorage()
orderService := application.NewOrderService(storage)
orderHandler := NewOrderHandler(orderService)*/

func provideOrderStorage() application.OrderStorage { return mock_postgres.NewOrderStorage() }
func provideOrderService(s application.OrderStorage) http_server.OrderProcessor {
	return application.NewOrderService(s)
}
func provideOrderHandler(port string, s http_server.OrderProcessor) *http_server.Server {
	return http_server.NewServer(port, s)
}
