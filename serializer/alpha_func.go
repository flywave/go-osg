package serializer

import (
	"github.com/flywave/go-osg"
	"github.com/flywave/go-osg/model"
)

func getComparisonFunc(obj interface{}) interface{} {
	return obj.(*model.AlphaFunc).ComparisonFunc
}

func setComparisonFunc(obj interface{}, fc interface{}) {
	obj.(*model.AlphaFunc).ComparisonFunc = *fc.(*int)
}

func getReferenceValue(obj interface{}) interface{} {
	return obj.(*model.AlphaFunc).ReferenceValue
}

func setReferenceValue(obj interface{}, fc interface{}) {
	obj.(*model.AlphaFunc).ReferenceValue = *fc.(*float32)
}

func init() {
	ser := osg.NewEnumSerializer("Function", getComparisonFunc, setComparisonFunc)
	ser.Add("NEVER", model.GLNEVER)
	ser.Add("LESS", model.GLLESS)
	ser.Add("EQUAL", model.GLEQUAL)
	ser.Add("LEQUAL", model.GLLEQUAL)
	ser.Add("GREATER", model.GLGREATER)
	ser.Add("NOTEQUAL", model.GLNOTEQUAL)
	ser.Add("GEQUAL", model.GLGEQUAL)
	ser.Add("GEQUAL", model.GLALWAYS)

	serf := osg.NewPropByValSerializer("ReferenceValue", false, getReferenceValue, setReferenceValue)

	fn := func() interface{} {
		al := model.NewAlphaFunc()
		return &al
	}
	wrap := osg.NewObjectWrapper("AlphaFunc", fn, "osg::Object osg::StateAttribute osg::AlphaFunc")
	wrap.AddSerializer(&ser, osg.RWENUM)
	wrap.AddSerializer(&serf, osg.RWFLOAT)
	osg.GetObjectWrapperManager().AddWrap(&wrap)
}
