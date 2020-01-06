package osg

import (
	"github.com/flywave/go-osg/model"
)

func getShadeMode(obj interface{}) interface{} {
	return &obj.(*model.ShadeModel).Mode
}

func setShadeMode(obj interface{}, val interface{}) {
	obj.(*model.ShadeModel).Mode = val.(int32)
}

func init() {
	fn := func() interface{} {
		sm := model.NewShadeModel()
		return sm
	}
	wrap := NewObjectWrapper("ShadeModel", fn, "osg::Object osg::StateAttribute osg::ShadeModel")
	ser := NewEnumSerializer("Mode", getShadeMode, setShadeMode)
	ser.Add("FLAT", model.FLAT)
	ser.Add("SMOOTH", model.SMOOTH)
	wrap.AddSerializer(ser, RWENUM)
	GetObjectWrapperManager().AddWrap(wrap)
}
