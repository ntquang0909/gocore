package cache

import (
	"context"
	"fmt"
	"testing"
	"time"

	goredis "github.com/go-redis/redis/v8"
	"github.com/patrickmn/go-cache"
	"github.com/stretchr/testify/assert"
	"github.com/thaitanloi365/gocore/cache/memory"
	"github.com/thaitanloi365/gocore/cache/redis"
)

func TestRedisCache(t *testing.T) {
	var redisCache Cache = redis.New(&redis.Config{
		Namespace: "redis_test",
		Options: &goredis.Options{
			Addr: "localhost:6379",
			DB:   0,
		},
	})

	type Author struct {
		Name string
	}

	var author = Author{
		Name: "Loi",
	}

	var err = redisCache.Set("test", &author, time.Hour)
	assert.NoError(t, err)

	var result Author
	err = redisCache.Get("test", &result)
	assert.NoError(t, err)
	assert.Equal(t, author, result)

	err = redisCache.Delete("test")
	assert.NoError(t, err)

	var result2 Author
	err = redisCache.Set("test", &author, time.Hour)
	assert.NoError(t, err)

	err = redisCache.Get("test", &result2)
	assert.NoError(t, err)

	for i := 0; i < 5; i++ {
		redisCache.Set(fmt.Sprintf("test1_%d", i), &author, cache.NoExpiration)
		redisCache.SetWithDefault(fmt.Sprintf("test2_%d", i), &author)
		redisCache.SetWithContextDefault(context.Background(), fmt.Sprintf("test3_%d", i), &author)
	}

	var list = redisCache.GetAllItemsWithContext(context.Background())
	for _, item := range list {
		fmt.Println("list", item)

	}

	var list2 = redisCache.GetAllItemsWithContext(context.Background(), "test2")
	for _, item := range list2 {
		fmt.Println("test2", item)

	}

	redisCache.Clear("test2")

	var list22 = redisCache.GetAllItemsWithContext(context.Background(), "test2")
	for _, item := range list22 {
		fmt.Println("test22", item)

	}

	var list3 = redisCache.GetAllItemsWithContext(context.Background(), "test3")
	for _, item := range list3 {
		fmt.Println("test3", item)

	}
	for i := 0; i < 5; i++ {
		var key = fmt.Sprintf("test_%d", i)
		var result Author
		err = redisCache.Get(key, &result)
		if err != nil {
			redisCache.Logger().Printf("Error %v\n", err)
			continue
		}
		redisCache.Logger().Printf("Key %s value = %v \n", key, result)
	}

	for i := 0; i < 5; i++ {
		var key = fmt.Sprintf("test_%d", i)
		var result Author
		err = redisCache.Get(key, &result)
		if err != nil {
			redisCache.Logger().Printf("Error %v\n", err)
			continue
		}
		redisCache.Logger().Printf("Key %s value = %v \n", key, result)
	}

	redisCache.Logger().Printf("%v\n", result2)
}

func TestMemCache(t *testing.T) {
	var memCache Cache = memory.New(&memory.Config{
		Namespace: "redis_test",
	})

	type Author struct {
		Name string
	}

	var author = Author{
		Name: "Loi",
	}

	var err = memCache.Set("test", &author, time.Hour)
	assert.NoError(t, err)

	var result Author
	err = memCache.Get("test", &result)
	assert.NoError(t, err)
	assert.Equal(t, author, result)

	err = memCache.Delete("test")
	assert.NoError(t, err)

	var result2 Author
	err = memCache.Set("test", &author, time.Hour)
	assert.NoError(t, err)

	err = memCache.Get("test", &result2)
	assert.NoError(t, err)

	for i := 0; i < 5; i++ {
		memCache.Set(fmt.Sprintf("test_%d", i), &author, time.Hour)
	}

	var list = memCache.GetAllItemsWithContext(context.Background())
	for _, item := range list {
		fmt.Println("list", item)

	}

	for i := 0; i < 5; i++ {
		var key = fmt.Sprintf("test_%d", i)
		var result Author
		err = memCache.Get(key, &result)
		if err != nil {
			memCache.Logger().Printf("Error %v\n", err)
			continue
		}
		memCache.Logger().Printf("Key %s value = %v \n", key, result)
	}

	memCache.Clear()
	for i := 0; i < 5; i++ {
		var key = fmt.Sprintf("test_%d", i)
		var result Author
		err = memCache.Get(key, &result)
		if err != nil {
			memCache.Logger().Printf("Error %v\n", err)
			continue
		}
		memCache.Logger().Printf("Key %s value = %v \n", key, result)
	}

	memCache.Logger().Printf("%v\n", result2)
}
