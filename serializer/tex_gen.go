package serializer

import (
	"github.com/flywave/go-osg"
	"github.com/flywave/go-osg/model"
)

func checkPlaneS(tex interface{}) bool {
	return true
}

func readPlaneS(is *osg.OsgIstream, obj interface{}) {
	pf := model.Planef{}
	is.Read(&pf)
}

func writePlaneS(os *osg.OsgOstream, obj interface{}) {
}

func checkPlaneT(tex interface{}) bool {
	return true
}

func readPlaneT(is *osg.OsgIstream, obj interface{}) {
	pf := model.Planef{}
	is.Read(&pf)
}

func writePlaneT(os *osg.OsgOstream, obj interface{}) {
}

func checkPlaneR(tex interface{}) bool {
	return true
}

func readPlaneR(is *osg.OsgIstream, obj interface{}) {
	pf := model.Planef{}
	is.Read(&pf)
}

func writePlaneR(os *osg.OsgOstream, obj interface{}) {
}

func checkPlaneQ(tex interface{}) bool {
	return true
}

func readPlaneQ(is *osg.OsgIstream, obj interface{}) {
	pf := model.Planef{}
	is.Read(&pf)
}

func writePlaneQ(os *osg.OsgOstream, obj interface{}) {
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
	wrap := osg.NewObjectWrapper2("TexGen", "flywave::osg::texgen", fn, "osg::Object osg::StateAttribute osg::TexGen")
	ser1 := osg.NewEnumSerializer("Mode", getTXMode, setTXMode)
	ser1.Add("OBJECTLINEAR", model.GLOBJECTLINEAR)
	ser1.Add("EYELINEAR", model.GLEYELINEAR)
	ser1.Add("SPHEREMAP", model.GLSPHEREMAP)
	ser1.Add("NORMALMAP", model.GLNORMALMAP)
	ser1.Add("REFLECTIONMAP", model.GLREFLECTIONMAP)
	wrap.AddSerializer(&ser1, osg.RWENUM)

	ser2 := osg.NewUserSerializer("PlaneS", checkPlaneT, readPlaneR, writePlaneQ)
	ser3 := osg.NewUserSerializer("PlaneT", checkPlaneT, readPlaneR, writePlaneQ)
	ser4 := osg.NewUserSerializer("PlaneR", checkPlaneT, readPlaneR, writePlaneQ)
	ser5 := osg.NewUserSerializer("PlaneQ", checkPlaneT, readPlaneR, writePlaneQ)
	wrap.AddSerializer(&ser2, osg.RWUSER)
	wrap.AddSerializer(&ser3, osg.RWUSER)
	wrap.AddSerializer(&ser4, osg.RWUSER)
	wrap.AddSerializer(&ser5, osg.RWUSER)

	osg.GetObjectWrapperManager().AddWrap(&wrap)
}
