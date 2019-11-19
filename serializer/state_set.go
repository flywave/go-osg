package serializer

import (
	"strings"

	"github.com/flywave/go-osg"
	"github.com/flywave/go-osg/model"
)

func readValue(is *osg.OsgIstream) int {
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

func readModes(is *osg.OsgIstream, mdlist model.ModeListType) {
	size := is.ReadSize()
	if size > 0 {
		for i := 0; i < size; i++ {
			var md int
			is.Read(&md)
			val := readValue(is)
			mdlist[md] = val
		}
		is.Read(is.ENDBRACKET)
	}
}

func readAttributes(is *osg.OsgIstream, attr model.AttributeListType) {
	size := is.ReadSize()
	if size > 0 {
		is.Read(is.BEGINBRACKET)
		for i := 0; i < size; i++ {
			ob := is.ReadObject(nil)
			is.PROPERTY.Name = "Value"
			is.Read(is.PROPERTY)
			val := readValue(is)
			if model.IsBaseOfStateAttribute(ob) {
				sa := ob.(model.StateAttributeInterface)
				rp := model.RefAttributePair{ob, val} //TODO
				attr[sa.GetType()] = &rp
			}
		}
		is.Read(is.ENDBRACKET)
	}
}

func writeValue(os *osg.OsgOstream, val int) {
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

func writeModes(os *osg.OsgOstream, mdlist model.ModeListType) {
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

func writeAttributes(os *osg.OsgOstream, attr model.AttributeListType) {
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

func readModeList(is *osg.OsgIstream, val interface{}) {
	ss := val.(*model.StateSet)
	readModes(is, ss.ModeList)
}

func writeModeList(os *osg.OsgOstream, obj interface{}) {
	ss := obj.(*model.StateSet)
	writeModes(os, ss.ModeList)
}

func checkAttributeList(obj interface{}) bool {
	ss := obj.(*model.StateSet)
	return len(ss.AttributeList) > 0
}

func readAttributeList(is *osg.OsgIstream, obj interface{}) {
	ss := obj.(*model.StateSet)
	readAttributes(is, ss.AttributeList)
}

func writeAttributeList(os *osg.OsgOstream, obj interface{}) {
	ss := obj.(*model.StateSet)
	writeAttributes(os, ss.AttributeList)
}

func checkTextureModeList(obj interface{}) bool {
	ss := obj.(*model.StateSet)
	return len(ss.TextureModeList) > 0
}

func readTextureModeList(is *osg.OsgIstream, obj interface{}) {
	ss := obj.(*model.StateSet)
	size := is.ReadSize()
	is.Read(is.BEGINBRACKET)
	if size > 0 {
		is.PROPERTY.Name = "Data"
		for i := 0; i < size; i++ {
			tmp := make(model.ModeListType)
			readModes(is, tmp)
			for k, v := range tmp {
				ss.SetTextureMode(i, k, v)
			}
		}
		is.Read(is.ENDBRACKET)
	}
}

func writeTextureModeList(os *osg.OsgOstream, obj interface{}) {
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

func readTextureAttributeList(is *osg.OsgIstream, obj interface{}) {
	ss := obj.(*model.StateSet)
	size := is.ReadSize()
	is.Read(is.BEGINBRACKET)
	if size > 0 {
		is.PROPERTY.Name = "Data"
		for i := 0; i < size; i++ {
			tmp := make(model.AttributeListType)
			readAttributes(is, tmp)
			for _, v := range tmp {
				ss.SetTextureAttribute(i, v.First, v.Second)
			}
		}
		is.Read(is.ENDBRACKET)
	}
}

func writeTextureAttributeList(os *osg.OsgOstream, obj interface{}) {
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

func readUniformList(is *osg.OsgIstream, obj interface{}) {
	size := is.ReadSize()
	is.Read(is.BEGINBRACKET)
	if size > 0 {
		is.PROPERTY.Name = "Value"
		for i := 0; i < size; i++ {
			ob := is.ReadObject(nil)
			is.Read(is.PROPERTY)
			if model.IsBaseOfUniform(ob) {
				readValue(is) //ignore
			}
		}
	}
	is.Read(is.ENDBRACKET)
}

func writeUniformList(os *osg.OsgOstream, obj interface{}) {
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

func readDefineList(is *osg.OsgIstream, obj interface{}) {
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

func writeDefineList(os *osg.OsgOstream, obj interface{}) {
	ss := obj.(*model.StateSet)
	size := len(ss.DefineList)
	os.Write(size)
	os.Write(os.BEGINBRACKET)
	os.Write(os.CRLF)
	for k, v := range ss.DefineList {
		os.WriteWrappedString(k)
		os.WriteWrappedString(v.First)
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
	wrap := osg.NewObjectWrapper2("StateSet", "flywave::osg::stateset", fn, "osg::Object osg::StateSet")
	ser1 := osg.NewUserSerializer("ModeList", checkModeList, readModeList, writeModeList)
	ser2 := osg.NewUserSerializer("AttributeList", checkAttributeList, readAttributeList, writeAttributeList)
	ser3 := osg.NewUserSerializer("TextureModeList", checkTextureModeList, readTextureModeList, writeTextureModeList)
	ser4 := osg.NewUserSerializer("TextureAttributeList", checkTextureAttributeList, readTextureAttributeList, writeTextureAttributeList)
	ser5 := osg.NewUserSerializer("UniformList", checkUniformList, readUniformList, writeUniformList)
	ser6 := osg.NewEnumSerializer("RenderingHint", getRenderingHint, setRenderingHint)
	ser7 := osg.NewEnumSerializer("RenderBinMode", getRenderBinMode, setRenderBinMode)
	ser8 := osg.NewPropByValSerializer("BinNumber", false, getBinNumber, setBinNumber)
	ser9 := osg.NewStringSerializer("BinName", getBinName, setBinName)
	ser10 := osg.NewEnumSerializer("NestRenderBins", getNestRenderBins, setNestRenderBins)
	ser11 := osg.NewObjectSerializer("UpdateCallback", getUpdateCallbackSS, setUpdateCallbackSS)
	ser12 := osg.NewObjectSerializer("EventCallback", getEventCallbackSS, setEventCallbackSS)

	wrap.AddSerializer(&ser1, osg.RWUSER)
	wrap.AddSerializer(&ser2, osg.RWUSER)
	wrap.AddSerializer(&ser3, osg.RWUSER)
	wrap.AddSerializer(&ser4, osg.RWUSER)
	wrap.AddSerializer(&ser5, osg.RWUSER)

	wrap.AddSerializer(&ser6, osg.RWINT)
	wrap.AddSerializer(&ser7, osg.RWINT)
	wrap.AddSerializer(&ser8, osg.RWINT)
	wrap.AddSerializer(&ser9, osg.RWSTRING)
	wrap.AddSerializer(&ser10, osg.RWBOOL)
	wrap.AddSerializer(&ser11, osg.RWOBJECT)
	wrap.AddSerializer(&ser12, osg.RWOBJECT)
	osg.GetObjectWrapperManager().AddWrap(&wrap)
}
