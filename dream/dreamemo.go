package dream

import (
	"flag"
	"github.com/B1NARY-GR0UP/dreamemo/app/server"
	"github.com/B1NARY-GR0UP/dreamemo/guidance"
	"github.com/B1NARY-GR0UP/dreamemo/memo"
	"github.com/B1NARY-GR0UP/dreamemo/strategy/eliminate/lru"
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

// ParseFlag quick start
// -addrs=:7246,:7247,:7248
// -addrs=:7247,:7246,:7248
// -addrs=:7248,:7246,:7247
// hint: first element is local instance
// TODO: 提供一个解析 flag 的函数，返回数组，包含地址配置
func ParseFlag() {
	flag.StringVar(&addrsFlag, addrsFlagName, addrsFlagDefaultValue, addrsFlagHint)
	flag.Parse()
}

// StandAlone in order to help user quick start
// StandAlone uses following default options
// StandAlone is in standalone mode, listen on :7246
// protocol             => protobuf
// eliminate strategy   => lru
// distributed strategy => consistent hash
// source               => redis
func StandAlone(opts ...Option) (*server.Engine, *guidance.Group) {
	// engine layer
	e := server.NewEngine()
	// eliminate layer
	l := lru.NewLRUCore()
	// memo layer
	m := memo.NewMemo(l)
	// guidance layer
	g := guidance.NewGroup(m, e)
	return e, g
}
