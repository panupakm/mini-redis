// mini-redis project server main.go
package server

import (
	"fmt"
	"net"
	"os"

	"github.com/panupakm/mini-redis/lib/payload"
	"github.com/panupakm/mini-redis/lib/ping"
)

const (
	Protocal = "tcp"
)

type Server struct {
	host, port string
	conn       net.Conn
	listener   net.Listener
}

func NewServer(host, port string) *Server {
	return &Server{
		host: host,
		port: port,
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
		go s.processClient(connection)
	}
}

func (s *Server) processClient(conn net.Conn) {
	var cmdstr payload.String
	cmdstr.ReaderFrom(conn)
	fmt.Println("Command:", cmdstr)
	switch cmdstr {
	case ping.Code:
		s.handlePing(conn)
	}
}

func (s *Server) handlePing(conn net.Conn) {

	msg := ping.NewPing(conn)

	pong := payload.String(msg.String())
	_, err := pong.WriterTo(conn)
	if err != nil {
		fmt.Println("Error writing:", err.Error())
	}
}
