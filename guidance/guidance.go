package guidance

import (
	"context"
	"fmt"
	"sync"

	"github.com/B1NARY-GR0UP/dreamemo/common/singleflight"
	"github.com/B1NARY-GR0UP/dreamemo/common/util"
	"github.com/B1NARY-GR0UP/dreamemo/loadbalance"
	"github.com/B1NARY-GR0UP/dreamemo/memo"
	"github.com/B1NARY-GR0UP/dreamemo/protocol/protobuf"
	"github.com/B1NARY-GR0UP/dreamemo/protocol/thrift"
	"github.com/B1NARY-GR0UP/dreamemo/strategy/eliminate"
	"github.com/B1NARY-GR0UP/inquisitor/core"
)

// guidance is a runtime object that maintain by dreamemo
var guidance = struct {
	sync.RWMutex
	groups map[string]*Group
}{
	groups: make(map[string]*Group),
}

type Group struct {
	options *Options
	memo    *memo.Memo
	engine  loadbalance.LoadBalancer
	sf      singleflight.SingleFlight
}

// NewGroup will not return a group pointer, use GetGroup function directly
func NewGroup(memo *memo.Memo, engine loadbalance.LoadBalancer, opts ...Option) {
	guidance.Lock()
	defer guidance.Unlock()
	options := newOptions(opts...)
	g := &Group{
		options: options,
		memo:    memo,
		engine:  engine,
		sf:      &singleflight.Group{},
	}
	guidance.groups[options.name] = g
}

// GetGroup return correspond group related to the name
func GetGroup(name string) *Group {
	guidance.RLock()
	defer guidance.RUnlock()
	return guidance.groups[name]
}

func (g *Group) Name() string {
	return g.options.name
}

func (g *Group) Get(ctx context.Context, key string) (memo.ByteView, error) {
	if key == "" {
		return memo.ByteView{}, fmt.Errorf("key is null")
	}
	if v, ok := g.memo.Get(eliminate.Key(key)); ok {
		core.Info("---DREAMEMO--- Core hit")
		return v, nil
	}
	return g.load(ctx, key)
}

func (g *Group) load(ctx context.Context, key string) (memo.ByteView, error) {
	bv, err := g.sf.Do(key, func() (any, error) {
		if g.engine != nil {
			if ins, ok := g.engine.Pick(key); ok {
				core.Info("---DREAMEMO--- Get from other instance")
				value, err := g.getFromInstance(ctx, ins, key)
				if err != nil {
					return memo.ByteView{}, err
				}
				g.populateMemo(key, value)
				return value, nil
			}
		}
		core.Info("---DREAMEMO--- Get locally")
		return g.getLocally(ctx, key)
	})
	if err != nil {
		return memo.ByteView{}, err
	}
	return bv.(memo.ByteView), nil
}

func (g *Group) getLocally(ctx context.Context, key string) (memo.ByteView, error) {
	bytes, err := g.options.getter.Get(ctx, key)
	if err != nil {
		return memo.ByteView{}, err
	}
	value := memo.ByteView{
		B: util.CopyBytes(bytes),
	}
	g.populateMemo(key, value)
	return value, nil
}

func (g *Group) getFromInstance(ctx context.Context, instance loadbalance.Instance, key string) (memo.ByteView, error) {
	if g.options.thrift {
		req := &thrift.GetRequest{
			Group: g.options.name,
			Key:   key,
		}
		resp := &thrift.GetResponse{}
		err := instance.Get(ctx, req, resp)
		if err != nil {
			return memo.ByteView{}, err
		}
		return memo.ByteView{
			B: resp.Value,
		}, nil
	} else {
		req := &protobuf.GetRequest{
			Group: g.options.name,
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
}

func (g *Group) populateMemo(key string, value memo.ByteView) {
	g.memo.Add(eliminate.Key(key), value)
}
