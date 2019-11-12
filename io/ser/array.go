package ser

import (
	"github.com/flywave/go-osg/io"
	"github.com/flywave/go-osg/model"
)

func GetBinding(obj interface{}) interface{} {
	return obj.(*model.Array).Binding
}

func SetBinding(obj interface{}, fc interface{}) {
	obj.(*model.Array).Binding = *fc.(*int)
}

func GetNormalize(obj interface{}) interface{} {
	return obj.(*model.Array).Normalize
}

func SetNormalize(obj interface{}, fc interface{}) {
	obj.(*model.Array).Normalize = *fc.(*bool)
}

func GetPreserveDataType(obj interface{}) interface{} {
	return obj.(*model.Array).PreserveDataType
}

func SetPreserveDataType(obj interface{}, fc interface{}) {
	obj.(*model.Array).PreserveDataType = *fc.(*bool)
}

func init() {
	ser := io.NewEnumSerializer("Binding", GetBinding, SetBinding)
	ser.Add("BIND_UNDEFINED", model.BIND_UNDEFINED)
	ser.Add("BIND_OFF", model.BIND_OFF)
	ser.Add("BIND_OVERALL", model.BIND_OVERALL)
	ser.Add("BIND_PER_PRIMITIVE_SET", model.BIND_PER_PRIMITIVE_SET)
	ser.Add("BIND_PER_VERTEX", model.BIND_PER_VERTEX)

	serb1 := io.NewPropByValSerializer("Normalize", false, GetNormalize, SetNormalize)
	serb2 := io.NewPropByValSerializer("PreserveDataType", false, GetPreserveDataType, SetPreserveDataType)

	fn := func() interface{} {
		ay := model.NewArray()
		return &ay
	}
	wrap := io.NewObjectWrapper("Array", fn, "osg::Object osg::BufferData osg::Array")
	wrap.MarkSerializerAsAdded("osg::BufferData")
	wrap.AddSerializer(&ser, io.RW_ENUM)
	wrap.AddSerializer(&serb1, io.RW_BOOL)
	wrap.AddSerializer(&serb2, io.RW_BOOL)
	io.AddUpdateWrapperVersionProxy(&wrap, 147)
	io.GetObjectWrapperManager().AddWrap(&wrap)
}
