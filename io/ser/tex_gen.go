package ser

import (
	"github.com/flywave/go-osg/io"
	"github.com/flywave/go-osg/model"
)

func checkPlaneS(tex interface{}) bool {
	return true
}

func readPlaneS(is *io.OsgIstream, obj interface{}) {
	pf := model.Planef{}
	is.Read(&pf)
}

func writePlaneS(os *io.OsgOstream, obj interface{}) {
}

func checkPlaneT(tex interface{}) bool {
	return true
}

func readPlaneT(is *io.OsgIstream, obj interface{}) {
	pf := model.Planef{}
	is.Read(&pf)
}

func writePlaneT(os *io.OsgOstream, obj interface{}) {
}

func checkPlaneR(tex interface{}) bool {
	return true
}

func readPlaneR(is *io.OsgIstream, obj interface{}) {
	pf := model.Planef{}
	is.Read(&pf)
}

func writePlaneR(os *io.OsgOstream, obj interface{}) {
}

func checkPlaneQ(tex interface{}) bool {
	return true
}

func readPlaneQ(is *io.OsgIstream, obj interface{}) {
	pf := model.Planef{}
	is.Read(&pf)
}

func writePlaneQ(os *io.OsgOstream, obj interface{}) {
}

func getTXMode(obj interface{}) interface{} {
	tx := obj.(*model.TexGen)
	return &tx.Mode
}

func setTXMode(obj interface{}, val interface{}) {
	tx := obj.(*model.TexGen)
	tx.Mode = val.(int)
}

func init() {
	fn := func() interface{} {
		tg := model.NewTexGen()
		return &tg
	}
	wrap := io.NewObjectWrapper2("TexGen", "flywave::osg::texgen", fn, "osg::Object osg::StateAttribute osg::TexGen")
	ser1 := io.NewEnumSerializer("Mode", getTXMode, setTXMode)
	ser1.Add("OBJECT_LINEAR", model.GL_OBJECT_LINEAR)
	ser1.Add("EYE_LINEAR", model.GL_EYE_LINEAR)
	ser1.Add("SPHERE_MAP", model.GL_SPHERE_MAP)
	ser1.Add("NORMAL_MAP", model.GL_NORMAL_MAP)
	ser1.Add("REFLECTION_MAP", model.GL_REFLECTION_MAP)
	wrap.AddSerializer(&ser1, io.RW_ENUM)

	ser2 := io.NewUserSerializer("PlaneS", checkPlaneT, readPlaneR, writePlaneQ)
	ser3 := io.NewUserSerializer("PlaneT", checkPlaneT, readPlaneR, writePlaneQ)
	ser4 := io.NewUserSerializer("PlaneR", checkPlaneT, readPlaneR, writePlaneQ)
	ser5 := io.NewUserSerializer("PlaneQ", checkPlaneT, readPlaneR, writePlaneQ)
	wrap.AddSerializer(&ser2, io.RW_USER)
	wrap.AddSerializer(&ser3, io.RW_USER)
	wrap.AddSerializer(&ser4, io.RW_USER)
	wrap.AddSerializer(&ser5, io.RW_USER)

	io.GetObjectWrapperManager().AddWrap(&wrap)
}
