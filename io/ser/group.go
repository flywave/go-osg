package ser

import (
	"github.com/flywave/go-osg/io"
	"github.com/flywave/go-osg/model"
)

func ChildrenChecker(obj interface{}) bool {
	return true
}
func ChildrenReader(is *io.OsgIstream, obj interface{}) {}

func ChildrenWriter(os *io.OsgOstream, obj interface{}) {}

func init() {
	fn := func() interface{} {
		g := model.NewGroup()
		return &g
	}
	wrap := io.NewObjectWrapper("Group", fn, "osg::Object osg::Node osg::Group")
	ser := io.NewUserSerializer("Children", ChildrenChecker, ChildrenReader, ChildrenWriter)
	wrap.AddSerializer(&ser, io.RW_OBJECT)
	io.GetObjectWrapperManager().AddWrap(&wrap)
}
