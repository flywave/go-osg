package serializer

import (
	"github.com/flywave/go-osg"
	"github.com/flywave/go-osg/model"
)

func getMode(obj interface{}) interface{} {
	return obj.(*model.CullFace).Mode
}

func setMode(obj interface{}, fc interface{}) {
	obj.(*model.CullFace).Mode = *fc.(*int)
}

func init() {
	wrap := osg.NewObjectWrapper2("CullFace", " model.CullFace", nil, "osg::Object osg::BufferData")
	ser := osg.NewEnumSerializer("Mode", getMode, setMode)
	ser.Add("FRONT", model.GLFRONT)
	ser.Add("FRONT", model.GLBACK)
	ser.Add("FRONTANDBACK", model.GLFRONTANDBACK)
	wrap.AddSerializer(&ser, osg.RWENUM)

	osg.AddUpdateWrapperVersionProxy(&wrap, 147)
	osg.GetObjectWrapperManager().AddWrap(&wrap)
}
