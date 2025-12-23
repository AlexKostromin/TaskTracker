package pkg

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"
)

type Closer func() error

var closers = make([]Closer, 0, 10)

func Add(c Closer) {
	closers = append(closers, c)
}

func CloseOnSignalContext(ctx context.Context) {

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	<-quit
	log.Println("🛑 Получен сигнал завершения...")

	for _, c := range closers {
		err := c()
		if err != nil {
			log.Printf("❌ Ошибка при закрытии: %v\n", err)
		}
	}

	log.Println("✅ Все ресурсы закрыты")
}
