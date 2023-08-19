package cache

import (
	"context"
	"encoding/json"
	"time"

	"github.com/rueian/rueidis"
)

// Result is the result of the cache
type Result struct {
	Value []byte
	// TTLRemaining is the time to live remaining in seconds
	TTLRemaining int64
}

// Cache is the interface for the cache
type Cache interface {
	// Get gets the value from the cache
	// ttl is the time to live in seconds
	Get(ctx context.Context, key string, ttl int64) (*Result, error)
	// Set sets the value to the cache
	// ttl is the time to live in seconds
	Set(ctx context.Context, key string, value any, ttl int64) error
}

func New(client rueidis.Client) Cache {
	return &cache{
		client: client,
	}
}

type cache struct {
	client rueidis.Client
}

func (c *cache) Get(ctx context.Context, key string, ttl int64) (*Result, error) {
	command := c.client.B().
		Get().
		Key(key).
		Cache()

	res := c.client.DoCache(ctx, command, time.Second*time.Duration(ttl))
	if res.Error() != nil {
		if res.Error() == rueidis.Nil {
			return nil, nil
		}
		return nil, res.Error()
	}

	value, err := res.AsBytes()
	if err != nil {
		return nil, err
	}

	result := Result{
		TTLRemaining: res.CacheTTL(),
		Value:        value,
	}
	return &result, nil
}

func (c *cache) Set(ctx context.Context, key string, value any, ttl int64) error {
	bytes, err := json.Marshal(value)
	if err != nil {
		return err
	}

	command := c.client.B().
		Set().
		Key(key).
		Value(rueidis.BinaryString(bytes)).
		ExSeconds(ttl).
		Build()

	err = c.client.Do(ctx, command).Error()
	if err != nil {
		return err
	}
	return nil
}
