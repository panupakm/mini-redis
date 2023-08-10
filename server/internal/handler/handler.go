package handler

import (
	"net"

	"github.com/panupakm/miniredis/server/context"
)

type ImplHandler struct {
}

type Handler interface {
	HandleGet(conn net.Conn, ctx *context.Context) error
	HandlePing(conn net.Conn) error
	HandlePub(conn net.Conn, ctx *context.Context) error
	HandleSet(conn net.Conn, ctx *context.Context) error
	HandleSub(conn net.Conn, ctx *context.Context) error
}

func NewHandler() Handler {
	return &ImplHandler{}
}

func (h *ImplHandler) HandleGet(conn net.Conn, ctx *context.Context) error {
	return HandleGet(conn, ctx)
}

func (h *ImplHandler) HandlePing(conn net.Conn) error {
	return HandlePing(conn)
}

func (h *ImplHandler) HandlePub(conn net.Conn, ctx *context.Context) error {
	return HandlePub(conn, ctx)
}

func (h *ImplHandler) HandleSub(conn net.Conn, ctx *context.Context) error {
	return HandleSub(conn, ctx)
}

func (h *ImplHandler) HandleSet(conn net.Conn, ctx *context.Context) error {
	return HandleSet(conn, ctx)
}
