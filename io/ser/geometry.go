package ser

import (
	"github.com/flywave/go-osg/io"
	"github.com/flywave/go-osg/model"
)

var lookup *io.IntLookup

func readAttributeBinding(is *io.OsgIstream) int32 {
	if is.IsBinary() {
		var val int32
		is.Read(&val)
		return val
	} else {
		str := is.ReadString()
		return lookup.GetValue(str)
	}
}

func writeAttributeBinding(os *io.OsgOstream, val int32) {
	if os.IsBinary() {
		os.Write(val)
	} else {
		os.Write(lookup.GetString(val))
	}
}

func readArray(is *io.OsgIstream) *model.Array {
	is.PROPERTY.Name = "Array"
	hasArray := false
	var ary *model.Array
	is.Read(is.PROPERTY)
	is.Read(&hasArray)
	if hasArray {
		*ary = model.NewArray()
		is.Read(&ary)
	}

	is.PROPERTY.Name = "Indices"
	hasIndices := false
	is.Read(is.PROPERTY)
	is.Read(&hasIndices)
	if hasIndices {
		ary2 := model.NewArray()
		is.Read(&ary)
		if ary != nil {
			ary.Udc.UserData = &ary2
		}
	}

	is.PROPERTY.Name = "Binding"
	is.Read(is.PROPERTY)
	bd := readAttributeBinding(is)
	if ary != nil {
		ary.Binding = bd
	}

	is.PROPERTY.Name = "Normalize"
	normalizeValue := -1
	is.Read(is.PROPERTY)
	is.Read(&normalizeValue)
	if ary != nil {
		ary.Normalize = normalizeValue != 0
	}
	return ary
}

func writeArray(os *io.OsgOstream, ary *model.Array) {
	os.PROPERTY.Name = "Array"
	os.Write(os.PROPERTY)
	os.Write(ary != nil)
	var ary2 *model.Array
	if ary != nil {
		os.Write(ary)
		ary2 = ary.Udc.UserData.(*model.Array)
	} else {
		os.Write(os.CRLF)
	}

	os.PROPERTY.Name = "Indices"
	os.Write(os.PROPERTY)
	os.Write(ary2 != nil)
	if ary2 != nil {
		os.Write(ary2)
	} else {
		os.Write(os.CRLF)
	}
	os.PROPERTY.Name = "Binding"
	os.Write(os.PROPERTY)
	writeAttributeBinding(os, ary.Binding)
}

func checkVertexData(geom interface{}) bool {
	return geom.(*model.Geometry).VertexArray != nil
}
func readVertexData(is *io.OsgIstream, g interface{}) {
	geom := g.(*model.Geometry)
	is.Read(is.BEGIN_BRACKET)
	is.Read(is.CRLF)
	ary := readArray(is)
	geom.VertexArray = ary
	is.Read(is.END_BRACKET)
	is.Read(is.CRLF)
}
func writeVertexData(os *io.OsgOstream, g interface{}) {
	geom := g.(*model.Geometry)
	os.Write(os.BEGIN_BRACKET)
	os.Write(os.CRLF)
	writeArray(os, geom.VertexArray)
	os.Write(os.END_BRACKET)
	os.Write(os.CRLF)
}

func checkNormalData(g interface{}) bool {
	geom := g.(*model.Geometry)
	return geom.NormalArray != nil
}
func readNormalData(is *io.OsgIstream, g interface{}) {
	geom := g.(*model.Geometry)
	is.Read(is.BEGIN_BRACKET)
	is.Read(is.CRLF)
	ary := readArray(is)
	geom.NormalArray = ary
	is.Read(is.END_BRACKET)
	is.Read(is.CRLF)
}
func writeNormalData(os *io.OsgOstream, g interface{}) {
	geom := g.(*model.Geometry)
	os.Write(os.BEGIN_BRACKET)
	os.Write(os.CRLF)
	writeArray(os, geom.NormalArray)
	os.Write(os.END_BRACKET)
	os.Write(os.CRLF)
}

func checkColorData(g interface{}) bool {
	geom := g.(*model.Geometry)
	return geom.ColorArray != nil
}
func readColorData(is *io.OsgIstream, g interface{}) {
	geom := g.(*model.Geometry)
	is.Read(is.BEGIN_BRACKET)
	is.Read(is.CRLF)
	ary := readArray(is)
	geom.ColorArray = ary
	is.Read(is.END_BRACKET)
	is.Read(is.CRLF)
}
func writeColorData(os *io.OsgOstream, g interface{}) {
	geom := g.(*model.Geometry)
	os.Write(os.BEGIN_BRACKET)
	os.Write(os.CRLF)
	writeArray(os, geom.ColorArray)
	os.Write(os.END_BRACKET)
	os.Write(os.CRLF)
}

func checkSecondaryColorData(g interface{}) bool {
	geom := g.(*model.Geometry)
	return geom.SecondaryColorArray != nil
}
func readSecondaryColorData(is *io.OsgIstream, g interface{}) {
	geom := g.(*model.Geometry)
	is.Read(is.BEGIN_BRACKET)
	is.Read(is.CRLF)
	ary := readArray(is)
	geom.SecondaryColorArray = ary
	is.Read(is.END_BRACKET)
	is.Read(is.CRLF)
}
func writeSecondaryColorData(os *io.OsgOstream, g interface{}) {
	geom := g.(*model.Geometry)
	os.Write(os.BEGIN_BRACKET)
	os.Write(os.CRLF)
	writeArray(os, geom.SecondaryColorArray)
	os.Write(os.END_BRACKET)
	os.Write(os.CRLF)
}

func checkFogCoordData(g interface{}) bool {
	geom := g.(*model.Geometry)
	return geom.FogCoordArray != nil
}
func readFogCoordData(is *io.OsgIstream, g interface{}) {
	geom := g.(*model.Geometry)
	is.Read(is.BEGIN_BRACKET)
	is.Read(is.CRLF)
	ary := readArray(is)
	geom.FogCoordArray = ary
	is.Read(is.END_BRACKET)
	is.Read(is.CRLF)
}
func writeFogCoordData(os *io.OsgOstream, g interface{}) {
	geom := g.(*model.Geometry)
	os.Write(os.BEGIN_BRACKET)
	os.Write(os.CRLF)
	writeArray(os, geom.FogCoordArray)
	os.Write(os.END_BRACKET)
	os.Write(os.CRLF)
}

func checkTexCoordData(g interface{}) bool {
	geom := g.(*model.Geometry)
	return geom.TexCoordArrayList != nil
}

func readTexCoordData(is *io.OsgIstream, g interface{}) {
	geom := g.(*model.Geometry)
	var size int = 0
	is.Read(&size)
	is.Read(is.BEGIN_BRACKET)
	index := 0
	is.PROPERTY.Name = "Data"
	geom.TexCoordArrayList = make([]*model.Array, size, size)
	for {
		if index == size {
			break
		}
		size--
		is.Read(is.PROPERTY)
		is.Read(is.BEGIN_BRACKET)

		ay := model.NewArray()
		is.Read(&ay)
		geom.TexCoordArrayList[index] = &ay
		is.Read(is.END_BRACKET)
	}
	is.Read(is.END_BRACKET)
}

func writeTexCoordData(os *io.OsgOstream, g interface{}) {
	geom := g.(*model.Geometry)
	length := len(geom.TexCoordArrayList)
	os.Write(length)
	os.Write(os.BEGIN_BRACKET)
	os.Write(os.CRLF)
	os.PROPERTY.Name = "Data"
	for _, ary := range geom.TexCoordArrayList {
		os.Write(os.PROPERTY)
		os.Write(os.BEGIN_BRACKET)
		os.Write(os.CRLF)
		os.Write(ary)
		os.Write(os.END_BRACKET)
		os.Write(os.CRLF)
	}
	os.Write(os.END_BRACKET)
	os.Write(os.CRLF)
}

func checkVertexAttribData(g interface{}) bool {
	geom := g.(*model.Geometry)
	return len(geom.VertexAttribList) > 0
}

func readVertexAttribData(is *io.OsgIstream, g interface{}) {
	geom := g.(*model.Geometry)
	var size int = 0
	is.Read(&size)
	is.Read(is.BEGIN_BRACKET)
	index := 0
	is.PROPERTY.Name = "Data"
	geom.VertexAttribList = make([]*model.Array, size, size)
	for {
		if index == size {
			break
		}
		size--
		is.Read(is.PROPERTY)
		is.Read(is.BEGIN_BRACKET)

		ay := model.NewArray()
		is.Read(&ay)
		geom.SetVertexAttribArray(index, &ay, model.BIND_UNDEFINED)
		is.Read(is.END_BRACKET)
	}
	is.Read(is.END_BRACKET)
}

func writeVertexAttribData(os *io.OsgOstream, g interface{}) {
	geom := g.(*model.Geometry)
	length := len(geom.VertexAttribList)
	os.Write(length)
	os.Write(os.BEGIN_BRACKET)
	os.Write(os.CRLF)
	os.PROPERTY.Name = "Data"
	for _, ary := range geom.VertexAttribList {
		os.Write(os.PROPERTY)
		os.Write(os.BEGIN_BRACKET)
		os.Write(os.CRLF)
		os.Write(ary)
		os.Write(os.END_BRACKET)
		os.Write(os.CRLF)
	}
	os.Write(os.END_BRACKET)
	os.Write(os.CRLF)
}

func checkFastPathHint(g interface{}) bool {
	return false
}

func readFastPathHint(is *io.OsgIstream, g interface{}) {
	value := false
	if !is.IsBinary() {
		is.Read(&value)
	}
}

func writeFastPathHint(os *io.OsgOstream, g interface{}) {
}

func getPrimitiveSetList(obj interface{}) interface{} {
	return obj.(*model.Geometry).Primitives
}

func setPrimitiveSetList(obj interface{}, pro interface{}) {
	obj.(*model.Geometry).AddPrimitiveSet(pro.(*model.PrimitiveSet))
}

func getTexCoordArrayList(obj interface{}) interface{} {
	return obj.(*model.Geometry).TexCoordArrayList
}

func setTexCoordArrayList(obj interface{}, pro interface{}) {
	pl := obj.(*model.Geometry).TexCoordArrayList
	pl = append(pl, pro.(*model.Array))
}

func getVertexAttribArrayList(obj interface{}) interface{} {
	return obj.(*model.Geometry).VertexAttribList
}

func setVertexAttribArrayList(obj interface{}, pro interface{}) {
	pl := obj.(*model.Geometry).VertexAttribList
	pl = append(pl, pro.(*model.Array))
}

func getVertexData(obj interface{}) interface{} {
	return obj.(*model.Geometry).VertexArray
}

func setVertexData(obj interface{}, pro interface{}) {
	obj.(*model.Geometry).VertexArray = pro.(*model.Array)
}

func getNormalData(obj interface{}) interface{} {
	return obj.(*model.Geometry).NormalArray
}

func setNormalData(obj interface{}, pro interface{}) {
	obj.(*model.Geometry).NormalArray = pro.(*model.Array)

}

func getColorData(obj interface{}) interface{} {
	return obj.(*model.Geometry).ColorArray
}

func setColorData(obj interface{}, pro interface{}) {
	obj.(*model.Geometry).ColorArray = pro.(*model.Array)

}

func getSecondaryColorData(obj interface{}) interface{} {
	return obj.(*model.Geometry).SecondaryColorArray

}

func setSecondaryColorData(obj interface{}, pro interface{}) {
	obj.(*model.Geometry).SecondaryColorArray = pro.(*model.Array)

}

func getFogCoordData(obj interface{}) interface{} {
	return obj.(*model.Geometry).FogCoordArray

}
func setFogCoordData(obj interface{}, pro interface{}) {
	obj.(*model.Geometry).FogCoordArray = pro.(*model.Array)
}

func init() {
	lk := io.NewIntLookup()
	lookup = &lk
	lookup.Add("BIND_OFF", 0)
	lookup.Add("BIND_OVERALL", 1)
	lookup.Add("BIND_PER_PRIMITIVE_SET", 2)
	lookup.Add("BIND_PER_PRIMITIVE", 3)
	lookup.Add("BIND_PER_VERTEX", 4)

	fn := func() interface{} {
		g := model.NewGeometry()
		return &g
	}

	wrap := io.NewObjectWrapper("Geometry", fn, "osg::Object osg::Node osg::Drawable osg::Geometry")
	io.AddUpdateWrapperVersionProxy(&wrap, 154)
	wrap.MarkSerializerAsAdded("osg::Node")
	vser := io.NewVectorSerializer("PrimitiveSetList", io.RW_OBJECT, 0, getPrimitiveSetList, setPrimitiveSetList)
	wrap.AddSerializer(&vser, io.RW_VECTOR)

	ser1 := io.NewUserSerializer("VertexData", checkVertexData, readVertexData, writeVertexData)
	ser2 := io.NewUserSerializer("NormalData", checkNormalData, readNormalData, writeNormalData)
	ser3 := io.NewUserSerializer("ColorData", checkColorData, readColorData, writeColorData)
	ser4 := io.NewUserSerializer("SecondaryColorData", checkSecondaryColorData, readSecondaryColorData, writeSecondaryColorData)
	ser5 := io.NewUserSerializer("FogCoordData", checkFogCoordData, readFogCoordData, writeFogCoordData)
	ser6 := io.NewUserSerializer("TexCoordData", checkTexCoordData, readTexCoordData, writeTexCoordData)
	ser7 := io.NewUserSerializer("VertexAttribData", checkVertexAttribData, readVertexAttribData, writeVertexAttribData)
	ser8 := io.NewUserSerializer("FastPathHint", checkFastPathHint, readFastPathHint, writeFastPathHint)

	wrap.AddSerializer(&ser1, io.RW_USER)
	wrap.AddSerializer(&ser2, io.RW_USER)
	wrap.AddSerializer(&ser3, io.RW_USER)
	wrap.AddSerializer(&ser4, io.RW_USER)
	wrap.AddSerializer(&ser5, io.RW_USER)
	wrap.AddSerializer(&ser6, io.RW_USER)
	wrap.AddSerializer(&ser7, io.RW_USER)
	wrap.AddSerializer(&ser8, io.RW_USER)

	io.AddUpdateWrapperVersionProxy(&wrap, 112)
	wrap.MarkSerializerAsRemoved("VertexData")
	wrap.MarkSerializerAsRemoved("NormalData")
	wrap.MarkSerializerAsRemoved("ColorData")
	wrap.MarkSerializerAsRemoved("SecondaryColorData")
	wrap.MarkSerializerAsRemoved("FogCoordData")
	wrap.MarkSerializerAsRemoved("TexCoordData")
	wrap.MarkSerializerAsRemoved("VertexAttribData")
	wrap.MarkSerializerAsRemoved("FastPathHint")

	ser11 := io.NewObjectSerializer("VertexData", getVertexData, setVertexData)
	ser21 := io.NewObjectSerializer("NormalData", getNormalData, setNormalData)
	ser31 := io.NewObjectSerializer("ColorData", getColorData, setColorData)
	ser41 := io.NewObjectSerializer("SecondaryColorData", getSecondaryColorData, setSecondaryColorData)
	ser51 := io.NewObjectSerializer("FogCoordData", getFogCoordData, setFogCoordData)

	wrap.AddSerializer(&ser11, io.RW_OBJECT)
	wrap.AddSerializer(&ser21, io.RW_OBJECT)
	wrap.AddSerializer(&ser31, io.RW_OBJECT)
	wrap.AddSerializer(&ser41, io.RW_OBJECT)
	wrap.AddSerializer(&ser51, io.RW_OBJECT)

	vser2 := io.NewVectorSerializer("TexCoordArrayList", io.RW_OBJECT, 0, getTexCoordArrayList, setTexCoordArrayList)
	vser3 := io.NewVectorSerializer("VertexAttribArrayList", io.RW_OBJECT, 0, getVertexAttribArrayList, setVertexAttribArrayList)
	wrap.AddSerializer(&vser2, io.RW_VECTOR)
	wrap.AddSerializer(&vser3, io.RW_VECTOR)
	io.GetObjectWrapperManager().AddWrap(&wrap)
}
