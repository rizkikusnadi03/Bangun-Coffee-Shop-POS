package redis

import (
	"context"
	"fmt"

	goredis "github.com/redis/go-redis/v9"
)

// NewRedis menginisialisasi client Redis, mengatur DB index ke 0, dan memverifikasi koneksi dengan ping.
func NewRedis(addr, password string) (*goredis.Client, error) {
	client := goredis.NewClient(&goredis.Options{
		Addr:     addr,
		Password: password,
		DB:       0,
	})

	ctx := context.Background()
	if err := client.Ping(ctx).Err(); err != nil {
		_ = client.Close()
		return nil, fmt.Errorf("failed to ping redis server at %s: %w", addr, err)
	}

	return client, nil
}
