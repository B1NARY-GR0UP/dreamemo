package memo

import (
	"sync"

	"github.com/B1NARY-GR0UP/dreamemo/strategy/eliminate"
	"github.com/B1NARY-GR0UP/inquisitor/core"
)

// Memo ensures that all accesses to the guidance are concurrency safe
// Memo is an encapsulation of the strategy layer
type Memo struct {
	sync.RWMutex
	memo eliminate.ICore
}

// NewMemo return a new memo base on eliminate layer
func NewMemo(core eliminate.ICore) *Memo {
	return &Memo{
		memo: core,
	}
}

// Add to strategy.ICore
func (m *Memo) Add(key eliminate.Key, value ByteView) {
	m.Lock()
	defer m.Unlock()
	m.memo.Add(key, value)
}

// Get from strategy.ICore
func (m *Memo) Get(key eliminate.Key) (ByteView, bool) {
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
func (m *Memo) Remove(key eliminate.Key) {
	m.Lock()
	defer m.Unlock()
	if m.memo == nil {
		core.Info("---DREAMEMO--- Core is empty")
		return
	}
	m.memo.Remove(key)
}

func (m *Memo) Clear() {
	m.Lock()
	defer m.Unlock()
	if m.memo == nil {
		core.Info("---DREAMEMO--- Core is empty")
		return
	}
	m.memo.Clear()
}
