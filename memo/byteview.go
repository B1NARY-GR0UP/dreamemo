package memo

import "github.com/B1NARY-GR0UP/dreamemo/util"

// ByteView data will store in ByteView form
type ByteView struct {
	B []byte
	// TODO: decide Add s string field
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
