package redis

import (
	"context"
	"fmt"
	"testing"
)

func TestRedisClient_Set(t *testing.T) {
	ctx := context.Background()
	redis, err := New()
	if err != nil {
		t.Fatal(err)
	}

	defer redis.Close()

	_, err = redis.Set(ctx, "key", "value", 0)
	if err != nil {
		t.Fatal(err)
	}

	val, err := redis.Get(ctx, "key")
	if err != nil {
		t.Fatal(err)
	}

	t.Log(fmt.Sprintf("success get value: %s", val))

	redis.FlushAll(ctx)
}

func TestRedisClient_HSet(t *testing.T) {
	ctx := context.Background()
	redis, err := New()
	if err != nil {
		t.Fatal(err)
	}
	defer redis.Close()

	_, err = redis.HSet(ctx, "key", "field", "value")
	if err != nil {
		t.Fatal(err)
	}

	_, err = redis.HSet(ctx, "key", "field2", "value2")
	if err != nil {
		t.Fatal(err)
	}

	val, err := redis.HGetAll(ctx, "key")
	if err != nil {
		t.Fatal(err)
	}

	t.Log(fmt.Sprintf("success get value: %s", val))

	redis.FlushAll(ctx)
}
