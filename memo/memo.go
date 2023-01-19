package memo

import (
	"github.com/B1NARY-GR0UP/dreamemo/strategy"
	"sync"
)

// Memo ensures that all accesses to the core are concurrency safe
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
