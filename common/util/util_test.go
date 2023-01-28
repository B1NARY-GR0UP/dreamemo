package util

import (
	"fmt"
	"testing"
)

func TestParseRequestURL(t *testing.T) {
	segments := ParseRequestURL("/_dreamemo/hello/dreamemo/", "/_dreamemo/")
	fmt.Println(segments)
}
