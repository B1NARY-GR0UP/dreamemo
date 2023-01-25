package server

import (
	"github.com/B1NARY-GR0UP/dreamemo/api"
	"github.com/B1NARY-GR0UP/dreamemo/app/client"
	"github.com/B1NARY-GR0UP/dreamemo/strategy/distributed"
	"github.com/B1NARY-GR0UP/dreamemo/util"
	"net/http"
)

// Engine TODO: use listenAndServe to start server
type Engine struct {
	// TODO: choose lb method according to users option
	ins     distributed.Instance
	options *Options
	clients map[string]*client.Client
}

func NewEngine(opts ...Option) *Engine {
	options := NewOptions(opts...)
	e := &Engine{
		options: options,
	}
	return e
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
	matchedGroup := api.GetGroup(groupName)
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
