package util

import (
	"github.com/B1NARY-GR0UP/inquisitor/core"
	"strings"
)

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
		core.Warnf("Request URL is invalid: %v", reqPath)
		return nil
	}
	idx := strings.LastIndex(reqPath, "/")
	if idx == len(reqPath)-1 {
		reqPath = reqPath[:len(reqPath)-1]
	}
	segments := strings.Split(reqPath[len(basePath):], "/")
	if len(segments) != 2 {
		core.Warnf("Request URL is invalid: %v", reqPath)
		return nil
	}
	return segments
}
