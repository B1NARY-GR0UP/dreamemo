package memo

import (
	"github.com/B1NARY-GR0UP/dreamemo/strategy"
	"github.com/B1NARY-GR0UP/inquisitor/core"
	"sync"
)

// Memo ensures that all accesses to the core are concurrency safe
// Memo is an encapsulation of the strategy layer
type Memo struct {
	sync.RWMutex
	memo strategy.Memo
}

// Add to strategy.Memo
func (m *Memo) Add(key strategy.Key, value ByteView) {
	m.Lock()
	defer m.Unlock()
	// TODO: 这里需要做类型判断（使用 Option 或者其他形式来获取使用的 CoreMemo 类型）
	// TODO: 实现懒加载
	m.memo.Add(key, value)
}

// Get from strategy.Memo
func (m *Memo) Get(key strategy.Key) (ByteView, bool) {
	m.RLock()
	defer m.RUnlock()
	if m.memo == nil {
		return ByteView{}, false
	}
	value, ok := m.memo.Get(key)
	if !ok {
		return ByteView{}, false
	}
	return value.(ByteView), true
}

// Remove memo entity
func (m *Memo) Remove(key strategy.Key) {
	// TODO: api layer should ensure that the key must not be null
	m.Lock()
	defer m.Unlock()
	if m.memo == nil {
		core.Info("Memo is Empty")
		return
	}
	m.memo.Remove(key)
}

func (m *Memo) Clear() {
	m.Lock()
	defer m.Unlock()
	if m.memo == nil {
		core.Info("Memo is Empty")
		return
	}
	m.memo.Clear()
}
