package redis

import (
	"blog-api/pkg/storage"
	"context"
	"time"

	"github.com/redis/go-redis/v9"
)

var ErrRecordNotFound = redis.Nil

type Storage struct {
	rdb       *redis.Client
	namespace string
	keyFunc   storage.KeyFunc
}

func New(rdb *redis.Client, keyFunc storage.KeyFunc) (storage.Storage, error) {
	s := &Storage{
		rdb:     rdb,
		keyFunc: storage.DefaultKeyFunc,
	}

	if keyFunc != nil {
		s.keyFunc = keyFunc
	}

	return s, nil
}

func (s *Storage) WithNamespace(namespace string) storage.Storage {
	if s.namespace != "" {
		namespace = s.namespace + "_" + namespace
	}
	return &Storage{
		rdb:       s.rdb,
		namespace: namespace,
		keyFunc:   s.keyFunc,
	}
}

func (s *Storage) Set(ctx context.Context, key string, value interface{}, ttl time.Duration) error {
	key = s.keyFunc(s.namespace, key)
	return s.rdb.Set(ctx, key, value, ttl).Err()
}

func (s *Storage) Get(ctx context.Context, key string) (string, error) {
	key = s.keyFunc(s.namespace, key)
	return s.rdb.Get(ctx, key).Result()
}

func (s *Storage) Delete(ctx context.Context, key string) error {
	key = s.keyFunc(s.namespace, key)
	return s.rdb.Del(ctx, key).Err()
}

func (s *Storage) Exists(ctx context.Context, key string) (bool, error) {
	key = s.keyFunc(s.namespace, key)
	val, err := s.rdb.Exists(ctx, key).Result()
	if err != nil {
		return false, err
	}
	return val == 1, nil
}
