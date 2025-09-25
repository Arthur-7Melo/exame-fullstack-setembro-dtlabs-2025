package redis

import (
	"context"
	"os"
	"time"

	"github.com/go-redis/redis/v8"
)

type Client struct {
	*redis.Client
}

func NewRedisClient() *Client {
	redisURL := os.Getenv("REDIS_URL")
	if redisURL == "" {
		redisURL = "redis://redis:6379"
	}

	opts, err := redis.ParseURL(redisURL)
	if err != nil {
		panic(err)
	}

	rdb := redis.NewClient(opts)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := rdb.Ping(ctx).Err(); err != nil {
		panic(err)
	}

	return &Client{rdb}
}

func (c *Client) Publish(ctx context.Context, channel string, message interface{}) error {
	return c.Client.Publish(ctx, channel, message).Err()
}

func (c *Client) Subscribe(ctx context.Context, channels ...string) *redis.PubSub {
	return c.Client.Subscribe(ctx, channels...)
}