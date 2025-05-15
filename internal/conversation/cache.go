package conversation

import (
	"sync"
	"time"
)

// CacheEntry represents a cached conversation with expiration
type CacheEntry struct {
	Messages    []Message
	LastUpdated time.Time
	ExpiresAt   time.Time
}

// Cache implements a thread-safe cache with TTL
type Cache struct {
	mu      sync.RWMutex
	entries map[string]*CacheEntry
	ttl     time.Duration
}

var (
	cache     *Cache
	cacheOnce sync.Once
)

// InitCache initializes the conversation cache
func InitCache(ttl time.Duration) {
	cacheOnce.Do(func() {
		cache = &Cache{
			entries: make(map[string]*CacheEntry),
			ttl:     ttl,
		}
		go cache.cleanupLoop()
	})
}

// Get retrieves a conversation from cache
func (c *Cache) Get(userID string) ([]Message, bool) {
	c.mu.RLock()
	defer c.mu.RUnlock()

	if entry, exists := c.entries[userID]; exists {
		if time.Now().Before(entry.ExpiresAt) {
			return entry.Messages, true
		}
		// Entry expired, remove it
		delete(c.entries, userID)
	}
	return nil, false
}

// Set stores a conversation in cache
func (c *Cache) Set(userID string, messages []Message) {
	c.mu.Lock()
	defer c.mu.Unlock()

	now := time.Now()
	c.entries[userID] = &CacheEntry{
		Messages:    messages,
		LastUpdated: now,
		ExpiresAt:   now.Add(c.ttl),
	}
}

// cleanupLoop periodically removes expired entries
func (c *Cache) cleanupLoop() {
	ticker := time.NewTicker(5 * time.Minute)
	defer ticker.Stop()

	for range ticker.C {
		c.mu.Lock()
		now := time.Now()
		for userID, entry := range c.entries {
			if now.After(entry.ExpiresAt) {
				delete(c.entries, userID)
			}
		}
		c.mu.Unlock()
	}
}
