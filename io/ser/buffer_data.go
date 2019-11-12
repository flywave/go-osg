package ser

import (
	"github.com/flywave/go-osg/io"
	"github.com/flywave/go-osg/model"
)

func init() {
	wrap := io.NewObjectWrapper2("BufferData", " model.BufferData", nil, "osg::Object osg::StateAttribute osg::CullFace")
	ser := io.NewEnumSerializer("Mode")
	ser.Add("FRONT", model.GL_FRONT)
	ser.Add("FRONT", model.GL_BACK)
	ser.Add("FRONT_AND_BACK", model.GL_FRONT_AND_BACK)
	wrap.AddSerializer(&ser, io.RW_ENUM)
	io.GetObjectWrapperManager().AddWrap(&wrap)
}
