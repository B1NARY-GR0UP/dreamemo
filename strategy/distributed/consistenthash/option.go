package consistenthash

import "hash/crc32"

type Option struct {
	F func(o *Options)
}

type Options struct {
	HashFunc          HashFunc
	ReplicationFactor int
}

func newOptions(opts ...Option) *Options {
	options := &Options{
		HashFunc:          crc32.ChecksumIEEE,
		ReplicationFactor: 10,
	}
	options.apply(opts...)
	return options
}

func (o *Options) apply(opts ...Option) {
	for _, opt := range opts {
		opt.F(o)
	}
}

func WithHashFunc(hashFunc HashFunc) Option {
	return Option{F: func(o *Options) {
		o.HashFunc = hashFunc
	}}
}

func WithReplicationFactor(num int) Option {
	return Option{F: func(o *Options) {
		o.ReplicationFactor = num
	}}
}
