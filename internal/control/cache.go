package control

import (
	"context"
	"errors"
	"time"

	"github.com/redis/go-redis/v9"
)

// Cache is the small get/set surface the server needs. *redisCache implements
// it against Redis; noopCache (when REDIS_URL is unset or Redis is down)
// makes every call a miss so all endpoints degrade to their uncached path —
// nothing breaks without Redis.
type Cache interface {
	Get(ctx context.Context, key string) (string, bool)
	Set(ctx context.Context, key, val string, ttl time.Duration)
	Del(ctx context.Context, key string)
	Ping(ctx context.Context) error
}

type redisCache struct{ client *redis.Client }

// ConnectCache dials Redis and verifies with PING.
func ConnectCache(ctx context.Context, url string) (Cache, error) {
	opt, err := redis.ParseURL(url)
	if err != nil {
		return nil, err
	}
	client := redis.NewClient(opt)
	pingCtx, cancel := context.WithTimeout(ctx, 3*time.Second)
	defer cancel()
	if err := client.Ping(pingCtx).Err(); err != nil {
		_ = client.Close()
		return nil, err
	}
	return &redisCache{client: client}, nil
}

func (c *redisCache) Get(ctx context.Context, key string) (string, bool) {
	val, err := c.client.Get(ctx, key).Result()
	if errors.Is(err, redis.Nil) || err != nil {
		return "", false
	}
	return val, true
}

func (c *redisCache) Set(ctx context.Context, key, val string, ttl time.Duration) {
	_ = c.client.Set(ctx, key, val, ttl).Err()
}

func (c *redisCache) Del(ctx context.Context, key string) {
	_ = c.client.Del(ctx, key).Err()
}

func (c *redisCache) Ping(ctx context.Context) error {
	return c.client.Ping(ctx).Err()
}

type noopCache struct{}

func (noopCache) Get(context.Context, string) (string, bool) { return "", false }
func (noopCache) Set(context.Context, string, string, time.Duration) {}
func (noopCache) Del(context.Context, string)                {}
func (noopCache) Ping(context.Context) error                 { return errors.New("cache not configured") }

// cache key space
const (
	cacheKeyHealthServices = "talon:health:services"
	cacheKeyAnalyzePrefix  = "talon:analyze:"
	cacheKeySessionPrefix  = "talon:sess:"
)

// TTLs — tuned so the dashboard feels instant without going stale.
const (
	healthServicesTTL = 5 * time.Second
	analyzeTTL        = 24 * time.Hour
	sessionCacheTTL   = 60 * time.Second
)
