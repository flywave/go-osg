package osg

import (
	"strings"

	"github.com/flywave/go-osg/model"
)

func checkWRAPS(obj interface{}) bool {
	return true
}

func readWRAPS(is *OsgIstream, obj interface{}) {
	tex := obj.(*model.Texture)
	mode := model.NewObjectGlenum()
	is.Read(mode)
	tex.SetWrap(model.WRAPS, int(mode.Value))
}

func writeWRAPS(os *OsgOstream, obj interface{}) {
	tex := obj.(*model.Texture)
	os.Write(tex.GetWrap(model.WRAPS))
	os.Write(os.CRLF)
}

func checkWRAPT(obj interface{}) bool {
	return true
}

func readWRAPT(is *OsgIstream, obj interface{}) {
	tex := obj.(*model.Texture)
	mode := model.NewObjectGlenum()
	is.Read(mode)
	tex.SetWrap(model.WRAPT, int(mode.Value))
}

func writeWRAPT(os *OsgOstream, obj interface{}) {
	tex := obj.(*model.Texture)
	os.Write(tex.GetWrap(model.WRAPT))
	os.Write(os.CRLF)
}

func checkWRAPR(obj interface{}) bool {
	return true
}

func readWRAPR(is *OsgIstream, obj interface{}) {
	tex := obj.(*model.Texture)
	mode := model.NewObjectGlenum()
	is.Read(mode)
	tex.SetWrap(model.WRAPR, int(mode.Value))
}

func writeWRAPR(os *OsgOstream, obj interface{}) {
	tex := obj.(*model.Texture)
	os.Write(tex.GetWrap(model.WRAPR))
	os.Write(os.CRLF)
}

func checkMINFILTER(obj interface{}) bool {
	return true
}

func readMINFILTER(is *OsgIstream, obj interface{}) {
	tex := obj.(*model.Texture)
	mode := model.NewObjectGlenum()
	is.Read(mode)
	tex.SetFilter(model.MINFILTER, int(mode.Value))
}

func writeMINFILTER(os *OsgOstream, obj interface{}) {
	tex := obj.(*model.Texture)
	tex.GetFilter(model.MINFILTER)
	os.Write(os.CRLF)
}

func checkMAGFILTER(obj interface{}) bool {
	return true
}

func readMAGFILTER(is *OsgIstream, obj interface{}) {
	tex := obj.(*model.Texture)
	mode := model.NewObjectGlenum()
	is.Read(mode)
	tex.SetFilter(model.MAGFILTER, int(mode.Value))
}

func writeMAGFILTER(os *OsgOstream, obj interface{}) {
	tex := obj.(*model.Texture)
	tex.GetFilter(model.MAGFILTER)
	os.Write(os.CRLF)
}

func checkSourceFormat(obj interface{}) bool {
	tex := obj.(*model.Texture)
	return tex.SourceFormat != model.GLZERO
}

func readSourceFormat(is *OsgIstream, obj interface{}) {
	tex := obj.(*model.Texture)
	mode := model.NewObjectGlenum()
	is.Read(mode)
	tex.SourceFormat = int(mode.Value)
}
func writeSourceFormat(os *OsgOstream, obj interface{}) {
	tex := obj.(*model.Texture)
	os.Write(tex.SourceFormat)
	os.Write(os.CRLF)
}

func checkSourceType(obj interface{}) bool {
	tex := obj.(*model.Texture)
	return tex.SourceType != model.GLZERO
}

func readSourceType(is *OsgIstream, obj interface{}) {
	tex := obj.(*model.Texture)
	mode := model.NewObjectGlenum()
	is.Read(mode)
	tex.SourceType = int(mode.Value)
}
func writeSourceType(os *OsgOstream, obj interface{}) {
	tex := obj.(*model.Texture)
	os.Write(tex.SourceType)
	os.Write(os.CRLF)
}

func checkInternalFormat(obj interface{}) bool {
	tex := obj.(*model.Texture)
	return tex.InternalFormat == model.USEUSERDEFINEDFORMAT
}

func readInternalFormat(is *OsgIstream, obj interface{}) {
	tex := obj.(*model.Texture)
	mode := model.NewObjectGlenum()
	is.Read(mode)
	if tex.InternalFormat == model.USEUSERDEFINEDFORMAT {
		tex.InternalFormat = int(mode.Value)
	}
}

func writeInternalFormat(os *OsgOstream, obj interface{}) {
	tex := obj.(*model.Texture)
	if os.IsBinary() && tex.InternalFormatMode != model.USEUSERDEFINEDFORMAT {
		os.Write(model.GLZERO)
		os.Write(os.CRLF)
	} else {
		os.Write(tex.InternalFormatMode)
		os.Write(os.CRLF)
	}
}

type dummyImageAttachment struct {
	unit    int
	level   int
	layered bool
	layer   int
	access  int
	format  int
}

var attachment1 dummyImageAttachment
var attachment2 dummyImageAttachment

func checkImageAttachment(obj interface{}) bool {
	return false
}

func readImageAttachment(is *OsgIstream, obj interface{}) {
	is.Read(&attachment1.unit)
	is.Read(&attachment1.level)
	is.Read(&attachment1.layered)
	is.Read(&attachment1.layer)
	is.Read(&attachment1.access)
	is.Read(&attachment1.format)
}

func writeImageAttachment(os *OsgOstream, obj interface{}) {
	os.Write(attachment2.unit)
	os.Write(attachment2.level)
	os.Write(attachment2.layered)
	os.Write(attachment2.layer)
	os.Write(attachment2.access)
	os.Write(attachment2.format)
}

func checkSwizzle(obj interface{}) bool {
	return true
}

func swizzleToCharacter(swizzle int, defaultCharacter byte) byte {
	switch swizzle {
	case model.GLRED:
		return 'R'
	case model.GLGREEN:
		return 'G'
	case model.GLBLUE:
		return 'B'
	case model.GLALPHA:
		return 'A'
	case model.GLZERO:
		return '0'
	case model.GLONE:
		return '1'
	}
	return defaultCharacter
}

func characterToSwizzle(character byte, defaultSwizzle int) int {
	switch character {
	case 'R':
		return model.GLRED
	case 'G':
		return model.GLGREEN
	case 'B':
		return model.GLBLUE
	case 'A':
		return model.GLALPHA
	case '0':
		return model.GLZERO
	case '1':
		return model.GLONE
	}
	return defaultSwizzle
}

func swizzleToString(swizzle [4]int) string {
	var bd strings.Builder
	bd.WriteByte(swizzleToCharacter(swizzle[0], 'R'))
	bd.WriteByte(swizzleToCharacter(swizzle[1], 'G'))
	bd.WriteByte(swizzleToCharacter(swizzle[2], 'B'))
	bd.WriteByte(swizzleToCharacter(swizzle[3], 'A'))
	return bd.String()
}

func stringToSwizzle(str string) [4]int {
	var swizzle [4]int
	swizzle[0] = characterToSwizzle(str[0], model.GLRED)
	swizzle[1] = characterToSwizzle(str[1], model.GLGREEN)
	swizzle[2] = characterToSwizzle(str[2], model.GLBLUE)
	swizzle[3] = characterToSwizzle(str[3], model.GLALPHA)

	return swizzle
}

func readSwizzle(is *OsgIstream, obj interface{}) {
	var swizzleString string
	is.Read(&swizzleString)
	tex := obj.(*model.Texture)
	tex.Swizzle = stringToSwizzle(swizzleString)
}

func writeSwizzle(os *OsgOstream, obj interface{}) {
	tex := obj.(*model.Texture)
	str := swizzleToString(tex.Swizzle)
	os.Write(&str)
	os.Write(os.CRLF)
}

func getMaxAnisotropy(obj interface{}) interface{} {
	tex := obj.(*model.Texture)
	return &tex.MaxAnisotropy
}

func setMaxAnisotropy(obj interface{}, val interface{}) {
	tex := obj.(*model.Texture)
	tex.MaxAnisotropy = val.(float32)
}

func getUseHardwareMipMapGeneration(obj interface{}) interface{} {
	tex := obj.(*model.Texture)
	return &tex.UseHardwareMipmapGeneration
}

func setUseHardwareMipMapGeneration(obj interface{}, val interface{}) {
	tex := obj.(*model.Texture)
	tex.UseHardwareMipmapGeneration = val.(bool)
}

func getUnRefImageDataAfterApply(obj interface{}) interface{} {
	tex := obj.(*model.Texture)
	return &tex.UnRefImageDataAfterApply
}

func setUnRefImageDataAfterApply(obj interface{}, val interface{}) {
	tex := obj.(*model.Texture)
	tex.UnRefImageDataAfterApply = val.(bool)
}

func getClientStorageHint(obj interface{}) interface{} {
	tex := obj.(*model.Texture)
	return &tex.ClientStorageHint
}

func setClientStorageHint(obj interface{}, val interface{}) {
	tex := obj.(*model.Texture)
	tex.ClientStorageHint = val.(bool)
}

func getResizeNonPowerOfTwoHint(obj interface{}) interface{} {
	tex := obj.(*model.Texture)
	return &tex.ResizeNonPowerOfTwoHint
}

func setResizeNonPowerOfTwoHint(obj interface{}, val interface{}) {
	tex := obj.(*model.Texture)
	tex.ResizeNonPowerOfTwoHint = val.(bool)
}

func getBorderColor(obj interface{}) interface{} {
	tex := obj.(*model.Texture)
	return &tex.BorderColor
}

func setBorderColor(obj interface{}, val interface{}) {
	tex := obj.(*model.Texture)
	tex.BorderColor = val.([4]float64)
}

func getBorderWidth(obj interface{}) interface{} {
	tex := obj.(*model.Texture)
	return &tex.BorderWidth
}

func setBorderWidth(obj interface{}, val interface{}) {
	tex := obj.(*model.Texture)
	tex.BorderWidth = val.(int)
}

func getInternalFormatMode(obj interface{}) interface{} {
	tex := obj.(*model.Texture)
	return &tex.InternalFormatMode
}

func setInternalFormatMode(obj interface{}, val interface{}) {
	tex := obj.(*model.Texture)
	tex.InternalFormatMode = val.(uint32)
}

func getShadowComparison(obj interface{}) interface{} {
	tex := obj.(*model.Texture)
	return &tex.UseShadowComparison
}

func setShadowComparison(obj interface{}, val interface{}) {
	tex := obj.(*model.Texture)
	tex.UseShadowComparison = val.(bool)
}

func getShadowCompareFunc(obj interface{}) interface{} {
	tex := obj.(*model.Texture)
	return &tex.ShadowCompareFunc
}

func setShadowCompareFunc(obj interface{}, val interface{}) {
	tex := obj.(*model.Texture)
	tex.ShadowCompareFunc = val.(int)
}

func getShadowTextureMode(obj interface{}) interface{} {
	tex := obj.(*model.Texture)
	return &tex.ShadowTextureMode
}

func setShadowTextureMode(obj interface{}, val interface{}) {
	tex := obj.(*model.Texture)
	tex.ShadowTextureMode = val.(int)
}
func getShadowAmbient(obj interface{}) interface{} {
	tex := obj.(*model.Texture)
	return &tex.ShadowAmbient
}

func setShadowAmbient(obj interface{}, val interface{}) {
	tex := obj.(*model.Texture)
	tex.ShadowAmbient = val.(float32)
}

func getMinLOD(obj interface{}) interface{} {
	tex := obj.(*model.Texture)
	return &tex.Minlod
}

func setMinLOD(obj interface{}, val interface{}) {
	tex := obj.(*model.Texture)
	tex.Minlod = val.(float32)
}

func getMaxLOD(obj interface{}) interface{} {
	tex := obj.(*model.Texture)
	return &tex.Maxlod
}

func setMaxLOD(obj interface{}, val interface{}) {
	tex := obj.(*model.Texture)
	tex.Maxlod = val.(float32)
}

func getLODBias(obj interface{}) interface{} {
	tex := obj.(*model.Texture)
	return &tex.Lodbias
}

func setLODBias(obj interface{}, val interface{}) {
	tex := obj.(*model.Texture)
	tex.Lodbias = val.(float32)
}

func init() {
	wrap := NewObjectWrapper("Texture", nil, "osg::Object osg::StateAttribute osg::Texture")
	ser1 := NewUserSerializer("WRAP_S", checkWRAPS, readWRAPS, writeWRAPS)
	ser2 := NewUserSerializer("WRAP_T", checkWRAPT, readWRAPT, writeWRAPT)
	ser3 := NewUserSerializer("WRAP_R", checkWRAPR, readWRAPR, writeWRAPR)
	ser4 := NewUserSerializer("MIN_FILTER", checkMINFILTER, readMINFILTER, writeMINFILTER)
	ser5 := NewUserSerializer("MAG_FILTER", checkMAGFILTER, readMAGFILTER, writeMAGFILTER)
	ser6 := NewPropByValSerializer("MaxAnisotropy", false, getMaxAnisotropy, setMaxAnisotropy)
	ser7 := NewPropByValSerializer("UseHardwareMipmapGeneration", false, getUseHardwareMipMapGeneration, setUseHardwareMipMapGeneration)
	ser8 := NewPropByValSerializer("UnRefImageDataAfterApply", false, getUnRefImageDataAfterApply, setUnRefImageDataAfterApply)
	ser9 := NewPropByValSerializer("ClientStorageHint", false, getClientStorageHint, setClientStorageHint)
	ser10 := NewPropByValSerializer("ResizeNonPowerOfTwoHint", false, getResizeNonPowerOfTwoHint, setResizeNonPowerOfTwoHint)
	// tydata := [4]float64{}
	ser11 := NewPropByValSerializer("BorderColor", false, getBorderColor, setBorderColor)
	ser12 := NewPropByValSerializer("BorderWidth", false, getBorderWidth, setBorderWidth)
	ser13 := NewEnumSerializer("InternalFormatMode", getInternalFormatMode, setInternalFormatMode)
	ser13.Add("USEIMAGEDATAFORMAT", model.USEIMAGEDATAFORMAT)
	ser13.Add("USEUSERDEFINEDFORMAT", model.USEUSERDEFINEDFORMAT)
	ser13.Add("USEARBCOMPRESSION", model.USEARBCOMPRESSION)
	ser13.Add("USES3TCDXT1COMPRESSION", model.USES3TCDXT1COMPRESSION)
	ser13.Add("USES3TCDXT3COMPRESSION", model.USES3TCDXT3COMPRESSION)
	ser13.Add("USES3TCDXT5COMPRESSION", model.USES3TCDXT5COMPRESSION)
	ser13.Add("USEPVRTC2BPPCOMPRESSION", model.USEPVRTC2BPPCOMPRESSION)
	ser13.Add("USEPVRTC4BPPCOMPRESSION", model.USEPVRTC4BPPCOMPRESSION)
	ser13.Add("USEETCCOMPRESSION", model.USEETCCOMPRESSION)
	ser13.Add("USERGTC1COMPRESSION", model.USERGTC1COMPRESSION)
	ser13.Add("USERGTC2COMPRESSION", model.USERGTC2COMPRESSION)
	ser13.Add("USES3TCDXT1cCOMPRESSION", model.USES3TCDXT1cCOMPRESSION)
	ser13.Add("USES3TCDXT1aCOMPRESSION", model.USES3TCDXT1aCOMPRESSION)

	ser14 := NewUserSerializer("InternalFormat", checkInternalFormat, readInternalFormat, writeInternalFormat)
	ser15 := NewUserSerializer("SourceFormat", checkSourceFormat, readSourceFormat, writeSourceFormat)
	ser16 := NewUserSerializer("SourceType", checkSourceType, readSourceType, writeSourceType)
	ser17 := NewPropByValSerializer("ShadowComparison", false, getShadowComparison, setShadowComparison)

	ser18 := NewEnumSerializer("ShadowCompareFunc", getShadowCompareFunc, setShadowCompareFunc)
	ser18.Add("NEVER", model.GLNEVER)
	ser18.Add("LESS", model.GLLESS)
	ser18.Add("EQUAL", model.GLEQUAL)
	ser18.Add("LEQUAL", model.GLLEQUAL)
	ser18.Add("GREATER", model.GLGREATER)
	ser18.Add("NOTEQUAL", model.GLNOTEQUAL)
	ser18.Add("GEQUAL", model.GLGEQUAL)
	ser18.Add("ALWAYS", model.GLALWAYS)

	ser19 := NewEnumSerializer("ShadowTextureMode", getShadowTextureMode, setShadowTextureMode)
	ser19.Add("LUMINANCE", model.GLLUMINANCE)
	ser19.Add("INTENSITY", model.GLINTENSITY)
	ser19.Add("ALPHA", model.GLALPHA)
	ser19.Add("NONE", model.NONE)

	ser20 := NewPropByValSerializer("ShadowAmbient", false, getShadowAmbient, setShadowAmbient)

	wrap.AddSerializer(ser1, RWUSER)
	wrap.AddSerializer(ser2, RWUSER)
	wrap.AddSerializer(ser3, RWUSER)
	wrap.AddSerializer(ser4, RWUSER)
	wrap.AddSerializer(ser5, RWUSER)
	wrap.AddSerializer(ser6, RWFLOAT)
	wrap.AddSerializer(ser7, RWBOOL)
	wrap.AddSerializer(ser8, RWBOOL)
	wrap.AddSerializer(ser9, RWBOOL)
	wrap.AddSerializer(ser10, RWBOOL)
	wrap.AddSerializer(ser11, RWVEC4D)
	wrap.AddSerializer(ser12, RWINT)
	wrap.AddSerializer(ser13, RWENUM)
	wrap.AddSerializer(ser14, RWUSER)
	wrap.AddSerializer(ser15, RWUSER)
	wrap.AddSerializer(ser16, RWUSER)
	wrap.AddSerializer(ser17, RWBOOL)
	wrap.AddSerializer(ser18, RWENUM)
	wrap.AddSerializer(ser19, RWENUM)
	wrap.AddSerializer(ser20, RWFLOAT)
	{
		uv := AddUpdateWrapperVersionProxy(wrap, 95)
		ser21 := NewUserSerializer("ImageAttachment", checkImageAttachment, readImageAttachment, writeImageAttachment)
		wrap.AddSerializer(ser21, RWUSER)
		uv.SetLastVersion()
	}

	{
		uv := AddUpdateWrapperVersionProxy(wrap, 154)
		wrap.MarkSerializerAsRemoved("ImageAttachment")
		uv.SetLastVersion()
	}

	{
		uv := AddUpdateWrapperVersionProxy(wrap, 98)
		ser22 := NewUserSerializer("Swizzle", checkSwizzle, readSwizzle, writeSwizzle)
		wrap.AddSerializer(ser22, RWUSER)
		uv.SetLastVersion()
	}

	{
		uv := AddUpdateWrapperVersionProxy(wrap, 155)
		ser23 := NewPropByValSerializer("MinLOD", false, getMinLOD, setMinLOD)
		ser24 := NewPropByValSerializer("MaxLOD", false, getMaxLOD, setMaxLOD)
		ser25 := NewPropByValSerializer("LODBias", false, getLODBias, setLODBias)
		wrap.AddSerializer(ser23, RWFLOAT)
		wrap.AddSerializer(ser24, RWFLOAT)
		wrap.AddSerializer(ser25, RWFLOAT)
		uv.SetLastVersion()
	}

	GetObjectWrapperManager().AddWrap(wrap)
}
