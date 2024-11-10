package render

import (
	"image"
	"sync"

	"github.com/golang/freetype/truetype"
	"github.com/oakmound/oak/v4/alg/intgeom"
)

// DefaultCache is the receiver for package level sprites, sheets, and font loading operations.
var DefaultCache = NewCache()

// Cache is a simple image data cache
type Cache struct {
	imageLock    sync.RWMutex
	loadedImages map[string]*image.RGBA

	sheetLock    sync.RWMutex
	loadedSheets map[string]*Sheet

	fontLock    sync.RWMutex
	loadedFonts map[string]*truetype.Font
}

// NewCache returns an empty Cache
func NewCache() *Cache {
	return &Cache{
		loadedImages: make(map[string]*image.RGBA),
		loadedSheets: make(map[string]*Sheet),
		loadedFonts:  make(map[string]*truetype.Font),
	}
}

// ClearAll will remove all elements from a Cache
func (c *Cache) ClearAll() {
	c.imageLock.Lock()
	c.sheetLock.Lock()
	c.fontLock.Lock()
	c.loadedImages = make(map[string]*image.RGBA)
	c.loadedSheets = make(map[string]*Sheet)
	c.loadedFonts = make(map[string]*truetype.Font)
	c.fontLock.Unlock()
	c.sheetLock.Unlock()
	c.imageLock.Unlock()
}

// Clear will remove elements matching the given key from the Cache.
func (c *Cache) Clear(key string) {
	c.imageLock.Lock()
	c.sheetLock.Lock()
	c.fontLock.Lock()
	delete(c.loadedImages, key)
	delete(c.loadedSheets, key)
	delete(c.loadedFonts, key)
	c.fontLock.Unlock()
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

// LoadComplexSheet calls LoadComplexSheet on the Default Cache.
func LoadComplexSheet(file string, opts ...Option) (*Sheet, error) {
	return DefaultCache.LoadComplexSheet(file, opts...)
}

// GetFont calls GetFont on the Default Cache.
func GetFont(file string) (*truetype.Font, error) {
	return DefaultCache.GetFont(file)
}

// LoadFont calls LoadFont on the Default Cache.
func LoadFont(file string) (*truetype.Font, error) {
	return DefaultCache.LoadFont(file)
}
