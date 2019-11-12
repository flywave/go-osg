package ser

import (
	"github.com/flywave/go-osg/io"
	"github.com/flywave/go-osg/model"
)

func init() {
	ser := io.NewEnumSerializer("Binding")
	ser.Add("BIND_UNDEFINED", model.BIND_UNDEFINED)
	ser.Add("BIND_OFF", model.BIND_OFF)
	ser.Add("BIND_OVERALL", model.BIND_OVERALL)
	ser.Add("BIND_PER_PRIMITIVE_SET", model.BIND_PER_PRIMITIVE_SET)
	ser.Add("BIND_PER_VERTEX", model.BIND_PER_VERTEX)

	serb1 := io.NewPropByValSerializer("Normalize", false)
	serb2 := io.NewPropByValSerializer("PreserveDataType", false)

	fn := func() interface{} {
		ay := model.NewArray()
		ser.EnumValue = &ay.Binding
		serb1.Prop = &ay.Normalize
		serb2.Prop = &ay.PreserveDataType
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
