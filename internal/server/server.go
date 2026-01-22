package server

import (
	"bufio"
	"context"
	"io"
	"log"
	"net"
	"strings"

	"github.com/aptolon/kv-store/internal/storage"
)

type Server struct {
	addr    string
	storage storage.Storage
	ready   chan string
}

func NewServer(addr string, storage storage.Storage) *Server {
	return &Server{
		addr:    addr,
		storage: storage,
		ready:   make(chan string, 1),
	}
}

func (s *Server) Start(ctx context.Context) error {
	listener, err := net.Listen("tcp", s.addr)
	if err != nil {
		return err
	}
	defer listener.Close()

	s.ready <- listener.Addr().String()

	log.Printf("server started on port %s", s.addr)
	for {
		select {
		case <-ctx.Done():
			return nil
		default:
			conn, err := listener.Accept()
			if err != nil {
				continue
			}
			go s.handleConn(conn)
		}
	}
}

func (s *Server) handleConn(conn net.Conn) {
	defer conn.Close()
	log.Printf("client connected on port %s", s.addr)
	reader := bufio.NewReader(conn)
	writer := bufio.NewWriter(conn)
	for {
		line, err := reader.ReadString('\n')

		if err == io.EOF {
			return
		}
		if err != nil {
			log.Println(err)
			return
		}
		line = strings.TrimSpace(line)
		resp := s.handleCommand(line)
		writer.WriteString(resp + "\n")
		writer.Flush()
	}
}

func (s *Server) handleCommand(line string) string {
	parts := strings.Fields(line)
	if len(parts) == 0 {
		return "ERROR empty command"
	}
	cmd := strings.ToUpper(parts[0])
	switch cmd {
	case "SET":
		if len(parts) != 3 {
			return "ERROR invalid arguments"
		}
		key := parts[1]
		value := parts[2]
		err := s.storage.Set(key, []byte(value))
		if err != nil {
			return "ERROR internal error"
		}
		return "OK"
	case "GET":
		if len(parts) != 2 {
			return "ERROR invalid arguments"
		}
		key := parts[1]
		value, err := s.storage.Get(key)
		if err != nil {
			return "ERROR internal error"
		}
		if value == nil {
			return "NULL"
		}
		return "VALUE " + string(value)
	case "DEL":
		if len(parts) != 2 {
			return "ERROR invalid arguments"
		}
		key := parts[1]
		err := s.storage.Delete(key)
		if err != nil {
			return "ERROR internal error"
		}
		return "OK"
	default:
		return "ERROR invalid command"
	}
}
