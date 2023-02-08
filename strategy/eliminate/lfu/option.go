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

package lfu

import "github.com/B1NARY-GR0UP/dreamemo/strategy/eliminate"

type Option func(o *Options)

type Options struct {
	MaxSize   int
	OnEvicted eliminate.EvictFunc
}

var defaultOptions = Options{
	MaxSize:   0,
	OnEvicted: nil,
}

func newOptions(opts ...Option) *Options {
	options := &Options{
		MaxSize:   defaultOptions.MaxSize,
		OnEvicted: defaultOptions.OnEvicted,
	}
	options.apply(opts...)
	return options
}

func (o *Options) apply(opts ...Option) {
	for _, opt := range opts {
		opt(o)
	}
}

func WithMaxSize(size int) Option {
	return func(o *Options) {
		o.MaxSize = size
	}
}

func WithEvictFunc(evict eliminate.EvictFunc) Option {
	return func(o *Options) {
		o.OnEvicted = evict
	}
}
