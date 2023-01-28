package lru

import (
	"container/list"

	"github.com/B1NARY-GR0UP/dreamemo/strategy/eliminate"
)

var _ eliminate.ICore = (*Core)(nil)

// Core is not safe under concurrent scene
type Core struct {
	eliminate.Core
	store map[eliminate.Key]*list.Element
	list  *list.List
}

// NewLRUCore will new a strategy object based on LRU algorithm
// TODO: use functional option pattern?
func NewLRUCore(maxSize int, onEvicted eliminate.EvictFunc) *Core {
	return &Core{
		Core: eliminate.Core{
			MaxSize:   maxSize,
			OnEvicted: onEvicted,
		},
		store: make(map[eliminate.Key]*list.Element),
		list:  list.New(),
	}
}

func (c *Core) Add(key eliminate.Key, value eliminate.Value) {
	if c.store == nil {
		c.store = make(map[eliminate.Key]*list.Element)
		c.list = list.New()
	}
	if ele, ok := c.store[key]; ok {
		c.list.MoveToFront(ele)
		ele.Value.(*eliminate.Entity).Value = value
		return
	}
	ele := c.list.PushFront(&eliminate.Entity{
		Key:   key,
		Value: value,
	})
	c.store[key] = ele
	c.UsedSize++
	for c.MaxSize != 0 && c.MaxSize < c.UsedSize {
		c.RemoveOldest()
	}
}

func (c *Core) Get(key eliminate.Key) (eliminate.Value, bool) {
	if c.store == nil {
		return nil, false
	}
	if ele, ok := c.store[key]; ok {
		c.list.MoveToFront(ele)
		return ele.Value.(*eliminate.Entity).Value, true
	}
	return nil, false
}

func (c *Core) Remove(key eliminate.Key) {
	if c.store == nil {
		return
	}
	if ele, ok := c.store[key]; ok {
		c.removeElement(ele)
	}
}

func (c *Core) Clear() {
	if c.OnEvicted != nil {
		for _, ele := range c.store {
			entity := ele.Value.(*eliminate.Entity)
			c.OnEvicted(entity.Key, entity.Value)
		}
	}
	c.list = nil
	c.store = nil
	c.UsedSize = 0
}

func (c *Core) Name() string {
	return "lru"
}

func (c *Core) RemoveOldest() {
	if c.store == nil {
		return
	}
	ele := c.list.Back()
	if ele != nil {
		c.removeElement(ele)
	}
}

func (c *Core) removeElement(ele *list.Element) {
	c.list.Remove(ele)
	entry := ele.Value.(*eliminate.Entity)
	delete(c.store, entry.Key)
	c.UsedSize--
	if c.UsedSize < 0 {
		panic("UsedSize must greater than or equal to 0")
	}
	if c.OnEvicted != nil {
		c.OnEvicted(entry.Key, entry.Value)
	}
}
