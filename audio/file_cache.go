package audio

import (
	"path/filepath"
	"sync"

	"github.com/oakmound/oak/v3/audio/pcm"
)

// DefaultCache is the receiver for package level loading operations.
var DefaultCache = NewCache()

// Cache is a simple audio data cache
type Cache struct {
	mu   sync.RWMutex
	data map[string]*BytesReader
}

// NewCache returns an empty Cache
func NewCache() *Cache {
	return &Cache{
		data: make(map[string]*BytesReader),
	}
}

// ClearAll will remove all elements from a Cache
func (c *Cache) ClearAll() {
	c.mu.Lock()
	c.data = make(map[string]*BytesReader)
	c.mu.Unlock()
}

// Clear will remove elements matching the given key from the Cache.
func (c *Cache) Clear(key string) {
	c.mu.Lock()
	delete(c.data, key)
	c.mu.Unlock()
}

func (c *Cache) setLoaded(file string, r pcm.Reader) {
	// This ReadAll and .Copy() on Cache.Read ensure that multiple loads from the cache do not
	// change the data that will be read on future reads.
	br := ReadAll(r)
	c.mu.Lock()
	c.data[file] = br
	c.data[filepath.Base(file)] = br
	c.mu.Unlock()
}

// Load calls Load on the Default Cache.
func Load(file string) (pcm.Reader, error) {
	return DefaultCache.Load(file)
}

// Get calls Get on the Default Cache.
func Get(file string) (pcm.Reader, error) {
	return DefaultCache.Get(file)
}
