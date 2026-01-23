package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/aptolon/kv-store/internal/server"
	"github.com/aptolon/kv-store/internal/storage"
)

func main() {
	ctx, cancel := signal.NotifyContext(
		context.Background(),
		os.Interrupt,
		syscall.SIGTERM,
	)
	defer cancel()

	store := storage.NewMemoryStorage()
	serv := server.NewServer(":8080", store)
	go func() {
		if err := serv.Start(ctx); err != nil {
			log.Printf("server stopped with error: %v", err)
			cancel()
		}
	}()

	<-ctx.Done()
	log.Println("shutdown signal received")
}
