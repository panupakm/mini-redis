// client to connect to mini redis server
package client

import (
	"fmt"
	"net"

	"github.com/panupakm/mini-redis/lib/payload"
)

type Client struct {
	conn net.Conn
	addr string
}

const (
	Protocal = "tcp"
)

func NewClient() *Client {
	return &Client{}
}

func (c *Client) Connect(addr string) error {
	conn, err := net.Dial(Protocal, addr)
	if err != nil {
		return err
	}
	c.conn = conn
	c.addr = addr
	return nil
}

func (c *Client) Close() error {
	return c.conn.Close()
}

func (c *Client) Ping(msg string) (int64, error) {
	pl := payload.String("ping")
	var n int64 = 0
	o, err := pl.WriterTo(c.conn)
	n += o
	if err != nil {
		return n, err
	}

	pl = payload.String(msg)
	o, err = pl.WriterTo(c.conn)
	if err != nil {
		return 0, err
	}
	return n + int64(0), nil
}

func (c *Client) write(value string) error {
	_, err := c.conn.Write([]byte(value))
	if err != nil {
		fmt.Println("Error writing:", err.Error())
		return err
	}

	return nil
}

func (c *Client) Read() (string, error) {
	buffer := make([]byte, 1024)
	mLen, err := c.conn.Read(buffer)
	if err != nil {
		fmt.Println("Error reading:", err.Error())
		return "", err
	}
	return string(buffer[:mLen]), nil
}

func (c *Client) ReadString() (string, error) {
	var str payload.String
	_, err := str.ReaderFrom(c.conn)
	if err != nil {
		return "", err
	}
	return str.String(), nil
}
