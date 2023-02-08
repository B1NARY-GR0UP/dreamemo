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
