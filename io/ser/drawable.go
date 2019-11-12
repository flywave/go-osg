package ser

import "github.com/flywave/go-osg/io"

func DrawableChecker(obj interface{}) bool {
	return true
}
func DrawableReader(is *io.OsgIstream, obj interface{}) {}
func DrawableWriter(os *io.OsgOstream, obj interface{}) {}

func init() {
	wrap := io.NewObjectWrapper("Drawable", nil, "osg::Object osg::Node osg::Drawable")
	io.AddUpdateWrapperVersionProxy(&wrap, 154)
	wrap.MarkSerializerAsAdded("osg::Node")

	ser := io.NewObjectSerializer("StateSet")

	ser_user := io.NewUserSerializer("InitialBound", DrawableChecker, DrawableReader, DrawableWriter)

}
