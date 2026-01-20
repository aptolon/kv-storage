package storage

import (
	"sync"
	"testing"
)

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

func TestStorageDeleteKey(t *testing.T) {
	store := newTestStorage()

	key := "123"
	value := []byte("456")

	err := store.Set(key, value)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	err = store.Delete(key)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	got, err := store.Get(key)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got != nil {
		t.Fatalf("expected nil, got %q", got)
	}
}
func TestStorageIsolatedSet(t *testing.T) {
	store := newTestStorage()
	key := "123"
	value := []byte("456")
	valcopy := make([]byte, len(value))
	copy(valcopy, value)

	err := store.Set(key, value)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	value[0] = '9'

	got, err := store.Get(key)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if string(got) != string(valcopy) {
		t.Fatalf("expected %q, got %q", valcopy, got)
	}

}

func TestStorageIsolatedGet(t *testing.T) {
	store := newTestStorage()
	key := "123"
	value := []byte("456")

	err := store.Set(key, value)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	got1, err := store.Get(key)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	got1[0] = '9'

	got2, err := store.Get(key)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if string(got2) != string(value) {
		t.Fatalf("expected %q, got %q", value, got2)
	}

}

func TestStorageConcurrentSetNoRace(t *testing.T) {
	store := newTestStorage()

	value := []byte("abc")
	key := "def"
	wg := &sync.WaitGroup{}
	workers := 1000
	errCh := make(chan error, workers)
	for range workers {
		wg.Add(1)
		go func() {
			defer wg.Done()
			err := store.Set(key, value)
			if err != nil {
				errCh <- err
			}

		}()
	}
	wg.Wait()
	close(errCh)
	for err := range errCh {
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

	}

}

func TestStorageConcurrentGetNoRace(t *testing.T) {
	store := newTestStorage()

	value := []byte("abc")
	key := "def"

	err := store.Set(key, value)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	wg := &sync.WaitGroup{}
	workers := 1000
	errCh := make(chan error, workers)
	valCh := make(chan []byte, workers)
	for range workers {
		wg.Add(1)
		go func() {
			defer wg.Done()
			got, err := store.Get(key)
			if err != nil {
				errCh <- err
			}

			valCh <- got

		}()
	}
	wg.Wait()
	close(errCh)
	close(valCh)
	for err := range errCh {
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
	}
	for got := range valCh {
		if string(got) != string(value) {
			t.Fatalf("expected %q, got %q", value, got)
		}
	}

}
func TestStorageConcurrentSetGetNoRace(t *testing.T) {
	store := newTestStorage()

	var workers = 1000
	value := []byte("abc")
	key := "def"

	var wg sync.WaitGroup
	errCh := make(chan error, workers*2)

	for range workers {
		wg.Add(1)
		go func() {
			defer wg.Done()
			if err := store.Set(key, value); err != nil {
				errCh <- err
			}
		}()
	}

	for range workers {
		wg.Add(1)
		go func() {
			defer wg.Done()

			_, err := store.Get(key)
			if err != nil {
				errCh <- err
			}
		}()
	}

	wg.Wait()
	close(errCh)

	for err := range errCh {
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
	}
}
