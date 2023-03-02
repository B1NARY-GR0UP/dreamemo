// Copyright 2023 BINARY Members
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
// http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.
//

package lfu

import (
	"sort"

	"github.com/B1NARY-GR0UP/dreamemo/strategy/eliminate"
	"github.com/B1NARY-GR0UP/inquisitor/core"
)

var _ eliminate.ICore = (*Core)(nil)

// Core is not safe under concurrent scene
type Core struct {
	eliminate.Core
	store    map[eliminate.Key]entity
	entities entities
}

type (
	entity struct {
		eliminate.Entity
		frequency int
	}
	entities []entity
)

// NewLFUCore will new a strategy object based on LFU algorithm
func NewLFUCore(opts ...Option) *Core {
	options := newOptions(opts...)
	return &Core{
		Core: eliminate.Core{
			MaxSize:   options.MaxSize,
			OnEvicted: options.OnEvicted,
		},
		store:    make(map[eliminate.Key]entity),
		entities: make(entities, 0),
	}
}

func (c *Core) Add(key eliminate.Key, value eliminate.Value) {
	if c.store == nil {
		c.store = make(map[eliminate.Key]entity)
		c.entities = make(entities, 0)
	}
	if e, ok := c.store[key]; ok {
		e.frequency++
		e.Value = value
		return
	}
	ent := entity{
		Entity: eliminate.Entity{
			Key:   key,
			Value: value,
		},
		frequency: 1,
	}
	c.store[key] = ent
	c.entities = append(c.entities, ent)
	for c.MaxSize != 0 && c.MaxSize < c.UsedSize {
		c.RemoveLowFrequency()
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

func (c *Core) RemoveLowFrequency() {
	if c.store == nil {
		return
	}
	sort.Sort(c.entities)
	removedEntity := c.entities[0]
	c.entities = c.entities[1:]
	delete(c.store, removedEntity.Key)
	c.UsedSize--
	if c.UsedSize < 0 {
		core.Error("---DREAMEMO--- UsedSize must greater than or equal to 0")
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

func (e entities) Len() int {
	return len(e)
}

func (e entities) Less(i, j int) bool {
	return e[i].frequency < e[j].frequency
}

func (e entities) Swap(i, j int) {
	e[i], e[j] = e[j], e[i]
}
