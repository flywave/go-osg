package serializer

import "github.com/flywave/go-osg"

func init() {
	wrap := osg.NewObjectWrapper2("BufferData", " model.BufferData", nil, "osg::Object osg::StateAttribute osg::CullFace")
	osg.GetObjectWrapperManager().AddWrap(&wrap)
}
