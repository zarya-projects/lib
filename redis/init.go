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
	opts   Options
}

type Base interface {
	Connect() error
	Close()
	HealthCheck() error
}

func New(opts Options, logger *slog.Logger) *SRedis {
	return &SRedis{opts: opts, logger: logger}
}

func (r *SRedis) Connect() error {
	var err error

	switch {
	case r.opts.Addr == "":
		return errors.New("redis: addr is required")
	case r.opts.Password == "":
		return errors.New("redis: password is required")
	case r.opts.DB == 0:
		return errors.New("redis: db is required")
	}

	if r.opts.MaxRetries == 0 {
		r.opts.MaxRetries = 5
	}
	if r.opts.TimeRetry == 0 {
		r.opts.TimeRetry = 5 * time.Second
	}

	for i := 0; i < r.opts.MaxRetries; i++ {
		r.client = redis.NewClient(&redis.Options{
			Addr:     r.opts.Addr,
			Password: r.opts.Password,
			DB:       r.opts.DB,
		})

		if err = r.client.Ping(context.Background()).Err(); err != nil {
			time.Sleep(r.opts.TimeRetry)
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
