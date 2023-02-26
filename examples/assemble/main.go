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

package main

import (
	"context"
	"fmt"
	"net/http"

	"github.com/B1NARY-GR0UP/dreamemo/app"
	"github.com/B1NARY-GR0UP/dreamemo/app/server"
	"github.com/B1NARY-GR0UP/dreamemo/common/util"
	"github.com/B1NARY-GR0UP/dreamemo/guidance"
	"github.com/B1NARY-GR0UP/dreamemo/memo"
	"github.com/B1NARY-GR0UP/dreamemo/source"
	"github.com/B1NARY-GR0UP/dreamemo/strategy/eliminate/lfu"
	log "github.com/B1NARY-GR0UP/inquisitor/core"
	"github.com/B1NARY-GR0UP/piano/core"
	"github.com/B1NARY-GR0UP/piano/core/bin"
)

var db = map[string]string{
	"binary": "dreamemo",
	"hello":  "world",
	"ping":   "pong",
}

func getFromDB(_ context.Context, key string) ([]byte, error) {
	log.Info("Get from DB")
	if v, ok := db[key]; ok {
		return []byte(v), nil
	}
	return nil, fmt.Errorf("key %v is not exist", key)
}

// go run . --addrs=http://localhost:7246,http://localhost:7247,http://localhost:7248 --api
// go run . --addrs=http://localhost:7247,http://localhost:7248,http://localhost:7246
// go run . --addrs=http://localhost:7248,http://localhost:7246,http://localhost:7247
// curl localhost:8080/hello?key=ping
func main() {
	addrs, api := util.ParseFlags()
	e := server.NewEngine(app.WithHostAddr(addrs[0]), app.WithThrift0())
	e.RegisterNodes(addrs...)
	l := lfu.NewLFUCore()
	m := memo.NewMemo(l)
	guidance.NewGroup(m, e, guidance.WithGroupName("hello"), guidance.WithThrift1(), guidance.WithGetter(source.GetterFunc(getFromDB)))
	if api {
		go func() {
			p := bin.Default(core.WithHostAddr(":8080"))
			p.GET("/hello", func(ctx context.Context, pk *core.PianoKey) {
				key := pk.Query("key")
				g := guidance.GetGroup("hello")
				value, _ := g.Get(ctx, key)
				pk.JSON(http.StatusOK, core.M{
					key: value.String(),
				})
			})
			p.Play()
		}()
	}
	e.Run()
}
