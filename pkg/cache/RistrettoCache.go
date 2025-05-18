package cache

import (
	"time"

	"github.com/dgraph-io/ristretto/v2"
)

var _ CacheMethod = (*Cache)(nil) // 确保 Cache 实现了 CacheMethod 接口的所有方法，否则会报错，提示没有实现所有方法，无法编译通过。

type Cache struct {
	cache         *ristretto.Cache[string, any]
	cacheInstance CacheInstance
	// ttl   time.Duration
}

func NewCache(ttl time.Duration, maxCost int64) (*Cache, error) {
	cache, err := ristretto.NewCache(&ristretto.Config[string, any]{
		NumCounters: 10 * maxCost, // 推荐: 10x maxCost
		MaxCost:     maxCost,      // 总体可用空间大小 (cost 的单位是任意的)
		BufferItems: 64,           // 推荐值
	})
	if err != nil {
		return nil, err
	}

	return &Cache{
		cache: cache,
		cacheInstance: CacheInstance{
			TTL: ttl,
		},
	}, nil
}

func (r *Cache) Get(key string) (any, bool) {
	return r.cache.Get(key)
}

func (r *Cache) Set(key string, val any) {
	// Cost = 1 是最简单估算方式；可按需自定义
	r.cache.SetWithTTL(key, val, 1, r.cacheInstance.TTL)
}
