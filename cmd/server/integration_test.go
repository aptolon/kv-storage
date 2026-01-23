package main

import (
	"context"
	"os"
	"testing"

	"github.com/aptolon/kv-store/internal/persistence"
	"github.com/aptolon/kv-store/internal/storage"
	"github.com/jackc/pgx/v5"
)

func TestPersistenceRestart(t *testing.T) {
	ctx := t.Context()

	db := os.Getenv("DATABASE_URL")

	if db == "" {
		t.Fatal("DATABASE_URL not set")
	}
	conn, err := pgx.Connect(ctx, db)

	if err != nil {
		t.Fatalf("postgres connect error: %v", err)
	}
	defer conn.Close(ctx)

	nameTable := "kv_snapshot_test"

	if err := persistence.CreateSnapshotTable(ctx, conn, nameTable); err != nil {
		t.Fatalf("create table error: %v", err)
	}
	repo := persistence.NewPostgresSnapshotRepository(conn, nameTable)

	data, err := repo.Load(ctx)
	if err != nil {
		t.Fatalf("load snapshot error: %v", err)
	}

	store := storage.NewMemoryStorage(data)

	key := "123"
	value := []byte("456")

	err = store.Set(key, value)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if err := repo.Save(
		context.Background(),
		store.Snapshot(),
	); err != nil {
		t.Fatalf("snapshot save error: %v", err)
	}

	data2, err := repo.Load(ctx)
	if err != nil {
		t.Fatalf("failed to load snapshot after restart: %v", err)
	}

	store2 := storage.NewMemoryStorage(data2)

	got, err := store2.Get(key)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if string(got) != string(value) {
		t.Fatalf("expected %q, got %q", value, got)
	}
	if err := persistence.DropSnapshotTable(ctx, conn, nameTable); err != nil {
		t.Fatalf("drop table error: %v", err)
	}
}
