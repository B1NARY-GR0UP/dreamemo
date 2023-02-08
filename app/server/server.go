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

package server

import (
	"context"
	"net/http"
	"strings"
	"sync"

	"github.com/B1NARY-GR0UP/dreamemo/app"
	"github.com/B1NARY-GR0UP/dreamemo/app/client"
	"github.com/B1NARY-GR0UP/dreamemo/common/util"
	"github.com/B1NARY-GR0UP/dreamemo/guidance"
	"github.com/B1NARY-GR0UP/dreamemo/loadbalance"
	"github.com/B1NARY-GR0UP/dreamemo/protocol/protobuf"
	pthrift "github.com/B1NARY-GR0UP/dreamemo/protocol/thrift"
	"github.com/B1NARY-GR0UP/dreamemo/strategy/distributed"
	"github.com/B1NARY-GR0UP/inquisitor/core"
	"github.com/apache/thrift/lib/go/thrift"
	"google.golang.org/protobuf/proto"
)

var _ loadbalance.LoadBalancer = (*Engine)(nil)

// Engine server engine of each instance
type Engine struct {
	sync.Mutex
	options *app.Options
	// addr of local instance
	self string
	// an instance only holds its addr
	instances distributed.Instance
	clients   map[string]*client.Client
	// TODO: (check needed)
	// TODO: add to options
	Transport func(context.Context) http.RoundTripper
}

func NewEngine(opts ...app.Option) *Engine {
	options := app.NewOptions(opts...)
	e := &Engine{
		// TODO: may cause bug, need secondly check
		options:   options,
		self:      options.Addr,
		instances: options.Strategy,
	}
	return e
}

// Run is used to start cluster, should not be used in standalone mode
func (e *Engine) Run() {
	core.Infof("---DREAMEMO--- Server is listening on %v", e.options.Addr)
	err := http.ListenAndServe(e.options.Addr, e)
	if err != nil {
		core.Errorf("---DREAMEMO--- Server started failed: %v", err)
	}
}

// RegisterInstances instance should be a valid addr e.g. localhost:7246 localhost:7247 localhost:7248
func (e *Engine) RegisterInstances(insts ...string) {
	e.Lock()
	defer e.Unlock()
	e.instances.Add(insts...)
	e.clients = make(map[string]*client.Client, len(insts))
	for _, instance := range insts {
		e.clients[instance] = &client.Client{
			Options:   e.options,
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
	if !strings.Contains(ins, e.self) {
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
	// TODO: add context field to Engine, use sync.Pool to optimize (will be supported)
	byteView, err := matchedGroup.Get(req.Context(), key)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	if e.options.Thrift {
		serializer := thrift.NewTSerializer()
		body, err := serializer.Write(context.Background(), &pthrift.GetResponse{Value: byteView.ByteSlice()})
		if err != nil {
			core.Warnf("---DREAMEMO--- thrift serialize err: %v", err)
		}
		w.Header().Set("Content-Type", "application/octet-stream")
		_, _ = w.Write(body)
	} else {
		body, err := proto.Marshal(&protobuf.GetResponse{Value: byteView.ByteSlice()})
		if err != nil {
			core.Warnf("---DREAMEMO--- Protobuf marshal err: %v", err)
		}
		w.Header().Set("Content-Type", "application/octet-stream")
		_, _ = w.Write(body)
	}
}
