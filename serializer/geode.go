package serializer

import (
	"github.com/flywave/go-osg"
	"github.com/flywave/go-osg/model"
)

func DrawablesChecker(obj interface{}) bool {
	g := obj.(*model.Geode)
	return len(g.Children) > 0
}
func DrawableReader(is *osg.OsgIstream, obj interface{}) {
	g := obj.(*model.Geode)
	var size int = 0
	is.Read(&size)
	is.Read(is.BEGINBRACKET)
	for i := 0; i < size; i++ {
		ob := is.ReadObject(nil)
		if model.IsBaseOfDrawable(ob) {
			g.AddChild(ob)
		}
	}
	is.Read(is.ENDBRACKET)
}

func DrawableWriter(os *osg.OsgOstream, obj interface{}) {
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
		return &gd
	}
	wrap := osg.NewObjectWrapper("Geode", fn, "osg::Object osg::Node osg::Geode")
	ser := osg.NewUserSerializer("Drawables", DrawablesChecker, DrawableReader, DrawableWriter)
	wrap.AddSerializer(&ser, osg.RWOBJECT)
	osg.GetObjectWrapperManager().AddWrap(&wrap)
}
