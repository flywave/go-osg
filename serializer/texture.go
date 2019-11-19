package serializer

import (
	"strings"

	"github.com/flywave/go-osg"
	"github.com/flywave/go-osg/model"
)

func checkWRAPS(obj interface{}) bool {
	return true
}

func readWRAPS(is *osg.OsgIstream, obj interface{}) {
	tex := obj.(*model.Texture)
	var mode int
	is.Read(&mode)
	tex.SetWrap(model.WRAPS, mode)
}

func writeWRAPS(os *osg.OsgOstream, obj interface{}) {
	tex := obj.(*model.Texture)
	os.Write(tex.GetWrap(model.WRAPS))
	os.Write(os.CRLF)
}

func checkWRAPT(obj interface{}) bool {
	return true
}

func readWRAPT(is *osg.OsgIstream, obj interface{}) {
	tex := obj.(*model.Texture)
	var mode int
	is.Read(&mode)
	tex.SetWrap(model.WRAPT, mode)
}

func writeWRAPT(os *osg.OsgOstream, obj interface{}) {
	tex := obj.(*model.Texture)
	os.Write(tex.GetWrap(model.WRAPT))
	os.Write(os.CRLF)
}

func checkWRAPR(obj interface{}) bool {
	return true
}

func readWRAPR(is *osg.OsgIstream, obj interface{}) {
	tex := obj.(*model.Texture)
	var mode int
	is.Read(&mode)
	tex.SetWrap(model.WRAPR, mode)
}

func writeWRAPR(os *osg.OsgOstream, obj interface{}) {
	tex := obj.(*model.Texture)
	os.Write(tex.GetWrap(model.WRAPR))
	os.Write(os.CRLF)
}

func checkMINFILTER(obj interface{}) bool {
	return true
}

func readMINFILTER(is *osg.OsgIstream, obj interface{}) {
	tex := obj.(*model.Texture)
	var mode int
	tex.SetFilter(model.MINFILTER, mode)
}

func writeMINFILTER(os *osg.OsgOstream, obj interface{}) {
	tex := obj.(*model.Texture)
	tex.GetFilter(model.MINFILTER)
	os.Write(os.CRLF)
}

func checkMAGFILTER(obj interface{}) bool {
	return true
}

func readMAGFILTER(is *osg.OsgIstream, obj interface{}) {
	tex := obj.(*model.Texture)
	var mode int
	tex.SetFilter(model.MAGFILTER, mode)
}

func writeMAGFILTER(os *osg.OsgOstream, obj interface{}) {
	tex := obj.(*model.Texture)
	tex.GetFilter(model.MAGFILTER)
	os.Write(os.CRLF)
}

func checkSourceFormat(obj interface{}) bool {
	tex := obj.(*model.Texture)
	return tex.SourceFormat != model.GLZERO
}

func readSourceFormat(is *osg.OsgIstream, obj interface{}) {
	tex := obj.(*model.Texture)
	var mode int
	is.Read(mode)
	tex.SourceFormat = mode
}
func writeSourceFormat(os *osg.OsgOstream, obj interface{}) {
	tex := obj.(*model.Texture)
	os.Write(tex.SourceFormat)
	os.Write(os.CRLF)
}

func checkSourceType(obj interface{}) bool {
	tex := obj.(*model.Texture)
	return tex.SourceType != model.GLZERO
}

func readSourceType(is *osg.OsgIstream, obj interface{}) {
	tex := obj.(*model.Texture)
	var mode int
	is.Read(mode)
	tex.SourceType = mode
}
func writeSourceType(os *osg.OsgOstream, obj interface{}) {
	tex := obj.(*model.Texture)
	os.Write(tex.SourceType)
	os.Write(os.CRLF)
}

func checkInternalFormat(obj interface{}) bool {
	tex := obj.(*model.Texture)
	return tex.InternalFormat == model.USEUSERDEFINEDFORMAT
}

func readInternalFormat(is *osg.OsgIstream, obj interface{}) {
	tex := obj.(*model.Texture)
	var mode int
	is.Read(&mode)
	if tex.InternalFormat == model.USEUSERDEFINEDFORMAT {
		tex.InternalFormat = mode
	}
}

func writeInternalFormat(os *osg.OsgOstream, obj interface{}) {
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

func readImageAttachment(is *osg.OsgIstream, obj interface{}) {
	is.Read(attachment1.unit)
	is.Read(attachment1.level)
	is.Read(attachment1.layered)
	is.Read(attachment1.layer)
	is.Read(attachment1.access)
	is.Read(attachment1.format)
}

func writeImageAttachment(os *osg.OsgOstream, obj interface{}) {
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

func readSwizzle(is *osg.OsgIstream, obj interface{}) {
	var swizzleString string
	is.Read(&swizzleString)
	tex := obj.(*model.Texture)
	tex.Swizzle = stringToSwizzle(swizzleString)
}

func writeSwizzle(os *osg.OsgOstream, obj interface{}) {
	tex := obj.(*model.Texture)
	os.Write(swizzleToString(tex.Swizzle))
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
	tex.InternalFormatMode = val.(model.InternalFormatMode)
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
	wrap := osg.NewObjectWrapper("Texture", nil, "osg::Object osg::StateAttribute osg::Texture")
	ser1 := osg.NewUserSerializer("WRAPS", checkWRAPS, readWRAPS, writeWRAPS)
	ser2 := osg.NewUserSerializer("WRAPT", checkWRAPT, readWRAPT, writeWRAPT)
	ser3 := osg.NewUserSerializer("WRAPR", checkWRAPR, readWRAPR, writeWRAPR)
	ser4 := osg.NewUserSerializer("MINFILTER", checkMINFILTER, readMINFILTER, writeMINFILTER)
	ser5 := osg.NewUserSerializer("MAGFILTER", checkMAGFILTER, readMAGFILTER, writeMAGFILTER)
	ser6 := osg.NewPropByValSerializer("MaxAnisotropy", false, getMaxAnisotropy, setMaxAnisotropy)
	ser7 := osg.NewPropByValSerializer("UseHardwareMipmapGeneration", false, getUseHardwareMipMapGeneration, setUseHardwareMipMapGeneration)
	ser8 := osg.NewPropByValSerializer("UnRefImageDataAfterApply", false, getUnRefImageDataAfterApply, setUnRefImageDataAfterApply)
	ser9 := osg.NewPropByValSerializer("ClientStorageHint", false, getClientStorageHint, setClientStorageHint)
	ser10 := osg.NewPropByValSerializer("ResizeNonPowerOfTwoHint", false, getResizeNonPowerOfTwoHint, setResizeNonPowerOfTwoHint)
	var tydata float64
	ser11 := osg.NewVectorSerializer("BorderColor", osg.RWDOUBLE, &tydata, getBorderColor, setBorderColor)
	ser12 := osg.NewPropByValSerializer("BorderWidth", false, getBorderWidth, setBorderWidth)
	ser13 := osg.NewEnumSerializer("InternalFormatMode", getInternalFormatMode, setInternalFormatMode)
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

	ser14 := osg.NewUserSerializer("InternalFormat", checkInternalFormat, readInternalFormat, writeInternalFormat)
	ser15 := osg.NewUserSerializer("SourceFormat", checkSourceFormat, readSourceFormat, writeSourceFormat)
	ser16 := osg.NewUserSerializer("SourceType", checkSourceType, readSourceType, writeSourceType)
	ser17 := osg.NewPropByValSerializer("ShadowComparison", false, getShadowComparison, setShadowComparison)

	ser18 := osg.NewEnumSerializer("ShadowCompareFunc", getShadowCompareFunc, setShadowCompareFunc)
	ser18.Add("NEVER", model.GLNEVER)
	ser18.Add("LESS", model.GLLESS)
	ser18.Add("EQUAL", model.GLEQUAL)
	ser18.Add("LEQUAL", model.GLLEQUAL)
	ser18.Add("GREATER", model.GLGREATER)
	ser18.Add("NOTEQUAL", model.GLNOTEQUAL)
	ser18.Add("GEQUAL", model.GLGEQUAL)
	ser18.Add("ALWAYS", model.GLALWAYS)

	ser19 := osg.NewEnumSerializer("ShadowTextureMode", getShadowTextureMode, setShadowTextureMode)
	ser19.Add("LUMINANCE", model.GLLUMINANCE)
	ser19.Add("INTENSITY", model.GLINTENSITY)
	ser19.Add("ALPHA", model.GLALPHA)
	ser19.Add("NONE", model.NONE)

	ser20 := osg.NewPropByValSerializer("ShadowAmbient", false, getShadowAmbient, setShadowAmbient)
	osg.AddUpdateWrapperVersionProxy(&wrap, 95)
	ser21 := osg.NewUserSerializer("ImageAttachment", checkImageAttachment, readImageAttachment, writeImageAttachment)

	osg.AddUpdateWrapperVersionProxy(&wrap, 98)
	wrap.MarkSerializerAsRemoved("ImageAttachment")

	osg.AddUpdateWrapperVersionProxy(&wrap, 98)
	ser22 := osg.NewUserSerializer("Swizzle", checkSwizzle, readSwizzle, writeSwizzle)

	osg.AddUpdateWrapperVersionProxy(&wrap, 155)
	ser23 := osg.NewPropByValSerializer("MinLOD", false, getMinLOD, setMinLOD)
	ser24 := osg.NewPropByValSerializer("MaxLOD", false, getMaxLOD, setMaxLOD)
	ser25 := osg.NewPropByValSerializer("LODBias", false, getLODBias, setLODBias)

	wrap.AddSerializer(&ser1, osg.RWUSER)
	wrap.AddSerializer(&ser2, osg.RWUSER)
	wrap.AddSerializer(&ser3, osg.RWUSER)
	wrap.AddSerializer(&ser4, osg.RWUSER)
	wrap.AddSerializer(&ser5, osg.RWUSER)
	wrap.AddSerializer(&ser6, osg.RWFLOAT)
	wrap.AddSerializer(&ser7, osg.RWBOOL)
	wrap.AddSerializer(&ser8, osg.RWBOOL)
	wrap.AddSerializer(&ser9, osg.RWBOOL)
	wrap.AddSerializer(&ser10, osg.RWBOOL)
	wrap.AddSerializer(&ser11, osg.RWVEC4D)
	wrap.AddSerializer(&ser12, osg.RWINT)
	wrap.AddSerializer(&ser13, osg.RWENUM)
	wrap.AddSerializer(&ser14, osg.RWUSER)
	wrap.AddSerializer(&ser15, osg.RWUSER)
	wrap.AddSerializer(&ser16, osg.RWUSER)
	wrap.AddSerializer(&ser17, osg.RWBOOL)
	wrap.AddSerializer(&ser18, osg.RWENUM)
	wrap.AddSerializer(&ser19, osg.RWENUM)
	wrap.AddSerializer(&ser20, osg.RWFLOAT)
	wrap.AddSerializer(&ser21, osg.RWUSER)
	wrap.AddSerializer(&ser22, osg.RWUSER)
	wrap.AddSerializer(&ser23, osg.RWFLOAT)
	wrap.AddSerializer(&ser24, osg.RWFLOAT)
	wrap.AddSerializer(&ser25, osg.RWFLOAT)
	osg.GetObjectWrapperManager().AddWrap(&wrap)
}
