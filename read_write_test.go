package osg

import (
	"testing"
)

func TestHello(t *testing.T) {
	rw := NewReadWrite()
	rw.ReadNode("test_data/simpleroom.osgt", nil)
}
