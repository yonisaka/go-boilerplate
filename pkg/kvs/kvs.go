package kvs

//go:generate rm -f ./kvs_mock.go
//go:generate mockgen -destination kvs_mock.go -package kvs -mock_names Client=GoMockClient -source kvs.go

import (
	"context"
	"time"
)

// Session represents an authenticated peer session.
type Session struct {
	ID     string `json:"id"`
	Handle string `json:"name"`
}

// Client is an interface for KVS cache.
type Client interface {
	// Set sets the value with expiration time.
	Set(ctx context.Context, key string, value interface{}, expire time.Duration) (interface{}, error)
	// HSet sets the value with expiration time.
	HSet(ctx context.Context, key string, field string, value interface{}) (interface{}, error)
	// Get gets the value by the given key.
	Get(ctx context.Context, key string) (interface{}, error)
	// HGetAll gets the value by the given key.
	HGetAll(ctx context.Context, key string) (interface{}, error)
	// Close closes the connection of KVS client.
	Close() error
	// FlushAll flushes all keys.
	FlushAll(ctx context.Context) error
}
