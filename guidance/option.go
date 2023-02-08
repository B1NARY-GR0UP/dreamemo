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

package guidance

import (
	"github.com/B1NARY-GR0UP/dreamemo/common/constant"
	"github.com/B1NARY-GR0UP/dreamemo/source"
	"github.com/B1NARY-GR0UP/dreamemo/source/redis"
)

type Option func(o *Options)

type Options struct {
	name   string
	thrift bool
	getter source.Getter
}

var defaultOptions = Options{
	name:   constant.DefaultGroupName,
	thrift: false,
	getter: redis.NewSource(),
}

func newOptions(opts ...Option) *Options {
	options := &Options{
		name:   defaultOptions.name,
		thrift: defaultOptions.thrift,
		getter: defaultOptions.getter,
	}
	options.apply(opts...)
	return options
}

func (o *Options) apply(opts ...Option) {
	for _, opt := range opts {
		opt(o)
	}
}

// WithGroupName define name for group
func WithGroupName(name string) Option {
	return func(o *Options) {
		o.name = name
	}
}

// WithThrift1 must be consistent will app.WithThrift0
func WithThrift1() Option {
	return func(o *Options) {
		o.thrift = true
	}
}

// WithGetter define backend datasource
func WithGetter(getter source.Getter) Option {
	return func(o *Options) {
		o.getter = getter
	}
}
