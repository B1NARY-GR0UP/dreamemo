package api

import (
	"context"
	"fmt"
	"github.com/B1NARY-GR0UP/dreamemo/memo"
	"github.com/B1NARY-GR0UP/dreamemo/strategy/eliminate"
	"github.com/B1NARY-GR0UP/dreamemo/util"
	"github.com/B1NARY-GR0UP/inquisitor/core"
	"sync"
)

type Group struct {
	name   string
	getter Getter
	memo   memo.Memo
}

var (
	mu     sync.RWMutex
	groups = make(map[string]*Group)
)

// NewGroup
// TODO: add cacheBytes; related to lazy init todo
func NewGroup(name string, getter Getter) *Group {
	if getter == nil {
		panic("Getter must not be nil")
	}
	mu.Lock()
	defer mu.Unlock()
	g := &Group{
		name:   name,
		getter: getter,
		memo:   memo.Memo{},
	}
	groups[name] = g
	return g
}

// GetGroup return correspond group related to the name
func GetGroup(name string) *Group {
	mu.RLock()
	defer mu.RUnlock()
	return groups[name]
}

func (g *Group) Name() string {
	return g.name
}

func (g *Group) Get(ctx context.Context, key string) (memo.ByteView, error) {
	// TODO: How to use context
	// TODO: refer to groupcache to improve logic
	if key == "" {
		return memo.ByteView{}, fmt.Errorf("key is null")
	}
	if v, ok := g.memo.Get(eliminate.Key(key)); ok {
		core.Info("memo hit")
		return v, nil
	}
	return g.load(ctx, key)
}

func (g *Group) load(ctx context.Context, key string) (memo.ByteView, error) {
	// TODO: refer to groupcache to improve logic
	return g.getLocally(ctx, key)
}

func (g *Group) getLocally(ctx context.Context, key string) (memo.ByteView, error) {
	bytes, err := g.getter.Get(ctx, key)
	if err != nil {
		return memo.ByteView{}, err
	}
	// TODO: refer to groupcache to improve logic
	value := memo.ByteView{
		B: util.CopyBytes(bytes),
	}
	g.populateMemo(key, value)
	return value, nil
}

func (g *Group) populateMemo(key string, value memo.ByteView) {
	// TODO: refer to groupcache to improve logic
	g.memo.Add(eliminate.Key(key), value)
}
