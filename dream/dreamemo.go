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

package dream

import (
	"github.com/B1NARY-GR0UP/dreamemo/app"
	"github.com/B1NARY-GR0UP/dreamemo/app/server"
	"github.com/B1NARY-GR0UP/dreamemo/guidance"
	"github.com/B1NARY-GR0UP/dreamemo/memo"
	"github.com/B1NARY-GR0UP/dreamemo/source"
	"github.com/B1NARY-GR0UP/dreamemo/strategy/eliminate/lru"
)

// StandAlone start in standalone mode
// uses following default options:
// protocol             => protobuf
// eliminate strategy   => lru
// distributed strategy => consistent hash
func StandAlone(getter source.Getter) {
	c := lru.NewLRUCore()
	m := memo.NewMemo(c)
	guidance.NewGroup(m, nil, guidance.WithGetter(getter))
}

// Cluster start in cluster mode
// uses following default options:
// protocol             => protobuf
// eliminate strategy   => lru
// distributed strategy => consistent hash
func Cluster(getter source.Getter, addrs ...string) *server.Engine {
	e := server.NewEngine(app.WithHostAddr(addrs[0]))
	e.RegisterInstances(addrs...)
	c := lru.NewLRUCore()
	m := memo.NewMemo(c)
	guidance.NewGroup(m, e, guidance.WithGetter(getter))
	return e
}
