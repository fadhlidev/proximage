package cache

import (
	"container/list"
	"sync"
	"time"
)

type entry struct {
	key       string
	data      []byte
	expiresAt time.Time
}

type Cache struct {
	mu    sync.Mutex
	cap   int
	ttl   time.Duration
	items map[string]*list.Element
	lru   *list.List
}

func New(capacity int, ttl time.Duration) *Cache {
	return &Cache{
		cap:   capacity,
		ttl:   ttl,
		items: make(map[string]*list.Element),
		lru:   list.New(),
	}
}

func (c *Cache) Get(key string) ([]byte, bool) {
	c.mu.Lock()
	defer c.mu.Unlock()

	el, ok := c.items[key]
	if !ok {
		return nil, false
	}

	e := el.Value.(*entry)
	if time.Now().After(e.expiresAt) {
		c.lru.Remove(el)
		delete(c.items, key)
		return nil, false
	}

	c.lru.MoveToFront(el)
	return e.data, true
}

func (c *Cache) Set(key string, data []byte) {
	c.mu.Lock()
	defer c.mu.Unlock()

	if el, ok := c.items[key]; ok {
		c.lru.MoveToFront(el)
		el.Value.(*entry).data = data
		el.Value.(*entry).expiresAt = time.Now().Add(c.ttl)
		return
	}

	if c.lru.Len() >= c.cap {
		// evict oldest
		back := c.lru.Back()
		if back != nil {
			c.lru.Remove(back)
			delete(c.items, back.Value.(*entry).key)
		}
	}

	e := &entry{key: key, data: data, expiresAt: time.Now().Add(c.ttl)}
	el := c.lru.PushFront(e)
	c.items[key] = el
}
