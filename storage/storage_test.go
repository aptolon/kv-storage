package storage

import "testing"

func newTestStorage() Storage {
	return NewMemoryStorage()
}

func TestStorageSetGet(t *testing.T) {
	store := newTestStorage()
	key := "123"
	value := []byte("456")

	err := store.Set(key, value)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	got, err := store.Get(key)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if string(got) != string(value) {
		t.Fatalf("expected %q, got %q", value, got)
	}
}

func TestStorageSetOverride(t *testing.T) {
	store := newTestStorage()
	key := "123"
	value := []byte("456")

	err := store.Set(key, value)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	value = []byte("789")
	err = store.Set(key, value)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	got, err := store.Get(key)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if string(got) != string(value) {
		t.Fatalf("expected %q, got %q", value, got)
	}
}

func TestStorageGetMissingKey(t *testing.T) {
	store := newTestStorage()

	got, err := store.Get("miss")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got != nil {
		t.Fatalf("expected nil, got %q", got)
	}
}

func TestStorageDeleteMissingKey(t *testing.T) {
	store := newTestStorage()

	err := store.Delete("miss")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}
