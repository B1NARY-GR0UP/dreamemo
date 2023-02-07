package main

/*
$ curl "http://localhost:9999/api?key=Tom"
630

$ curl "http://localhost:9999/api?key=kkk"
kkk not exist
*/

import (
	"context"
	"flag"
	"fmt"
	"github.com/B1NARY-GR0UP/dreamemo/app"
	"github.com/B1NARY-GR0UP/dreamemo/app/server"
	"github.com/B1NARY-GR0UP/dreamemo/guidance"
	"github.com/B1NARY-GR0UP/dreamemo/memo"
	"github.com/B1NARY-GR0UP/dreamemo/source"
	"github.com/B1NARY-GR0UP/dreamemo/strategy/eliminate/lru"
	"log"
	"net/http"
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
	g.RegisterEngine(engine)
	log.Println("dreamemo is running at", addr)
	engine.Run()
}

func startAPIServer(apiAddr string, g *guidance.Group) {
	http.Handle("/api", http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			key := r.URL.Query().Get("key")
			view, err := g.Get(context.Background(), key)
			fmt.Println("api server", view)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			w.Header().Set("Content-Type", "application/octet-stream")
			w.Write(view.ByteSlice())
		}))
	log.Println("fontend server is running at", apiAddr)
	log.Fatal(http.ListenAndServe(apiAddr[7:], nil))
}

func main() {
	var port int
	var api bool
	flag.IntVar(&port, "port", 8001, "Geecache server port")
	flag.BoolVar(&api, "api", false, "Start a api server?")
	flag.Parse()

	apiAddr := "http://localhost:9999"
	addrMap := map[int]string{
		8001: "http://localhost:8001",
		8002: "http://localhost:8002",
		8003: "http://localhost:8003",
	}

	var addrs []string
	for _, v := range addrMap {
		addrs = append(addrs, v)
	}

	gee := createGroup()
	// 只有 8003 才启动 api server
	if api {
		go startAPIServer(apiAddr, gee)
	}
	startDreamemoServer(addrMap[port], addrs, gee)
}
