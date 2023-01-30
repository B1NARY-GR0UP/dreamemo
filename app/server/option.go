package server

const (
	defaultBasePath = "/_dreamemo/"
	defaultHostAddr = ":7246"
)

type Option func(o *Options)

type Options struct {
	// TODO: add JSON option
	BasePath string
	Addr     string
}

func newOptions(opts ...Option) *Options {
	options := &Options{
		BasePath: defaultBasePath,
		Addr:     defaultHostAddr,
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
