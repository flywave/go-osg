package osg

import (
	"testing"
)

func TestReadNode(t *testing.T) {
	rw := NewReadWrite()
	rw.ReadNode("/Volumes/Projection/flywave/src/tests/flywave/osg/skydome.osgt", nil)
}
