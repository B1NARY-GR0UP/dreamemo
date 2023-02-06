package guidance

import (
	"github.com/B1NARY-GR0UP/dreamemo/common/constant"
	"github.com/B1NARY-GR0UP/dreamemo/source"
	"github.com/B1NARY-GR0UP/dreamemo/source/redis"
)

type Option func(o *Options)

type Options struct {
	Name   string
	Getter source.Getter
}

var defaultOptions = Options{
	Name:   constant.DefaultGroupName,
	Getter: redis.NewSource(),
}

func newOptions(opts ...Option) *Options {
	options := &Options{
		Name:   defaultOptions.Name,
		Getter: defaultOptions.Getter,
	}
	options.apply(opts...)
	return options
}

func (o *Options) apply(opts ...Option) {
	for _, opt := range opts {
		opt(o)
	}
}

// WithName define name for group
func WithName(name string) Option {
	return func(o *Options) {
		o.Name = name
	}
}

// WithGetter define backend datasource
func WithGetter(getter source.Getter) Option {
	return func(o *Options) {
		o.Getter = getter
	}
}
