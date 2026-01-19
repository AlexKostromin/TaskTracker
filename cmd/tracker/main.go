package main

import (
	"log"

	"github.com/AlexKostromin/TaskTracker/internal/tracker/bootstrap"
)

func main() {
	// Создание приложения
	app := bootstrap.Configure()

	// Запуск приложения
	if err := app.Run(); err != nil {
		log.Fatalf("❌ Ошибка запуска приложения: %v", err)
	}
}
