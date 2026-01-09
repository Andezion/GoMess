package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/Andezion/GoMess/server"
)

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	srv := server.NewServer(ctx)

	go srv.Run()

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)

	errChan := make(chan error, 1)
	go func() {
		if err := srv.Listen(":8080"); err != nil {
			errChan <- err
		}
	}()

	fmt.Println("Сервер мессенджера запущен на порту 8080")
	fmt.Println("Нажмите Ctrl+C для остановки")

	select {
	case <-sigChan:
		fmt.Println("\nПолучен сигнал остановки...")
		cancel()
	case err := <-errChan:
		fmt.Printf("Ошибка сервера: %v\n", err)
		cancel()
	}

	fmt.Println("Сервер остановлен")
}
