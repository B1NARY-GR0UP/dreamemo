package common

import "sync"

type SingleFlight interface {
	Do(key string, fn func() (any, error)) (any, error)
}

type (
	Group struct {
		mu sync.Mutex
		m  map[string]*call
	}
	call struct {
		wg  sync.WaitGroup
		val any
		err error
	}
)

func (g *Group) Do(key string, fn func() (any, error)) (any, error) {
	g.mu.Lock()
	if g.m == nil {
		g.m = make(map[string]*call)
	}
	if c, ok := g.m[key]; ok {
		g.mu.Unlock()
		c.wg.Wait()
		return c.val, c.err
	}
	c := new(call)
	c.wg.Add(1)
	g.m[key] = c
	g.mu.Unlock()

	c.val, c.err = fn()
	c.wg.Done()

	g.mu.Lock()
	delete(g.m, key)
	g.mu.Unlock()

	return c.val, c.err
}
