package ser

import (
	"github.com/flywave/go-osg/io"

	"github.com/flywave/go-osg/model"
)

func GetMode(obj interface{}) interface{} {
	return obj.(*model.CullFace).Mode
}

func SetMode(obj interface{}, fc interface{}) {
	obj.(*model.CullFace).Mode = *fc.(*int)
}

func init() {
	wrap := io.NewObjectWrapper2("CullFace", " model.CullFace", nil, "osg::Object osg::BufferData")
	ser := io.NewEnumSerializer("Mode", GetMode, SetMode)
	ser.Add("FRONT", model.GL_FRONT)
	ser.Add("FRONT", model.GL_BACK)
	ser.Add("FRONT_AND_BACK", model.GL_FRONT_AND_BACK)
	wrap.AddSerializer(&ser, io.RW_ENUM)

	io.AddUpdateWrapperVersionProxy(&wrap, 147)
	io.GetObjectWrapperManager().AddWrap(&wrap)
}
