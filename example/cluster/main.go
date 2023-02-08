package main

import (
	"context"
	"fmt"
	"net/http"

	"github.com/B1NARY-GR0UP/dreamemo/common/constant"
	"github.com/B1NARY-GR0UP/dreamemo/common/util"
	"github.com/B1NARY-GR0UP/dreamemo/dream"
	"github.com/B1NARY-GR0UP/dreamemo/guidance"
	"github.com/B1NARY-GR0UP/dreamemo/source"
	log "github.com/B1NARY-GR0UP/inquisitor/core"
	"github.com/B1NARY-GR0UP/piano/core"
	"github.com/B1NARY-GR0UP/piano/core/bin"
)

var db = map[string]string{
	"binary": "dreamemo",
	"hello":  "world",
	"foo":    "bar",
}

func getFromDB(_ context.Context, key string) ([]byte, error) {
	log.Info("Get from DB")
	if v, ok := db[key]; ok {
		return []byte(v), nil
	}
	return nil, fmt.Errorf("key %v is not exist", key)
}

func main() {
	addrs, api := util.ParseFlags()
	e := dream.Cluster(source.GetterFunc(getFromDB), addrs...)
	if api {
		go func() {
			p := bin.Default(core.WithHostAddr(":8080"))
			p.GET("/hello", func(ctx context.Context, pk *core.PianoKey) {
				key := pk.Query("key")
				g := guidance.GetGroup(constant.DefaultGroupName)
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
