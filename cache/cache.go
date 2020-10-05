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
	Get(key string, value interface{}) error
	Set(key string, value interface{}, expiration time.Duration) error
	SetWithDefault(key string, value interface{}) error
	Delete(keys ...string) error
	Clear()
	GetAllItemsWithContext(ctx context.Context) []types.Item
	GetWithContext(ctx context.Context, key string, src interface{}) error
	SetWithContext(ctx context.Context, key string, value interface{}, expiration time.Duration) error
	SetWithContextDefault(ctx context.Context, key string, value interface{}) error
	DeleteWithContext(ctx context.Context, keys ...string) error
	ClearWithContext(ctx context.Context)
	Logger() types.Logger
}
