package redis

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
	"github.com/yonisaka/go-boilerplate/pkg/kvs"
)

type redisClient struct {
	*redis.Client

	network         string
	addr            string
	password        string
	maxRetries      int
	minRetryBackoff time.Duration
	maxRetryBackoff time.Duration
	db              int
	dialTimeout     time.Duration
	readTimeout     time.Duration
	writeTimeout    time.Duration
	poolSize        int
	minIdleConns    int
}

// New returns KVS interface implementations.
func New(opts ...Option) (kvs.Client, error) {
	r := new(redisClient)

	for _, opt := range append(defaultOptions, opts...) {
		if err := opt(r); err != nil {
			return nil, fmt.Errorf("failed to apply option: %w", err)
		}
	}

	//nolint:exhaustivestruct
	r.Client = redis.NewClient(&redis.Options{
		Network:         r.network,
		Addr:            r.addr,
		Password:        r.password,
		DB:              r.db,
		MaxRetries:      r.maxRetries,
		MinRetryBackoff: r.minRetryBackoff,
		MaxRetryBackoff: r.maxRetryBackoff,
		DialTimeout:     r.dialTimeout,
		ReadTimeout:     r.readTimeout,
		WriteTimeout:    r.writeTimeout,
		PoolSize:        r.poolSize,
		MinIdleConns:    r.minIdleConns,
	})

	return r, nil
}

func (r *redisClient) Set(ctx context.Context, key string, value interface{}, expire time.Duration) (interface{}, error) {
	val, err := r.Client.Set(ctx, key, value, expire).Result()
	if err != nil {
		return nil, fmt.Errorf("failed to execute set command of redis. key: %v, value: %v: %w", key, value, err)
	}

	return val, nil
}

func (r *redisClient) HSet(ctx context.Context, key string, field string, value interface{}) (interface{}, error) {
	val, err := r.Client.HSet(ctx, key, field, value).Result()
	if err != nil {
		return nil, fmt.Errorf("failed to execute set command of redis. key: %v, field: %v, value: %v: %w", key, field, value, err)
	}

	return val, nil
}

func (r *redisClient) Get(ctx context.Context, key string) (interface{}, error) {
	val, err := r.Client.Get(ctx, key).Result()
	if errors.Is(err, redis.Nil) {
		return nil, fmt.Errorf("failed to find a value by the key `%s`: %w", key, err)
	} else if err != nil {
		return nil, fmt.Errorf("failed to execute get command. key: %v: %w", key, err)
	}

	return val, nil
}

func (r *redisClient) HGetAll(ctx context.Context, key string) (interface{}, error) {
	val, err := r.Client.HGetAll(ctx, key).Result()
	if errors.Is(err, redis.Nil) {
		return nil, fmt.Errorf("failed to find a value by the key `%s`: %w", key, err)
	} else if err != nil {
		return nil, fmt.Errorf("failed to execute get command. key: %v: %w", key, err)
	}

	return val, nil
}

func (r *redisClient) Close() error {
	if err := r.Client.Close(); err != nil {
		return fmt.Errorf("failed to close redis connection: %w", err)
	}

	return nil
}

func (r *redisClient) FlushAll(ctx context.Context) error {
	err := r.Client.FlushAll(ctx).Err()

	return err
}
