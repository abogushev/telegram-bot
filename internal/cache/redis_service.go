package cache

import (
	"context"
	"fmt"
	"time"

	"github.com/go-redis/redis/v9"
)

type RedisService struct {
	client *redis.Client
}

var ErrNotFound = redis.Nil

func NewRedisCache(host string, port int) (*RedisService, error) {
	options, err := redis.ParseURL(fmt.Sprintf("redis://%s:%v", host, port))
	if err != nil {
		return nil, err
	}
	return &RedisService{client: redis.NewClient(options)}, nil
}

func (r *RedisService) Set(ctx context.Context, key, value string, ttl time.Duration) error {
	return r.client.Set(ctx, key, value, ttl).Err()
}

func (r *RedisService) Get(ctx context.Context, key string) (string, error) {
	return r.client.Get(ctx, key).Result()
}

func (r *RedisService) Delete(ctx context.Context, key string) error {
	return r.client.Del(ctx, key).Err()
}
