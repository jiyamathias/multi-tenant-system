package controller

import (
	"context"
	"time"
)

// SetValue sets and writes value into redis
func (c *Controller) SetValue(ctx context.Context, key string, value interface{}, ttl time.Duration) error {
	return c.redis.SetValue(ctx, key, value, ttl)
}

// GetStringValue retrieves the value of a key from inside redis as string
func (c *Controller) GetStringValue(ctx context.Context, key string) (string, error) {
	return c.redis.GetStringValue(ctx, key)
}
