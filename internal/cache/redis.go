package cache

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"time"

	"github.com/redis/go-redis/v9"
)


const opTimeout = 200 * time.Millisecond

// Cache is a thin wrapper around the Redis client.
type Cache struct {
	client *redis.Client
}

func New(addr, password string, db int) *Cache {
	client := redis.NewClient(&redis.Options{
		Addr:         addr,
		Password:     password,
		DB:           db,
		DialTimeout:  opTimeout,
		ReadTimeout:  opTimeout,
		WriteTimeout: opTimeout,
		MaxRetries:   -1, // -1 turns retries off, so a dead Redis fails fast
	})
	return &Cache{client: client}
}

// Ping checks the connection. Used once at startup.
func (c *Cache) Ping() error {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	return c.client.Ping(ctx).Err()
}

func (c *Cache) Close() error {
	return c.client.Close()
}

// GetJSON reads a key and decodes it into dest.
// It returns true only on a HIT. Any problem (key missing, Redis down, bad
// JSON) returns false, so the caller simply asks the database.
func (c *Cache) GetJSON(key string, dest any) bool {
	ctx, cancel := context.WithTimeout(context.Background(), opTimeout)
	defer cancel()

	data, err := c.client.Get(ctx, key).Bytes()
	if err != nil {
		// redis.Nil just means "key does not exist" - a normal miss, not an error.
		if !errors.Is(err, redis.Nil) {
			log.Printf("cache: GET %s failed: %v", key, err)
		}
		return false
	}
	if err := json.Unmarshal(data, dest); err != nil {
		log.Printf("cache: bad JSON in %s: %v", key, err)
		return false
	}
	return true
}

// SetJSON stores a value as JSON with a TTL. Errors are only logged:
// failing to cache must never fail the request.
func (c *Cache) SetJSON(key string, value any, ttl time.Duration) {
	data, err := json.Marshal(value)
	if err != nil {
		log.Printf("cache: cannot encode %s: %v", key, err)
		return
	}
	ctx, cancel := context.WithTimeout(context.Background(), opTimeout)
	defer cancel()
	if err := c.client.Set(ctx, key, data, ttl).Err(); err != nil {
		log.Printf("cache: SET %s failed: %v", key, err)
	}
}

// Delete removes keys (invalidation). If it fails, the TTL will clean up later.
func (c *Cache) Delete(keys ...string) {
	ctx, cancel := context.WithTimeout(context.Background(), opTimeout)
	defer cancel()
	if err := c.client.Del(ctx, keys...).Err(); err != nil {
		log.Printf("cache: DEL %v failed: %v", keys, err)
	}
}

// Key builders in one place, so a typo in a key can never create two
// different keys for the same thing.
func AuthorKey(id int) string { return fmt.Sprintf("author:%d", id) }
func BookKey(id int) string   { return fmt.Sprintf("book:%d", id) }