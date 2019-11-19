package serializer

import (
	"github.com/flywave/go-osg"
	"github.com/flywave/go-osg/model"
)

func ChildrenChecker(obj interface{}) bool {
	g := obj.(*model.Group)
	return len(g.Children) > 0
}
func ChildrenReader(is *osg.OsgIstream, obj interface{}) {
	g := obj.(*model.Group)
	var size int = 0
	is.Read(&size)
	is.Read(is.BEGINBRACKET)
	for i := 0; i < size; i++ {
		ob := is.ReadObject(nil)
		if ob != nil {
			g.AddChild(ob)
		}
	}
	is.Read(is.ENDBRACKET)
}

func ChildrenWriter(os *osg.OsgOstream, obj interface{}) {
	g := obj.(*model.Group)
	size := len(g.Children)
	os.Write(size)
	os.Write(os.BEGINBRACKET)
	for i := 0; i < size; i++ {
		os.Write(g.Children[i])
	}
	os.Write(os.ENDBRACKET)
	os.Write(os.CRLF)
}

func init() {
	fn := func() interface{} {
		g := model.NewGroup()
		return &g
	}
	wrap := osg.NewObjectWrapper("Group", fn, "osg::Object osg::Node osg::Group")
	ser := osg.NewUserSerializer("Children", ChildrenChecker, ChildrenReader, ChildrenWriter)
	wrap.AddSerializer(&ser, osg.RWOBJECT)
	osg.GetObjectWrapperManager().AddWrap(&wrap)
}
