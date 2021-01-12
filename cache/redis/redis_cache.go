package redis

import (
	"context"
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	"github.com/go-redis/redis/v8"
	jsoniter "github.com/json-iterator/go"
	"github.com/pkg/errors"
	"github.com/thaitanloi365/gocore/cache/types"
)

const name = "redis"

// Errors
var (
	ErrKeyNotFound = errors.New("Key not found")
)

var instance *Client

// Config config
type Config struct {
	RedisAddress      string
	RedisDBMode       int
	RedisPassword     string
	Namespace         string
	Logger            types.Logger
	DefaultExpiration time.Duration
}

// Client client
type Client struct {
	rdb       *redis.Client
	config    *Config
	namespace string
	logger    types.Logger
}

// New get the redis client
func New(config *Config) *Client {
	var rdb = redis.NewClient(&redis.Options{
		Addr:       config.RedisAddress,
		DB:         config.RedisDBMode,
		Password:   config.RedisPassword,
		MaxRetries: 3,
	})

	var instance = &Client{
		config:    config,
		namespace: "gocore_redis_cache",
		rdb:       rdb,
		logger:    log.New(os.Stdout, "\r\n", 0),
	}

	if config.Namespace != "" {
		instance.namespace = config.Namespace
	}

	if config.Logger != nil {
		instance.logger = config.Logger
	}

	for i := 0; i < 10; i++ {
		if err := instance.rdb.Ping(context.Background()).Err(); err != nil {
			instance.logger.Printf("[%d/%d] Connect to Redis error: %v\n", i, 10, err)
		}
		time.Sleep(1)
	}

	return instance

}

// GetInstance get instance
func GetInstance() *Client {
	if instance == nil {
		panic("You must call New() first")
	}

	return instance
}

// Type get type
func (client *Client) Type() string {
	return name
}

// Logger get logger
func (client *Client) Logger() types.Logger {
	return client.logger
}

// Get get key
func (client *Client) Get(key string, value interface{}) error {
	return client.GetWithContext(context.Background(), key, value)
}

// GetAllKeysWithContext get all keys with context
func (client *Client) GetAllKeysWithContext(ctx context.Context, prefix ...string) []string {
	var ns = ""
	if len(prefix) > 0 {
		ns = prefix[0]
	}
	var keys = []string{}
	var iter = client.rdb.Scan(ctx, 0, fmt.Sprintf("%s_%s*", client.namespace, ns), 0).Iterator()
	for iter.Next(ctx) {
		var key = iter.Val()
		keys = append(keys, key)
	}

	return keys
}

// GetAllKeys get all keys
func (client *Client) GetAllKeys(prefix ...string) []string {
	return client.GetAllKeysWithContext(context.Background(), prefix...)
}

// GetAllItems get all items
func (client *Client) GetAllItems(prefix ...string) (list []types.Item) {
	return client.GetAllItemsWithContext(context.Background(), prefix...)
}

// GetAllItemsWithContext get key
func (client *Client) GetAllItemsWithContext(ctx context.Context, prefix ...string) (list []types.Item) {
	var ns = ""
	if len(prefix) > 0 {
		ns = prefix[0]
	}

	var iter = client.rdb.Scan(ctx, 0, fmt.Sprintf("%s_%s*", client.namespace, ns), 0).Iterator()
	for iter.Next(ctx) {
		var key = iter.Val()
		val, err := client.rdb.Get(ctx, key).Result()
		if err == nil {
			var item = types.Item{
				Key:   key,
				Value: val,
			}
			list = append(list, item)

		}
	}

	return
}

// GetWithContext get key
func (client *Client) GetWithContext(ctx context.Context, key string, value interface{}) error {
	var k = client.Key(key)
	val, err := client.rdb.Get(ctx, k).Result()

	if err != nil {
		if err == redis.Nil {
			return ErrKeyNotFound
		}
		return err
	}

	err = jsoniter.Unmarshal([]byte(val), value)
	if err != nil {
		client.logger.Printf("Unmarshal entity with key = %s error: %v\n", key, err)
		return err
	}
	return nil
}

// Set set key
func (client *Client) Set(key string, value interface{}, expiration time.Duration) error {
	return client.SetWithContext(context.Background(), key, value, expiration)
}

// SetWithDefault set key with default  expiration
func (client *Client) SetWithDefault(key string, value interface{}) error {

	return client.SetWithContextDefault(context.Background(), key, value)
}

// SetWithContext set key with context
func (client *Client) SetWithContext(ctx context.Context, key string, value interface{}, expiration time.Duration) error {
	cacheEntry, err := jsoniter.Marshal(value)
	if err != nil {
		client.logger.Printf("Marshal entity with key = %s error: %v\n", key, err)
		return err
	}
	var k = client.Key(key)
	err = client.rdb.Set(ctx, k, cacheEntry, expiration).Err()
	if err != nil {
		client.logger.Printf("Set value with key = %s error: %v\n", key, err)
		return err
	}
	return nil
}

// SetWithContextDefault set key with context and default expiration
func (client *Client) SetWithContextDefault(ctx context.Context, key string, value interface{}) error {
	return client.SetWithContext(ctx, key, value, client.config.DefaultExpiration)
}

// Delete delete by key
func (client *Client) Delete(keys ...string) error {
	return client.DeleteWithContext(context.Background(), keys...)
}

// DeleteWithContext delete by key with context
func (client *Client) DeleteWithContext(ctx context.Context, keys ...string) error {
	var listKey = []string{}
	for _, key := range keys {
		listKey = append(listKey, client.Key(key))
	}
	var err = client.rdb.Del(ctx, listKey...).Err()
	if err != nil {
		client.logger.Printf("Delete keys = %v error: %v\n", keys, err)
		return err
	}
	return nil
}

// Clear clear all records
func (client *Client) Clear(prefix ...string) {
	client.ClearWithContext(context.Background(), prefix...)

}

// ClearWithContext clear all records with context
func (client *Client) ClearWithContext(ctx context.Context, prefix ...string) {
	var ns = ""
	if len(prefix) > 0 {
		ns = prefix[0]
	}
	var iter = client.rdb.Scan(ctx, 0, fmt.Sprintf("%s_%s*", client.namespace, ns), 0).Iterator()
	for iter.Next(ctx) {
		var key = iter.Val()
		if strings.HasPrefix(key, client.Key(ns)) {
			var err = client.rdb.Del(ctx, key).Err()
			if err != nil {
				client.logger.Printf("Clear key = %s error: %v\n", key, err)
			}
		}
	}

}

// RedisClient get redis client
func (client *Client) RedisClient() *redis.Client {
	return client.rdb

}

// Key get full key
func (client *Client) Key(k string) string {
	if client.namespace == "" {
		return k
	}

	var key = fmt.Sprintf("%s_%s", client.namespace, k)
	return key
}
