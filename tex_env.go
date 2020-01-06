package osg

import (
	"github.com/flywave/go-osg/model"
)

func getEnvMode(obj interface{}) interface{} {
	tx := obj.(*model.TexEnv)
	return &tx.Mode
}

func setEnvMode(obj interface{}, val interface{}) {
	tx := obj.(*model.TexEnv)
	tx.Mode = val.(int32)
}

func getEnvColor(obj interface{}) interface{} {
	tx := obj.(*model.TexEnv)
	return &tx.Color
}

func setEnvColor(obj interface{}, val interface{}) {
	tx := obj.(*model.TexEnv)
	tx.Color = val.([4]float32)
}

func init() {
	fn := func() interface{} {
		tg := model.NewTexEnv()
		return tg
	}
	wrap := NewObjectWrapper("TexEnv", fn, "osg::Object osg::StateAttribute osg::TexEnv")
	ser1 := NewEnumSerializer("Mode", getEnvMode, setEnvMode)
	ser1.Add("DECAL", model.GLDECAL)
	ser1.Add("MODULATE", model.GLMODULATE)
	ser1.Add("BLEND", model.GLBLEND)
	ser1.Add("REPLACE", model.GLREPLACE)
	ser1.Add("ADD", model.GLADD)
	wrap.AddSerializer(ser1, RWENUM)

	ser2 := NewPropByRefSerializer("Color", getEnvColor, setEnvColor)
	wrap.AddSerializer(ser2, RWVEC4F)

}
