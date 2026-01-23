package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/aptolon/kv-store/internal/persistence"
	"github.com/aptolon/kv-store/internal/server"
	"github.com/aptolon/kv-store/internal/storage"
	"github.com/jackc/pgx/v5"
)

func main() {
	ctx, cancel := signal.NotifyContext(
		context.Background(),
		os.Interrupt,
		syscall.SIGTERM,
	)
	defer cancel()

	db := os.Getenv("DATABASE_URL")

	if db == "" {
		log.Fatal("DATABASE_URL not set")
	}
	conn, err := pgx.Connect(ctx, db)

	if err != nil {
		log.Fatalf("postgres connect error: %v", err)
	}
	defer conn.Close(ctx)

	nameTable := "kv_snapshot"
	if err := persistence.CreateSnapshotTable(ctx, conn, nameTable); err != nil {
		log.Fatalf("create table error: %v", err)
	}
	repo := persistence.NewPostgresSnapshotRepository(conn, nameTable)

	data, err := repo.Load(ctx)
	if err != nil {
		log.Fatalf("load snapshot error: %v", err)
	}

	store := storage.NewMemoryStorage(data)

	port := os.Getenv("SERV_PORT")
	serv := server.NewServer(port, store)
	go func() {
		if err := serv.Start(ctx); err != nil {
			log.Printf("server stopped with error: %v", err)
			cancel()
		}
	}()

	<-ctx.Done()
	if err := repo.Save(
		context.Background(),
		store.Snapshot(),
	); err != nil {
		log.Printf("snapshot save error: %v", err)
	}
	log.Println("shutdown signal received")
}
