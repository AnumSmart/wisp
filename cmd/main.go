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

	//канал для ошибок запуска сервера
	serverErr := make(chan error, 1)

	// Запуск сервера в отдельной горутине
	go func() {
		if err := server.Run(); err != nil {
			serverErr <- err
		}
	}()

	// Ожидаем сигнал завершения или ошибку сервера
	select {
	case <-exit:
		log.Println("Shutting down server...")
	case err := <-serverErr:
		log.Printf("Server error: %v", err)
	}

	// Завершаем работу с таймаутом
	shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer shutdownCancel()

	if err := server.Shutdown(shutdownCtx); err != nil {
		log.Fatalf("Server shutdown error: %v", err)
	}

	log.Println("Server exited properly")
}
