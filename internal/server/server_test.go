package server

import (
	"fmt"
	"strings"
	"testing"

	"github.com/aptolon/kv-store/internal/storage"
)

func newTestServer() *Server {
	return NewServer(
		":0",
		storage.NewMemoryStorage(),
	)
}
func TestHandleCommandSet(t *testing.T) {
	s := newTestServer()

	key := "123"
	value := "456"

	resp := s.handleCommand(fmt.Sprintf("SET %s %s", key, value))
	expectedResp := "OK"
	if resp != expectedResp {
		t.Fatalf("expected %q, got %q", expectedResp, resp)
	}

	respVal, _ := s.storage.Get(key)
	if string(respVal) != string(value) {
		t.Fatalf("expected value %q, got %q", value, respVal)
	}
}

func TestHandleCommandGet(t *testing.T) {
	s := newTestServer()

	key := "123"
	value := []byte("456")

	s.storage.Set(key, value)

	resp := s.handleCommand("GET " + key)
	gotResp := "VALUE " + string(value)
	if resp != gotResp {
		t.Fatalf("expected %q, got %q", resp, gotResp)
	}
}

func TestHandleCommandDel(t *testing.T) {
	s := newTestServer()

	key := "123"
	value := []byte("456")

	s.storage.Set(key, value)

	resp := s.handleCommand("DEL " + key)
	gotResp := "OK"
	if resp != gotResp {
		t.Fatalf("expected %q, got %q", resp, gotResp)
	}

	val, _ := s.storage.Get(key)
	if val != nil {
		t.Fatalf("expected nil, got %q", val)
	}
}

func TestHandleCommandGetMissingKey(t *testing.T) {
	s := newTestServer()

	key := "123"

	resp := s.handleCommand("GET " + key)
	gotResp := "NULL"
	if resp != gotResp {
		t.Fatalf("expected %q, got %q", resp, gotResp)
	}
}

func TestHandleInvalidCommand(t *testing.T) {
	s := newTestServer()

	tests := []string{
		"SET a",
		"GET",
		"GET a b c d",
		"DEL",
		"SET a 1 2",
		"",
		"\n",
		"   ",
		"a a a",
	}

	for _, cmd := range tests {
		resp := s.handleCommand(cmd)
		if !strings.HasPrefix(resp, "ERROR") {
			t.Fatalf("cmd %q: expected ERROR, got %q", cmd, resp)
		}
	}
}
