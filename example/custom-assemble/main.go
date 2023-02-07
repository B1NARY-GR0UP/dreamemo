package main

import (
	"context"
	"fmt"
	"log"
	"net/http"

	"github.com/B1NARY-GR0UP/dreamemo/app"
	"github.com/B1NARY-GR0UP/dreamemo/app/server"
	"github.com/B1NARY-GR0UP/dreamemo/guidance"
	"github.com/B1NARY-GR0UP/dreamemo/memo"
	"github.com/B1NARY-GR0UP/dreamemo/source"
	"github.com/B1NARY-GR0UP/dreamemo/strategy/eliminate/lfu"
)

var db = map[string]string{
	"binary": "dreamemo",
	"hello":  "world",
	"foo":    "bar",
}

func main() {
	addrs := []string{"http://localhost:8001", "http://localhost:8002", "http://localhost:8003"}
	e := server.NewEngine(app.WithHostAddr("localhost:8003"), app.WithThriftI())
	e.RegisterInstances(addrs...)
	l := lfu.NewLFUCore()
	m := memo.NewMemo(l)
	guidance.NewGroup(m, e, guidance.WithName("color"), guidance.WithGetter(source.GetterFunc(func(ctx context.Context, key string) ([]byte, error) {
		log.Println("[SlowDB] search key", key)
		if v, ok := db[key]; ok {
			return []byte(v), nil
		}
		return nil, fmt.Errorf("%s not exist", key)
	})), guidance.WithThriftII())
	go startAPIServer("localhost:9999", guidance.GetGroup("color"))
	e.Run()
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
			_, _ = w.Write(view.ByteSlice())
		}))
	log.Println("api server is running at", apiAddr)
	log.Fatal(http.ListenAndServe(apiAddr, nil))
}
