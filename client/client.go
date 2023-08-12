// Package client implementation streamss requests from client to server.
package client

import (
	"crypto/tls"
	"fmt"
	"net"

	"github.com/panupakm/miniredis/payload"
	"github.com/panupakm/miniredis/request"
)

type Dialer interface {
	Dial(network, addr string, config *tls.Config) (net.Conn, error)
}

type ImplDialer struct {
	Dialer
}

func (*ImplDialer) Dial(network, addr string, config *tls.Config) (net.Conn, error) {
	if config != nil {
		return tls.Dial(network, addr, config)
	}
	return net.Dial(network, addr)
}

type Client struct {
	conn   net.Conn
	addr   string
	dial   Dialer
	config *tls.Config
}

type ResultChannel struct {
	*payload.Result
	Err error
}

const (
	protocal = "tcp"
)

func NewClient() *Client {
	return &Client{
		dial: &ImplDialer{},
	}
}

func (c *Client) Connect(addr string, config *tls.Config) error {
	conn, err := c.dial.Dial(protocal, addr, config)
	if err != nil {
		return err
	}
	c.conn = conn
	c.addr = addr
	c.config = config
	return nil
}

func (c *Client) Close() error {
	if c.conn == nil {
		return fmt.Errorf("connection is not established")
	}
	return c.conn.Close()
}

func (c *Client) Ping(msg string) (chan ResultChannel, error) {
	ch := make(chan ResultChannel)
	pl := payload.String(request.PingCode)
	var n int64 = 0
	o, err := pl.WriteTo(c.conn)
	n += o
	if err != nil {
		return nil, err
	}

	pl = payload.String(msg)
	if _, err := pl.WriteTo(c.conn); err != nil {
		return nil, err
	}

	go func() {
		r, err := c.readResult()
		ch <- ResultChannel{
			Result: r,
			Err:    err,
		}
	}()

	return ch, nil
}

func (c *Client) Sub(topic string) (*Subsriber, chan ResultChannel, error) {
	ch := make(chan ResultChannel)
	pl := payload.String(request.SubCode)
	_, err := pl.WriteTo(c.conn)
	if err != nil {
		return nil, nil, err
	}

	pl = payload.String(topic)
	_, err = pl.WriteTo(c.conn)
	if err != nil {
		return nil, nil, err
	}

	go func() {
		r, err := c.readResult()
		ch <- ResultChannel{Result: r, Err: err}
	}()

	return &Subsriber{
		messages: make(chan payload.SubMsg),
		conn:     c.conn,
	}, ch, nil
}

func (c *Client) SetString(key string, value string) (chan ResultChannel, error) {
	pl := payload.String(request.SetCode)
	if _, err := pl.WriteTo(c.conn); err != nil {
		return nil, err
	}

	pl = payload.String(key)
	if _, err := pl.WriteTo(c.conn); err != nil {
		return nil, err
	}

	pl = payload.String(value)
	if _, err := pl.WriteTo(c.conn); err != nil {
		return nil, err
	}

	ch := make(chan ResultChannel)
	go func() {
		r, err := c.readResult()
		ch <- ResultChannel{Result: r, Err: err}
	}()
	return ch, nil
}

func (c *Client) Get(key string) (chan ResultChannel, error) {
	pl := payload.String(request.GetCode)
	if _, err := pl.WriteTo(c.conn); err != nil {
		return nil, err
	}

	pl = payload.String(key)
	if _, err := pl.WriteTo(c.conn); err != nil {
		return nil, err
	}

	ch := make(chan ResultChannel)
	go func() {
		r, err := c.readResult()
		ch <- ResultChannel{Result: r, Err: err}
	}()
	return ch, nil
}

func (c *Client) PubString(topic string, msg string) (chan ResultChannel, error) {
	err := request.PubStringTo(c.conn, topic, msg)
	if err != nil {
		return nil, err
	}

	ch := make(chan ResultChannel)
	go func() {
		r, err := c.readResult()
		ch <- ResultChannel{Result: r, Err: err}
	}()
	return ch, nil
}

func (c *Client) readResult() (*payload.Result, error) {
	var result payload.Result
	if _, err := result.ReadFrom(c.conn); err != nil {
		return nil, err
	}
	return &result, nil
}
