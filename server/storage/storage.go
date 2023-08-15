// Package db provides functions for interacting with a database.
package storage

import (
	"github.com/panupakm/miniredis/internal/payload"
	"github.com/panupakm/miniredis/server/storage/internal"
)

type Storage interface {
	Set(key string, value payload.General) error
	Get(key string) (payload.General, error)
}

func NewDefaultStorage() Storage {
	return internal.NewStorage()
}
