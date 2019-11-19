package serializer

import (
	"github.com/flywave/go-osg"
	"github.com/flywave/go-osg/model"
)

var lookup *osg.IntLookup

func readAttributeBinding(is *osg.OsgIstream) int32 {
	if is.IsBinary() {
		var val int32
		is.Read(&val)
		return val
	} else {
		str := is.ReadString()
		return lookup.GetValue(str)
	}
}

func writeAttributeBinding(os *osg.OsgOstream, val int32) {
	if os.IsBinary() {
		os.Write(val)
	} else {
		os.Write(lookup.GetString(val))
	}
}

func readArray(is *osg.OsgIstream) *model.Array {
	is.PROPERTY.Name = "Array"
	hasArray := false
	var ary *model.Array
	is.Read(is.PROPERTY)
	is.Read(&hasArray)
	if hasArray {
		ary = is.ReadArray()
	}

	is.PROPERTY.Name = "Indices"
	hasIndices := false
	is.Read(is.PROPERTY)
	is.Read(&hasIndices)
	if hasIndices {
		ary2 := is.ReadArray()
		if ary != nil {
			ary.Udc.UserData = ary2
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

func writeArray(os *osg.OsgOstream, ary *model.Array) {
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
func readVertexData(is *osg.OsgIstream, g interface{}) {
	geom := g.(*model.Geometry)
	is.Read(is.BEGINBRACKET)
	is.Read(is.CRLF)
	ary := readArray(is)
	geom.VertexArray = ary
	is.Read(is.ENDBRACKET)
	is.Read(is.CRLF)
}
func writeVertexData(os *osg.OsgOstream, g interface{}) {
	geom := g.(*model.Geometry)
	os.Write(os.BEGINBRACKET)
	os.Write(os.CRLF)
	writeArray(os, geom.VertexArray)
	os.Write(os.ENDBRACKET)
	os.Write(os.CRLF)
}

func checkNormalData(g interface{}) bool {
	geom := g.(*model.Geometry)
	return geom.NormalArray != nil
}
func readNormalData(is *osg.OsgIstream, g interface{}) {
	geom := g.(*model.Geometry)
	is.Read(is.BEGINBRACKET)
	is.Read(is.CRLF)
	ary := readArray(is)
	geom.NormalArray = ary
	is.Read(is.ENDBRACKET)
	is.Read(is.CRLF)
}
func writeNormalData(os *osg.OsgOstream, g interface{}) {
	geom := g.(*model.Geometry)
	os.Write(os.BEGINBRACKET)
	os.Write(os.CRLF)
	writeArray(os, geom.NormalArray)
	os.Write(os.ENDBRACKET)
	os.Write(os.CRLF)
}

func checkColorData(g interface{}) bool {
	geom := g.(*model.Geometry)
	return geom.ColorArray != nil
}
func readColorData(is *osg.OsgIstream, g interface{}) {
	geom := g.(*model.Geometry)
	is.Read(is.BEGINBRACKET)
	is.Read(is.CRLF)
	ary := readArray(is)
	geom.ColorArray = ary
	is.Read(is.ENDBRACKET)
	is.Read(is.CRLF)
}
func writeColorData(os *osg.OsgOstream, g interface{}) {
	geom := g.(*model.Geometry)
	os.Write(os.BEGINBRACKET)
	os.Write(os.CRLF)
	writeArray(os, geom.ColorArray)
	os.Write(os.ENDBRACKET)
	os.Write(os.CRLF)
}

func checkSecondaryColorData(g interface{}) bool {
	geom := g.(*model.Geometry)
	return geom.SecondaryColorArray != nil
}
func readSecondaryColorData(is *osg.OsgIstream, g interface{}) {
	geom := g.(*model.Geometry)
	is.Read(is.BEGINBRACKET)
	is.Read(is.CRLF)
	ary := readArray(is)
	geom.SecondaryColorArray = ary
	is.Read(is.ENDBRACKET)
	is.Read(is.CRLF)
}
func writeSecondaryColorData(os *osg.OsgOstream, g interface{}) {
	geom := g.(*model.Geometry)
	os.Write(os.BEGINBRACKET)
	os.Write(os.CRLF)
	writeArray(os, geom.SecondaryColorArray)
	os.Write(os.ENDBRACKET)
	os.Write(os.CRLF)
}

func checkFogCoordData(g interface{}) bool {
	geom := g.(*model.Geometry)
	return geom.FogCoordArray != nil
}
func readFogCoordData(is *osg.OsgIstream, g interface{}) {
	geom := g.(*model.Geometry)
	is.Read(is.BEGINBRACKET)
	is.Read(is.CRLF)
	ary := readArray(is)
	geom.FogCoordArray = ary
	is.Read(is.ENDBRACKET)
	is.Read(is.CRLF)
}
func writeFogCoordData(os *osg.OsgOstream, g interface{}) {
	geom := g.(*model.Geometry)
	os.Write(os.BEGINBRACKET)
	os.Write(os.CRLF)
	writeArray(os, geom.FogCoordArray)
	os.Write(os.ENDBRACKET)
	os.Write(os.CRLF)
}

func checkTexCoordData(g interface{}) bool {
	geom := g.(*model.Geometry)
	return geom.TexCoordArrayList != nil
}

func readTexCoordData(is *osg.OsgIstream, g interface{}) {
	geom := g.(*model.Geometry)
	var size int = 0
	is.Read(&size)
	is.Read(is.BEGINBRACKET)
	index := 0
	is.PROPERTY.Name = "Data"
	geom.TexCoordArrayList = make([]*model.Array, size, size)
	for {
		if index == size {
			break
		}
		size--
		is.Read(is.PROPERTY)
		is.Read(is.BEGINBRACKET)

		ay := is.ReadArray()
		geom.TexCoordArrayList[index] = ay
		is.Read(is.ENDBRACKET)
	}
	is.Read(is.ENDBRACKET)
}

func writeTexCoordData(os *osg.OsgOstream, g interface{}) {
	geom := g.(*model.Geometry)
	length := len(geom.TexCoordArrayList)
	os.Write(length)
	os.Write(os.BEGINBRACKET)
	os.Write(os.CRLF)
	os.PROPERTY.Name = "Data"
	for _, ary := range geom.TexCoordArrayList {
		os.Write(os.PROPERTY)
		os.Write(os.BEGINBRACKET)
		os.Write(os.CRLF)
		os.Write(ary)
		os.Write(os.ENDBRACKET)
		os.Write(os.CRLF)
	}
	os.Write(os.ENDBRACKET)
	os.Write(os.CRLF)
}

func checkVertexAttribData(g interface{}) bool {
	geom := g.(*model.Geometry)
	return len(geom.VertexAttribList) > 0
}

func readVertexAttribData(is *osg.OsgIstream, g interface{}) {
	geom := g.(*model.Geometry)
	var size int = 0
	is.Read(&size)
	is.Read(is.BEGINBRACKET)
	index := 0
	is.PROPERTY.Name = "Data"
	geom.VertexAttribList = make([]*model.Array, size, size)
	for {
		if index == size {
			break
		}
		size--
		is.Read(is.PROPERTY)
		is.Read(is.BEGINBRACKET)

		ay := is.ReadArray()
		geom.SetVertexAttribArray(index, ay, model.BINDUNDEFINED)
		is.Read(is.ENDBRACKET)
	}
	is.Read(is.ENDBRACKET)
}

func writeVertexAttribData(os *osg.OsgOstream, g interface{}) {
	geom := g.(*model.Geometry)
	length := len(geom.VertexAttribList)
	os.Write(length)
	os.Write(os.BEGINBRACKET)
	os.Write(os.CRLF)
	os.PROPERTY.Name = "Data"
	for _, ary := range geom.VertexAttribList {
		os.Write(os.PROPERTY)
		os.Write(os.BEGINBRACKET)
		os.Write(os.CRLF)
		os.Write(ary)
		os.Write(os.ENDBRACKET)
		os.Write(os.CRLF)
	}
	os.Write(os.ENDBRACKET)
	os.Write(os.CRLF)
}

func checkFastPathHint(g interface{}) bool {
	return false
}

func readFastPathHint(is *osg.OsgIstream, g interface{}) {
	value := false
	if !is.IsBinary() {
		is.Read(&value)
	}
}

func writeFastPathHint(os *osg.OsgOstream, g interface{}) {
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
	lk := osg.NewIntLookup()
	lookup = &lk
	lookup.Add("BINDOFF", 0)
	lookup.Add("BINDOVERALL", 1)
	lookup.Add("BINDPERPRIMITIVESET", 2)
	lookup.Add("BINDPERPRIMITIVE", 3)
	lookup.Add("BINDPERVERTEX", 4)

	fn := func() interface{} {
		g := model.NewGeometry()
		return &g
	}

	wrap := osg.NewObjectWrapper("Geometry", fn, "osg::Object osg::Node osg::Drawable osg::Geometry")
	osg.AddUpdateWrapperVersionProxy(&wrap, 154)
	wrap.MarkSerializerAsAdded("osg::Node")
	ps := model.NewPrimitiveSet()
	vser := osg.NewVectorSerializer("PrimitiveSetList", osg.RWOBJECT, &ps, getPrimitiveSetList, setPrimitiveSetList)
	wrap.AddSerializer(&vser, osg.RWVECTOR)

	ser1 := osg.NewUserSerializer("VertexData", checkVertexData, readVertexData, writeVertexData)
	ser2 := osg.NewUserSerializer("NormalData", checkNormalData, readNormalData, writeNormalData)
	ser3 := osg.NewUserSerializer("ColorData", checkColorData, readColorData, writeColorData)
	ser4 := osg.NewUserSerializer("SecondaryColorData", checkSecondaryColorData, readSecondaryColorData, writeSecondaryColorData)
	ser5 := osg.NewUserSerializer("FogCoordData", checkFogCoordData, readFogCoordData, writeFogCoordData)
	ser6 := osg.NewUserSerializer("TexCoordData", checkTexCoordData, readTexCoordData, writeTexCoordData)
	ser7 := osg.NewUserSerializer("VertexAttribData", checkVertexAttribData, readVertexAttribData, writeVertexAttribData)
	ser8 := osg.NewUserSerializer("FastPathHint", checkFastPathHint, readFastPathHint, writeFastPathHint)

	wrap.AddSerializer(&ser1, osg.RWUSER)
	wrap.AddSerializer(&ser2, osg.RWUSER)
	wrap.AddSerializer(&ser3, osg.RWUSER)
	wrap.AddSerializer(&ser4, osg.RWUSER)
	wrap.AddSerializer(&ser5, osg.RWUSER)
	wrap.AddSerializer(&ser6, osg.RWUSER)
	wrap.AddSerializer(&ser7, osg.RWUSER)
	wrap.AddSerializer(&ser8, osg.RWUSER)

	osg.AddUpdateWrapperVersionProxy(&wrap, 112)
	wrap.MarkSerializerAsRemoved("VertexData")
	wrap.MarkSerializerAsRemoved("NormalData")
	wrap.MarkSerializerAsRemoved("ColorData")
	wrap.MarkSerializerAsRemoved("SecondaryColorData")
	wrap.MarkSerializerAsRemoved("FogCoordData")
	wrap.MarkSerializerAsRemoved("TexCoordData")
	wrap.MarkSerializerAsRemoved("VertexAttribData")
	wrap.MarkSerializerAsRemoved("FastPathHint")

	ser11 := osg.NewObjectSerializer("VertexData", getVertexData, setVertexData)
	ser21 := osg.NewObjectSerializer("NormalData", getNormalData, setNormalData)
	ser31 := osg.NewObjectSerializer("ColorData", getColorData, setColorData)
	ser41 := osg.NewObjectSerializer("SecondaryColorData", getSecondaryColorData, setSecondaryColorData)
	ser51 := osg.NewObjectSerializer("FogCoordData", getFogCoordData, setFogCoordData)

	wrap.AddSerializer(&ser11, osg.RWOBJECT)
	wrap.AddSerializer(&ser21, osg.RWOBJECT)
	wrap.AddSerializer(&ser31, osg.RWOBJECT)
	wrap.AddSerializer(&ser41, osg.RWOBJECT)
	wrap.AddSerializer(&ser51, osg.RWOBJECT)

	vser2 := osg.NewVectorSerializer("TexCoordArrayList", osg.RWOBJECT, 0, getTexCoordArrayList, setTexCoordArrayList)
	vser3 := osg.NewVectorSerializer("VertexAttribArrayList", osg.RWOBJECT, 0, getVertexAttribArrayList, setVertexAttribArrayList)
	wrap.AddSerializer(&vser2, osg.RWVECTOR)
	wrap.AddSerializer(&vser3, osg.RWVECTOR)
	osg.GetObjectWrapperManager().AddWrap(&wrap)
}
