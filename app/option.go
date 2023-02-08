package app

import (
	"github.com/B1NARY-GR0UP/dreamemo/common/constant"
	"github.com/B1NARY-GR0UP/dreamemo/common/util"
	"github.com/B1NARY-GR0UP/dreamemo/strategy/distributed"
	"github.com/B1NARY-GR0UP/dreamemo/strategy/distributed/consistenthash"
)

type Option func(o *Options)

type Options struct {
	BasePath string
	Addr     string
	Strategy distributed.Instance
	Thrift   bool
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
