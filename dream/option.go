package dream

const (
	defaultGroupName = "binary"
	defaultHostAddr  = ":7246"
)

type Option func(o *Options)

type Options struct {
}

func newOptions(opts ...Option) *Options {
	options := &Options{}
	options.apply(opts...)
	return options
}

func (o *Options) apply(opts ...Option) {
	for _, opt := range opts {
		opt(o)
	}
}
