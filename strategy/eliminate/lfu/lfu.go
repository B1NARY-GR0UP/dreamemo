package lfu

import (
	"sort"

	"github.com/B1NARY-GR0UP/dreamemo/strategy/eliminate"
)

var _ eliminate.ICore = (*Core)(nil)

type Core struct {
	eliminate.Core
	store    map[eliminate.Key]Entity
	entities Entities
}

type (
	Entity struct {
		eliminate.Entity
		frequency int
	}
	Entities []Entity
)

// NewLFUCore will new a strategy object based on LFU algorithm
// TODO: use functional option pattern?
func NewLFUCore(maxSize int, onEvicted eliminate.EvictFunc) *Core {
	return &Core{
		Core: eliminate.Core{
			MaxSize:   maxSize,
			OnEvicted: onEvicted,
		},
		store:    make(map[eliminate.Key]Entity),
		entities: make(Entities, 0),
	}
}

func (c *Core) Add(key eliminate.Key, value eliminate.Value) {
	if c.store == nil {
		c.store = make(map[eliminate.Key]Entity)
		c.entities = make(Entities, 0)
	}
	if e, ok := c.store[key]; ok {
		e.frequency++
		e.Value = value
		return
	}
	entity := Entity{
		Entity: eliminate.Entity{
			Key:   key,
			Value: value,
		},
		frequency: 1,
	}
	c.store[key] = entity
	c.entities = append(c.entities, entity)
	for c.MaxSize != 0 && c.MaxSize < c.UsedSize {
		c.RemoveLowFrequency(key)
	}
}

func (c *Core) Get(key eliminate.Key) (eliminate.Value, bool) {
	if c.store == nil {
		return nil, false
	}
	if e, ok := c.store[key]; ok {
		e.frequency++
		return e.Value, true
	}
	return nil, false
}

func (c *Core) Remove(key eliminate.Key) {
	if c.store == nil {
		return
	}
	if e, ok := c.store[key]; ok {
		for i, entity := range c.entities {
			if entity.Key == key {
				copy(c.entities[i:], c.entities[i+1:])
				c.entities = c.entities[:len(c.entities)-1]
			}
		}
		delete(c.store, key)
		c.UsedSize--
		if c.UsedSize < 0 {
			panic("UsedSize must greater than or equal to 0")
		}
		if c.OnEvicted != nil {
			c.OnEvicted(e.Key, e.Value)
		}
	}
}

func (c *Core) RemoveLowFrequency(key eliminate.Key) {
	if c.store == nil {
		return
	}
	sort.Sort(c.entities)
	removedEntity := c.entities[0]
	c.entities = c.entities[1:]
	delete(c.store, key)
	c.UsedSize--
	if c.UsedSize < 0 {
		panic("UsedSize must greater than or equal to 0")
	}
	if c.OnEvicted != nil {
		c.OnEvicted(removedEntity.Key, removedEntity.Value)
	}
}

func (c *Core) Clear() {
	if c.OnEvicted != nil {
		for _, e := range c.store {
			c.OnEvicted(e.Key, e.Value)
		}
	}
	c.entities = nil
	c.store = nil
	c.UsedSize = 0
}

func (c *Core) Name() string {
	return "lfu"
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
