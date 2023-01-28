package singleflight

import "sync"

type SingleFlight interface {
	Do(key string, fn func() (any, error)) (any, error)
}

type (
	Group struct {
		sync.Mutex
		m map[string]*call
	}
	call struct {
		wg  sync.WaitGroup
		val any
		err error
	}
)

func (g *Group) Do(key string, fn func() (any, error)) (any, error) {
	g.Lock()
	if g.m == nil {
		g.m = make(map[string]*call)
	}
	if c, ok := g.m[key]; ok {
		g.Unlock()
		c.wg.Wait()
		return c.val, c.err
	}
	c := new(call)
	c.wg.Add(1)
	g.m[key] = c
	g.Unlock()

	c.val, c.err = fn()
	c.wg.Done()

	g.Lock()
	delete(g.m, key)
	g.Unlock()

	return c.val, c.err
}
