package redis

import (
	"context"
	"errors"
	"log/slog"
	"time"

	"github.com/redis/go-redis/v9"
)

type Options struct {
	Addr       string
	Password   string
	DB         int
	MaxRetries int
	TimeRetry  time.Duration
}

type SRedis struct {
	client *redis.Client
	logger *slog.Logger
}

type Base interface {
	Connect() error
	Close()
	HealthCheck() error
}

func (r *SRedis) Connect(opts Options) error {
	var err error

	switch {
	case opts.Addr == "":
		return errors.New("redis: addr is required")
	case opts.Password == "":
		return errors.New("redis: password is required")
	case opts.DB == 0:
		return errors.New("redis: db is required")
	}

	if opts.MaxRetries == 0 {
		opts.MaxRetries = 5
	}
	if opts.TimeRetry == 0 {
		opts.TimeRetry = 5 * time.Second
	}

	for i := 0; i < opts.MaxRetries; i++ {
		r.client = redis.NewClient(&redis.Options{
			Addr:     opts.Addr,
			Password: opts.Password,
			DB:       opts.DB,
		})

		if err = r.client.Ping(context.Background()).Err(); err != nil {
			time.Sleep(opts.TimeRetry)
			continue
		}

		break
	}

	return nil
}

func (r *SRedis) Close() error {
	return r.client.Close()
}

func (r *SRedis) HealthCheck() error {
	return r.client.Ping(context.Background()).Err()
}
