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

package eliminate

// ICore defines allowed operations of a memo
type ICore interface {
	Add(Key, Value)
	Get(Key) (Value, bool)
	Remove(Key)
	Clear()
	Name() string
}

type (
	Core struct {
		MaxSize   int // 0 => no limit
		UsedSize  int
		OnEvicted EvictFunc
	}
	EvictFunc func(key Key, value Value)
)

// Entity
type (
	Entity struct {
		Key   Key
		Value Value
	}
	Key   string
	Value any
)
