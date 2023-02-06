package consistenthash

import "hash/crc32"

type Option func(o *Options)

// Options TODO: options should be allowed to edit to user
type Options struct {
	HashFunc          HashFunc
	ReplicationFactor int
}

var defaultOptions = Options{
	HashFunc:          crc32.ChecksumIEEE,
	ReplicationFactor: 10,
}

func newOptions(opts ...Option) *Options {
	options := &Options{
		HashFunc:          defaultOptions.HashFunc,
		ReplicationFactor: defaultOptions.ReplicationFactor,
	}
	options.apply(opts...)
	return options
}

func (o *Options) apply(opts ...Option) {
	for _, opt := range opts {
		opt(o)
	}
}

func WithHashFunc(hashFunc HashFunc) Option {
	return func(o *Options) {
		o.HashFunc = hashFunc
	}
}

func WithReplicationFactor(factor int) Option {
	return func(o *Options) {
		o.ReplicationFactor = factor
	}
}
