package ser

import "github.com/flywave/go-osg/io"

func init() {
	wrap := io.NewObjectWrapper2("CullFace", " model.CullFace", nil, "osg::Object osg::BufferData")
	io.AddUpdateWrapperVersionProxy(&wrap, 147)
	io.GetObjectWrapperManager().AddWrap(&wrap)
}
