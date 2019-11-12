package ser

import (
	"github.com/flywave/go-osg/io"
)

func init() {
	wrap := io.NewObjectWrapper2("BufferData", " model.BufferData", nil, "osg::Object osg::StateAttribute osg::CullFace")
	io.GetObjectWrapperManager().AddWrap(&wrap)
}
