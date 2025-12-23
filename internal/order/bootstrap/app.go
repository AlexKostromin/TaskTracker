package bootstrap

import (
	"context"
	"log"

	"gitlab.com/godevs2/micro/internal/order/config"
	"gitlab.com/godevs2/micro/internal/order/infrastructure/http-server"
	"gitlab.com/godevs2/micro/pkg"
)

type App struct {
	server *http_server.Server
	config *config.Config
}

func New() *App {
	// Загрузка конфигурации
	cfg := config.Load()

	// Создание HTTP сервера

	storage := provideOrderStorage()
	orderService := provideOrderService(storage)
	orderHandler := provideOrderHandler(cfg.HTTPPort, orderService)

	app := &App{
		server: orderHandler,
		config: cfg,
	}

	// Регистрируем closer для graceful shutdown
	pkg.Add(app.gracefulShutdown)

	return app
}

func (a *App) Run() error {
	// Запуск HTTP сервера
	go func() {
		if err := a.server.Start(); err != nil {
			log.Printf("Ошибка запуска сервера: %v", err)
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

	// Останавливаем HTTP сервер
	if err := a.server.Shutdown(ctx); err != nil {
		log.Printf("❌ Ошибка при остановке сервера: %v\n", err)
		return err
	}

	// Здесь можно добавить закрытие других ресурсов:
	// - База данных
	// - Redis
	// - Другие соединения

	log.Println("✅ Приложение остановлено")
	return nil
}
