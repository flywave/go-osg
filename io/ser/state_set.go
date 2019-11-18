package ser

import (
	"strings"

	"github.com/flywave/go-osg/io"
	"github.com/flywave/go-osg/model"
)

func readValue(is *io.OsgIstream) int {
	var val int
	var str string
	if is.IsBinary() {
		is.Read(&val)
	} else {
		is.Read(&str)
		if strings.Contains(str, "OFF") {
			val = model.OFF
		}
		if strings.Contains(str, "ON") {
			val = model.ON
		}
		if strings.Contains(str, "OVERRIDE") {
			val = model.OVERRIDE
		}
		if strings.Contains(str, "PROTECTED") {
			val = model.PROTECTED
		}
		if strings.Contains(str, "INHERIT") {
			val = model.INHERIT
		}
	}
	return val
}

func readModes(is *io.OsgIstream, mdlist model.ModeListType) {
	size := is.ReadSize()
	if size > 0 {
		for i := 0; i < size; i++ {
			var md int
			is.Read(&md)
			val := readValue(is)
			mdlist[md] = val
		}
		is.Read(is.END_BRACKET)
	}
}

func readAttributes(is *io.OsgIstream, attr model.AttributeListType) {
	size := is.ReadSize()
	if size > 0 {
		is.Read(is.BEGIN_BRACKET)
		for i := 0; i < size; i++ {
			ob := is.ReadObject()
			is.PROPERTY.Name = "Value"
			is.Read(is.PROPERTY)
			val := readValue(is)
			if model.IsBaseOfStateAttribute(ob) {
				sa := ob.(model.StateAttributeInterface)
				rp := model.RefAttributePair{ob, val}
				attr[sa.GetType()] = &rp
			}
		}
		is.Read(is.END_BRACKET)
	}
}

func writeValue(os *io.OsgOstream, val int) {
	if os.IsBinary() {
		os.Write(val)
	} else {
		var str string = ""
		if val|model.ON != 0 {
			str += "|ON"
		}
		if val|model.OVERRIDE != 0 {
			str += "|OVERRIDE"
		}
		if val|model.ON != 0 {
			str += "|PROTECTED"
		}
		if val|model.ON != 0 {
			str += "|INHERIT"
		}
		if str != "" {
			os.Write(str)
		} else {
			os.Write("OFF")
		}
	}
}

func writeModes(os *io.OsgOstream, mdlist model.ModeListType) {
	size := len(mdlist)
	os.Write(size)
	if size > 0 {
		os.Write(os.BEGIN_BRACKET)
		os.Write(os.CRLF)
		for k, v := range mdlist {
			os.Write(k)
			os.Write(v)
			os.Write(os.CRLF)
		}
		os.Write(os.END_BRACKET)
		os.Write(os.CRLF)
	}
}

func writeAttributes(os *io.OsgOstream, attr model.AttributeListType) {
	size := len(attr)
	os.Write(size)
	if size > 0 {
		os.Write(os.BEGIN_BRACKET)
		os.Write(os.CRLF)
		os.PROPERTY.Name = "Value"
		for _, v := range attr {
			os.Write(v.First)
			os.Write(os.PROPERTY)
			writeValue(os, v.Second)
			os.Write(os.CRLF)
		}
		os.Write(os.END_BRACKET)
		os.Write(os.CRLF)
	}
}

func checkModeList(obj interface{}) bool {
	ss := obj.(*model.StateSet)
	return len(ss.ModeList) > 0
}

func readModeList(is *io.OsgIstream, val interface{}) {
	ss := val.(*model.StateSet)
	readModes(is, ss.ModeList)
}

func writeModeList(os *io.OsgOstream, obj interface{}) {
	ss := obj.(*model.StateSet)
	writeModes(os, ss.ModeList)
}

func checkAttributeList(obj interface{}) bool {
	ss := obj.(*model.StateSet)
	return len(ss.AttributeList) > 0
}

func readAttributeList(is *io.OsgIstream, obj interface{}) {
	ss := obj.(*model.StateSet)
	readAttributes(is, ss.AttributeList)
}

func writeAttributeList(os *io.OsgOstream, obj interface{}) {
	ss := obj.(*model.StateSet)
	writeAttributes(os, ss.AttributeList)
}

func checkTextureModeList(obj interface{}) bool {
	ss := obj.(*model.StateSet)
	return len(ss.TextureModeList) > 0
}

func readTextureModeList(is *io.OsgIstream, obj interface{}) {
	ss := obj.(*model.StateSet)
	size := is.ReadSize()
	is.Read(is.BEGIN_BRACKET)
	if size > 0 {
		is.PROPERTY.Name = "Data"
		for i := 0; i < size; i++ {
			tmp := make(model.ModeListType)
			readModes(is, tmp)
			for k, v := range tmp {
				ss.SetTextureMode(i, k, v)
			}
		}
		is.Read(is.END_BRACKET)
	}
}

func writeTextureModeList(os *io.OsgOstream, obj interface{}) {
	ss := obj.(*model.StateSet)
	size := len(ss.TextureModeList)
	os.Write(size)
	os.PROPERTY.Name = "Data"
	for _, tl := range ss.TextureModeList {
		os.Write(os.PROPERTY)
		writeModes(os, tl)
	}
	os.Write(os.BEGIN_BRACKET)
	os.Write(os.CRLF)
}

func checkTextureAttributeList(obj interface{}) bool {
	ss := obj.(*model.StateSet)
	return len(ss.TextureAttributeList) > 0
}

func readTextureAttributeList(is *io.OsgIstream, obj interface{}) {
	ss := obj.(*model.StateSet)
	size := is.ReadSize()
	is.Read(is.BEGIN_BRACKET)
	if size > 0 {
		is.PROPERTY.Name = "Data"
		for i := 0; i < size; i++ {
			tmp := make(model.AttributeListType)
			readAttributes(is, tmp)
			for _, v := range tmp {
				ss.SetTextureAttribute(i, v.First, v.Second)
			}
		}
		is.Read(is.END_BRACKET)
	}
}

func writeTextureAttributeList(os *io.OsgOstream, obj interface{}) {
	ss := obj.(*model.StateSet)
	size := len(ss.TextureAttributeList)
	os.Write(size)
	os.Write(os.BEGIN_BRACKET)
	os.Write(os.CRLF)
	os.PROPERTY.Name = "Data"
	for _, ta := range ss.TextureAttributeList {
		os.Write(os.PROPERTY)
		writeAttributes(os, ta)
	}
	os.Write(os.END_BRACKET)
	os.Write(os.CRLF)
}

func checkUniformList(obj interface{}) bool {
	ss := obj.(*model.StateSet)
	return len(ss.UniformList) > 0
}

func readUniformList(is *io.OsgIstream, obj interface{}) {
	size := is.ReadSize()
	is.Read(is.BEGIN_BRACKET)
	if size > 0 {
		is.PROPERTY.Name = "Value"
		for i := 0; i < size; i++ {
			obj := is.ReadObject()
			is.Read(is.PROPERTY)
			if model.IsBaseOfUniform(obj) {
				readValue(is) //ignore
			}
		}
	}
	is.Read(is.END_BRACKET)
}

func writeUniformList(os *io.OsgOstream, obj interface{}) {
	ss := obj.(*model.StateSet)
	size := len(ss.UniformList)
	os.Write(size)
	os.PROPERTY.Name = "Value"
	for _, l := range ss.UniformList {
		os.Write(l.First)
		os.Write(os.PROPERTY)
		writeValue(os, l.Second)
		os.Write(os.CRLF)
	}
	os.Write(os.BEGIN_BRACKET)
	os.Write(os.CRLF)
}

func checkDefineList(obj interface{}) bool {
	ss := obj.(*model.StateSet)
	return len(ss.DefineList) > 0
}

func readDefineList(is *io.OsgIstream, obj interface{}) {
	ss := obj.(*model.StateSet)
	size := is.ReadSize()
	var defineName, defineValue string
	for i := 0; i < size; i++ {
		is.ReadWrappedString(&defineName)
		is.ReadWrappedString(&defineValue)
		is.PROPERTY.Name = "Value"
		is.Read(is.PROPERTY)
		val := readValue(is)
		ss.SetDefine(defineName, defineValue, val)
	}
	is.Read(is.END_BRACKET)
}

func writeDefineList(os *io.OsgOstream, obj interface{}) {
	ss := obj.(*model.StateSet)
	size := len(ss.DefineList)
	os.Write(size)
	os.Write(os.BEGIN_BRACKET)
	os.Write(os.CRLF)
	for k, v := range ss.DefineList {
		os.WriteWrappedString(k)
		os.WriteWrappedString(v.First)
		os.PROPERTY.Name = "Value"
		os.Write(os.PROPERTY)
		writeValue(os, v.Second)
		os.Write(os.CRLF)
	}
	os.Write(os.END_BRACKET)
	os.Write(os.CRLF)
}

func getRenderingHint(obj interface{}) interface{} {
	ss := obj.(*model.StateSet)
	return &ss.RenderingHint
}

func setRenderingHint(obj interface{}, val interface{}) {
	ss := obj.(*model.StateSet)
	ss.RenderingHint = val.(int)
}

func getRenderBinMode(obj interface{}) interface{} {
	ss := obj.(*model.StateSet)
	return &ss.BinMode
}

func setRenderBinMode(obj interface{}, val interface{}) {
	ss := obj.(*model.StateSet)
	ss.BinMode = val.(int)
}

func getBinNumber(obj interface{}) interface{} {
	ss := obj.(*model.StateSet)
	return &ss.BinNum
}

func setBinNumber(obj interface{}, val interface{}) {
	ss := obj.(*model.StateSet)
	ss.BinNum = val.(int)
}

func getBinName(obj interface{}) interface{} {
	ss := obj.(*model.StateSet)
	return &ss.BinName
}

func setBinName(obj interface{}, val interface{}) {
	ss := obj.(*model.StateSet)
	ss.BinName = val.(string)
}

func getNestRenderBins(obj interface{}) interface{} {
	ss := obj.(*model.StateSet)
	return &ss.NestRenderBins
}

func setNestRenderBins(obj interface{}, val interface{}) {
	ss := obj.(*model.StateSet)
	ss.NestRenderBins = val.(bool)
}

func getUpdateCallbackSS(obj interface{}) interface{} {
	ss := obj.(*model.StateSet)
	return ss.UpdateCallback
}

func setUpdateCallbackSS(obj interface{}, val interface{}) {
	ss := obj.(*model.StateSet)
	ss.UpdateCallback = val.(*model.Callback)
}

func getEventCallbackSS(obj interface{}) interface{} {
	ss := obj.(*model.StateSet)
	return ss.EventCallback
}

func setEventCallbackSS(obj interface{}, val interface{}) {
	ss := obj.(*model.StateSet)
	ss.EventCallback = val.(*model.Callback)
}

func init() {
	fn := func() interface{} {
		ss := model.NewStateSet()
		return &ss
	}
	wrap := io.NewObjectWrapper2("StateSet", "flywave::osg::stateset", fn, "osg::Object osg::StateSet")
	ser1 := io.NewUserSerializer("ModeList", checkModeList, readModeList, writeModeList)
	ser2 := io.NewUserSerializer("AttributeList", checkAttributeList, readAttributeList, writeAttributeList)
	ser3 := io.NewUserSerializer("TextureModeList", checkTextureModeList, readTextureModeList, writeTextureModeList)
	ser4 := io.NewUserSerializer("TextureAttributeList", checkTextureAttributeList, readTextureAttributeList, writeTextureAttributeList)
	ser5 := io.NewUserSerializer("UniformList", checkUniformList, readUniformList, writeUniformList)
	ser6 := io.NewEnumSerializer("RenderingHint", getRenderingHint, setRenderingHint)
	ser7 := io.NewEnumSerializer("RenderBinMode", getRenderBinMode, setRenderBinMode)
	ser8 := io.NewPropByValSerializer("BinNumber", false, getBinNumber, setBinNumber)
	ser9 := io.NewStringSerializer("BinName", getBinName, setBinName)
	ser10 := io.NewEnumSerializer("NestRenderBins", getNestRenderBins, setNestRenderBins)
	ser11 := io.NewObjectSerializer("UpdateCallback", getUpdateCallbackSS, setUpdateCallbackSS)
	ser12 := io.NewObjectSerializer("EventCallback", getEventCallbackSS, setEventCallbackSS)

	wrap.AddSerializer(&ser1, io.RW_USER)
	wrap.AddSerializer(&ser2, io.RW_USER)
	wrap.AddSerializer(&ser3, io.RW_USER)
	wrap.AddSerializer(&ser4, io.RW_USER)
	wrap.AddSerializer(&ser5, io.RW_USER)

	wrap.AddSerializer(&ser6, io.RW_INT)
	wrap.AddSerializer(&ser7, io.RW_INT)
	wrap.AddSerializer(&ser8, io.RW_INT)
	wrap.AddSerializer(&ser9, io.RW_STRING)
	wrap.AddSerializer(&ser10, io.RW_BOOL)
	wrap.AddSerializer(&ser11, io.RW_OBJECT)
	wrap.AddSerializer(&ser12, io.RW_OBJECT)
	io.GetObjectWrapperManager().AddWrap(&wrap)
}
