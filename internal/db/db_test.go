// Package db provides functions for interacting with a database.
package db

import (
	"testing"

	"github.com/panupakm/miniredis/payload"
	"github.com/stretchr/testify/assert"
)

func TestNewDb(t *testing.T) {
	tests := []struct {
		name string
		want *Db
	}{
		{
			name: "new DB",
			want: &Db{
				pairs: make(map[string]payload.General),
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := NewDb()
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestDb_Set(t *testing.T) {
	type fields struct {
		pairs map[string]payload.General
	}
	type args struct {
		key   string
		value payload.General
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		{
			name: "set key",
			fields: fields{
				pairs: make(map[string]payload.General),
			},
			args: args{
				key:   "key",
				value: *payload.NewGeneral(payload.StringType, []byte("value")),
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			db := &Db{
				pairs: tt.fields.pairs,
			}
			err := db.Set(tt.args.key, tt.args.value)
			assert.NoError(t, err)
			assert.Equal(t, tt.args.value, db.pairs[tt.args.key])
		})
	}
}

func TestDb_Get(t *testing.T) {
	type fields struct {
		pairs map[string]payload.General
	}
	type args struct {
		key string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    payload.General
		wantErr bool
	}{
		{
			name: "get key",
			fields: fields{
				pairs: map[string]payload.General{
					"key": *payload.NewGeneral(payload.StringType, []byte("value")),
				},
			},
			args: args{
				key: "key",
			},
			want: *payload.NewGeneral(payload.StringType, []byte("value")),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			db := &Db{
				pairs: tt.fields.pairs,
			}
			got, err := db.Get(tt.args.key)
			assert.NoError(t, err)
			assert.Equal(t, tt.want, got)
		})
	}
}
