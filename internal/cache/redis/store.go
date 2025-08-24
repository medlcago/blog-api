package redis

import (
	"blog-api/config"
	"blog-api/internal/cache"
	"context"
	"time"

	"github.com/redis/go-redis/v9"
)

var ErrRecordNotFound = redis.Nil

type Storage struct {
	rdb        *redis.Client
	namespace  string
	keyBuilder *cache.KeyBuilder
}

func New(redisCfg config.RedisConfig) (*Storage, error) {
	rdb := redis.NewClient(&redis.Options{
		Addr:     redisCfg.Addr,
		Password: redisCfg.Password,
		DB:       redisCfg.DB,
	})

	if err := rdb.Ping(context.Background()).Err(); err != nil {
		return nil, err
	}

	return &Storage{
		rdb:        rdb,
		keyBuilder: cache.NewKeyBuilder("", ":"),
	}, nil
}

func (s *Storage) WithNamespace(namespace string) cache.Storage {
	if s.namespace != "" {
		namespace = s.namespace + "_" + namespace
	}
	return &Storage{
		rdb:        s.rdb,
		namespace:  namespace,
		keyBuilder: cache.NewKeyBuilder(namespace, ":"),
	}
}

func (s *Storage) Set(ctx context.Context, key string, value interface{}, ttl time.Duration) error {
	key = s.keyBuilder.Build(key)
	return s.rdb.Set(ctx, key, value, ttl).Err()
}

func (s *Storage) Get(ctx context.Context, key string) (string, error) {
	key = s.keyBuilder.Build(key)
	return s.rdb.Get(ctx, key).Result()
}

func (s *Storage) Delete(ctx context.Context, key string) error {
	key = s.keyBuilder.Build(key)
	return s.rdb.Del(ctx, key).Err()
}

func (s *Storage) Exists(ctx context.Context, key string) (bool, error) {
	key = s.keyBuilder.Build(key)
	val, err := s.rdb.Exists(ctx, key).Result()
	if err != nil {
		return false, err
	}
	return val == 1, nil
}
