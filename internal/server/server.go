package server

import (
	"bufio"
	"context"
	"io"
	"log"
	"net"
	"strings"
	"sync"

	"github.com/aptolon/kv-store/internal/storage"
)

type Server struct {
	addr     string
	storage  storage.Storage
	listener net.Listener
	ready    chan string
	wg       *sync.WaitGroup
}

func NewServer(addr string, storage storage.Storage) *Server {
	return &Server{
		addr:    addr,
		storage: storage,
		ready:   make(chan string, 1),
		wg:      &sync.WaitGroup{},
	}
}

func (s *Server) Start(ctx context.Context) error {
	var err error
	s.listener, err = net.Listen("tcp", s.addr)
	if err != nil {
		return err
	}
	defer s.listener.Close()

	s.ready <- s.listener.Addr().String()

	log.Printf("server started on port %s", s.addr)

	go func() {
		<-ctx.Done()
		log.Println("server stopped gracefully")
		s.listener.Close()
	}()

	for {
		conn, err := s.listener.Accept()
		if err != nil {
			return nil
		}
		s.wg.Add(1)
		go s.handleConn(ctx, conn)
	}
}

func (s *Server) handleConn(ctx context.Context, conn net.Conn) {
	defer conn.Close()
	log.Printf("client connected on port %s", s.addr)
	reader := bufio.NewReader(conn)
	writer := bufio.NewWriter(conn)
	for {
		select {
		case <-ctx.Done():
			return
		default:
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
