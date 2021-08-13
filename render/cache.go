package render

import (
	"image"
	"sync"

	"github.com/oakmound/oak/v3/alg/intgeom"
)

// The DefaultCache is the receiver for package level sprite and sheet loading operations.
var DefaultCache = NewCache()

// Cache is a simple image data cache
type Cache struct {
	imageLock    sync.RWMutex
	loadedImages map[string]*image.RGBA

	sheetLock    sync.RWMutex
	loadedSheets map[string]*Sheet
}

// NewCache returns an empty Cache
func NewCache() *Cache {
	return &Cache{
		loadedImages: make(map[string]*image.RGBA),
		loadedSheets: make(map[string]*Sheet),
	}
}

// ClearAll will remove all elements from a Cache
func (c *Cache) ClearAll() {
	c.imageLock.Lock()
	c.sheetLock.Lock()
	c.loadedImages = make(map[string]*image.RGBA)
	c.loadedSheets = make(map[string]*Sheet)
	c.sheetLock.Unlock()
	c.imageLock.Unlock()
}

// Clear will remove elements matching the given key from the Cache.
func (c *Cache) Clear(key string) {
	c.imageLock.Lock()
	c.sheetLock.Lock()
	delete(c.loadedImages, key)
	delete(c.loadedSheets, key)
	c.sheetLock.Unlock()
	c.imageLock.Unlock()
}

// GetSprite calls GetSprite on the Default Cache.
func GetSprite(file string) (*Sprite, error) {
	return DefaultCache.GetSprite(file)
}

// LoadSprite calls LoadSprite on the Default Cache.
func LoadSprite(file string) (*Sprite, error) {
	return DefaultCache.LoadSprite(file)
}

// GetSheet calls GetSheet on the Default Cache.
func GetSheet(file string) (*Sheet, error) {
	return DefaultCache.GetSheet(file)
}

// LoadSheet calls LoadSheet on the Default Cache.
func LoadSheet(file string, cellSize intgeom.Point2) (*Sheet, error) {
	return DefaultCache.LoadSheet(file, cellSize)
}
