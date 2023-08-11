// Package server provies functionality to start and stop the mini redis server.
package server

import (
	"crypto/tls"
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
	config     *tls.Config
	closechan  chan struct{}
}

func NewServer(host string, port uint, db *db.Db, ps *pubsub.PubSub, config *tls.Config) *Server {
	return &Server{
		host:      host,
		port:      fmt.Sprint(port),
		db:        db,
		ps:        ps,
		handler:   handler.NewHandler(),
		config:    config,
		closechan: make(chan struct{}),
	}
}

func (s *Server) Close() error {

	err := s.listener.Close()
	close(s.closechan)
	return err
}

func (s *Server) ListenAndServe() error {

	listener, err := func() (net.Listener, error) {
		if s.config == nil {
			fmt.Printf(("Unsecure listening on %s:%s\n"), s.host, s.port)
			return net.Listen(Protocal, fmt.Sprintf("%s:%s", s.host, s.port))
		}
		fmt.Printf(("Secure listening on %s:%s\n"), s.host, s.port)
		return tls.Listen(Protocal, fmt.Sprintf("%s:%s", s.host, s.port), s.config)
	}()
	s.listener = listener
	if err != nil {
		fmt.Println("Error listening:", err.Error())
		return err
	}
	fmt.Println("Waiting for client...")

	disconnect := make(chan net.Conn)
	go func() {
		for {
			connection, err := s.listener.Accept()
			if err != nil {
				fmt.Println("Error accepting:", err.Error())
				break
			}
			fmt.Println("client connected")
			go processClient(connection, context.NewContext(s.db, s.ps), s.handler, disconnect)
		}
	}()

	for {
		select {
		case disconn := <-disconnect:
			fmt.Println("Deallocating resource for disconnect connection")
			s.removeConnection(disconn)
		case <-s.closechan:
			break
		}
	}
}

func (s *Server) removeConnection(conn net.Conn) {
	s.ps.Unsub(conn)
}

func processClient(conn net.Conn, ctx *context.Context, handler handler.Handler, disconnect chan net.Conn) {

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
	if disconnect != nil {
		disconnect <- conn
	}
	fmt.Println("client closed")
}
