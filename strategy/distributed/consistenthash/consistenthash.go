package consistenthash

import (
	"fmt"
	"github.com/B1NARY-GR0UP/dreamemo/strategy/distributed"
	"sort"
)

var _ distributed.Instance = (*Hash)(nil)

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
	// same key will get same hash so this ensures that a picked instance won't pick another instance
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
