package ser

import (
	"github.com/flywave/go-osg/io"
	"github.com/flywave/go-osg/model"
)

func DrawablesChecker(obj interface{}) bool {
	g := obj.(*model.Geode)
	return len(g.Children) > 0
}
func DrawableReader(is *io.OsgIstream, obj interface{}) {
	g := obj.(*model.Geode)
	var size int = 0
	is.Read(&size)
	is.Read(is.BEGIN_BRACKET)
	for i := 0; i < size; i++ {
		obj = is.ReadObject()
		if model.IsBaseOfDrawable(obj) {
			g.AddChild(obj)
		}
	}
	is.Read(is.END_BRACKET)
}

func DrawableWriter(os *io.OsgOstream, obj interface{}) {
	g := obj.(*model.Geode)
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
		gd := model.NewGeode()
		return &gd
	}
	wrap := io.NewObjectWrapper("Geode", fn, "osg::Object osg::Node osg::Geode")
	ser := io.NewUserSerializer("Drawables", DrawablesChecker, DrawableReader, DrawableWriter)
	wrap.AddSerializer(&ser, io.RW_OBJECT)
	io.GetObjectWrapperManager().AddWrap(&wrap)
}
