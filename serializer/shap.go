package serializer

import "github.com/flywave/go-osg"

func init() {
	wrap := osg.NewObjectWrapper("Shape", nil, "osg::Object osg::Shape")
	osg.GetObjectWrapperManager().AddWrap(&wrap)
}
