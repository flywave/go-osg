package osg

import (
	"fmt"
	"testing"
)

func TestReadNode(t *testing.T) {
	rw := NewReadWrite()
	res := rw.ReadNode("test_data/skydome.osgt", nil)
	obj := res.GetNode()
	if obj == nil {
		fmt.Println("....")
	}
}
