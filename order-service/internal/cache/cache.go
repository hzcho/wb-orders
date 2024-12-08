package cache

import (
	"container/list"
	"sync"
	"time"

	"wb-orders/internal/domain/model"
)

type ICache interface {
	Get(key string) (model.Order, bool)
	SetMany(items map[string]model.Order)
	Set(key string, value model.Order)
}

type Cache struct {
	mu        sync.Mutex
	data      map[string]*list.Element
	evictList *list.List
	capacity  int
}

type entry struct {
	key       string
	value     model.Order
	timestamp time.Time
}

func NewCache(capacity int) *Cache {
	if capacity <= 0 {
		panic("capacity must be greater than zero")
	}
	return &Cache{
		data:      make(map[string]*list.Element, capacity),
		evictList: list.New(),
		capacity:  capacity,
	}
}

func (c *Cache) SetMany(items map[string]model.Order) {
	c.mu.Lock()
	defer c.mu.Unlock()
	for key, value := range items {
		if elem, exists := c.data[key]; exists {
			c.evictList.MoveToFront(elem)
			elem.Value.(*entry).value = value
			elem.Value.(*entry).timestamp = time.Now()
			continue
		}
		if c.evictList.Len() >= c.capacity {
			c.evict()
		}
		e := &entry{key: key, value: value, timestamp: time.Now()}
		elem := c.evictList.PushFront(e)
		c.data[key] = elem
	}
}

func (c *Cache) Set(key string, value model.Order) {
	c.mu.Lock()
	defer c.mu.Unlock()
	if elem, exists := c.data[key]; exists {
		c.evictList.MoveToFront(elem)
		elem.Value.(*entry).value = value
		elem.Value.(*entry).timestamp = time.Now()
		return
	}
	if c.evictList.Len() >= c.capacity {
		c.evict()
	}
	e := &entry{key: key, value: value, timestamp: time.Now()}
	elem := c.evictList.PushFront(e)
	c.data[key] = elem
}

func (c *Cache) Get(key string) (model.Order, bool) {
	c.mu.Lock()
	defer c.mu.Unlock()
	if elem, exists := c.data[key]; exists {
		c.evictList.MoveToFront(elem)
		elem.Value.(*entry).timestamp = time.Now()
		return elem.Value.(*entry).value, true
	}
	return model.Order{}, false
}

func (c *Cache) evict() {
	if elem := c.evictList.Back(); elem != nil {
		c.evictList.Remove(elem)
		delete(c.data, elem.Value.(*entry).key)
	}
}
