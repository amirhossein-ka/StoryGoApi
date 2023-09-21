package redis

import (
	"StoryGoAPI/config"
	"context"
	"fmt"
	"github.com/go-redis/cache/v9"
	"time"

	"github.com/redis/go-redis/v9"
)

type CacheImpl struct {
	rdb   *redis.Client
	cache *cache.Cache
	cfg   config.Redis
}

func NewCache(cfg config.Redis) (*CacheImpl, error) {
	rdb := redis.NewClient(&redis.Options{
		Addr:     cfg.Addr,
		Password: cfg.Password,
		Username: cfg.Username,
		DB:       cfg.DB,
		OnConnect: func(ctx context.Context, cn *redis.Conn) error {
			if err := cn.Ping(ctx).Err(); err != nil {
				return err
			}
			return nil
		},
	})

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()
	if err := rdb.Ping(ctx).Err(); err != nil {
		return nil, fmt.Errorf("failed to connect to redis: %w", err)
	}
	cache_ := cache.New(&cache.Options{
		Redis: rdb,
		//LocalCache: cache.NewTinyLFU(2048, cfg.ExpTime),
	})

	return &CacheImpl{rdb: rdb, cfg: cfg, cache: cache_}, nil
}

func (c *CacheImpl) Close() error {
	return c.rdb.Close()
}

func (c *CacheImpl) addToList(ctx context.Context, key string, keys ...string) error {
	if err := c.rdb.LPush(ctx, key, keys).Err(); err != nil {
		return err
	}
	return c.rdb.ExpireAt(ctx, key, time.Now().Add(c.cfg.ExpTime)).Err()
}

// getValues retrieve all values from redis list and remove the list.
func (c *CacheImpl) getValues(ctx context.Context, key string) ([]string, error) {
	result, err := c.rdb.LRange(ctx, key, 0, -1).Result()
	if err != nil {
		return nil, err
	}
	if err = c.rdb.Del(ctx, key).Err(); err != nil {
		return nil, err
	}

	return result, nil
}

// Set adds value to redis cache.
// example usage:
//
//	// note the pointer obj
//	CacheImpl.Set(ctx, "my_key", &obj, func(*cache.Item) (any, error) {
//	    return &Object{...}, nil // somehow get the object (from a db) and return it
//	})
//
// and obj is the returned object from the `do` function set in params
func (c *CacheImpl) Set(ctx context.Context, key, listKey string, value any, do func(item *cache.Item) (any, error)) error {
	if value == nil {
		return fmt.Errorf("value must not be nil")
	}
	var item = &cache.Item{
		Ctx:   ctx,
		Key:   key,
		Value: value,
		TTL:   c.cfg.ExpTime,
	}
	if do != nil {
		item.Do = do
		if err := c.cache.Once(item); err != nil {
			return err
		}
	}

	if err := c.cache.Set(item); err != nil {
		return err
	}

	if err := c.addToList(ctx, listKey, key); err != nil {
		return err
	}
	return nil
}

// Get puts the value returned from cache to value.
// pass `value` as pointer.
func (c *CacheImpl) Get(ctx context.Context, key string, value any) error {
	if value == nil {
		return fmt.Errorf("value must not be nil")
	}
	return c.cache.Get(ctx, key, value)
}

func (c *CacheImpl) InvalidateCache(ctx context.Context, listKey string) error {
	cacheKeys, err := c.getValues(ctx, listKey)
	if err != nil {
		return err
	}
	for i := 0; i < len(cacheKeys); i++ {
		err = c.cache.Delete(ctx, cacheKeys[i])
		if err != nil {
			return err
		}
	}
	return nil
}

func (c *CacheImpl) BlackList(ctx context.Context, jwtToken string) error {
	return c.rdb.SetNX(ctx, jwtToken, 1, c.cfg.BlacklistExpTime).Err()
}

func (c *CacheImpl) Exists(ctx context.Context, jwtToken string) bool {
	ok, err := c.rdb.Get(ctx, jwtToken).Bool()
	if err != nil {
		return false
	}
	return ok
}

func (c *CacheImpl) Delete(ctx context.Context, key string) error {
	return c.cache.Delete(ctx, key)
}
