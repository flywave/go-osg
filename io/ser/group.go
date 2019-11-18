package ser

import (
	"github.com/flywave/go-osg/io"
	"github.com/flywave/go-osg/model"
)

func ChildrenChecker(obj interface{}) bool {
	g := obj.(*model.Group)
	return len(g.Children) > 0
}
func ChildrenReader(is *io.OsgIstream, obj interface{}) {
	g := obj.(*model.Group)
	var size int = 0
	is.Read(&size)
	is.Read(is.BEGIN_BRACKET)
	for i := 0; i < size; i++ {
		ob := is.ReadObject(nil)
		if ob != nil {
			g.AddChild(ob)
		}
	}
	is.Read(is.END_BRACKET)
}

func ChildrenWriter(os *io.OsgOstream, obj interface{}) {
	g := obj.(*model.Group)
	size := len(g.Children)
	os.Write(size)
	os.Write(os.BEGIN_BRACKET)
	for i := 0; i < size; i++ {
		os.Write(g.Children[i])
	}
	os.Write(os.END_BRACKET)
	os.Write(os.CRLF)
}

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
