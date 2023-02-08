package util

import (
	"flag"
	"strings"

	"github.com/B1NARY-GR0UP/inquisitor/core"
)

// ParseFlags e.g.
// --addrs=http://localhost:7246,http://localhost:7247,http://localhost:7248 --api
// --addrs=http://localhost:7247,http://localhost:7248,http://localhost:7246
// --addrs=http://localhost:7248,http://localhost:7246,http://localhost:7247
// hint: first element is local instance
func ParseFlags() (addrs []string, api bool) {
	var addrsFlag string
	var apiFlag bool
	flag.StringVar(&addrsFlag, "addrs", "", "instances addresses")
	flag.BoolVar(&apiFlag, "api", false, "start api or not")
	flag.Parse()
	return strings.Split(addrsFlag, ","), apiFlag
}

func CopyBytes(b []byte) []byte {
	c := make([]byte, len(b))
	copy(c, b)
	return c
}

// ParseRequestURL
// Default Request Form: host:port/_dreamemo/group/key
// segments[0]: group
// segments[1]: key
func ParseRequestURL(reqPath, basePath string) []string {
	if !strings.HasPrefix(reqPath, basePath) {
		core.Warnf("---DREAMEMO--- Request URL is Invalid: %v", reqPath)
		return nil
	}
	idx := strings.LastIndex(reqPath, "/")
	if idx == len(reqPath)-1 {
		reqPath = reqPath[:len(reqPath)-1]
	}
	segments := strings.Split(reqPath[len(basePath):], "/")
	if len(segments) != 2 {
		core.Warnf("---DREAMEMO--- Request URL is Invalid: %v", reqPath)
		return nil
	}
	return segments
}

func StandardizeAddr(addr string) string {
	segments := strings.Split(addr, "://")
	length := len(segments)
	if length == 1 {
		return segments[0]
	}
	if length == 2 {
		return segments[1]
	}
	return ""
}
