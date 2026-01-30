package redis

import (
	"context"
	"time"

	goredis "github.com/go-redis/redis/v8"
)

func New(addr string) *goredis.Client {
	return goredis.NewClient(&goredis.Options{
		Addr:         addr,
		DialTimeout:  5 * time.Second,
		ReadTimeout:  2 * time.Second,
		WriteTimeout: 2 * time.Second,
	})
}

func Ping(ctx context.Context, rdb *goredis.Client) error {
	return rdb.Ping(ctx).Err()
}
