// miniredis project server main.go
package server

import (
	"fmt"
	"io"
	"net"
	"os"

	"github.com/panupakm/miniredis"
	"github.com/panupakm/miniredis/internal/db"
	"github.com/panupakm/miniredis/internal/server/handler"
	"github.com/panupakm/miniredis/lib/cmd"
	"github.com/panupakm/miniredis/lib/payload"
)

const (
	Protocal = "tcp"
)

type Server struct {
	host, port string
	conn       net.Conn
	listener   net.Listener
	db         *db.Db
}

func NewServer(host, port string, db *db.Db) *Server {
	return &Server{
		host: host,
		port: port,
		db:   db,
	}
}

func (s *Server) Close() error {
	return s.listener.Close()
}

func (s *Server) Start() {
	listener, err := net.Listen(Protocal, fmt.Sprintf("%s:%s", s.host, s.port))
	s.listener = listener
	if err != nil {
		fmt.Println("Error listening:", err.Error())
		os.Exit(1)
	}

	fmt.Printf(("Listening on %s:%s\n"), s.host, s.port)
	fmt.Println("Waiting for client...")

	for {
		connection, err := s.listener.Accept()
		if err != nil {
			fmt.Println("Error accepting:", err.Error())
			continue
		}
		fmt.Println("client connected")
		go s.processClient(connection, miniredis.NewContext(s.db))
	}
}

// TODO: handle close connections
func (s *Server) processClient(conn net.Conn, ctx *miniredis.Context) {

	for {
		var cmdstr payload.String
		_, err := cmdstr.ReaderFrom(conn)
		if err != nil {
			if err == io.EOF {
				fmt.Printf("Client %s has closed connection\n", conn.RemoteAddr())
			} else {
				fmt.Printf("Unexpected error: %s\n", err)
			}
			break
		}
		switch cmdstr {
		case cmd.PingCode:
			handler.HandlePing(conn, ctx)
		case cmd.SetCode:
			handler.HandleSet(conn, ctx)
		case cmd.GetCode:
			handler.HandleGet(conn, ctx)
		}
	}
	fmt.Println("client closed")
}
