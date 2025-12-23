package main

import (
	"log"

	"gitlab.com/godevs2/micro/internal/order/bootstrap"
)

func main() {
	app := bootstrap.Configure()

	// Запускаем приложение
	log.Println("🚀 Запуск Order Service...")
	if err := app.Run(); err != nil {

		log.Fatalf("❌ Ошибка запуска приложения: %v", err)
	}
}
