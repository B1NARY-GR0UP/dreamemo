package dream

const (
	defaultGroupName = "binary"
)

type Option func(o *Options)

type Options struct {
	GroupName string
}

func newOptions(opts ...Option) *Options {
	options := &Options{
		GroupName: defaultGroupName,
	}
	options.apply(opts...)
	return options
}

func (o *Options) apply(opts ...Option) {
	for _, opt := range opts {
		opt(o)
	}
}

func WithGroupName(name string) Option {
	return func(o *Options) {
		o.GroupName = name
	}
}
