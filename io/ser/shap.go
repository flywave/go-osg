package ser

import "github.com/flywave/go-osg/io"

func init() {
	wrap := io.NewObjectWrapper("Shape", nil, "osg::Object osg::Shape")
	io.GetObjectWrapperManager().AddWrap(&wrap)
}
