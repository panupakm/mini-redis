package handler

import (
	"io"

	"github.com/panupakm/miniredis/server/context"
)

type ImplHandler struct {
}

type Handler interface {
	HandleGet(rw io.ReadWriter, ctx *context.Context) error
	HandlePing(rw io.ReadWriter) error
	HandlePub(rw io.ReadWriter, ctx *context.Context) error
	HandleSet(rw io.ReadWriter, ctx *context.Context) error
	HandleSub(rw io.ReadWriter, ctx *context.Context) error
}

func NewHandler() Handler {
	return &ImplHandler{}
}

func (h *ImplHandler) HandleGet(rw io.ReadWriter, ctx *context.Context) error {
	return HandleGet(rw, ctx)
}

func (h *ImplHandler) HandlePing(rw io.ReadWriter) error {
	return HandlePing(rw)
}

func (h *ImplHandler) HandlePub(rw io.ReadWriter, ctx *context.Context) error {
	return HandlePub(rw, ctx)
}

func (h *ImplHandler) HandleSub(rw io.ReadWriter, ctx *context.Context) error {
	return HandleSub(rw, ctx)
}

func (h *ImplHandler) HandleSet(rw io.ReadWriter, ctx *context.Context) error {
	return HandleSet(rw, ctx)
}
