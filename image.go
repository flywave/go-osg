package osg

import "github.com/flywave/go-osg/model"

func getFileName(obj interface{}) interface{} {
	return &obj.(*model.Image).FileName
}

func setFileName(obj interface{}, val interface{}) {
	obj.(*model.Image).FileName = val.(string)
}

func getWriteHint(obj interface{}) interface{} {
	return &obj.(*model.Image).WriteHint
}

func setWriteHint(obj interface{}, val interface{}) {
	obj.(*model.Image).WriteHint = val.(int32)
}

func getAllocationMode(obj interface{}) interface{} {
	return &obj.(*model.Image).AllocationMode
}

func setAllocationMode(obj interface{}, val interface{}) {
	obj.(*model.Image).AllocationMode = val.(int32)
}

func getOrigin(obj interface{}) interface{} {
	return &obj.(*model.Image).Origin
}

func setOrigin(obj interface{}, val interface{}) {
	obj.(*model.Image).Origin = val.(int32)
}

func getInternalTextureFormat(obj interface{}) interface{} {
	return &obj.(*model.Image).InternalTextureFormat
}

func setInternalTextureFormat(obj interface{}, val interface{}) {
	obj.(*model.Image).InternalTextureFormat = val.(int32)
}

func getDataType(obj interface{}) interface{} {
	return &obj.(*model.Image).DataType
}

func setDataType(obj interface{}, val interface{}) {
	obj.(*model.Image).DataType = val.(int32)
}

func getPixelFormat(obj interface{}) interface{} {
	return &obj.(*model.Image).PixelFormat
}

func setPixelFormat(obj interface{}, val interface{}) {
	obj.(*model.Image).PixelFormat = val.(int32)
}

func getRowLength(obj interface{}) interface{} {
	return &obj.(*model.Image).RowLength
}

func setRowLength(obj interface{}, val interface{}) {
	obj.(*model.Image).RowLength = val.(int32)
}

func getPacking(obj interface{}) interface{} {
	return &obj.(*model.Image).Packing
}

func setPacking(obj interface{}, val interface{}) {
	obj.(*model.Image).Packing = val.(int32)
}
func init() {
	fn := func() interface{} {
		g := model.NewImage()
		return g
	}
	wrap := NewObjectWrapper("Image", fn, "osg::Object osg::BufferData osg::Image")
	{
		uv := AddUpdateWrapperVersionProxy(wrap, 154)
		wrap.MarkSerializerAsAdded("osg::BufferData")
		uv.SetLastVersion()
	}

	uv := AddUpdateWrapperVersionProxy(wrap, 112)
	ser := NewStringSerializer("FileName", getFileName, setFileName)

	ser2 := NewEnumSerializer("WriteHint", getWriteHint, setWriteHint)
	ser2.Add("NOPREFERENCE", model.NOPREFERENCE)
	ser2.Add("STOREINLINE", model.STOREINLINE)
	ser2.Add("EXTERNALFILE", model.NOPREFERENCE)

	ser3 := NewEnumSerializer("AllocationMode", getAllocationMode, setAllocationMode)
	ser3.Add("NODELETE", model.NODELETE)
	ser3.Add("USENEWDELETE", model.USENEWDELETE)
	ser3.Add("EXTERNALFILE", model.USEMALLOCFREE)

	ser4 := NewGlenumSerializer("InternalTextureFormat", getInternalTextureFormat, setInternalTextureFormat)
	ser5 := NewGlenumSerializer("DataType", getDataType, setDataType)
	ser6 := NewGlenumSerializer("PixelFormat", getPixelFormat, setPixelFormat)
	ser7 := NewPropByValSerializer("RowLength", false, getRowLength, setRowLength)
	ser8 := NewPropByValSerializer("Packing", false, getPacking, setPacking)

	ser9 := NewEnumSerializer("Origin", getOrigin, setOrigin)
	ser9.Add("BOTTOMLEFT", model.BOTTOMLEFT)
	ser9.Add("TOPLEFT", model.TOPLEFT)
	uv.SetLastVersion()

	wrap.AddSerializer(ser, RWSTRING)
	wrap.AddSerializer(ser2, RWENUM)
	wrap.AddSerializer(ser3, RWENUM)
	wrap.AddSerializer(ser4, RWGLENUM)
	wrap.AddSerializer(ser5, RWGLENUM)
	wrap.AddSerializer(ser6, RWGLENUM)
	wrap.AddSerializer(ser7, RWINT)
	wrap.AddSerializer(ser8, RWINT)
	wrap.AddSerializer(ser9, RWENUM)
	GetObjectWrapperManager().AddWrap(wrap)
}
