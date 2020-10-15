package cache

import (
	"context"
	"time"

	"github.com/thaitanloi365/gocore/cache/types"
)

// NoExpiration no expiration
const NoExpiration time.Duration = -1

// Cache cache
type Cache interface {
	Type() string
	Key(k string) string
	Get(key string, value interface{}) error
	Set(key string, value interface{}, expiration time.Duration) error
	SetWithDefault(key string, value interface{}) error
	Delete(keys ...string) error
	Clear(prefix ...string)
	ClearWithContext(ctx context.Context, prefix ...string)
	GetAllKeys(prefix ...string) []string
	GetAllKeysWithContext(ctx context.Context, prefix ...string) []string
	GetAllItems(prefix ...string) []types.Item
	GetAllItemsWithContext(ctx context.Context, prefix ...string) []types.Item
	GetWithContext(ctx context.Context, key string, src interface{}) error
	SetWithContext(ctx context.Context, key string, value interface{}, expiration time.Duration) error
	SetWithContextDefault(ctx context.Context, key string, value interface{}) error
	DeleteWithContext(ctx context.Context, keys ...string) error
	Logger() types.Logger
}
