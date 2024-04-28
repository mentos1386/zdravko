package kv

import "time"

type KeyValueStore interface {
	Close() error

	Get(key string) ([]byte, error)
	Set(key string, value []byte, ttl time.Duration) error
	Increment(key string) (int, error)
	Delete(key string) error
	Keys(prefix string) ([]string, error)
}
