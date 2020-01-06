package osg

import (
	"github.com/flywave/go-osg/model"
)

func checkPlaneS(tex interface{}) bool {
	return true
}

func readPlaneS(is *OsgIstream, obj interface{}) {
	pf := model.Planef{}
	is.Read(&pf)
}

func writePlaneS(os *OsgOstream, obj interface{}) {
}

func checkPlaneT(tex interface{}) bool {
	return true
}

func readPlaneT(is *OsgIstream, obj interface{}) {
	pf := model.Planef{}
	is.Read(&pf)
}

func writePlaneT(os *OsgOstream, obj interface{}) {
}

func checkPlaneR(tex interface{}) bool {
	return true
}

func readPlaneR(is *OsgIstream, obj interface{}) {
	pf := model.Planef{}
	is.Read(&pf)
}

func writePlaneR(os *OsgOstream, obj interface{}) {
}

func checkPlaneQ(tex interface{}) bool {
	return true
}

func readPlaneQ(is *OsgIstream, obj interface{}) {
	pf := model.Planef{}
	is.Read(&pf)
}

func writePlaneQ(os *OsgOstream, obj interface{}) {
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
		return tg
	}
	wrap := NewObjectWrapper("TexGen", fn, "osg::Object osg::StateAttribute osg::TexGen")
	ser1 := NewEnumSerializer("Mode", getTXMode, setTXMode)
	ser1.Add("OBJECTLINEAR", model.GLOBJECTLINEAR)
	ser1.Add("EYELINEAR", model.GLEYELINEAR)
	ser1.Add("SPHEREMAP", model.GLSPHEREMAP)
	ser1.Add("NORMALMAP", model.GLNORMALMAP)
	ser1.Add("REFLECTIONMAP", model.GLREFLECTIONMAP)
	wrap.AddSerializer(ser1, RWENUM)

	ser2 := NewUserSerializer("PlaneS", checkPlaneT, readPlaneR, writePlaneQ)
	ser3 := NewUserSerializer("PlaneT", checkPlaneT, readPlaneR, writePlaneQ)
	ser4 := NewUserSerializer("PlaneR", checkPlaneT, readPlaneR, writePlaneQ)
	ser5 := NewUserSerializer("PlaneQ", checkPlaneT, readPlaneR, writePlaneQ)
	wrap.AddSerializer(ser2, RWUSER)
	wrap.AddSerializer(ser3, RWUSER)
	wrap.AddSerializer(ser4, RWUSER)
	wrap.AddSerializer(ser5, RWUSER)

	GetObjectWrapperManager().AddWrap(wrap)
}
