package osg

import (
	"github.com/flywave/go-osg/model"
)

func getMode(obj interface{}) interface{} {
	return &obj.(*model.CullFace).Mode
}

func setMode(obj interface{}, fc interface{}) {
	obj.(*model.CullFace).Mode = int(fc.(int32))
}

func init() {
	wrap := NewObjectWrapper("CullFace", nil, "osg::Object osg::StateAttribute osg::CullFace")
	ser := NewEnumSerializer("Mode", getMode, setMode)
	ser.Add("FRONT", model.GLFRONT)
	ser.Add("BACK", model.GLBACK)
	ser.Add("FRONTANDBACK", model.GLFRONTANDBACK)
	wrap.AddSerializer(ser, RWENUM)
	GetObjectWrapperManager().AddWrap(wrap)
}
