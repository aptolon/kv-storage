package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/aptolon/kv-store/internal/server"
	"github.com/aptolon/kv-store/internal/storage"
)

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, os.Interrupt, syscall.SIGTERM)

	store := storage.NewMemoryStorage()
	serv := server.NewServer(":8080", store)
	go func() {
		if err := serv.Start(ctx); err != nil {
			log.Printf("server stopped with error: %v", err)
			cancel()
		}
	}()

	<-sigCh
	log.Println("shutdown signal received")
	time.Sleep(time.Second)

}
