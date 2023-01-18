package cache

import "sort"

type LFUCache struct {
	Cache
	cache    map[Key]LFUEntity
	entities Entities
}

type (
	LFUEntity struct {
		Entity
		frequency int
	}
	Entities []LFUEntity
)

// NewLFU will new a cache object based on LFU algorithm
// TODO: use functional option pattern?
func NewLFU(maxSize int, onEvicted EvictFunc) *LFUCache {
	return &LFUCache{
		Cache: Cache{
			MaxSize:   maxSize,
			OnEvicted: onEvicted,
		},
		cache:    make(map[Key]LFUEntity),
		entities: make(Entities, 0),
	}
}

func (c *LFUCache) Add(key Key, value Value) {
	if c.cache == nil {
		c.cache = make(map[Key]LFUEntity)
		c.entities = make(Entities, 0)
	}
	if e, ok := c.cache[key]; ok {
		e.frequency++
		e.Value = value
		return
	}
	entity := LFUEntity{
		Entity: Entity{
			Key:   key,
			Value: value,
		},
		frequency: 1,
	}
	c.cache[key] = entity
	c.entities = append(c.entities, entity)
	for c.MaxSize != 0 && c.MaxSize < c.UsedSize {
		c.RemoveLowFrequency(key)
	}
}

func (c *LFUCache) Get(key Key) (Value, bool) {
	if c.cache == nil {
		return nil, false
	}
	if e, ok := c.cache[key]; ok {
		e.frequency++
		return e.Value, true
	}
	return nil, false
}

func (c *LFUCache) Remove(key Key) {
	if c.cache == nil {
		return
	}
	if e, ok := c.cache[key]; ok {
		for i, entity := range c.entities {
			if entity.Key == key {
				copy(c.entities[i:], c.entities[i+1:])
				c.entities = c.entities[:len(c.entities)-1]
			}
		}
		delete(c.cache, key)
		c.UsedSize--
		if c.UsedSize < 0 {
			panic("UsedSize must greater than or equal to 0")
		}
		if c.OnEvicted != nil {
			c.OnEvicted(e.Key, e.Value)
		}
	}
}

func (c *LFUCache) RemoveLowFrequency(key Key) {
	if c.cache == nil {
		return
	}
	sort.Sort(c.entities)
	removedEntity := c.entities[0]
	c.entities = c.entities[1:]
	delete(c.cache, key)
	c.UsedSize--
	if c.UsedSize < 0 {
		panic("UsedSize must greater than or equal to 0")
	}
	if c.OnEvicted != nil {
		c.OnEvicted(removedEntity.Key, removedEntity.Value)
	}
}

func (c *LFUCache) Clear() {
	if c.OnEvicted != nil {
		for _, e := range c.cache {
			c.OnEvicted(e.Key, e.Value)
		}
	}
	c.entities = nil
	c.cache = nil
	c.UsedSize = 0
}

func (e Entities) Len() int {
	return len(e)
}

func (e Entities) Less(i, j int) bool {
	return e[i].frequency < e[j].frequency
}

func (e Entities) Swap(i, j int) {
	e[i], e[j] = e[j], e[i]
}
