package ser

import (
	"github.com/flywave/go-osg/io"
	"github.com/flywave/go-osg/model"
)

func DrawablesChecker(obj interface{}) bool {
	return true
}
func DrawableReader(is *io.OsgIstream, obj interface{}) {}

func DrawableWriter(os *io.OsgOstream, obj interface{}) {}

func init() {
	fn := func() interface{} {
		gd := model.NewGeode()
		return &gd
	}
	wrap := io.NewObjectWrapper("Geode", fn, "osg::Object osg::Node osg::Geode")
	ser := io.NewUserSerializer("Drawables", DrawablesChecker, DrawableReader, DrawableWriter)
	wrap.AddSerializer(&ser, io.RW_OBJECT)
	io.GetObjectWrapperManager().AddWrap(&wrap)
}
