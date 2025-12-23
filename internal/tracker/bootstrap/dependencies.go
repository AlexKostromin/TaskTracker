package bootstrap

import (
	http_server "gitlab.com/godevs2/micro/internal/order/infrastructure/http-server"
	"gitlab.com/godevs2/micro/internal/tracker/application"
	"gitlab.com/godevs2/micro/internal/tracker/infrastructure/mock-postgres"
)

func provideTrackerStorage() application.TrackerStorage {
	return mock_postgres.NewTrackerStorage()
}
func provideTrackerService(s application.TrackerStorage) http_server.TrackerProcessor {
	return application.NewTrackerService(s)
}
func provideTrackerHandler(port string, s http_server.TrackerProcessor) *http_server.OrderProcessor {
	return http_server.NewServer(port, s)
}
