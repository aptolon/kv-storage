package server

import (
	"bufio"
	"context"
	"fmt"
	"net"
	"sync"
	"testing"
	"time"

	"github.com/aptolon/kv-store/internal/storage"
)

func TestTCPStressSingleConnection(t *testing.T) {
	ctx, cancel := context.WithTimeout(t.Context(), 5*time.Second)
	defer cancel()
	data := make(map[string][]byte)
	store := storage.NewMemoryStorage(data)
	srv := NewServer(":0", store)

	go func() {
		if err := srv.Start(ctx); err != nil {
			t.Errorf("server error: %v", err)
		}
	}()
	var addr string
	select {
	case addr = <-srv.ready:
	case <-time.After(time.Second):
		t.Fatalf("server didn't started")
	}

	conn, err := net.Dial("tcp", addr)

	if err != nil {
		t.Fatalf("failed to connect: %v", err)
	}
	defer conn.Close()

	reader := bufio.NewReader(conn)
	writer := bufio.NewWriter(conn)

	var rng = 1000

	for i := range rng {
		key := fmt.Sprintf("key_%d", i)
		value := fmt.Sprintf("value_%d", i)
		// SET
		fmt.Fprintf(writer, "SET %s %s\n", key, value)
		writer.Flush()

		resp, _ := reader.ReadString('\n')
		expectedResp := "OK\n"
		if resp != expectedResp {
			t.Fatalf("expected %q, got %q", expectedResp, resp)
		}

		// GET
		fmt.Fprintf(writer, "GET %s\n", key)
		writer.Flush()

		resp, _ = reader.ReadString('\n')
		expectedResp = fmt.Sprintf("VALUE %s\n", value)
		if resp != expectedResp {
			t.Fatalf("expected %q, got %q", expectedResp, resp)
		}
	}

}

func TestTCPStressConcurrentClients(t *testing.T) {
	ctx, cancel := context.WithTimeout(t.Context(), 5*time.Second)
	defer cancel()

	data := make(map[string][]byte)
	store := storage.NewMemoryStorage(data)
	srv := NewServer(":0", store)

	go func() {
		if err := srv.Start(ctx); err != nil {
			t.Errorf("server error: %v", err)
		}
	}()
	var addr string
	select {
	case addr = <-srv.ready:
	case <-time.After(time.Second):
		t.Fatalf("server didn't started")
	}
	workers := 10
	errCh := make(chan string, workers*3)
	wg := &sync.WaitGroup{}

	for idWorker := range workers {
		wg.Add(1)

		go func() {
			defer wg.Done()
			conn, err := net.Dial("tcp", addr)

			if err != nil {
				errCh <- fmt.Sprintf("failed to connect: %v", err)
			}
			defer conn.Close()

			reader := bufio.NewReader(conn)
			writer := bufio.NewWriter(conn)

			key := fmt.Sprintf("key_%d", idWorker)
			value := fmt.Sprintf("value_%d", idWorker)

			// SET
			fmt.Fprintf(writer, "SET %s %s\n", key, value)
			writer.Flush()

			resp, _ := reader.ReadString('\n')
			expectedResp := "OK\n"
			if resp != expectedResp {
				errCh <- fmt.Sprintf("expected %q, got %q", expectedResp, resp)
			}

			// GET
			fmt.Fprintf(writer, "GET %s\n", key)
			writer.Flush()

			resp, _ = reader.ReadString('\n')
			expectedResp = fmt.Sprintf("VALUE %s\n", value)
			if resp != expectedResp {
				errCh <- fmt.Sprintf("expected %q, got %q", expectedResp, resp)
			}
		}()
	}
	wg.Wait()
	close(errCh)

	for err := range errCh {
		t.Fatalf("%v", err)
	}
}
