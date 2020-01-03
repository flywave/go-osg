package osg

import (
	"fmt"
	"testing"
)

func TestReadNode(t *testing.T) {
	rw := NewReadWrite()
	res := rw.ReadNode("test_data/Tile_+003_+003_L18_000.osgb", nil)
	obj := res.GetNode()
	if obj == nil {
		fmt.Println("....")
	}

}
