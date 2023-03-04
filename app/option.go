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

package app

import (
	"context"
	"net/http"
	"time"

	"github.com/B1NARY-GR0UP/dreamemo/common/constant"
	"github.com/B1NARY-GR0UP/dreamemo/common/util"
	"github.com/B1NARY-GR0UP/dreamemo/strategy/distributed"
	"github.com/B1NARY-GR0UP/dreamemo/strategy/distributed/consistenthash"
)

type Option func(o *Options)

type Options struct {
	BasePath    string
	Addr        string
	Strategy    distributed.Dispatcher
	Thrift      bool
	Transport   func(context.Context) http.RoundTripper
	DetectDelay time.Duration
}

var defaultOptions = Options{
	BasePath:    constant.DefaultBasePath,
	Addr:        constant.DefaultStandAloneAddr,
	Strategy:    consistenthash.New(),
	Thrift:      false,
	Transport:   nil,
	DetectDelay: time.Second * 60,
}

// NewOptions used for both server and client
func NewOptions(opts ...Option) *Options {
	options := &Options{
		BasePath: defaultOptions.BasePath,
		Addr:     defaultOptions.Addr,
		// TODO: support more distributed strategy (will be supported)
		Strategy:    defaultOptions.Strategy,
		Thrift:      defaultOptions.Thrift,
		Transport:   defaultOptions.Transport,
		DetectDelay: defaultOptions.DetectDelay,
	}
	options.apply(opts...)
	return options
}

func (o *Options) apply(opts ...Option) {
	for _, opt := range opts {
		opt(o)
	}
}

// WithBasePath used to define base path of server
func WithBasePath(path string) Option {
	return func(o *Options) {
		o.BasePath = path
	}
}

// WithHostAddr used to define host addr server listens
func WithHostAddr(addr string) Option {
	addr = util.StandardizeAddr(addr)
	return func(o *Options) {
		o.Addr = addr
	}
}

// WithThrift0 must be consistent will guidance.WithThriftII
func WithThrift0() Option {
	return func(o *Options) {
		o.Thrift = true
	}
}

// WithTransport define your own http transport
func WithTransport(tpt func(context.Context) http.RoundTripper) Option {
	return func(o *Options) {
		o.Transport = tpt
	}
}

// WithDetectDelay define the delay of starting heartbeat detect
func WithDetectDelay(t time.Duration) Option {
	return func(o *Options) {
		o.DetectDelay = t
	}
}
