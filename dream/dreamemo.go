package dream

import (
	"flag"
	"net/http"
	"strings"

	"github.com/B1NARY-GR0UP/dreamemo/app/server"
	"github.com/B1NARY-GR0UP/dreamemo/guidance"
	"github.com/B1NARY-GR0UP/dreamemo/memo"
	"github.com/B1NARY-GR0UP/dreamemo/source/redis"
	"github.com/B1NARY-GR0UP/dreamemo/strategy/eliminate/lfu"
)

// Dreamemo Have we ever been sober
type Dreamemo struct {
	options *Options
}

const (
	// addrsFlag
	addrsFlagName         = "addrs"
	addrsFlagDefaultValue = ":7246"
	addrsFlagHint         = "instances addresses"
	// TODO: add more flags
)

var addrsFlag string

func initFlag() {
	// quick start
	// -addrs=:7246,:7247,:7248
	// -addrs=:7247,:7246,:7248
	// -addrs=:7248,:7246,:7247
	// hint: first element is local instance
	flag.StringVar(&addrsFlag, addrsFlagName, addrsFlagDefaultValue, addrsFlagHint)
	flag.Parse()
}

// Default in order to help user quick start
// Default uses following default options
// protocol             => protobuf
// eliminate strategy   => lru
// distributed strategy => consistent hash
// source               => redis
func Default(opts ...Option) {
	initFlag()
	addrs := strings.Split(addrsFlag, ",")
	options := newOptions(opts...)
	// TODO: group should be a field of engine to init
	m := memo.NewMemo(lfu.NewLFUCore(-1, nil))
	group := guidance.NewGroup(m, options.GroupName, redis.NewSource())
	// NewEngine(NewGroup(NewMemo(NewCore(opts), opts), opts), opts)
	engine := server.NewEngine(group, server.WithHostAddr(addrs[0]))
	engine.Set(addrs...)
	_ = http.ListenAndServe(addrs[0], engine)
}
