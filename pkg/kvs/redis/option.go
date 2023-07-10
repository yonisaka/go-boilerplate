package redis

import (
	"errors"
	"fmt"
	"time"
)

// Option configures Redis client.
type Option func(r *redisClient) error

var defaultOptions = []Option{
	WithNetwork("tcp"),
	WithAddr("127.0.0.1:6379"),
	WithDB(0),
	WithMaxRetries(3),
	WithMinRetryBackoff("8ms"),
	WithMaxRetryBackoff("1s"),
	WithDialTimeout("5s"),
	WithWriteTimeout("3s"),
	WithReadTimeout("3s"),
	WithPoolSize(10),
	WithMinIdleConns(7),
}

// WithNetwork returns an option that network.
func WithNetwork(str string) Option {
	return func(r *redisClient) error {
		if len(str) == 0 {
			return errors.New("failed to set redis.network")
		}

		r.network = str

		return nil
	}
}

// WithAddr returns an option that set address.
func WithAddr(addr string) Option {
	return func(r *redisClient) error {
		if len(addr) == 0 {
			return errors.New("failed to set redis.addr")
		}

		r.addr = addr

		return nil
	}
}

// WithPassword returns an option that set password.
func WithPassword(pass string) Option {
	return func(r *redisClient) error {
		if len(pass) == 0 {
			return errors.New("failed to set redis.password")
		}

		r.password = pass

		return nil
	}
}

// WithDB returns an option that set db.
func WithDB(db int) Option {
	return func(r *redisClient) error {
		if db < 0 {
			return fmt.Errorf("failed to set redis.db: %d", db)
		}

		r.db = db

		return nil
	}
}

// WithMaxRetries returns an option that set max retries.
func WithMaxRetries(n int) Option {
	return func(r *redisClient) error {
		if n < 0 {
			return fmt.Errorf("failed to set redis.maxRetries: %d", n)
		}

		r.maxRetries = n

		return nil
	}
}

// WithMinRetryBackoff returns an option that set min retry backoff.
func WithMinRetryBackoff(str string) Option {
	return func(r *redisClient) error {
		d, err := parseDuration(str)
		if err != nil {
			return fmt.Errorf("failed to set redis.minRetryBackoff: %w", err)
		}

		r.minRetryBackoff = d

		return nil
	}
}

// WithMaxRetryBackoff returns an option that set max retry backoff.
func WithMaxRetryBackoff(str string) Option {
	return func(r *redisClient) error {
		d, err := parseDuration(str)
		if err != nil {
			return fmt.Errorf("failed to set redis.maxRetryBackoff: %w", err)
		}

		r.maxRetryBackoff = d

		return nil
	}
}

// WithDialTimeout returns an option that set dial timeout.
func WithDialTimeout(str string) Option {
	return func(r *redisClient) error {
		d, err := parseDuration(str)
		if err != nil {
			return fmt.Errorf("failed to set redis.dialTimeout: %s: %w", str, err)
		}

		if d < 0 {
			return fmt.Errorf("failed to set redis.dialTimeout: %d", d)
		}

		r.dialTimeout = d

		return nil
	}
}

// WithReadTimeout returns an option that set read timeout.
func WithReadTimeout(str string) Option {
	return func(r *redisClient) error {
		d, err := parseDuration(str)
		if err != nil {
			return fmt.Errorf("failed to set redis.readTimeout: %w", err)
		}

		r.readTimeout = d

		return nil
	}
}

// WithWriteTimeout returns an option that set write timeout.
func WithWriteTimeout(str string) Option {
	return func(r *redisClient) error {
		d, err := parseDuration(str)
		if err != nil {
			return fmt.Errorf("failed to set redis.writeTimeout: %w", err)
		}

		r.writeTimeout = d

		return nil
	}
}

// WithPoolSize returns an option that set pool size.
func WithPoolSize(n int) Option {
	return func(r *redisClient) error {
		if n < 0 {
			return fmt.Errorf("failed to set redis.poolSize %d", n)
		}

		r.poolSize = n

		return nil
	}
}

// WithMinIdleConns returns an option that set min idle conns.
func WithMinIdleConns(n int) Option {
	return func(r *redisClient) error {
		if n < 0 {
			return fmt.Errorf("failed to set redis.minIdleConns: %d", n)
		}

		r.minIdleConns = n

		return nil
	}
}

func parseDuration(str string) (time.Duration, error) {
	if len(str) == 0 {
		return 0, errors.New("empty string")
	}

	d, err := time.ParseDuration(str)
	if err != nil {
		return 0, fmt.Errorf("failed to parse: %w", err)
	}

	if d < 0 {
		return -1, nil
	}

	return d, nil
}
