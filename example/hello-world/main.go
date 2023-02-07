package main

import (
	"context"
	"net/http"

	"github.com/B1NARY-GR0UP/dreamemo/common/constant"
	"github.com/B1NARY-GR0UP/dreamemo/dream"
	"github.com/B1NARY-GR0UP/dreamemo/guidance"
	"github.com/B1NARY-GR0UP/piano/core"
	"github.com/B1NARY-GR0UP/piano/core/bin"
)

func main() {
	e := dream.StandAlone()
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
	e.Run()
}
