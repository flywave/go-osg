package osg

import (
	"fmt"
	"testing"

	"github.com/flywave/go-osg/model"
)

func TestReadNode(t *testing.T) {
	rw := NewReadWrite()
	res := rw.ReadNode("test_data/osgb/Data/Tile_+003_+003/Tile_+003_+003.osgb", nil)
	obj := res.GetNode()
	if obj == nil {
		fmt.Println("....")
	}

	vst := model.NewNodeVisitor()
	obj.Accept(vst)
}
