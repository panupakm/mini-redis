package db

import (
	"fmt"

	"github.com/panupakm/miniredis/lib/payload"
)

type TypeBuffer struct {
	Typ payload.ValueType
	Buf []byte
}

type Db struct {
	Map map[string]TypeBuffer
}

func NewDb() *Db {
	return &Db{
		Map: make(map[string]TypeBuffer),
	}
}

func NewTypeBuffer(typ payload.ValueType, buf []byte) TypeBuffer {
	return TypeBuffer{
		Typ: typ,
		Buf: buf,
	}
}

func (db *Db) Set(key string, value TypeBuffer) error {
	db.Map[key] = value
	return nil
}

func (db *Db) Get(key string) (TypeBuffer, error) {
	v, ok := db.Map[key]
	if !ok {
		return TypeBuffer{}, fmt.Errorf("key %s not found", key)
	}
	return v, nil
}
