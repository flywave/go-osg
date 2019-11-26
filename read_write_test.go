package osg

import (
	"testing"
)

func TestReadNode(t *testing.T) {
	rw := NewReadWrite()
	rw.ReadNode("test_data/skydome.osgt", nil)
}
