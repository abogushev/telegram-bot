package cache

import (
	"context"
	"fmt"
	"time"

	"github.com/go-redis/redis/v9"
)

type redisService struct {
	client *redis.Client
	ctx    context.Context
}

var ErrNotFound = redis.Nil

func NewRedisCache(ctx context.Context, host string, port int) (*redisService, error) {
	options, err := redis.ParseURL(fmt.Sprintf("redis://%s:%v", host, port))
	if err != nil {
		return nil, err
	}
	return &redisService{client: redis.NewClient(options), ctx: ctx}, nil
}

func (r *redisService) Set(key, value string, ttl time.Duration) error {
	return r.client.Set(r.ctx, key, value, ttl).Err()
}

func (r *redisService) Get(key string) (string, error) {
	return r.client.Get(r.ctx, key).Result()
}

func (r *redisService) Delete(key string) error {
	return r.client.Del(r.ctx, key).Err()
}
