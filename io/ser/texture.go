package ser

import (
	"strings"

	"github.com/flywave/go-osg/io"
	"github.com/flywave/go-osg/model"
)

func checkWRAPS(obj interface{}) bool {
	return true
}

func readWRAPS(is *io.OsgIstream, obj interface{}) {
	tex := obj.(*model.Texture)
	var mode int
	is.Read(&mode)
	tex.SetWrap(model.WRAP_S, mode)
}

func writeWRAPS(os *io.OsgOstream, obj interface{}) {
	tex := obj.(*model.Texture)
	os.Write(tex.GetWrap(model.WRAP_S))
	os.Write(os.CRLF)
}

func checkWRAPT(obj interface{}) bool {
	return true
}

func readWRAPT(is *io.OsgIstream, obj interface{}) {
	tex := obj.(*model.Texture)
	var mode int
	is.Read(&mode)
	tex.SetWrap(model.WRAP_T, mode)
}

func writeWRAPT(os *io.OsgOstream, obj interface{}) {
	tex := obj.(*model.Texture)
	os.Write(tex.GetWrap(model.WRAP_T))
	os.Write(os.CRLF)
}

func checkWRAPR(obj interface{}) bool {
	return true
}

func readWRAPR(is *io.OsgIstream, obj interface{}) {
	tex := obj.(*model.Texture)
	var mode int
	is.Read(&mode)
	tex.SetWrap(model.WRAP_R, mode)
}

func writeWRAPR(os *io.OsgOstream, obj interface{}) {
	tex := obj.(*model.Texture)
	os.Write(tex.GetWrap(model.WRAP_R))
	os.Write(os.CRLF)
}

func checkMINFILTER(obj interface{}) bool {
	return true
}

func readMINFILTER(is *io.OsgIstream, obj interface{}) {
	tex := obj.(*model.Texture)
	var mode int
	tex.SetFilter(model.MIN_FILTER, mode)
}

func writeMINFILTER(os *io.OsgOstream, obj interface{}) {
	tex := obj.(*model.Texture)
	tex.GetFilter(model.MIN_FILTER)
	os.Write(os.CRLF)
}

func checkMAGFILTER(obj interface{}) bool {
	return true
}

func readMAGFILTER(is *io.OsgIstream, obj interface{}) {
	tex := obj.(*model.Texture)
	var mode int
	tex.SetFilter(model.MAG_FILTER, mode)
}

func writeMAGFILTER(os *io.OsgOstream, obj interface{}) {
	tex := obj.(*model.Texture)
	tex.GetFilter(model.MAG_FILTER)
	os.Write(os.CRLF)
}

func checkSourceFormat(obj interface{}) bool {
	tex := obj.(*model.Texture)
	return tex.SourceFormat != model.GL_ZERO
}

func readSourceFormat(is *io.OsgIstream, obj interface{}) {
	tex := obj.(*model.Texture)
	var mode int
	is.Read(mode)
	tex.SourceFormat = mode
}
func writeSourceFormat(os *io.OsgOstream, obj interface{}) {
	tex := obj.(*model.Texture)
	os.Write(tex.SourceFormat)
	os.Write(os.CRLF)
}

func checkSourceType(obj interface{}) bool {
	tex := obj.(*model.Texture)
	return tex.SourceType != model.GL_ZERO
}

func readSourceType(is *io.OsgIstream, obj interface{}) {
	tex := obj.(*model.Texture)
	var mode int
	is.Read(mode)
	tex.SourceType = mode
}
func writeSourceType(os *io.OsgOstream, obj interface{}) {
	tex := obj.(*model.Texture)
	os.Write(tex.SourceType)
	os.Write(os.CRLF)
}

func checkInternalFormat(obj interface{}) bool {
	tex := obj.(*model.Texture)
	return tex.InternalFormat == model.USE_USER_DEFINED_FORMAT
}

func readInternalFormat(is *io.OsgIstream, obj interface{}) {
	tex := obj.(*model.Texture)
	var mode int
	is.Read(&mode)
	if tex.InternalFormat == model.USE_USER_DEFINED_FORMAT {
		tex.InternalFormat = mode
	}
}

func writeInternalFormat(os *io.OsgOstream, obj interface{}) {
	tex := obj.(*model.Texture)
	if os.IsBinary() && tex.InternalFormatMode != model.USE_USER_DEFINED_FORMAT {
		os.Write(model.GL_ZERO)
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

func readImageAttachment(is *io.OsgIstream, obj interface{}) {
	is.Read(attachment1.unit)
	is.Read(attachment1.level)
	is.Read(attachment1.layered)
	is.Read(attachment1.layer)
	is.Read(attachment1.access)
	is.Read(attachment1.format)
}

func writeImageAttachment(os *io.OsgOstream, obj interface{}) {
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
	case model.GL_RED:
		return 'R'
	case model.GL_GREEN:
		return 'G'
	case model.GL_BLUE:
		return 'B'
	case model.GL_ALPHA:
		return 'A'
	case model.GL_ZERO:
		return '0'
	case model.GL_ONE:
		return '1'
	}
	return defaultCharacter
}

func characterToSwizzle(character byte, defaultSwizzle int) int {
	switch character {
	case 'R':
		return model.GL_RED
	case 'G':
		return model.GL_GREEN
	case 'B':
		return model.GL_BLUE
	case 'A':
		return model.GL_ALPHA
	case '0':
		return model.GL_ZERO
	case '1':
		return model.GL_ONE
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
	swizzle[0] = characterToSwizzle(str[0], model.GL_RED)
	swizzle[1] = characterToSwizzle(str[1], model.GL_GREEN)
	swizzle[2] = characterToSwizzle(str[2], model.GL_BLUE)
	swizzle[3] = characterToSwizzle(str[3], model.GL_ALPHA)

	return swizzle
}

func readSwizzle(is *io.OsgIstream, obj interface{}) {
	var swizzleString string
	is.Read(swizzleString)
	tex := obj.(*model.Texture)
	tex.Swizzle = stringToSwizzle(swizzleString)
}

func writeSwizzle(os *io.OsgOstream, obj interface{}) {
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
	wrap := io.NewObjectWrapper("Texture", nil, "osg::Object osg::StateAttribute osg::Texture")
	ser1 := io.NewUserSerializer("WRAPS", checkWRAPS, readWRAPS, writeWRAPS)
	ser2 := io.NewUserSerializer("WRAPT", checkWRAPT, readWRAPT, writeWRAPT)
	ser3 := io.NewUserSerializer("WRAPR", checkWRAPR, readWRAPR, writeWRAPR)
	ser4 := io.NewUserSerializer("MINFILTER", checkMINFILTER, readMINFILTER, writeMINFILTER)
	ser5 := io.NewUserSerializer("MAGFILTER", checkMAGFILTER, readMAGFILTER, writeMAGFILTER)
	ser6 := io.NewPropByValSerializer("MaxAnisotropy", false, getMaxAnisotropy, setMaxAnisotropy)
	ser7 := io.NewPropByValSerializer("UseHardwareMipmapGeneration", false, getUseHardwareMipMapGeneration, setUseHardwareMipMapGeneration)
	ser8 := io.NewPropByValSerializer("UnRefImageDataAfterApply", false, getUnRefImageDataAfterApply, setUnRefImageDataAfterApply)
	ser9 := io.NewPropByValSerializer("ClientStorageHint", false, getClientStorageHint, setClientStorageHint)
	ser10 := io.NewPropByValSerializer("ResizeNonPowerOfTwoHint", false, getResizeNonPowerOfTwoHint, setResizeNonPowerOfTwoHint)
	ser11 := io.NewVectorSerializer("BorderColor", io.RW_DOUBLE, 4, getBorderColor, setBorderColor)
	ser12 := io.NewPropByValSerializer("BorderWidth", false, getBorderWidth, setBorderWidth)
	ser13 := io.NewEnumSerializer("InternalFormatMode", getInternalFormatMode, setInternalFormatMode)
	ser13.Add("USE_IMAGE_DATA_FORMAT", model.USE_IMAGE_DATA_FORMAT)
	ser13.Add("USE_USER_DEFINED_FORMAT", model.USE_USER_DEFINED_FORMAT)
	ser13.Add("USE_ARB_COMPRESSION", model.USE_ARB_COMPRESSION)
	ser13.Add("USE_S3TC_DXT1_COMPRESSION", model.USE_S3TC_DXT1_COMPRESSION)
	ser13.Add("USE_S3TC_DXT3_COMPRESSION", model.USE_S3TC_DXT3_COMPRESSION)
	ser13.Add("USE_S3TC_DXT5_COMPRESSION", model.USE_S3TC_DXT5_COMPRESSION)
	ser13.Add("USE_PVRTC_2BPP_COMPRESSION", model.USE_PVRTC_2BPP_COMPRESSION)
	ser13.Add("USE_PVRTC_4BPP_COMPRESSION", model.USE_PVRTC_4BPP_COMPRESSION)
	ser13.Add("USE_ETC_COMPRESSION", model.USE_ETC_COMPRESSION)
	ser13.Add("USE_RGTC1_COMPRESSION", model.USE_RGTC1_COMPRESSION)
	ser13.Add("USE_RGTC2_COMPRESSION", model.USE_RGTC2_COMPRESSION)
	ser13.Add("USE_S3TC_DXT1c_COMPRESSION", model.USE_S3TC_DXT1c_COMPRESSION)
	ser13.Add("USE_S3TC_DXT1a_COMPRESSION", model.USE_S3TC_DXT1a_COMPRESSION)

	ser14 := io.NewUserSerializer("InternalFormat", checkInternalFormat, readInternalFormat, writeInternalFormat)
	ser15 := io.NewUserSerializer("SourceFormat", checkSourceFormat, readSourceFormat, writeSourceFormat)
	ser16 := io.NewUserSerializer("SourceType", checkSourceType, readSourceType, writeSourceType)
	ser17 := io.NewPropByValSerializer("ShadowComparison", false, getShadowComparison, setShadowComparison)

	ser18 := io.NewEnumSerializer("ShadowCompareFunc", getShadowCompareFunc, setShadowCompareFunc)
	ser18.Add("NEVER", model.GL_NEVER)
	ser18.Add("LESS", model.GL_LESS)
	ser18.Add("EQUAL", model.GL_EQUAL)
	ser18.Add("LEQUAL", model.GL_LEQUAL)
	ser18.Add("GREATER", model.GL_GREATER)
	ser18.Add("NOTEQUAL", model.GL_NOTEQUAL)
	ser18.Add("GEQUAL", model.GL_GEQUAL)
	ser18.Add("ALWAYS", model.GL_ALWAYS)

	ser19 := io.NewEnumSerializer("ShadowTextureMode", getShadowTextureMode, setShadowTextureMode)
	ser19.Add("LUMINANCE", model.GL_LUMINANCE)
	ser19.Add("INTENSITY", model.GL_INTENSITY)
	ser19.Add("ALPHA", model.GL_ALPHA)
	ser19.Add("NONE", model.NONE)

	ser20 := io.NewPropByValSerializer("ShadowAmbient", false, getShadowAmbient, setShadowAmbient)
	io.AddUpdateWrapperVersionProxy(&wrap, 95)
	ser21 := io.NewUserSerializer("ImageAttachment", checkImageAttachment, readImageAttachment, writeImageAttachment)

	io.AddUpdateWrapperVersionProxy(&wrap, 98)
	wrap.MarkSerializerAsRemoved("ImageAttachment")

	io.AddUpdateWrapperVersionProxy(&wrap, 98)
	ser22 := io.NewUserSerializer("Swizzle", checkSwizzle, readSwizzle, writeSwizzle)

	io.AddUpdateWrapperVersionProxy(&wrap, 155)
	ser23 := io.NewPropByValSerializer("MinLOD", false, getMinLOD, setMinLOD)
	ser24 := io.NewPropByValSerializer("MaxLOD", false, getMaxLOD, setMaxLOD)
	ser25 := io.NewPropByValSerializer("LODBias", false, getLODBias, setLODBias)

	wrap.AddSerializer(&ser1, io.RW_USER)
	wrap.AddSerializer(&ser2, io.RW_USER)
	wrap.AddSerializer(&ser3, io.RW_USER)
	wrap.AddSerializer(&ser4, io.RW_USER)
	wrap.AddSerializer(&ser5, io.RW_USER)
	wrap.AddSerializer(&ser6, io.RW_FLOAT)
	wrap.AddSerializer(&ser7, io.RW_BOOL)
	wrap.AddSerializer(&ser8, io.RW_BOOL)
	wrap.AddSerializer(&ser9, io.RW_BOOL)
	wrap.AddSerializer(&ser10, io.RW_BOOL)
	wrap.AddSerializer(&ser11, io.RW_VEC4D)
	wrap.AddSerializer(&ser12, io.RW_INT)
	wrap.AddSerializer(&ser13, io.RW_ENUM)
	wrap.AddSerializer(&ser14, io.RW_USER)
	wrap.AddSerializer(&ser15, io.RW_USER)
	wrap.AddSerializer(&ser16, io.RW_USER)
	wrap.AddSerializer(&ser17, io.RW_BOOL)
	wrap.AddSerializer(&ser18, io.RW_ENUM)
	wrap.AddSerializer(&ser19, io.RW_ENUM)
	wrap.AddSerializer(&ser20, io.RW_FLOAT)
	wrap.AddSerializer(&ser21, io.RW_USER)
	wrap.AddSerializer(&ser22, io.RW_USER)
	wrap.AddSerializer(&ser23, io.RW_FLOAT)
	wrap.AddSerializer(&ser24, io.RW_FLOAT)
	wrap.AddSerializer(&ser25, io.RW_FLOAT)
	io.GetObjectWrapperManager().AddWrap(&wrap)
}
