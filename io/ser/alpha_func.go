package ser

import (
	"github.com/flywave/go-osg/io"
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
	ser := io.NewEnumSerializer("Function", getComparisonFunc, setComparisonFunc)
	ser.Add("NEVER", model.GL_NEVER)
	ser.Add("LESS", model.GL_LESS)
	ser.Add("EQUAL", model.GL_EQUAL)
	ser.Add("LEQUAL", model.GL_LEQUAL)
	ser.Add("GREATER", model.GL_GREATER)
	ser.Add("NOTEQUAL", model.GL_NOTEQUAL)
	ser.Add("GEQUAL", model.GL_GEQUAL)
	ser.Add("GEQUAL", model.GL_ALWAYS)

	serf := io.NewPropByValSerializer("ReferenceValue", false, getReferenceValue, setReferenceValue)

	fn := func() interface{} {
		al := model.NewAlphaFunc()
		return &al
	}
	wrap := io.NewObjectWrapper("AlphaFunc", fn, "osg::Object osg::StateAttribute osg::AlphaFunc")
	wrap.AddSerializer(&ser, io.RW_ENUM)
	wrap.AddSerializer(&serf, io.RW_FLOAT)
	io.GetObjectWrapperManager().AddWrap(&wrap)
}
