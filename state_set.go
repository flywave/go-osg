package osg

import (
	"strings"

	"github.com/flywave/go-osg/model"
)

func readValue(is *OsgIstream) int32 {
	var val int32
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

func readModes(is *OsgIstream, mdlist model.ModeListType) {
	size := is.ReadSize()
	if size > 0 {
		is.Read(is.BEGINBRACKET)
		for i := 0; i < size; i++ {
			md := model.ObjectGlenum{}
			is.Read(&md)
			val := readValue(is)
			mdlist[md.Value] = val
		}
		is.Read(is.ENDBRACKET)
	}
}

func readAttributes(is *OsgIstream, attr model.AttributeListType) {
	size := is.ReadSize()
	if size > 0 {
		is.Read(is.BEGINBRACKET)
		for i := 0; i < size; i++ {
			ob := is.ReadObject(nil)
			is.PROPERTY.Name = "Value"
			is.Read(is.PROPERTY)
			val := readValue(is)
			sa, ok := ob.(model.StateAttributeInterface)
			if ok {
				rp := model.RefAttributePair{sa, val}
				attr[sa.GetType()] = &rp
			}
		}
		is.Read(is.ENDBRACKET)
	}
}

func writeValue(os *OsgOstream, val int32) {
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

func writeModes(os *OsgOstream, mdlist model.ModeListType) {
	size := len(mdlist)
	os.Write(size)
	if size > 0 {
		os.Write(os.BEGINBRACKET)
		os.Write(os.CRLF)
		for k, v := range mdlist {
			os.Write(k)
			os.Write(v)
			os.Write(os.CRLF)
		}
		os.Write(os.ENDBRACKET)
		os.Write(os.CRLF)
	}
}

func writeAttributes(os *OsgOstream, attr model.AttributeListType) {
	size := len(attr)
	os.Write(size)
	if size > 0 {
		os.Write(os.BEGINBRACKET)
		os.Write(os.CRLF)
		os.PROPERTY.Name = "Value"
		for _, v := range attr {
			os.Write(v.First)
			os.Write(os.PROPERTY)
			writeValue(os, v.Second)
			os.Write(os.CRLF)
		}
		os.Write(os.ENDBRACKET)
		os.Write(os.CRLF)
	}
}

func checkModeList(obj interface{}) bool {
	ss := obj.(*model.StateSet)
	return len(ss.ModeList) > 0
}

func readModeList(is *OsgIstream, val interface{}) {
	ss := val.(*model.StateSet)
	readModes(is, ss.ModeList)
}

func writeModeList(os *OsgOstream, obj interface{}) {
	ss := obj.(*model.StateSet)
	writeModes(os, ss.ModeList)
}

func checkAttributeList(obj interface{}) bool {
	ss := obj.(*model.StateSet)
	return len(ss.AttributeList) > 0
}

func readAttributeList(is *OsgIstream, obj interface{}) {
	ss := obj.(*model.StateSet)
	readAttributes(is, ss.AttributeList)
}

func writeAttributeList(os *OsgOstream, obj interface{}) {
	ss := obj.(*model.StateSet)
	writeAttributes(os, ss.AttributeList)
}

func checkTextureModeList(obj interface{}) bool {
	ss := obj.(*model.StateSet)
	return len(ss.TextureModeList) > 0
}

func readTextureModeList(is *OsgIstream, obj interface{}) {
	ss := obj.(*model.StateSet)
	size := is.ReadSize()
	is.Read(is.BEGINBRACKET)
	if size > 0 {
		is.PROPERTY.Name = "Data"
		for i := 0; i < size; i++ {
			is.Read(is.PROPERTY)
			tmp := make(model.ModeListType)
			readModes(is, tmp)
			for k, v := range tmp {
				ss.SetTextureMode(i, k, v)
			}
		}
		is.Read(is.ENDBRACKET)
	}
}

func writeTextureModeList(os *OsgOstream, obj interface{}) {
	ss := obj.(*model.StateSet)
	size := len(ss.TextureModeList)
	os.Write(size)
	os.PROPERTY.Name = "Data"
	for _, tl := range ss.TextureModeList {
		os.Write(os.PROPERTY)
		writeModes(os, tl)
	}
	os.Write(os.BEGINBRACKET)
	os.Write(os.CRLF)
}

func checkTextureAttributeList(obj interface{}) bool {
	ss := obj.(*model.StateSet)
	return len(ss.TextureAttributeList) > 0
}

func readTextureAttributeList(is *OsgIstream, obj interface{}) {
	ss := obj.(*model.StateSet)
	size := is.ReadSize()
	is.Read(is.BEGINBRACKET)
	if size > 0 {
		is.PROPERTY.Name = "Data"
		for i := 0; i < size; i++ {
			is.Read(is.PROPERTY)
			tmp := make(model.AttributeListType)
			readAttributes(is, tmp)
			for _, v := range tmp {
				ss.SetTextureAttribute(i, v.First, v.Second)
			}
		}
		is.Read(is.ENDBRACKET)
	}
}

func writeTextureAttributeList(os *OsgOstream, obj interface{}) {
	ss := obj.(*model.StateSet)
	size := len(ss.TextureAttributeList)
	os.Write(size)
	os.Write(os.BEGINBRACKET)
	os.Write(os.CRLF)
	os.PROPERTY.Name = "Data"
	for _, ta := range ss.TextureAttributeList {
		os.Write(os.PROPERTY)
		writeAttributes(os, ta)
	}
	os.Write(os.ENDBRACKET)
	os.Write(os.CRLF)
}

func checkUniformList(obj interface{}) bool {
	ss := obj.(*model.StateSet)
	return len(ss.UniformList) > 0
}

func readUniformList(is *OsgIstream, obj interface{}) {
	size := is.ReadSize()
	is.Read(is.BEGINBRACKET)
	if size > 0 {
		is.PROPERTY.Name = "Value"
		for i := 0; i < size; i++ {
			is.Read(is.PROPERTY)
			ob := is.ReadObject(nil)
			is.Read(is.PROPERTY)
			if model.IsBaseOfUniform(ob) {
				readValue(is) //ignore
			}
		}
	}
	is.Read(is.ENDBRACKET)
}

func writeUniformList(os *OsgOstream, obj interface{}) {
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
	os.Write(os.BEGINBRACKET)
	os.Write(os.CRLF)
}

func checkDefineList(obj interface{}) bool {
	ss := obj.(*model.StateSet)
	return len(ss.DefineList) > 0
}

func readDefineList(is *OsgIstream, obj interface{}) {
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
	is.Read(is.ENDBRACKET)
}

func writeDefineList(os *OsgOstream, obj interface{}) {
	ss := obj.(*model.StateSet)
	size := len(ss.DefineList)
	os.Write(size)
	os.Write(os.BEGINBRACKET)
	os.Write(os.CRLF)
	for k, v := range ss.DefineList {
		os.Write(&k)
		os.Write(&v.First)
		os.PROPERTY.Name = "Value"
		os.Write(os.PROPERTY)
		writeValue(os, v.Second)
		os.Write(os.CRLF)
	}
	os.Write(os.ENDBRACKET)
	os.Write(os.CRLF)
}

func getRenderingHint(obj interface{}) interface{} {
	ss := obj.(*model.StateSet)
	return &ss.RenderingHint
}

func setRenderingHint(obj interface{}, val interface{}) {
	ss := obj.(*model.StateSet)
	ss.RenderingHint = val.(int32)
}

func getRenderBinMode(obj interface{}) interface{} {
	ss := obj.(*model.StateSet)
	return &ss.BinMode
}

func setRenderBinMode(obj interface{}, val interface{}) {
	ss := obj.(*model.StateSet)
	ss.BinMode = val.(int32)
}

func getBinNumber(obj interface{}) interface{} {
	ss := obj.(*model.StateSet)
	return &ss.BinNum
}

func setBinNumber(obj interface{}, val interface{}) {
	ss := obj.(*model.StateSet)
	ss.BinNum = val.(int32)
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
		return ss
	}
	wrap := NewObjectWrapper("StateSet", fn, "osg::Object osg::StateSet")
	ser1 := NewUserSerializer("ModeList", checkModeList, readModeList, writeModeList)
	ser2 := NewUserSerializer("AttributeList", checkAttributeList, readAttributeList, writeAttributeList)
	ser3 := NewUserSerializer("TextureModeList", checkTextureModeList, readTextureModeList, writeTextureModeList)
	ser4 := NewUserSerializer("TextureAttributeList", checkTextureAttributeList, readTextureAttributeList, writeTextureAttributeList)
	ser5 := NewUserSerializer("UniformList", checkUniformList, readUniformList, writeUniformList)
	ser6 := NewEnumSerializer("RenderingHint", getRenderingHint, setRenderingHint)
	ser7 := NewEnumSerializer("RenderBinMode", getRenderBinMode, setRenderBinMode)
	ser8 := NewPropByValSerializer("BinNumber", false, getBinNumber, setBinNumber)
	ser9 := NewStringSerializer("BinName", getBinName, setBinName)
	ser10 := NewEnumSerializer("NestRenderBins", getNestRenderBins, setNestRenderBins)
	ser11 := NewObjectSerializer("UpdateCallback", getUpdateCallbackSS, setUpdateCallbackSS)
	ser12 := NewObjectSerializer("EventCallback", getEventCallbackSS, setEventCallbackSS)

	wrap.AddSerializer(ser1, RWUSER)
	wrap.AddSerializer(ser2, RWUSER)
	wrap.AddSerializer(ser3, RWUSER)
	wrap.AddSerializer(ser4, RWUSER)
	wrap.AddSerializer(ser5, RWUSER)

	wrap.AddSerializer(ser6, RWINT)
	wrap.AddSerializer(ser7, RWINT)
	wrap.AddSerializer(ser8, RWINT)
	wrap.AddSerializer(ser9, RWSTRING)
	wrap.AddSerializer(ser10, RWBOOL)
	wrap.AddSerializer(ser11, RWOBJECT)
	wrap.AddSerializer(ser12, RWOBJECT)
	GetObjectWrapperManager().AddWrap(wrap)
	{
		uv := AddUpdateWrapperVersionProxy(wrap, 151)
		wrap.MarkSerializerAsAdded("DefineList")
		uv.SetLastVersion()
	}
}
