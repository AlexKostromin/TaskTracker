package bootstrap

import (
	"context"
	"log"

	"github.com/AlexKostromin/TaskTracker/internal/tracker/config"
	http_server "github.com/AlexKostromin/TaskTracker/internal/tracker/infrastructure/http-server"
	"github.com/AlexKostromin/TaskTracker/pkg"
)

type App struct {
	server *http_server.Server
	config *config.Config
}

func New() *App {
	cfg := config.Load()

	db, closeDB, err := provideDB(cfg)
	if err != nil {
		log.Fatalf("❌ Ошибка подключения к БД: %v", err)
	}
	pkg.Add(closeDB)

	if err := provideMigrations(db, cfg); err != nil {
		log.Fatalf("❌ Ошибка применения миграций: %v", err)
	}

	storage := provideTrackerStorage(db)
	trackerService := provideTrackerService(storage)
	trackerHandler := provideTrackerHandler(cfg.HTTPPort, trackerService)

	app := &App{
		server: trackerHandler,
		config: cfg,
	}

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
