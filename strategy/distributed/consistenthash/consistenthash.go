package consistenthash

import (
	"fmt"
	"sort"
)

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
	hash := h.options.HashFunc([]byte(key))
	idx := sort.Search(len(h.ring), func(i int) bool {
		return h.ring[i] >= hash
	})
	if idx >= len(h.ring) {
		idx = 0
	}
	return h.nodes[h.ring[idx]]
}

func (h *Hash) Name() string {
	return "consistenthash"
}
