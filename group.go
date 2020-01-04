package osg

import (
	"github.com/flywave/go-osg/model"
)

func ChildrenChecker(obj interface{}) bool {
	g := obj.(*model.Group)
	return len(g.Children) > 0
}
func ChildrenReader(is *OsgIstream, obj interface{}) {
	g := obj.(model.GroupInterface)
	size := is.ReadSize()
	is.Read(is.BEGINBRACKET)
	for i := 0; i < int(size); i++ {
		ob := is.ReadObject(nil)
		if ob != nil {
			g.AddChild(ob.(model.NodeInterface))
		}
	}
	is.Read(is.ENDBRACKET)
}

func ChildrenWriter(os *OsgOstream, obj interface{}) {
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
		return g
	}
	wrap := NewObjectWrapper("Group", fn, "osg::Object osg::Node osg::Group")
	ser := NewUserSerializer("Children", ChildrenChecker, ChildrenReader, ChildrenWriter)
	wrap.AddSerializer(ser, RWOBJECT)
	GetObjectWrapperManager().AddWrap(wrap)
}
