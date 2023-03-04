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

import (
	"fmt"
	"sort"

	"github.com/B1NARY-GR0UP/dreamemo/common/util"
	"github.com/B1NARY-GR0UP/dreamemo/strategy/distributed"
	"github.com/B1NARY-GR0UP/inquisitor/core"
)

var _ distributed.Dispatcher = (*Hash)(nil)

type Hash struct {
	ring    []uint32
	nodes   map[uint32]string
	options *Options
}

type HashFunc func(data []byte) uint32

func New(opts ...Option) *Hash {
	options := newOptions(opts...)
	return &Hash{
		ring:    make([]uint32, 0),
		nodes:   make(map[uint32]string, 0),
		options: options,
	}
}

func (h *Hash) Add(nodes ...string) {
	for _, node := range nodes {
		for i := 0; i < h.options.ReplicationFactor; i++ {
			hash := h.options.HashFunc([]byte(fmt.Sprintf("%s%d", node, i)))
			h.nodes[hash] = node
			h.ring = append(h.ring, hash)
		}
	}
	sort.Slice(h.ring, func(i, j int) bool {
		return h.ring[i] < h.ring[j]
	})
}

func (h *Hash) Get(key string) string {
	if len(h.nodes) == 0 {
		return ""
	}
	// same key will get same hash so this ensures that a picked node won't pick another node
	hash := h.options.HashFunc([]byte(key))
	idx := sort.Search(len(h.ring), func(i int) bool {
		return h.ring[i] >= hash
	})
	if idx >= len(h.ring) {
		idx = 0
	}
	return h.nodes[h.ring[idx]]
}

func (h *Hash) Remove(key string) {
	if key == "" {
		core.Warn("---DREAMEMO--- Key should not be empty")
		return
	}
	for i := 0; i < h.options.ReplicationFactor; i++ {
		hash := h.options.HashFunc([]byte(fmt.Sprintf("%s%d", key, i)))
		delete(h.nodes, hash)
		idx := util.SearchUint32s(h.ring, hash)
		if idx == -1 {
			core.Error("---DREAMEMO--- Down node removed failed")
			return
		}
		copy(h.ring[idx:], h.ring[idx+1:])
		h.ring = h.ring[:len(h.ring)-1]
	}
}

func (h *Hash) Name() string {
	return "consistenthash"
}
