package server

import (
	"bufio"
	"fmt"
	"net"
	"testing"

	"github.com/aptolon/kv-store/internal/storage"
)

func TestTCPSetGet(t *testing.T) {
	ctx := t.Context()

	store := storage.NewMemoryStorage()
	srv := NewServer(":0", store)

	go func() {
		if err := srv.Start(ctx); err != nil {
			t.Errorf("server error: %v", err)
		}
	}()
	addr := <-srv.ready

	conn, err := net.Dial("tcp", addr)
	if err != nil {
		t.Fatalf("failed to connect: %v", err)
	}
	defer conn.Close()

	reader := bufio.NewReader(conn)
	writer := bufio.NewWriter(conn)

	key := "123"
	value := "456"

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
