package goe

import (
	context "context"
	"fmt"
	"time"

	"github.com/gomodule/redigo/redis"
)

type Store interface {
	// Add appends provided data into the store.
	Add(ctx context.Context, createdAt time.Time, id string, data []byte) error
}

// NewRedisStore returns a Store implementation that is using a Redis instance
// as the storage engine.
func NewRedisStore(addr string) Store {
	return redisStore{
		pool: &redis.Pool{
			MaxIdle:     3,
			IdleTimeout: 240 * time.Second,
			Dial: func() (redis.Conn, error) {
				return redis.Dial("tcp", addr)
			},
		},
	}
}

type redisStore struct {
	pool *redis.Pool
}

func (s redisStore) Add(ctx context.Context, createdAt time.Time, id string, data []byte) error {
	c, err := s.pool.GetContext(ctx)
	if err != nil {
		return fmt.Errorf("acquire pool connection: %w", err)
	}
	defer c.Close()

	payload := fmt.Sprintf("%d\t%s\t%x", createdAt.UnixNano(), id, data)
	if _, err := c.Do("LPUSH", "randomer", payload); err != nil {
		return fmt.Errorf("redis lpush: %w", err)
	}
	return nil
}
