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

// quick start
// -addrs=:7246,:7247,:7248
// -addrs=:7247,:7246,:7248
// -addrs=:7248,:7246,:7247
// hint: first element is local instance
func initFlag() {
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
func StandAlone(opts ...Option) *guidance.Group {
	// TODO: 虽然是默认配置，但是每层的小配置是需要允许用户修改的
	// eliminate layer
	l := lru.NewLRUCore()
	// memo layer
	m := memo.NewMemo(l)
	// engine layer
	e := server.NewEngine()
	// guidance layer
	guidance.NewGroup(m)
	e.Run()
	return guidance.GetGroup("binary")
}
