package avitoredis

import (
	"context"
	"time"

	"github.com/redis/go-redis/v9"
)



type Cache struct {
	rdb *redis.Client
}

func NewCache(
	rd *redis.Client,
) *Cache {
	return &Cache{
		rdb: rd,
	}
}

func (c *Cache) Set(
	ctx context.Context, key string, data []byte, ttl time.Duration,
) ([]byte, error) {
	return c.rdb.Set(ctx, key, data, ttl).Bytes()
}

func (c *Cache) Get(
	ctx context.Context, key string,
) ([]byte, error) {
	return c.rdb.Get(ctx, key).Bytes()
}

func (c *Cache) Delete(
	ctx context.Context, key string,
) error {
	return c.rdb.Del(ctx, key).Err()
}
