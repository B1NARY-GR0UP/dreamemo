package util

import (
	"github.com/B1NARY-GR0UP/inquisitor/core"
	"strings"
)

// RespFlag is used to judge response type
// default value is 0, means response in thrift or protobuf
// if change it to 1, means response in text, JSON or other optional type
var RespFlag int64 = 0

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
		core.Warnf("[DREAMEMO] Request URL is Invalid: %v", reqPath)
		return nil
	}
	idx := strings.LastIndex(reqPath, "/")
	if idx == len(reqPath)-1 {
		reqPath = reqPath[:len(reqPath)-1]
	}
	segments := strings.Split(reqPath[len(basePath):], "/")
	if len(segments) != 2 {
		core.Warnf("[DREAMEMO] Request URL is Invalid: %v", reqPath)
		return nil
	}
	return segments
}
