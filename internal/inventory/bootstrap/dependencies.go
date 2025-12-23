package bootstrap

import (
	application "gitlab.com/godevs2/micro/internal/inventory/application"
	grpc_server "gitlab.com/godevs2/micro/internal/inventory/infrastructure/grpc-server"
	mock "gitlab.com/godevs2/micro/internal/inventory/infrastructure/mock"
)

func provideInventoryStorage() application.InventoryStorage {
	return mock.NewInventoryStorage()
}
func provideInventoryService(s application.InventoryStorage) grpc_server.InventoryProcessor {
	return application.NewInventoryService(s)
}
func provideInventoryHandler(port string, s grpc_server.InventoryProcessor) *grpc_server.Server {
	return grpc_server.NewServer(port, s)
}
