package redis

import (
	"context"
	"fmt"
	"time"

	redis "github.com/go-redis/redis/v8"
	"github.com/rs/zerolog"

	"codematic/pkg/environment"
	"codematic/pkg/helper"
)

const packageNameRedis = "store.redis"

type (
	// Redis store object
	Redis struct {
		env             *environment.Env
		logger          zerolog.Logger
		connectionError error
		client          *redis.Client
		url             string
	}
)

// NewRedis creates a new Redis object as a KeyValue instance
func NewRedis(e *environment.Env, z zerolog.Logger, url string) *KvStore {
	log := z.With().Str(helper.LogStrPackageLevel, packageNameRedis).Logger()
	r := &Redis{
		env:    e,
		logger: log,
		url:    url,
	}
	// connect to the storage
	r.connectionError = r.Connect()
	if r.connectionError != nil {
		fmt.Printf("There was an error connecting to redis: %v", r.connectionError.Error())
	}
	k := KvStore(r)
	return &k
}

// Connect handles connection to data source server implementation
func (r *Redis) Connect() error {
	ctx := context.Background()

	opt, err := redis.ParseURL(r.url)
	if err != nil {
		r.logger.Err(err).Msg("unable to parse redis server url")
		r.connectionError = err
		return err
	}

	rdb := redis.NewClient(opt)
	st := rdb.Ping(ctx)
	if err := st.Err(); err != nil {
		r.logger.Err(err).Msg("connection to redis server failed")
		r.connectionError = err
		return err
	}
	r.logger.Info().Msg("[success] connected to redis server")
	r.connectionError = nil // connection did NOT fail.a
	r.client = rdb          // set the client
	return nil
}

// GetValue retrieves the value // GetValue retrieves the value of a key from inside redis
func (r *Redis) GetValue(ctx context.Context, key string, result interface{}) error {
	if r.connectionError != nil {
		// attempt to reconnect
		err := r.Connect()
		if err != nil {
			return ErrConnectionToSourceFailed
		}
	}
	// hopefully the connection to the store is okay
	err := r.client.Get(ctx, key).Scan(result)
	if err != nil {
		// connecting issue or not able to retrieve from server
		r.logger.Err(err).Str("key", key).Msg(ErrFailedToRetrieveValue.Error())
		return ErrFailedToRetrieveValue
	}
	return nil
}

// GetStringValue retrieves the value of a key from inside redis as string
func (r *Redis) GetStringValue(ctx context.Context, key string) (string, error) {
	if r.connectionError != nil {
		// attempt to reconnect
		err := r.Connect()
		if err != nil {
			return "", ErrConnectionToSourceFailed
		}
	}

	// hopefully the connection to the store is okay
	return r.client.Get(ctx, key).Val(), nil
}

// SetValue sets and writes value into redis
func (r *Redis) SetValue(ctx context.Context, key string, value interface{}, ttl time.Duration) error {
	if r.connectionError != nil {
		// attempt to reconnect
		err := r.Connect()
		if err != nil {
			return ErrConnectionToSourceFailed
		}
	}

	// hopefully the connection to the store is okay
	res := r.client.Set(ctx, key, value, ttl)
	if res.Err() != nil {
		return res.Err()
	}

	return nil
}

// DeleteValue will delete redis key
func (r *Redis) DeleteValue(ctx context.Context, key string) error {
	res := r.client.Del(ctx, key)
	if res.Err() != nil {
		return res.Err()
	}
	return nil
}
