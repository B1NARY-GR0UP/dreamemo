package cache

import "container/list"

// LRUCache is not safe under concurrent scene
type LRUCache struct {
	Cache
	cache map[Key]*list.Element
	list  *list.List
}

// NewLRU will new a cache object based on LRU algorithm
// TODO: use functional option pattern?
func NewLRU(maxSize int, onEvicted EvictFunc) *LRUCache {
	return &LRUCache{
		Cache: Cache{
			MaxSize:   maxSize,
			OnEvicted: onEvicted,
		},
		cache: make(map[Key]*list.Element),
		list:  list.New(),
	}
}

func (c *LRUCache) Add(key Key, value Value) {
	if c.cache == nil {
		c.cache = make(map[Key]*list.Element)
		c.list = list.New()
	}
	if ele, ok := c.cache[key]; ok {
		c.list.MoveToFront(ele)
		ele.Value.(*Entity).Value = value
		return
	}
	ele := c.list.PushFront(&Entity{
		Key:   key,
		Value: value,
	})
	c.cache[key] = ele
	c.UsedSize++
	for c.MaxSize != 0 && c.MaxSize < c.UsedSize {
		c.RemoveOldest()
	}
}

func (c *LRUCache) Get(key Key) (Value, bool) {
	if c.cache == nil {
		return nil, false
	}
	if ele, ok := c.cache[key]; ok {
		c.list.MoveToFront(ele)
		return ele.Value.(*Entity).Value, true
	}
	return nil, false
}

func (c *LRUCache) Remove(key Key) {
	if c.cache == nil {
		return
	}
	if ele, ok := c.cache[key]; ok {
		c.removeElement(ele)
	}
}

func (c *LRUCache) Clear() {
	if c.OnEvicted != nil {
		for _, ele := range c.cache {
			entity := ele.Value.(*Entity)
			c.OnEvicted(entity.Key, entity.Value)
		}
	}
	c.list = nil
	c.cache = nil
	c.UsedSize = 0
}

func (c *LRUCache) RemoveOldest() {
	if c.cache == nil {
		return
	}
	ele := c.list.Back()
	if ele != nil {
		c.removeElement(ele)
	}
}

func (c *LRUCache) removeElement(ele *list.Element) {
	c.list.Remove(ele)
	entry := ele.Value.(*Entity)
	delete(c.cache, entry.Key)
	c.UsedSize--
	if c.UsedSize < 0 {
		panic("UsedSize must greater than or equal to 0")
	}
	if c.OnEvicted != nil {
		c.OnEvicted(entry.Key, entry.Value)
	}
}
