package server

import (
	"bufio"
	"context"
	"fmt"
	"net"
	"strings"
	"testing"
	"time"

	"github.com/aptolon/kv-store/internal/storage"
)

func TestTCPSetGet(t *testing.T) {
	ctx := t.Context()
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
func TestTCPGetMissingKey(t *testing.T) {
	ctx := t.Context()

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

	key := "123"

	// GET
	fmt.Fprintf(writer, "GET %s\n", key)
	writer.Flush()

	resp, _ := reader.ReadString('\n')
	expectedResp := "NULL\n"
	if resp != expectedResp {
		t.Fatalf("expected %q, got %q", expectedResp, resp)
	}

}

func TestTCPInvalidCommand(t *testing.T) {
	ctx := t.Context()

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

	tests := []string{
		"SET a\n",
		"GET\n",
		"GET a b c d\n",
		"DEL\n",
		"SET a 1 2\n",
		"\n",
		"\n\n",
		"   \n",
		"a a a\n",
	}
	for _, cmd := range tests {
		writer.WriteString(cmd)
		writer.Flush()
		resp, _ := reader.ReadString('\n')
		if !strings.HasPrefix(resp, "ERROR") {
			t.Fatalf("cmd %q: expected ERROR, got %q", cmd, resp)
		}

	}

}

func TestTCPMultiCommand(t *testing.T) {
	ctx := t.Context()

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

	key1 := "123"
	value1 := "456"
	key2 := "789"
	value2 := "123"

	// SET
	fmt.Fprintf(writer, "SET %s %s\n", key1, value1)
	writer.Flush()

	resp, _ := reader.ReadString('\n')
	expectedResp := "OK\n"
	if resp != expectedResp {
		t.Fatalf("expected %q, got %q", expectedResp, resp)
	}
	// SET
	fmt.Fprintf(writer, "SET %s %s\n", key2, value2)
	writer.Flush()

	resp, _ = reader.ReadString('\n')
	expectedResp = "OK\n"
	if resp != expectedResp {
		t.Fatalf("expected %q, got %q", expectedResp, resp)
	}

	// GET
	fmt.Fprintf(writer, "GET %s\n", key1)
	writer.Flush()

	resp, _ = reader.ReadString('\n')
	expectedResp = fmt.Sprintf("VALUE %s\n", value1)
	if resp != expectedResp {
		t.Fatalf("expected %q, got %q", expectedResp, resp)
	}

	// GET
	fmt.Fprintf(writer, "GET %s\n", key2)
	writer.Flush()

	resp, _ = reader.ReadString('\n')
	expectedResp = fmt.Sprintf("VALUE %s\n", value2)
	if resp != expectedResp {
		t.Fatalf("expected %q, got %q", expectedResp, resp)
	}

}

func TestTCPGracefulShutdown(t *testing.T) {
	ctx, cancel := context.WithTimeout(t.Context(), 5*time.Second)
	defer cancel()

	data := make(map[string][]byte)
	store := storage.NewMemoryStorage(data)
	srv := NewServer(":0", store)

	serverErr := make(chan error, 1)

	go func() {
		serverErr <- srv.Start(ctx)
	}()

	var addr string
	select {
	case addr = <-srv.ready:
	case <-time.After(time.Second):
		t.Fatalf("server didn't start")
	}

	conn, err := net.Dial("tcp", addr)
	if err != nil {
		t.Fatalf("failed to connect: %v", err)
	}

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
	cancel()

	select {
	case err := <-serverErr:
		if err != nil && err != context.Canceled {
			t.Fatalf("server returned error: %v", err)
		}
	case <-time.After(time.Second):
		t.Fatalf("server didn't shutdown")
	}

	_, err = net.DialTimeout("tcp", addr, 200*time.Millisecond)
	if err == nil {
		t.Fatalf("expected connection отказ после shutdown")
	}
}
