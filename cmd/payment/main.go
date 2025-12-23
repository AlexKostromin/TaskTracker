package main

import (
	"gitlab.com/godevs2/micro/internal/payment/bootstrap"
	"log"
)

func main() {
	// Создание приложения
	app := bootstrap.Configure()

	// Запуск приложения
	if err := app.Run(); err != nil {
		log.Fatalf("❌ Ошибка запуска приложения: %v", err)
	}
}
