package osg

import (
	"github.com/flywave/go-osg/model"
)

func DrawablesChecker(obj interface{}) bool {
	g := obj.(*model.Geode)
	return len(g.Children) > 0
}
func DrawableReader(is *OsgIstream, obj interface{}) {
	g := obj.(*model.Geode)
	size := is.ReadSize()
	is.Read(is.BEGINBRACKET)
	for i := 0; i < int(size); i++ {
		ob := is.ReadObject(nil)
		_, ok := ob.(model.DrawableInterface)
		if ok {
			child := ob.(model.NodeInterface)
			g.AddChild(child)
		}
	}
	is.Read(is.ENDBRACKET)
}

func DrawableWriter(os *OsgOstream, obj interface{}) {
	g := obj.(*model.Geode)
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
		gd := model.NewGeode()
		return gd
	}
	wrap := NewObjectWrapper("Geode", fn, "osg::Object osg::Node osg::Geode")
	ser := NewUserSerializer("Drawables", DrawablesChecker, DrawableReader, DrawableWriter)
	wrap.AddSerializer(ser, RWOBJECT)
	GetObjectWrapperManager().AddWrap(wrap)
}
