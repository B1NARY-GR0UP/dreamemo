package guidance

import (
	"context"
	"fmt"
	"github.com/B1NARY-GR0UP/dreamemo/common/singleflight"
	"github.com/B1NARY-GR0UP/dreamemo/common/util"
	"github.com/B1NARY-GR0UP/dreamemo/loadbalance"
	"github.com/B1NARY-GR0UP/dreamemo/memo"
	"github.com/B1NARY-GR0UP/dreamemo/protocol/protobuf"
	"github.com/B1NARY-GR0UP/dreamemo/source"
	"github.com/B1NARY-GR0UP/dreamemo/strategy/eliminate"
	"github.com/B1NARY-GR0UP/inquisitor/core"
	"sync"
)

type Group struct {
	name   string
	getter source.Getter
	memo   *memo.Memo
	lbr    loadbalance.LoadBalancer
	sf     singleflight.SingleFlight
}

var syncGroups = struct {
	sync.RWMutex
	groups map[string]*Group
}{
	groups: make(map[string]*Group),
}

// NewGroup
// TODO: add cacheBytes; related to lazy init todo
func NewGroup(memo *memo.Memo, name string, getter source.Getter) *Group {
	if getter == nil {
		panic("Getter must not be nil")
	}
	syncGroups.Lock()
	defer syncGroups.Unlock()
	g := &Group{
		name:   name,
		getter: getter,
		memo:   memo,
		sf:     &singleflight.Group{},
		// TODO: initLoadBalancer according to user's options
	}
	syncGroups.groups[name] = g
	return g
}

// GetGroup return correspond group related to the name
func GetGroup(name string) *Group {
	syncGroups.RLock()
	defer syncGroups.RUnlock()
	return syncGroups.groups[name]
}

func (g *Group) Name() string {
	return g.name
}

func (g *Group) Get(ctx context.Context, key string) (memo.ByteView, error) {
	// TODO: How to use context
	// TODO: refer to groupcache to improve guidance
	if key == "" {
		return memo.ByteView{}, fmt.Errorf("key is null")
	}
	if v, ok := g.memo.Get(eliminate.Key(key)); ok {
		core.Info("[DREAMEMO] ICore Hit")
		return v, nil
	}
	return g.load(ctx, key)
}

func (g *Group) load(ctx context.Context, key string) (memo.ByteView, error) {
	// TODO: review guidance
	bv, err := g.sf.Do(key, func() (any, error) {
		if g.lbr != nil {
			if ins, ok := g.lbr.Pick(key); ok {
				if value, err := g.getFromInstance(ctx, ins, key); err != nil {
					return value, nil
				}
			}
		}
		core.Info("[DREAMEMO] Get Locally")
		return g.getLocally(ctx, key)
	})
	if err != nil {
		return bv.(memo.ByteView), nil
	}
	return memo.ByteView{}, err
}

func (g *Group) getLocally(ctx context.Context, key string) (memo.ByteView, error) {
	bytes, err := g.getter.Get(ctx, key)
	if err != nil {
		return memo.ByteView{}, err
	}
	// TODO: refer to groupcache to improve guidance
	value := memo.ByteView{
		B: util.CopyBytes(bytes),
	}
	g.populateMemo(key, value)
	return value, nil
}

func (g *Group) getFromInstance(ctx context.Context, instance loadbalance.Instance, key string) (memo.ByteView, error) {
	// TODO: support thrift
	req := &protobuf.GetRequest{
		Group: g.name,
		Key:   key,
	}
	resp := &protobuf.GetResponse{}
	err := instance.Get(ctx, req, resp)
	if err != nil {
		return memo.ByteView{}, err
	}
	return memo.ByteView{
		B: resp.Value,
	}, err
}

func (g *Group) populateMemo(key string, value memo.ByteView) {
	// TODO: refer to groupcache to improve guidance
	g.memo.Add(eliminate.Key(key), value)
}
