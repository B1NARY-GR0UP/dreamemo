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

package memo

import (
	"github.com/B1NARY-GR0UP/dreamemo/common/util"
)

// ByteView data will store in ByteView form
type ByteView struct {
	B []byte
}

// Len returns the length of the data
func (v ByteView) Len() int {
	return len(v.B)
}

// ByteSlice returns the copy of the data as a byte slice
func (v ByteView) ByteSlice() []byte {
	return util.CopyBytes(v.B)
}

// String returns the data as a string
func (v ByteView) String() string {
	return string(v.B)
}

// TODO: Add more methods
