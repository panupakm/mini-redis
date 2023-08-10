// Package db provides functions for interacting with a database.
package db

import (
	"fmt"

	"github.com/panupakm/miniredis/payload"
)

type Db struct {
	pairs map[string]payload.General
}

type SetGet interface {
	Set(key string, value payload.General) error
	Get(key string) (payload.General, error)
}

func NewDb() *Db {
	return &Db{
		pairs: make(map[string]payload.General),
	}
}

func (db *Db) Set(key string, value payload.General) error {
	db.pairs[key] = value
	return nil
}

func (db *Db) Get(key string) (payload.General, error) {
	v, ok := db.pairs[key]
	if !ok {
		return payload.General{}, fmt.Errorf("key %s not found", key)
	}
	return v, nil
}
