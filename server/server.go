// Package server provies functionality to start and stop the mini redis server.
package server

import (
	"bytes"
	"crypto/tls"
	"encoding/binary"
	"fmt"
	"io"
	"net"
	"os"

	"github.com/panupakm/miniredis/payload"
	cmd "github.com/panupakm/miniredis/request"
	"github.com/panupakm/miniredis/server/context"
	"github.com/panupakm/miniredis/server/internal/handler"
	"github.com/panupakm/miniredis/server/pubsub"
	"github.com/panupakm/miniredis/server/storage"
)

const (
	Protocal = "tcp"
)

type readWriter struct {
	io.Reader
	io.Writer
}

type Config struct {
	tls.Config
	PersistentPath string
}

func (c *Config) hasCertificates() bool {
	return len(c.Certificates) > 0
}

type Server struct {
	conn        net.Conn
	listener    net.Listener
	storage     storage.Storage
	pubsub      pubsub.PubSub
	handler     handler.Handler
	config      *Config
	closechan   chan struct{}
	persistFile *os.File
}

func NewServer(db storage.Storage, ps pubsub.PubSub) *Server {
	return &Server{
		storage:   db,
		pubsub:    ps,
		handler:   handler.NewHandler(),
		closechan: make(chan struct{}),
	}
}

func (s *Server) Close() error {
	if s.persistFile != nil {
		s.persistFile.Close()
	}
	err := s.listener.Close()
	close(s.closechan)

	return err
}

func (s *Server) restoreServer() error {
	if s.config == nil || len(s.config.PersistentPath) == 0 {
		fmt.Println("No file to restore server state from")
		return nil
	}

	fmt.Println("Try to restore server state from:", s.config.PersistentPath)
	f, err := os.OpenFile(s.config.PersistentPath, os.O_RDWR|os.O_APPEND|os.O_CREATE, 0600)
	if err != nil {
		return err
	}
	s.persistFile = f

	// Read the contents of the file.
	data, err := io.ReadAll(s.persistFile)
	if err != nil {
		return err
	}

	dataReader := bytes.NewReader(data)

	for {
		var size uint32
		if err := binary.Read(dataReader, binary.BigEndian, &size); err != nil {
			break
		}
		buff := make([]byte, size)
		dataReader.Read(buff)

		discardWriter := io.Discard
		processBytesCommand(readWriter{
			Reader: bytes.NewReader(buff),
			Writer: discardWriter,
		}, context.NewContext(s.storage, s.pubsub), s.handler)
	}

	return nil
}

func (s *Server) ListenAndServe(host string, port uint, config *Config) error {
	s.config = config
	listener, err := func() (net.Listener, error) {
		if s.config == nil || !s.config.hasCertificates() {
			fmt.Printf(("Unsecure listening on %s:%d\n"), host, port)
			return net.Listen(Protocal, net.JoinHostPort(host, fmt.Sprint(port)))
		}
		fmt.Printf(("Secure listening on %s:%d\n"), host, port)
		return tls.Listen(Protocal, net.JoinHostPort(host, fmt.Sprint(port)), &s.config.Config)
	}()
	s.listener = listener
	if err != nil {
		fmt.Println("Error listening:", err.Error())
		return err
	}

	s.restoreServer()
	fmt.Println("Waiting for client...")

	disconnect := make(chan net.Conn)
	change := make(chan []byte)
	go func() {
		for {
			connection, err := s.listener.Accept()
			if err != nil {
				fmt.Println("Error accepting:", err.Error())
				break
			}
			fmt.Println("client connected")
			go processClient(connection, context.NewContext(s.storage, s.pubsub), s.handler, disconnect, change)
		}
	}()

	for {
		select {
		case disconn := <-disconnect:
			fmt.Println("Deallocating resource for disconnect connection")
			s.removeConnection(disconn)
		case changeBytes := <-change:
			if s.persistFile != nil {
				lenBytes := make([]byte, 4)
				binary.BigEndian.PutUint32(lenBytes, uint32(len(changeBytes)))
				_, err := s.persistFile.Write(append(lenBytes, changeBytes...))
				if err != nil {
					fmt.Println("Error writing to file:", err.Error())
				}
			}
		case <-s.closechan:
		}
	}
}

func (s *Server) removeConnection(conn net.Conn) {
	s.pubsub.Unsub(conn)
}

func getPayloadCommandFromConn(conn net.Conn) ([]byte, error) {

	var cmdstr payload.String
	_, err := cmdstr.ReadFrom(conn)
	if err != nil {
		return nil, err
	}

	switch cmdstr {
	case cmd.PingCode:
		return cmd.PingReadFrom(conn).Bytes(), nil
	case cmd.SetCode:
		return cmd.SetReadFrom(conn).Bytes(), nil
	case cmd.GetCode:
		return cmd.GetReadFrom(conn).Bytes(), nil
	case cmd.SubCode:
		return cmd.SubReadFrom(conn).Bytes(), nil
	case cmd.PubCode:
		return cmd.PubReadFrom(conn).Bytes(), nil
	default:
		return nil, fmt.Errorf("unknown command: %s", cmdstr)
	}
}

func processClient(conn net.Conn, ctx *context.Context, handler handler.Handler, disconnect chan net.Conn, change chan []byte) {
	for {
		buffer, err := getPayloadCommandFromConn(conn)
		if err != nil {
			if err == io.EOF {
				fmt.Printf("Client %s has closed connection\n", conn.RemoteAddr())
			} else {
				fmt.Printf("Unexpected error: %s\n", err)
			}
			break
		}
		if change != nil {
			change <- buffer
		}
		processBytesCommand(readWriter{
			Reader: bytes.NewReader(buffer),
			Writer: conn,
		}, ctx, handler)
	}
	if disconnect != nil {
		disconnect <- conn
	}
	fmt.Println("client closed")
}

func processBytesCommand(rw io.ReadWriter, ctx *context.Context, handler handler.Handler) {
	var cmdstr payload.String
	_, err := cmdstr.ReadFrom(rw)
	if err != nil {
		fmt.Printf("Error reading command: %s\n", err.Error())
		return
	}
	switch cmdstr {
	case cmd.PingCode:
		handler.HandlePing(rw)
	case cmd.SetCode:
		handler.HandleSet(rw, ctx)
	case cmd.GetCode:
		handler.HandleGet(rw, ctx)
	case cmd.SubCode:
		handler.HandleSub(rw, ctx)
	case cmd.PubCode:
		handler.HandlePub(rw, ctx)
	}
}
