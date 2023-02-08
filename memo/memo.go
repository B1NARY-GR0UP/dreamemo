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
