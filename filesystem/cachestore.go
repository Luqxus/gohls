package filesystem

import (
	"context"

	"github.com/redis/go-redis/v9"
)

type CacheStore interface {
	retrieveVideoData(ctx context.Context, key string) ([]byte, error)
	setVideoData(ctx context.Context, key string, videoData []byte) error
}

type RedisStore struct {
	client *redis.Client
}

func NewRedisStore() *RedisStore {
	return &RedisStore{
		client: redis.NewClient(&redis.Options{
			Addr: "localhost:6379",
			DB:   0,
		}),
	}
}

// retrieve video data from redis cache
func (r *RedisStore) retrieveVideoData(ctx context.Context, key string) ([]byte, error) {
	return r.client.Get(ctx, key).Bytes()
}

// commit video data to redis cache
func (r *RedisStore) setVideoData(ctx context.Context, key string, videoData []byte) error {
	return r.client.Set(ctx, key, videoData, 0).Err()
}
