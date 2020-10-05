package memory

import (
	"context"
	"errors"
	"fmt"
	"log"
	"os"
	"reflect"
	"strings"
	"time"

	"github.com/patrickmn/go-cache"
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
	DefaultExpiration time.Duration
	CleanupInterval   time.Duration
	Logger            types.Logger
	Namespace         string
}

// Client client
type Client struct {
	cache     *cache.Cache
	config    *Config
	logger    types.Logger
	namespace string
}

// New init cache
func New(config *Config) *Client {
	return NewFrom(config, make(map[string]cache.Item))
}

// NewFrom init cache with default value
func NewFrom(config *Config, items map[string]cache.Item) *Client {
	var c = cache.NewFrom(config.DefaultExpiration, config.CleanupInterval, items)

	instance = &Client{
		cache:     c,
		config:    config,
		logger:    log.New(os.Stdout, "\r\n", 0),
		namespace: "gocore_memory_cache",
	}

	if config.Namespace != "" {
		instance.namespace = config.Namespace
	}

	if config.Logger != nil {
		instance.logger = config.Logger
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

// GetAllItemsWithContext get key
func (client *Client) GetAllItemsWithContext(ctx context.Context) (list []types.Item) {
	for key, value := range client.cache.Items() {
		var item = types.Item{
			Key:   key,
			Value: value.Object,
		}
		list = append(list, item)
	}
	return
}

// Get get key
func (client *Client) Get(key string, value interface{}) error {
	return client.GetWithContext(context.Background(), key, value)
}

// GetWithContext get key
func (client *Client) GetWithContext(ctx context.Context, key string, value interface{}) error {
	var k = client.getKey(key)
	v, found := client.cache.Get(k)
	if !found {
		client.logger.Printf("Key = %s is not found\n", k)
		return ErrKeyNotFound
	}

	var i = reflect.ValueOf(v)
	var o = reflect.ValueOf(value)

	o.Elem().Set(i.Elem())
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
	var k = client.getKey(key)
	client.cache.Set(k, value, expiration)
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
	for _, key := range keys {
		client.cache.Delete(client.getKey(key))
	}
	return nil
}

// Clear clear all records
func (client *Client) Clear() {
	client.ClearWithContext(context.Background())

}

// ClearWithContext clear all records with context
func (client *Client) ClearWithContext(ctx context.Context) {
	for key := range client.cache.Items() {
		if strings.HasPrefix(key, client.namespace) {
			client.cache.Delete(key)
		}
	}
}

// Client get redis client
func (client *Client) Client() *cache.Cache {
	return client.cache

}

func (client *Client) getKey(k string) string {
	if client.namespace == "" {
		return k
	}

	var key = fmt.Sprintf("%s_%s", client.namespace, k)
	return key
}
