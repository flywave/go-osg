package ser

import (
	"github.com/flywave/go-osg/io"
	"github.com/flywave/go-osg/model"
)

func getNumInstances(obj interface{}) interface{} {
	return &obj.(*model.PrimitiveSet).NumInstances
}

func setNumInstances(obj interface{}, val interface{}) {
	obj.(*model.PrimitiveSet).NumInstances = val.(int)
}

func getPrimitMode(obj interface{}) interface{} {
	return &obj.(*model.PrimitiveSet).Mode
}

func setPrimitMode(obj interface{}, val interface{}) {
	obj.(*model.PrimitiveSet).Mode = val.(uint)
}

func getDAFirst(obj interface{}) interface{} {
	return &obj.(*model.DrawArrays).First
}

func setDAFirst(obj interface{}, val interface{}) {
	obj.(*model.DrawArrays).First = val.(int)
}

func getDACount(obj interface{}) interface{} {
	return &obj.(*model.DrawArrays).Count
}

func setDACount(obj interface{}, val interface{}) {
	obj.(*model.DrawArrays).Count = val.(uint64)
}

func getDLFirst(obj interface{}) interface{} {
	return &obj.(*model.DrawArrays).First
}

func setDLFirst(obj interface{}, val interface{}) {
	obj.(*model.DrawArrays).First = val.(int)
}

func getDAData(obj interface{}) interface{} {
	return &obj.(*model.DrawArrayLengths).Data
}

func setDAData(obj interface{}, val interface{}) {
	obj.(*model.DrawArrayLengths).Data = val.([]uint64)
}

func getDADataUByte(obj interface{}) interface{} {
	return &obj.(*model.DrawElementsUByte).Data
}

func setDADataUByte(obj interface{}, val interface{}) {
	obj.(*model.DrawElementsUByte).Data = val.([]uint8)
}

func getDADataUShort(obj interface{}) interface{} {
	return &obj.(*model.DrawElementsUShort).Data
}

func setDADataUShort(obj interface{}, val interface{}) {
	obj.(*model.DrawElementsUShort).Data = val.([]uint16)
}

func getDADataUInt(obj interface{}) interface{} {
	return &obj.(*model.DrawElementsUInt).Data
}

func setDADataUInt(obj interface{}, val interface{}) {
	obj.(*model.DrawElementsUInt).Data = val.([]uint32)
}

func init() {
	fn1 := func() interface{} {
		prv := model.NewPrimitiveSet()
		return &prv
	}

	fn2 := func() interface{} {
		dl := model.NewDrawArrays()
		return &dl
	}

	fn3 := func() interface{} {
		dl := model.NewDrawArrayLengths()
		return &dl
	}

	fn4 := func() interface{} {
		dl := model.NewDrawElementsUByte()
		return &dl
	}

	fn5 := func() interface{} {
		dl := model.NewDrawElementsUShort()
		return &dl
	}

	fn6 := func() interface{} {
		dl := model.NewDrawElementsUInt()
		return &dl
	}

	wrap1 := io.NewObjectWrapper("PrimitiveSet", fn1, "osg::Object osg::BufferData osg::PrimitiveSet")
	io.AddUpdateWrapperVersionProxy(&wrap1, 147)
	wrap1.MarkSerializerAsAdded("osg::BufferData")
	ser1 := io.NewPropByValSerializer("NumInstances", false, getNumInstances, setNumInstances)
	ser2 := io.NewEnumSerializer("Mode", getPrimitMode, setPrimitMode)
	ser2.Add("POINTS", model.GL_POINTS)
	ser2.Add("LINES", model.GL_LINES)
	ser2.Add("LINE_STRIP", model.GL_LINE_STRIP)
	ser2.Add("LINE_LOOP", model.GL_LINE_LOOP)
	ser2.Add("TRIANGLES", model.GL_TRIANGLES)
	ser2.Add("TRIANGLE_STRIP", model.GL_TRIANGLE_STRIP)
	ser2.Add("TRIANGLE_FAN", model.GL_TRIANGLE_FAN)
	ser2.Add("QUADS", model.GL_QUADS)
	ser2.Add("QUAD_STRIP", model.GL_QUAD_STRIP)
	ser2.Add("POLYGON", model.GL_POLYGON)
	ser2.Add("LINES_ADJACENCY", model.GL_LINES_ADJACENCY)
	ser2.Add("LINE_STRIP_ADJACENCY", model.GL_LINE_STRIP_ADJACENCY)
	ser2.Add("TRIANGLES_ADJACENCY", model.GL_TRIANGLES_ADJACENCY)
	ser2.Add("TRIANGLE_STRIP_ADJACENCY", model.GL_TRIANGLE_STRIP_ADJACENCY)
	ser2.Add("PATCHES", model.GL_PATCHES)
	wrap1.AddSerializer(&ser1, io.RW_INT)
	wrap1.AddSerializer(&ser2, io.RW_ENUM)
	io.GetObjectWrapperManager().AddWrap(&wrap1)

	wrap2 := io.NewObjectWrapper("DrawArrays", fn2, "osg::Object osg::BufferData osg::PrimitiveSet osg::DrawArrays")
	io.AddUpdateWrapperVersionProxy(&wrap2, 147)
	wrap2.MarkSerializerAsAdded("osg::BufferData")
	ser3 := io.NewPropByValSerializer("First", false, getDAFirst, setDAFirst)
	ser4 := io.NewPropByValSerializer("Count", false, getDACount, setDACount)
	wrap2.AddSerializer(&ser3, io.RW_INT)
	wrap2.AddSerializer(&ser4, io.RW_UINT)
	io.GetObjectWrapperManager().AddWrap(&wrap2)

	wrap3 := io.NewObjectWrapper2("DrawArrayLengths", "flywave::osg::drawarraylengths", fn3, "osg::Object osg::BufferData osg::PrimitiveSet osg::DrawArrays")
	io.AddUpdateWrapperVersionProxy(&wrap2, 147)
	wrap3.MarkSerializerAsAdded("osg::BufferData")
	ser5 := io.NewPropByValSerializer("First", false, getDLFirst, setDLFirst)
	ser6 := io.NewVectorSerializer("Data", io.RW_OBJECT, 4, getDAData, setDAData)
	wrap3.AddSerializer(&ser5, io.RW_INT)
	wrap3.AddSerializer(&ser6, io.RW_VECTOR)
	io.GetObjectWrapperManager().AddWrap(&wrap3)

	wrap4 := io.NewObjectWrapper2("DrawElementsUByte", "flywave::osg::drawelementsubyte", fn4, "osg::Object osg::BufferData osg::PrimitiveSet osg::DrawElementsUByte")
	io.AddUpdateWrapperVersionProxy(&wrap4, 147)
	wrap4.MarkSerializerAsAdded("osg::BufferData")
	ser7 := io.NewVectorSerializer("Data", io.RW_OBJECT, 4, getDADataUByte, setDADataUByte)
	wrap4.AddSerializer(&ser7, io.RW_UCHAR)
	io.GetObjectWrapperManager().AddWrap(&wrap4)

	wrap5 := io.NewObjectWrapper2("DrawElementsUShort", "flywave::osg::drawelementsushort", fn5, "osg::Object osg::BufferData osg::PrimitiveSet osg::DrawElementsUShort")
	io.AddUpdateWrapperVersionProxy(&wrap5, 147)
	wrap5.MarkSerializerAsAdded("osg::BufferData")
	ser8 := io.NewVectorSerializer("Data", io.RW_OBJECT, 4, getDADataUShort, setDADataUShort)
	wrap5.AddSerializer(&ser8, io.RW_USHORT)
	io.GetObjectWrapperManager().AddWrap(&wrap5)

	wrap6 := io.NewObjectWrapper2("DrawElementsUInt", "flywave::osg::drawelementsuint", fn6, "osg::Object osg::BufferData osg::PrimitiveSet osg::DrawElementsUInt")
	io.AddUpdateWrapperVersionProxy(&wrap6, 147)
	wrap6.MarkSerializerAsAdded("osg::BufferData")
	ser9 := io.NewVectorSerializer("Data", io.RW_OBJECT, 4, getDADataUInt, setDADataUInt)
	wrap5.AddSerializer(&ser9, io.RW_UINT)
	io.GetObjectWrapperManager().AddWrap(&wrap6)
}
