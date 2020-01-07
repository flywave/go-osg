package osg

import (
	"bytes"
	"compress/zlib"
	"fmt"
	"io"
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

func TestCompress(t *testing.T) {
	buf := []byte{1, 2, 3, 4, 5, 6, 7, 7, 8, 8, 2, 3, 4, 5, 6, 7, 7, 8, 8, 2, 3, 4, 5, 6, 7, 7, 8, 8, 2, 3, 4, 5, 6, 7, 7, 8, 8, 2, 3, 4, 5, 6, 7, 7, 8, 8, 2, 3, 4, 5, 6, 7, 7, 8, 8, 2, 3, 4, 5, 6, 7, 7, 8, 8, 2, 3, 4, 5, 6, 7, 7, 8, 8, 2, 3, 4, 5, 6, 7, 7, 8, 8, 2, 3, 4, 5, 6, 7, 7, 8, 8, 2, 3, 4, 5, 6, 7, 7, 8, 8, 2, 3, 4, 5, 6, 7, 7, 8, 8}
	var bt []byte
	bf := bytes.NewBuffer(bt)
	w := zlib.NewWriter(bf)
	w.Write(buf)
	w.Close()
	fmt.Println(bf.Bytes())

	r, _ := zlib.NewReader(bf)
	var src1 []byte
	for {
		buf1 := make([]byte, 4096)
		n, _ := io.ReadFull(r, buf1)
		if n == 0 {
			break
		}
		src1 = append(src1, buf1[0:n]...)
	}
	fmt.Println(src1)
}
