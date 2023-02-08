// Copyright 2023 BINARY Members
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
// http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.
//

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

// CopyBytes copy bytes to a new []byte
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

// StandardizeAddr make addr standard
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
