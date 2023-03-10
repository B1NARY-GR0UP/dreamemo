// Copyright 2023 BINARY Members
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
// http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.
//

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

// WithHashFunc define your hash function
func WithHashFunc(hashFunc HashFunc) Option {
	return func(o *Options) {
		o.HashFunc = hashFunc
	}
}

// WithReplicationFactor define consistent hash replication factor
func WithReplicationFactor(factor int) Option {
	return func(o *Options) {
		o.ReplicationFactor = factor
	}
}
