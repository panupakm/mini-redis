// Package server provies functionality to start and stop the mini redis server.
package server

import (
	"fmt"
	"io"
	"net"

	"github.com/panupakm/miniredis/internal/db"
	"github.com/panupakm/miniredis/internal/pubsub"
	"github.com/panupakm/miniredis/payload"
	cmd "github.com/panupakm/miniredis/request"
	"github.com/panupakm/miniredis/server/context"
	"github.com/panupakm/miniredis/server/internal/handler"
)

const (
	Protocal = "tcp"
)

type Server struct {
	host, port string
	conn       net.Conn
	listener   net.Listener
	db         *db.Db
	ps         *pubsub.PubSub
	handler    handler.Handler
}

func NewServer(host string, port uint, db *db.Db, ps *pubsub.PubSub) *Server {
	return &Server{
		host:    host,
		port:    fmt.Sprint(port),
		db:      db,
		ps:      ps,
		handler: handler.NewHandler(),
	}
}

func (s *Server) Close() error {
	return s.listener.Close()
}

func (s *Server) ListenAndServe() error {
	listener, err := net.Listen(Protocal, fmt.Sprintf("%s:%s", s.host, s.port))
	s.listener = listener
	if err != nil {
		fmt.Println("Error listening:", err.Error())
		return err
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
		go processClient(connection, context.NewContext(s.db, s.ps), s.handler)
	}
}

func processClient(conn net.Conn, ctx *context.Context, handler handler.Handler) {

	for {
		var cmdstr payload.String
		_, err := cmdstr.ReadFrom(conn)
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
			handler.HandlePing(conn)
		case cmd.SetCode:
			handler.HandleSet(conn, ctx)
		case cmd.GetCode:
			handler.HandleGet(conn, ctx)
		case cmd.SubCode:
			handler.HandleSub(conn, ctx)
		case cmd.PubCode:
			handler.HandlePub(conn, ctx)
		}
	}
	fmt.Println("client closed")
}
