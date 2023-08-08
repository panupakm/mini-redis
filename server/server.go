// Package server provies functionality to start and stop the mini redis server.
package server

import (
	"fmt"
	"io"
	"net"

	"github.com/panupakm/miniredis"
	"github.com/panupakm/miniredis/internal/db"
	"github.com/panupakm/miniredis/internal/pubsub"
	"github.com/panupakm/miniredis/payload"
	cmd "github.com/panupakm/miniredis/request"
	"github.com/panupakm/miniredis/server/handler"
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
}

func NewServer(host, port string, db *db.Db, ps *pubsub.PubSub) *Server {
	return &Server{
		host: host,
		port: port,
		db:   db,
		ps:   ps,
	}
}

func (s *Server) Close() error {
	return s.listener.Close()
}

func (s *Server) Start() error {
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
		go processClient(connection, miniredis.NewContext(s.db, s.ps))
	}
}

func processClient(conn net.Conn, ctx *miniredis.Context) {

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
