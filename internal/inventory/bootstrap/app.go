package bootstrap

import (
	"context"
	"gitlab.com/godevs2/micro/internal/inventory/config"
	grpc_server "gitlab.com/godevs2/micro/internal/inventory/infrastructure/grpc-server"
	"log"

	"gitlab.com/godevs2/micro/pkg"
)

type App struct {
	grpcServer *grpc_server.Server
	config     *config.Config
}

func New() *App {
	// Загрузка конфигурации
	cfg := config.Load()

	// Создание gRPC сервера

	storage := provideInventoryStorage()
	inventoryService := provideInventoryService(storage)
	inventoryHandler := provideInventoryHandler(cfg.GRPCPort, inventoryService)

	app := &App{
		grpcServer: inventoryHandler,
		config:     cfg,
	}

	// Регистрируем closer для graceful shutdown
	pkg.Add(app.gracefulShutdown)

	return app
}

func (a *App) Run() error {
	// Запуск gRPC сервера
	go func() {
		if err := a.grpcServer.Start(); err != nil {
			log.Printf("Ошибка запуска gRPC сервера: %v", err)
		}
	}()

	// Ожидание сигнала завершения через closer
	ctx := context.Background()
	pkg.CloseOnSignalContext(ctx)

	return nil
}

// gracefulShutdown закрывает все ресурсы приложения
func (a *App) gracefulShutdown() error {
	log.Println("🛑 Завершение работы приложения...")

	// Создаем контекст с таймаутом для остановки сервера
	ctx, cancel := context.WithTimeout(context.Background(), a.config.ShutdownTimeout)
	defer cancel()

	// Останавливаем gRPC сервер
	if err := a.grpcServer.Shutdown(ctx); err != nil {
		log.Printf("❌ Ошибка при остановке gRPC сервера: %v\n", err)
		return err
	}

	log.Println("✅ Приложение остановлено")
	return nil
}
