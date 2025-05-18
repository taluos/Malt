package cache

import "time"

type CacheInstance struct {
	TTL time.Duration `json:"TTL" mapstructure:"TTL"` // minisecond
}

type CacheMethod interface {

	// Get returns the value for the given key.
	Get(key string) (any, bool)

	// Set sets the value for the given key.
	Set(key string, val any)
}
