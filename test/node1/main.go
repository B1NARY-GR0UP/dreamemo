package main

/*
$ curl "http://localhost:9999/api?key=Tom"
630

$ curl "http://localhost:9999/api?key=kkk"
kkk not exist
*/

import (
	"context"
	"fmt"
	"github.com/B1NARY-GR0UP/dreamemo/app"
	"github.com/B1NARY-GR0UP/dreamemo/app/server"
	"github.com/B1NARY-GR0UP/dreamemo/guidance"
	"github.com/B1NARY-GR0UP/dreamemo/memo"
	"github.com/B1NARY-GR0UP/dreamemo/source"
	"github.com/B1NARY-GR0UP/dreamemo/strategy/eliminate/lru"
	"log"
)

var db = map[string]string{
	"red":   "#FF0000",
	"green": "#00FF00",
	"blue":  "#0000FF",
}

func createGroup() *guidance.Group {
	l := lru.NewLRUCore()
	m := memo.NewMemo(l)
	guidance.NewGroup(m, guidance.WithName("color"), guidance.WithGetter(source.GetterFunc(func(ctx context.Context, key string) ([]byte, error) {
		log.Println("[SlowDB] search key", key)
		if v, ok := db[key]; ok {
			return []byte(v), nil
		}
		return nil, fmt.Errorf("%s not exist", key)
	})))
	return guidance.GetGroup("color")
}

func startDreamemoServer(addr string, addrs []string, g *guidance.Group) {
	engine := server.NewEngine(app.WithHostAddr(addr[7:]))
	engine.Register(addrs...)
	//g.RegisterEngine(engine)
	log.Println("dreamemo is running at", addr)
	engine.Run()
}

func main() {
	addrMap := map[int]string{
		8001: "http://localhost:8001",
		8002: "http://localhost:8002",
		8003: "http://localhost:8003",
	}
	var addrs []string
	for _, v := range addrMap {
		addrs = append(addrs, v)
	}
	g := createGroup()
	startDreamemoServer("http://localhost:8001", addrs, g)

}
