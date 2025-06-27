package pokecache

import (
  "time"
  "sync"
)

type cacheEntry struct {
  createdAt time.Time 
  val []byte
}

type Cache struct {
  mu sync.RWMutex
  entries map[string]cacheEntry
  timeout time.Duration
  done chan struct{}
}

func NewCache(timeout time.Duration) *Cache {
  return &Cache {
    entries: make(map[string]cacheEntry),
    timeout: timeout,
    done: make(chan struct{}),
  } 
} 

func (c *Cache) Add(key string, value []byte) {
  c.mu.Lock()
  defer c.mu.Unlock()

  c.entries[key] = cacheEntry{
    createdAt: time.Now(),
    val: value,
  }
}

func (c *Cache) Get(key string) ([]byte, bool) {
  c.mu.Lock()
  defer c.mu.Unlock()

  entry, exists := c.entries[key]
  if !exists {
    return nil, false
  }

  return entry.val, true
} 

func (c *Cache) reapLoop(interval time.Duration) {
  ticker := time.NewTicker(interval)
  defer ticker.Stop()

  for {
    select {
    case <- ticker.C:
      c.mu.Lock()
      for key, entry := range c.entries {
        if time.Since(entry.createdAt) > interval {
          delete(c.entries, key)
        }
      }
      c.mu.Unlock()
    case <- c.done:
      return
    }
  }
}
