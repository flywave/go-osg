package osg

import (
	"github.com/flywave/go-osg/model"
)

type hasReferenceFrame interface {
	GetReferenceFrame() int32
	SetReferenceFrame(int32)
}

func getReferenceFrame(obj interface{}) interface{} {
	switch t := obj.(type) {
	case *model.Transform:
		return &t.ReferenceFrame
	case *model.PositionAttitudeTransform:
		return &t.ReferenceFrame
	case *model.MatrixTransform:
		return &t.ReferenceFrame
	default:
		return nil
	}
}

func setReferenceFrame(obj interface{}, val interface{}) {
	v := int(val.(int32))
	switch t := obj.(type) {
	case *model.Transform:
		t.ReferenceFrame = v
	case *model.PositionAttitudeTransform:
		t.ReferenceFrame = v
	case *model.MatrixTransform:
		t.ReferenceFrame = v
	}
}

func init() {
	fn := func() interface{} {
		td := model.NewTransform()
		return td
	}
	wrap := NewObjectWrapper("Transform", fn, "osg::Object osg::Node osg::Group osg::Transform")
	ser1 := NewEnumSerializer("ReferenceFrame", getReferenceFrame, setReferenceFrame)
	ser1.Add("RELATIVERF", model.RELATIVERF)
	ser1.Add("ABSOLUTERF", model.ABSOLUTERF)
	ser1.Add("ABSOLUTERFINHERITVIEWPOINT", model.ABSOLUTERFINHERITVIEWPOINT)
	wrap.AddSerializer(ser1, RWENUM)
	GetObjectWrapperManager().AddWrap(wrap)
}
