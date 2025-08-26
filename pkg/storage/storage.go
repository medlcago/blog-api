package storage

import (
	"context"
	"time"
)

type Storage interface {
	WithNamespace(namespace string) Storage
	Set(ctx context.Context, key string, value interface{}, ttl time.Duration) error
	Get(ctx context.Context, key string) (string, error)
	Delete(ctx context.Context, key string) error
	Exists(ctx context.Context, key string) (bool, error)
}
