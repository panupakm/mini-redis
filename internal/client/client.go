// client to connect to mini redis server
package client

import (
	"net"

	"github.com/panupakm/miniredis/lib/cmd"
	"github.com/panupakm/miniredis/lib/payload"
)

type Client struct {
	conn net.Conn
	addr string
}

type ResultChannel struct {
	*payload.Result
	Str string
	Err error
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

func (c *Client) Ping(msg string) (chan ResultChannel, error) {
	ch := make(chan ResultChannel)
	pl := payload.String(cmd.PingCode)
	var n int64 = 0
	o, err := pl.WriteTo(c.conn)
	n += o
	if err != nil {
		return nil, err
	}

	pl = payload.String(msg)
	o, err = pl.WriteTo(c.conn)
	if err != nil {
		return nil, err
	}

	go func() {
		r, err := c.ReadResult()
		ch <- ResultChannel{Str: string(r.Buffer), Err: err}
	}()

	return ch, nil
}

func (c *Client) Sub(topic string) (chan ResultChannel, error) {
	ch := make(chan ResultChannel)
	pl := payload.String(cmd.SubCode)
	var n int64 = 0
	o, err := pl.WriteTo(c.conn)
	n += o
	if err != nil {
		return nil, err
	}

	pl = payload.String(topic)
	o, err = pl.WriteTo(c.conn)
	if err != nil {
		return nil, err
	}

	go func() {
		r, err := c.ReadResult()
		ch <- ResultChannel{Str: string(r.Buffer), Err: err}
	}()

	return ch, nil
}

func (c *Client) SetString(key string, value string) (chan ResultChannel, error) {
	pl := payload.String(cmd.SetCode)
	o, err := pl.WriteTo(c.conn)
	if err != nil {
		return nil, err
	}

	var n int64 = int64(o)
	pl = payload.String(key)
	o, err = pl.WriteTo(c.conn)
	if err != nil {
		return nil, err
	}

	n += o

	pl = payload.String(value)
	o, err = pl.WriteTo(c.conn)
	if err != nil {
		return nil, err
	}

	ch := make(chan ResultChannel)
	go func() {
		r, err := c.ReadResult()
		ch <- ResultChannel{Result: r, Err: err}
	}()
	return ch, nil
}

func (c *Client) Get(key string) (chan ResultChannel, error) {
	pl := payload.String(cmd.GetCode)
	_, err := pl.WriteTo(c.conn)
	if err != nil {
		return nil, err
	}

	pl = payload.String(key)
	_, err = pl.WriteTo(c.conn)
	if err != nil {
		return nil, err
	}

	ch := make(chan ResultChannel)
	go func() {
		r, err := c.ReadResult()
		ch <- ResultChannel{Result: r, Err: err}
	}()
	return ch, nil
}

func (c *Client) PubString(topic string, msg string) (chan ResultChannel, error) {
	pl := payload.String(cmd.PubCode)
	_, err := pl.WriteTo(c.conn)
	if err != nil {
		return nil, err
	}

	pl = payload.String(topic)
	_, err = pl.WriteTo(c.conn)
	if err != nil {
		return nil, err
	}

	pl = payload.String(msg)
	_, err = pl.WriteTo(c.conn)
	if err != nil {
		return nil, err
	}

	ch := make(chan ResultChannel)
	go func() {
		r, err := c.ReadResult()
		ch <- ResultChannel{Result: r, Err: err}
	}()
	return ch, nil
}

func (c *Client) ReadResult() (*payload.Result, error) {
	var result payload.Result
	_, err := result.ReadFrom(c.conn)
	if err != nil {
		return nil, err
	}
	return &result, nil
}
