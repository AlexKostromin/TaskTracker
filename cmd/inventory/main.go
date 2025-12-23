package main

import (
	"gitlab.com/godevs2/micro/internal/inventory/bootstrap"
	"log"
)

func main() {

	app := bootstrap.Configure()

	if err := app.Run(); err != nil {
		log.Fatalf("❌ Ошибка запуска приложения: %v", err)
	}
}
