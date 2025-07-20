// Package redis houses all the connections related to redis
package redis

import (
	"context"
	"errors"
	"time"
)

// KvStore interface
//
//go:generate mockgen -source redis.go -destination ./mock/mock_redis.go -package mock KvStore
type KvStore interface {
	GetValue(ctx context.Context, key string, result interface{}) error
	GetStringValue(ctx context.Context, key string) (string, error)
	SetValue(ctx context.Context, key string, value interface{}, ttl time.Duration) error
	DeleteValue(ctx context.Context, key string) error
	Connect() error
}

var (
	// ErrConnectionToSourceFailed if the connection to the data source cannot be established
	ErrConnectionToSourceFailed = errors.New("connection to  redis data source cannot be established")
	// ErrFailedToRetrieveValue if there is issue retrieving the value from source
	ErrFailedToRetrieveValue = errors.New("failed to retrieve the value from source")
)
