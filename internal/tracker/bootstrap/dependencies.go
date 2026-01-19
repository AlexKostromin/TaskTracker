package bootstrap

import (
	"database/sql"
	"time"

	"github.com/AlexKostromin/TaskTracker/internal/tracker/application"
	"github.com/AlexKostromin/TaskTracker/internal/tracker/config"
	http_server "github.com/AlexKostromin/TaskTracker/internal/tracker/infrastructure/http-server"
	"github.com/AlexKostromin/TaskTracker/internal/tracker/infrastructure/postgres"
	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/pressly/goose/v3"
)

func provideDB(cfg *config.Config) (*sql.DB, func() error, error) {
	db, err := sql.Open("pgx", cfg.DBURI)
	if err != nil {
		return nil, nil, err
	}

	// Базовые лимиты пула (можно расширить конфигом позже).
	db.SetMaxOpenConns(10)
	db.SetMaxIdleConns(10)
	db.SetConnMaxLifetime(30 * time.Minute)

	if err := db.Ping(); err != nil {
		_ = db.Close()
		return nil, nil, err
	}

	return db, db.Close, nil
}

func provideMigrations(db *sql.DB, cfg *config.Config) error {
	// goose определяет диалект по драйверу, но для уверенности зададим явно.
	if err := goose.SetDialect("postgres"); err != nil {
		return err
	}
	return goose.Up(db, cfg.MigrationsDir)
}

func provideTrackerStorage(db *sql.DB) application.TrackerStorage {
	return postgres.NewTrackerStorage(db)
}
func provideTrackerService(s application.TrackerStorage) http_server.TrackerProcessor {
	return application.NewTrackerService(s)
}
func provideTrackerHandler(port string, s http_server.TrackerProcessor) *http_server.Server {
	return http_server.NewServer(port, s)
}
