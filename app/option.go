package app

import (
	"github.com/B1NARY-GR0UP/dreamemo/common/constant"
	"github.com/B1NARY-GR0UP/dreamemo/strategy/distributed"
	"github.com/B1NARY-GR0UP/dreamemo/strategy/distributed/consistenthash"
)

type Option func(o *Options)

type Options struct {
	// TODO: add JSON option
	BasePath string
	Addr     string
	Strategy distributed.Instance
	// used for both server and client
	Thrift bool
}

var defaultOptions = Options{
	BasePath: constant.DefaultBasePath,
	Addr:     constant.DefaultStandAloneAddr,
	Strategy: consistenthash.New(),
	Thrift:   false,
}

func NewOptions(opts ...Option) *Options {
	options := &Options{
		BasePath: defaultOptions.BasePath,
		Addr:     defaultOptions.Addr,
		// TODO: support more distributed strategy (will be supported)
		Strategy: defaultOptions.Strategy,
		Thrift:   defaultOptions.Thrift,
	}
	options.apply(opts...)
	return options
}

func (o *Options) apply(opts ...Option) {
	for _, opt := range opts {
		opt(o)
	}
}

func WithBasePath(path string) Option {
	return func(o *Options) {
		o.BasePath = path
	}
}

func WithHostAddr(addr string) Option {
	return func(o *Options) {
		o.Addr = addr
	}
}

func WithThrift() Option {
	return func(o *Options) {
		o.Thrift = true
	}
}
