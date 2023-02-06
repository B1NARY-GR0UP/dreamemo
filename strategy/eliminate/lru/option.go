package lru

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

func WithEvictFunc(fn eliminate.EvictFunc) Option {
	return func(o *Options) {
		o.OnEvicted = fn
	}
}
