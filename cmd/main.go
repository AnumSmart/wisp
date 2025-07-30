package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func main() {
	// Создаем корневой контекст
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Инициализируем сервер
	server := NewServer(ctx)

	// механизм gracefull shutdown
	exit := make(chan os.Signal, 1)
	signal.Notify(exit, syscall.SIGINT, syscall.SIGTERM)

	// Запуск сервера в отдельной горутине
	go func() {
		if err := server.Run(); err != nil {
			log.Fatalf("Failed to start server: %v", err)
		}
	}()

	// Выход
	<-exit

	log.Println("Shutting down server...")

	// Завершаем работу с таймаутом
	shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer shutdownCancel()

	if err := server.Shutdown(shutdownCtx); err != nil {
		log.Fatalf("Server shutdown error: %v", err)
	}

	log.Println("Server exited properly")
}
