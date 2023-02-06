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
	"sync/atomic"
)

// guidance is a runtime object that maintain by dreamemo
var guidance = struct {
	sync.RWMutex
	groups map[string]*Group
}{
	groups: make(map[string]*Group),
}

type Group struct {
	memo   *memo.Memo
	name   string
	getter source.Getter
	lbr    loadbalance.LoadBalancer
	sf     singleflight.SingleFlight
}

// NewGroup will not return a group pointer, use GetGroup function directly
// TODO: add cacheBytes; related to lazy init todo
func NewGroup(memo *memo.Memo, opts ...Option) {
	guidance.Lock()
	defer guidance.Unlock()
	options := newOptions(opts...)
	g := &Group{
		memo:   memo,
		name:   options.Name,
		getter: options.Getter,
		sf:     &singleflight.Group{},
		// TODO: initLoadBalancer according to user's options
	}
	guidance.groups[options.Name] = g
}

// GetGroup return correspond group related to the name
func GetGroup(name string) *Group {
	guidance.RLock()
	defer guidance.RUnlock()
	return guidance.groups[name]
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
	flagChanged := atomic.CompareAndSwapInt64(&util.RespFlag, 0, 1)
	if !flagChanged {
		panic("Flag must be changed")
	}
	// TODO: use apex GetRequest
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
	}, nil
}

func (g *Group) populateMemo(key string, value memo.ByteView) {
	// TODO: refer to groupcache to improve guidance
	g.memo.Add(eliminate.Key(key), value)
}
