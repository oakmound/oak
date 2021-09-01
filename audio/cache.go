package audio

import (
	"path/filepath"
	"sync"
)

// DefaultCache is the receiver for package level loading operations.
var DefaultCache = NewCache()

// Cache is a simple audio data cache
type Cache struct {
	mu   sync.RWMutex
	data map[string]Data
}

// NewCache returns an empty Cache
func NewCache() *Cache {
	return &Cache{
		data: make(map[string]Data),
	}
}

// ClearAll will remove all elements from a Cache
func (c *Cache) ClearAll() {
	c.mu.Lock()
	c.data = make(map[string]Data)
	c.mu.Unlock()
}

// Clear will remove elements matching the given key from the Cache.
func (c *Cache) Clear(key string) {
	c.mu.Lock()
	delete(c.data, key)
	c.mu.Unlock()
}

func (c *Cache) setLoaded(file string, data Data) {
	c.mu.Lock()
	c.data[file] = data
	c.data[filepath.Base(file)] = data
	c.mu.Unlock()
}

// Load calls Load on the Default Cache.
func Load(file string) (Data, error) {
	return DefaultCache.Load(file)
}

// Get calls Get on the Default Cache.
func Get(file string) (Data, error) {
	return DefaultCache.Get(file)
}
