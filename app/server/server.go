package server

import (
	"context"
	"net/http"
	"sync"

	"github.com/B1NARY-GR0UP/dreamemo/app/client"
	"github.com/B1NARY-GR0UP/dreamemo/common/util"
	"github.com/B1NARY-GR0UP/dreamemo/guidance"
	"github.com/B1NARY-GR0UP/dreamemo/loadbalance"
	"github.com/B1NARY-GR0UP/dreamemo/strategy/distributed"
)

var _ loadbalance.LoadBalancer = (*Engine)(nil)

// Engine TODO: use listenAndServe to start server
type Engine struct {
	// TODO: choose lb method according to users option
	sync.Mutex
	// addr of local instance
	// TODO: 需要在实例化时赋值？？
	self    string
	options *Options
	// an instance only holds its addr
	// TODO: support more field to an instance e.g. tags
	instances distributed.Instance
	clients   map[string]*client.Client
	// TODO: how to use
	Transport func(context.Context) http.RoundTripper
}

func NewEngine(group *guidance.Group, opts ...Option) *Engine {
	options := newOptions(opts...)
	e := &Engine{
		options: options,
	}
	return e
}

func (e *Engine) Run() {
	_ = http.ListenAndServe(e.options.Addr, e)
}

// Set instance should be a valid addr e.g. localhost:7246 localhost:7247 localhost:7248
func (e *Engine) Set(instances ...string) {
	e.Lock()
	defer e.Unlock()
	// TODO: decide to use which kind of distributed strategy according to the option
	e.instances.Add(instances...)
	e.clients = make(map[string]*client.Client, len(instances))
	for _, instance := range instances {
		e.clients[instance] = &client.Client{
			BasePath:  instance + e.options.BasePath,
			Transport: e.Transport,
		}
	}
}

// Pick an instance according to the given key
func (e *Engine) Pick(key string) (loadbalance.Instance, bool) {
	e.Lock()
	defer e.Unlock()
	ins := e.instances.Get(key)
	if ins == "" {
		return nil, false
	}
	if ins != e.self {
		return e.clients[ins], true
	}
	return nil, false
}

// ServeHTTP implements the http.Handler interface
func (e *Engine) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	segments := util.ParseRequestURL(req.URL.Path, e.options.BasePath)
	if segments == nil {
		http.Error(w, "Bad Request", http.StatusBadRequest)
		return
	}
	groupName := segments[0]
	key := segments[1]
	matchedGroup := guidance.GetGroup(groupName)
	if matchedGroup == nil {
		http.Error(w, "No Such Group: "+groupName, http.StatusBadRequest)
		return
	}
	// TODO: use req.Context() or other context, refer to groupcache
	// TODO: add context field to Engine, use sync.Pool to optimize
	byteView, err := matchedGroup.Get(req.Context(), key)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	// TODO: use JSON, protobuf, thrift, ByteSlice, string according to user's option
	w.Header().Set("Content-Type", "application/octet-stream")
	_, _ = w.Write(byteView.ByteSlice())
}
