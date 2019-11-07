package model

import (
	"math"

	"github.com/ungerik/go3d/vec3"
)

type WriteHint uint32
type AllocationMode uint32
type Origin uint32

const (
	NO_PREFERENCE WriteHint = 0
	STORE_INLINE  WriteHint = 1
	EXTERNAL_FILE WriteHint = 2

	NO_DELETE       AllocationMode = 0
	USE_NEW_DELETE  AllocationMode = 1
	USE_MALLOC_FREE AllocationMode = 2

	BOTTOM_LEFT Origin = 0
	TOP_LEFT    Origin = 1
	IMAGE_T     string = "osg::Image"
)

type Image struct {
	BufferData
	FileName              string
	WriteHint             WriteHint
	Origin                Origin
	S                     int
	T                     int
	R                     int
	RowLength             int
	InternalTextureFormat int
	PixelFormat           int
	DataType              int
	Packing               uint
	PixelAspectRatio      float32

	AllocationMode AllocationMode
	Data           []uint8
}

func NewImage() Image {
	b := NewBufferData()
	b.Type = IMAGE_T
	return Image{BufferData: b, S: 0, T: 0, R: 0, RowLength: 0, InternalTextureFormat: 0, PixelFormat: 0, DataType: 0, Packing: 4, PixelAspectRatio: 1, AllocationMode: USE_NEW_DELETE}
}

func IsPackedType(t int) bool {
	switch t {
	case (GL_UNSIGNED_BYTE_3_3_2):
	case (GL_UNSIGNED_BYTE_2_3_3_REV):
	case (GL_UNSIGNED_SHORT_5_6_5):
	case (GL_UNSIGNED_SHORT_5_6_5_REV):
	case (GL_UNSIGNED_SHORT_4_4_4_4):
	case (GL_UNSIGNED_SHORT_4_4_4_4_REV):
	case (GL_UNSIGNED_SHORT_5_5_5_1):
	case (GL_UNSIGNED_SHORT_1_5_5_5_REV):
	case (GL_UNSIGNED_INT_8_8_8_8):
	case (GL_UNSIGNED_INT_8_8_8_8_REV):
	case (GL_UNSIGNED_INT_10_10_10_2):
	case (GL_UNSIGNED_INT_2_10_10_10_REV):
		return true
	}
	return false
}

func ComputePixelFormat(formt int) int {
	switch formt {
	case (GL_ALPHA16F_ARB):
	case (GL_ALPHA32F_ARB):
		return GL_ALPHA

	case (GL_LUMINANCE16F_ARB):
	case (GL_LUMINANCE32F_ARB):
		return GL_LUMINANCE

	case (GL_INTENSITY16F_ARB):
	case (GL_INTENSITY32F_ARB):
		return GL_INTENSITY

	case (GL_LUMINANCE_ALPHA16F_ARB):
	case (GL_LUMINANCE_ALPHA32F_ARB):
		return GL_LUMINANCE_ALPHA

	case (GL_RGB32F_ARB):
	case (GL_RGB16F_ARB):
		return GL_RGB

	case (GL_RGBA8):
	case (GL_RGBA16):
	case (GL_RGBA32F_ARB):
	case (GL_RGBA16F_ARB):
		return GL_RGBA

	case (GL_ALPHA8I_EXT):
	case (GL_ALPHA16I_EXT):
	case (GL_ALPHA32I_EXT):
	case (GL_ALPHA8UI_EXT):
	case (GL_ALPHA16UI_EXT):
	case (GL_ALPHA32UI_EXT):
		return GL_ALPHA_INTEGER_EXT

	case (GL_LUMINANCE8I_EXT):
	case (GL_LUMINANCE16I_EXT):
	case (GL_LUMINANCE32I_EXT):
	case (GL_LUMINANCE8UI_EXT):
	case (GL_LUMINANCE16UI_EXT):
	case (GL_LUMINANCE32UI_EXT):
		return GL_LUMINANCE_INTEGER_EXT

	case (GL_INTENSITY8I_EXT):
	case (GL_INTENSITY16I_EXT):
	case (GL_INTENSITY32I_EXT):
	case (GL_INTENSITY8UI_EXT):
	case (GL_INTENSITY16UI_EXT):
	case (GL_INTENSITY32UI_EXT):
		return GL_LUMINANCE_INTEGER_EXT

	case (GL_LUMINANCE_ALPHA8I_EXT):
	case (GL_LUMINANCE_ALPHA16I_EXT):
	case (GL_LUMINANCE_ALPHA32I_EXT):
	case (GL_LUMINANCE_ALPHA8UI_EXT):
	case (GL_LUMINANCE_ALPHA16UI_EXT):
	case (GL_LUMINANCE_ALPHA32UI_EXT):
		return GL_LUMINANCE_ALPHA_INTEGER_EXT

	case (GL_RGB32I_EXT):
	case (GL_RGB16I_EXT):
	case (GL_RGB8I_EXT):
	case (GL_RGB32UI_EXT):
	case (GL_RGB16UI_EXT):
	case (GL_RGB8UI_EXT):
		return GL_RGB_INTEGER_EXT

	case (GL_RGBA32I_EXT):
	case (GL_RGBA16I_EXT):
	case (GL_RGBA8I_EXT):
	case (GL_RGBA32UI_EXT):
	case (GL_RGBA16UI_EXT):
	case (GL_RGBA8UI_EXT):
		return GL_RGBA_INTEGER_EXT

	case (GL_DEPTH_COMPONENT16):
	case (GL_DEPTH_COMPONENT24):
	case (GL_DEPTH_COMPONENT32):
	case (GL_DEPTH_COMPONENT32F):
	case (GL_DEPTH_COMPONENT32F_NV):
		return GL_DEPTH_COMPONENT
	}
	return formt
}

func ComputeNumComponents(pixelFormat int) uint {
	switch pixelFormat {
	case (GL_COMPRESSED_RGB_S3TC_DXT1_EXT):
		return 3
	case (GL_COMPRESSED_RGBA_S3TC_DXT1_EXT):
		return 4
	case (GL_COMPRESSED_RGBA_S3TC_DXT3_EXT):
		return 4
	case (GL_COMPRESSED_RGBA_S3TC_DXT5_EXT):
		return 4
	case (GL_COMPRESSED_SIGNED_RED_RGTC1_EXT):
		return 1
	case (GL_COMPRESSED_RED_RGTC1_EXT):
		return 1
	case (GL_COMPRESSED_SIGNED_RED_GREEN_RGTC2_EXT):
		return 2
	case (GL_COMPRESSED_RED_GREEN_RGTC2_EXT):
		return 2
	case (GL_COMPRESSED_RGB_PVRTC_4BPPV1_IMG):
		return 3
	case (GL_COMPRESSED_RGB_PVRTC_2BPPV1_IMG):
		return 3
	case (GL_COMPRESSED_RGBA_PVRTC_4BPPV1_IMG):
		return 4
	case (GL_COMPRESSED_RGBA_PVRTC_2BPPV1_IMG):
		return 4
	case (GL_ETC1_RGB8_OES):
		return 3
	case (GL_COMPRESSED_RGB8_ETC2):
		return 3
	case (GL_COMPRESSED_SRGB8_ETC2):
		return 3
	case (GL_COMPRESSED_RGB8_PUNCHTHROUGH_ALPHA1_ETC2):
		return 4
	case (GL_COMPRESSED_SRGB8_PUNCHTHROUGH_ALPHA1_ETC2):
		return 4
	case (GL_COMPRESSED_RGBA8_ETC2_EAC):
		return 4
	case (GL_COMPRESSED_SRGB8_ALPHA8_ETC2_EAC):
		return 4
	case (GL_COMPRESSED_R11_EAC):
		return 1
	case (GL_COMPRESSED_SIGNED_R11_EAC):
		return 1
	case (GL_COMPRESSED_RG11_EAC):
		return 2
	case (GL_COMPRESSED_SIGNED_RG11_EAC):
		return 2
	case (GL_COLOR_INDEX):
		return 1
	case (GL_STENCIL_INDEX):
		return 1
	case (GL_DEPTH_COMPONENT):
		return 1
	case (GL_DEPTH_COMPONENT16):
		return 1
	case (GL_DEPTH_COMPONENT24):
		return 1
	case (GL_DEPTH_COMPONENT32):
		return 1
	case (GL_DEPTH_COMPONENT32F):
		return 1
	case (GL_DEPTH_COMPONENT32F_NV):
		return 1
	case (GL_RED):
		return 1
	case (GL_GREEN):
		return 1
	case (GL_BLUE):
		return 1
	case (GL_ALPHA):
		return 1
	case (GL_ALPHA8I_EXT):
		return 1
	case (GL_ALPHA8UI_EXT):
		return 1
	case (GL_ALPHA16I_EXT):
		return 1
	case (GL_ALPHA16UI_EXT):
		return 1
	case (GL_ALPHA32I_EXT):
		return 1
	case (GL_ALPHA32UI_EXT):
		return 1
	case (GL_ALPHA16F_ARB):
		return 1
	case (GL_ALPHA32F_ARB):
		return 1
	case (GL_R32F):
		return 1
	case (GL_RG):
		return 2
	case (GL_RG32F):
		return 2
	case (GL_RGB):
		return 3
	case (GL_BGR):
		return 3
	case (GL_RGB8I_EXT):
		return 3
	case (GL_RGB8UI_EXT):
		return 3
	case (GL_RGB16I_EXT):
		return 3
	case (GL_RGB16UI_EXT):
		return 3
	case (GL_RGB32I_EXT):
		return 3
	case (GL_RGB32UI_EXT):
		return 3
	case (GL_RGB16F_ARB):
		return 3
	case (GL_RGB32F_ARB):
		return 3
	case (GL_RGBA16F_ARB):
		return 4
	case (GL_RGBA32F_ARB):
		return 4
	case (GL_RGBA):
		return 4
	case (GL_BGRA):
		return 4
	case (GL_RGBA8):
		return 4
	case (GL_LUMINANCE):
		return 1
	case (GL_LUMINANCE4):
		return 1
	case (GL_LUMINANCE8):
		return 1
	case (GL_LUMINANCE12):
		return 1
	case (GL_LUMINANCE16):
		return 1
	case (GL_LUMINANCE8I_EXT):
		return 1
	case (GL_LUMINANCE8UI_EXT):
		return 1
	case (GL_LUMINANCE16I_EXT):
		return 1
	case (GL_LUMINANCE16UI_EXT):
		return 1
	case (GL_LUMINANCE32I_EXT):
		return 1
	case (GL_LUMINANCE32UI_EXT):
		return 1
	case (GL_LUMINANCE16F_ARB):
		return 1
	case (GL_LUMINANCE32F_ARB):
		return 1
	case (GL_LUMINANCE4_ALPHA4):
		return 2
	case (GL_LUMINANCE6_ALPHA2):
		return 2
	case (GL_LUMINANCE8_ALPHA8):
		return 2
	case (GL_LUMINANCE12_ALPHA4):
		return 2
	case (GL_LUMINANCE12_ALPHA12):
		return 2
	case (GL_LUMINANCE16_ALPHA16):
		return 2
	case (GL_INTENSITY):
		return 1
	case (GL_INTENSITY4):
		return 1
	case (GL_INTENSITY8):
		return 1
	case (GL_INTENSITY12):
		return 1
	case (GL_INTENSITY16):
		return 1
	case (GL_INTENSITY8UI_EXT):
		return 1
	case (GL_INTENSITY8I_EXT):
		return 1
	case (GL_INTENSITY16I_EXT):
		return 1
	case (GL_INTENSITY16UI_EXT):
		return 1
	case (GL_INTENSITY32I_EXT):
		return 1
	case (GL_INTENSITY32UI_EXT):
		return 1
	case (GL_INTENSITY16F_ARB):
		return 1
	case (GL_INTENSITY32F_ARB):
		return 1
	case (GL_LUMINANCE_ALPHA):
		return 2
	case (GL_LUMINANCE_ALPHA8I_EXT):
		return 2
	case (GL_LUMINANCE_ALPHA8UI_EXT):
		return 2
	case (GL_LUMINANCE_ALPHA16I_EXT):
		return 2
	case (GL_LUMINANCE_ALPHA16UI_EXT):
		return 2
	case (GL_LUMINANCE_ALPHA32I_EXT):
		return 2
	case (GL_LUMINANCE_ALPHA32UI_EXT):
		return 2
	case (GL_LUMINANCE_ALPHA16F_ARB):
		return 2
	case (GL_LUMINANCE_ALPHA32F_ARB):
		return 2
	case (GL_HILO_NV):
		return 2
	case (GL_DSDT_NV):
		return 2
	case (GL_DSDT_MAG_NV):
		return 3
	case (GL_DSDT_MAG_VIB_NV):
		return 4
	case (GL_RED_INTEGER_EXT):
		return 1
	case (GL_GREEN_INTEGER_EXT):
		return 1
	case (GL_BLUE_INTEGER_EXT):
		return 1
	case (GL_ALPHA_INTEGER_EXT):
		return 1
	case (GL_RGB_INTEGER_EXT):
		return 3
	case (GL_RGBA_INTEGER_EXT):
		return 4
	case (GL_BGR_INTEGER_EXT):
		return 3
	case (GL_BGRA_INTEGER_EXT):
		return 4
	case (GL_LUMINANCE_INTEGER_EXT):
		return 1
	case (GL_LUMINANCE_ALPHA_INTEGER_EXT):
		return 2
	case (GL_COMPRESSED_RGBA_ASTC_4x4_KHR):
		return 4
	case (GL_COMPRESSED_RGBA_ASTC_5x4_KHR):
		return 4
	case (GL_COMPRESSED_RGBA_ASTC_5x5_KHR):
		return 4
	case (GL_COMPRESSED_RGBA_ASTC_6x5_KHR):
		return 4
	case (GL_COMPRESSED_RGBA_ASTC_6x6_KHR):
		return 4
	case (GL_COMPRESSED_RGBA_ASTC_8x5_KHR):
		return 4
	case (GL_COMPRESSED_RGBA_ASTC_8x6_KHR):
		return 4
	case (GL_COMPRESSED_RGBA_ASTC_8x8_KHR):
		return 4
	case (GL_COMPRESSED_RGBA_ASTC_10x5_KHR):
		return 4
	case (GL_COMPRESSED_RGBA_ASTC_10x6_KHR):
		return 4
	case (GL_COMPRESSED_RGBA_ASTC_10x8_KHR):
		return 4
	case (GL_COMPRESSED_RGBA_ASTC_10x10_KHR):
		return 4
	case (GL_COMPRESSED_RGBA_ASTC_12x10_KHR):
		return 4
	case (GL_COMPRESSED_RGBA_ASTC_12x12_KHR):
		return 4
	case (GL_COMPRESSED_SRGB8_ALPHA8_ASTC_4x4_KHR):
		return 4
	case (GL_COMPRESSED_SRGB8_ALPHA8_ASTC_5x4_KHR):
		return 4
	case (GL_COMPRESSED_SRGB8_ALPHA8_ASTC_5x5_KHR):
		return 4
	case (GL_COMPRESSED_SRGB8_ALPHA8_ASTC_6x5_KHR):
		return 4
	case (GL_COMPRESSED_SRGB8_ALPHA8_ASTC_6x6_KHR):
		return 4
	case (GL_COMPRESSED_SRGB8_ALPHA8_ASTC_8x5_KHR):
		return 4
	case (GL_COMPRESSED_SRGB8_ALPHA8_ASTC_8x6_KHR):
		return 4
	case (GL_COMPRESSED_SRGB8_ALPHA8_ASTC_8x8_KHR):
		return 4
	case (GL_COMPRESSED_SRGB8_ALPHA8_ASTC_10x5_KHR):
		return 4
	case (GL_COMPRESSED_SRGB8_ALPHA8_ASTC_10x6_KHR):
		return 4
	case (GL_COMPRESSED_SRGB8_ALPHA8_ASTC_10x8_KHR):
		return 4
	case (GL_COMPRESSED_SRGB8_ALPHA8_ASTC_10x10_KHR):
		return 4
	case (GL_COMPRESSED_SRGB8_ALPHA8_ASTC_12x10_KHR):
		return 4
	case (GL_COMPRESSED_SRGB8_ALPHA8_ASTC_12x12_KHR):
		return 4
	}
	return 0
}

func ComputePixelSizeInBits(format int, t int) uint {
	switch format {
	case (GL_COMPRESSED_RGB_S3TC_DXT1_EXT):
		return 4
	case (GL_COMPRESSED_RGBA_S3TC_DXT1_EXT):
		return 4
	case (GL_COMPRESSED_RGBA_S3TC_DXT3_EXT):
		return 8
	case (GL_COMPRESSED_RGBA_S3TC_DXT5_EXT):
		return 8
	case (GL_COMPRESSED_SIGNED_RED_RGTC1_EXT):
		return 4
	case (GL_COMPRESSED_RED_RGTC1_EXT):
		return 4
	case (GL_COMPRESSED_SIGNED_RED_GREEN_RGTC2_EXT):
		return 8
	case (GL_COMPRESSED_RED_GREEN_RGTC2_EXT):
		return 8
	case (GL_COMPRESSED_RGB_PVRTC_4BPPV1_IMG):
		return 4
	case (GL_COMPRESSED_RGB_PVRTC_2BPPV1_IMG):
		return 2
	case (GL_COMPRESSED_RGBA_PVRTC_4BPPV1_IMG):
		return 4
	case (GL_COMPRESSED_RGBA_PVRTC_2BPPV1_IMG):
		return 2
	case (GL_ETC1_RGB8_OES):
		return 4
	case (GL_COMPRESSED_RGB8_ETC2):
		return 4
	case (GL_COMPRESSED_SRGB8_ETC2):
		return 4
	case (GL_COMPRESSED_RGB8_PUNCHTHROUGH_ALPHA1_ETC2):
		return 4
	case (GL_COMPRESSED_SRGB8_PUNCHTHROUGH_ALPHA1_ETC2):
		return 4
	case (GL_COMPRESSED_RGBA8_ETC2_EAC):
		return 8
	case (GL_COMPRESSED_SRGB8_ALPHA8_ETC2_EAC):
		return 8
	case (GL_COMPRESSED_R11_EAC):
		return 4
	case (GL_COMPRESSED_SIGNED_R11_EAC):
		return 4
	case (GL_COMPRESSED_RG11_EAC):
		return 8
	case (GL_COMPRESSED_SIGNED_RG11_EAC):
		return 8
	}

	switch format {
	case (GL_COMPRESSED_ALPHA):
	case (GL_COMPRESSED_LUMINANCE):
	case (GL_COMPRESSED_LUMINANCE_ALPHA):
	case (GL_COMPRESSED_INTENSITY):
	case (GL_COMPRESSED_RGB):
	case (GL_COMPRESSED_RGBA):
		return 0

	}
	switch format {
	case (GL_COMPRESSED_RGBA_ASTC_4x4_KHR):
	case (GL_COMPRESSED_RGBA_ASTC_5x4_KHR):
	case (GL_COMPRESSED_RGBA_ASTC_5x5_KHR):
	case (GL_COMPRESSED_RGBA_ASTC_6x5_KHR):
	case (GL_COMPRESSED_RGBA_ASTC_6x6_KHR):
	case (GL_COMPRESSED_RGBA_ASTC_8x5_KHR):
	case (GL_COMPRESSED_RGBA_ASTC_8x6_KHR):
	case (GL_COMPRESSED_RGBA_ASTC_8x8_KHR):
	case (GL_COMPRESSED_RGBA_ASTC_10x5_KHR):
	case (GL_COMPRESSED_RGBA_ASTC_10x6_KHR):
	case (GL_COMPRESSED_RGBA_ASTC_10x8_KHR):
	case (GL_COMPRESSED_RGBA_ASTC_10x10_KHR):
	case (GL_COMPRESSED_RGBA_ASTC_12x10_KHR):
	case (GL_COMPRESSED_RGBA_ASTC_12x12_KHR):
	case (GL_COMPRESSED_SRGB8_ALPHA8_ASTC_4x4_KHR):
	case (GL_COMPRESSED_SRGB8_ALPHA8_ASTC_5x4_KHR):
	case (GL_COMPRESSED_SRGB8_ALPHA8_ASTC_5x5_KHR):
	case (GL_COMPRESSED_SRGB8_ALPHA8_ASTC_6x5_KHR):
	case (GL_COMPRESSED_SRGB8_ALPHA8_ASTC_6x6_KHR):
	case (GL_COMPRESSED_SRGB8_ALPHA8_ASTC_8x5_KHR):
	case (GL_COMPRESSED_SRGB8_ALPHA8_ASTC_8x6_KHR):
	case (GL_COMPRESSED_SRGB8_ALPHA8_ASTC_8x8_KHR):
	case (GL_COMPRESSED_SRGB8_ALPHA8_ASTC_10x5_KHR):
	case (GL_COMPRESSED_SRGB8_ALPHA8_ASTC_10x6_KHR):
	case (GL_COMPRESSED_SRGB8_ALPHA8_ASTC_10x8_KHR):
	case (GL_COMPRESSED_SRGB8_ALPHA8_ASTC_10x10_KHR):
	case (GL_COMPRESSED_SRGB8_ALPHA8_ASTC_12x10_KHR):
	case (GL_COMPRESSED_SRGB8_ALPHA8_ASTC_12x12_KHR):
		{
			footprint := ComputeBlockFootprint(format)
			pixelsPerBlock := footprint[0] * footprint[1]
			bitsPerBlock := ComputeBlockSize(format, 0) // 16 x 8 = 128
			bitsPerPixel := bitsPerBlock / uint(pixelsPerBlock)
			if bitsPerBlock == bitsPerPixel*uint(pixelsPerBlock) {
				return bitsPerPixel
			}
			return 0
		}
		return 0
	}

	switch format {
	case (GL_LUMINANCE4):
		return 4
	case (GL_LUMINANCE8):
		return 8
	case (GL_LUMINANCE12):
		return 12
	case (GL_LUMINANCE16):
		return 16
	case (GL_LUMINANCE4_ALPHA4):
		return 8
	case (GL_LUMINANCE6_ALPHA2):
		return 8
	case (GL_LUMINANCE8_ALPHA8):
		return 16
	case (GL_LUMINANCE12_ALPHA4):
		return 16
	case (GL_LUMINANCE12_ALPHA12):
		return 24
	case (GL_LUMINANCE16_ALPHA16):
		return 32
	case (GL_INTENSITY4):
		return 4
	case (GL_INTENSITY8):
		return 8
	case (GL_INTENSITY12):
		return 12
	case (GL_INTENSITY16):
		return 16
	}
	return 0
}

func ComputeBlockFootprint(pixelFormat int) vec3.T {
	switch pixelFormat {
	case (GL_COMPRESSED_RGB_S3TC_DXT1_EXT):
	case (GL_COMPRESSED_RGBA_S3TC_DXT1_EXT):
	case (GL_COMPRESSED_RGBA_S3TC_DXT3_EXT):
	case (GL_COMPRESSED_RGBA_S3TC_DXT5_EXT):
		return vec3.T{4, 4, 4}

	case (GL_COMPRESSED_SIGNED_RED_RGTC1_EXT):
	case (GL_COMPRESSED_RED_RGTC1_EXT):
	case (GL_COMPRESSED_SIGNED_RED_GREEN_RGTC2_EXT):
	case (GL_COMPRESSED_RED_GREEN_RGTC2_EXT):
	case (GL_COMPRESSED_RGB_PVRTC_4BPPV1_IMG):
	case (GL_COMPRESSED_RGBA_PVRTC_4BPPV1_IMG):
	case (GL_ETC1_RGB8_OES):
	case (GL_COMPRESSED_RGB8_ETC2):
	case (GL_COMPRESSED_SRGB8_ETC2):
	case (GL_COMPRESSED_RGB8_PUNCHTHROUGH_ALPHA1_ETC2):
	case (GL_COMPRESSED_SRGB8_PUNCHTHROUGH_ALPHA1_ETC2):
	case (GL_COMPRESSED_RGBA8_ETC2_EAC):
	case (GL_COMPRESSED_SRGB8_ALPHA8_ETC2_EAC):
	case (GL_COMPRESSED_R11_EAC):
	case (GL_COMPRESSED_SIGNED_R11_EAC):
	case (GL_COMPRESSED_RG11_EAC):
	case (GL_COMPRESSED_SIGNED_RG11_EAC):
		return vec3.T{4, 4, 1}
	case (GL_COMPRESSED_RGB_PVRTC_2BPPV1_IMG):
	case (GL_COMPRESSED_RGBA_PVRTC_2BPPV1_IMG):
		return vec3.T{8, 4, 1} // no 3d texture support in pvrtc at all
	case (GL_COMPRESSED_RGBA_ASTC_4x4_KHR):
		return vec3.T{4, 4, 1}
	case (GL_COMPRESSED_RGBA_ASTC_5x4_KHR):
		return vec3.T{5, 4, 1}
	case (GL_COMPRESSED_RGBA_ASTC_5x5_KHR):
		return vec3.T{5, 5, 1}
	case (GL_COMPRESSED_RGBA_ASTC_6x5_KHR):
		return vec3.T{6, 5, 1}
	case (GL_COMPRESSED_RGBA_ASTC_6x6_KHR):
		return vec3.T{6, 6, 1}
	case (GL_COMPRESSED_RGBA_ASTC_8x5_KHR):
		return vec3.T{8, 5, 1}
	case (GL_COMPRESSED_RGBA_ASTC_8x6_KHR):
		return vec3.T{8, 6, 1}
	case (GL_COMPRESSED_RGBA_ASTC_8x8_KHR):
		return vec3.T{8, 8, 1}
	case (GL_COMPRESSED_RGBA_ASTC_10x5_KHR):
		return vec3.T{10, 5, 1}
	case (GL_COMPRESSED_RGBA_ASTC_10x6_KHR):
		return vec3.T{10, 6, 1}
	case (GL_COMPRESSED_RGBA_ASTC_10x8_KHR):
		return vec3.T{10, 8, 1}
	case (GL_COMPRESSED_RGBA_ASTC_10x10_KHR):
		return vec3.T{10, 10, 1}
	case (GL_COMPRESSED_RGBA_ASTC_12x10_KHR):
		return vec3.T{12, 10, 1}
	case (GL_COMPRESSED_RGBA_ASTC_12x12_KHR):
		return vec3.T{12, 12, 1}
	case (GL_COMPRESSED_SRGB8_ALPHA8_ASTC_4x4_KHR):
		return vec3.T{4, 4, 1}
	case (GL_COMPRESSED_SRGB8_ALPHA8_ASTC_5x4_KHR):
		return vec3.T{5, 4, 1}
	case (GL_COMPRESSED_SRGB8_ALPHA8_ASTC_5x5_KHR):
		return vec3.T{5, 5, 1}
	case (GL_COMPRESSED_SRGB8_ALPHA8_ASTC_6x5_KHR):
		return vec3.T{6, 5, 1}
	case (GL_COMPRESSED_SRGB8_ALPHA8_ASTC_6x6_KHR):
		return vec3.T{6, 6, 1}
	case (GL_COMPRESSED_SRGB8_ALPHA8_ASTC_8x5_KHR):
		return vec3.T{8, 5, 1}
	case (GL_COMPRESSED_SRGB8_ALPHA8_ASTC_8x6_KHR):
		return vec3.T{8, 6, 1}
	case (GL_COMPRESSED_SRGB8_ALPHA8_ASTC_8x8_KHR):
		return vec3.T{8, 8, 1}
	case (GL_COMPRESSED_SRGB8_ALPHA8_ASTC_10x5_KHR):
		return vec3.T{10, 5, 1}
	case (GL_COMPRESSED_SRGB8_ALPHA8_ASTC_10x6_KHR):
		return vec3.T{10, 6, 1}
	case (GL_COMPRESSED_SRGB8_ALPHA8_ASTC_10x8_KHR):
		return vec3.T{10, 8, 1}
	case (GL_COMPRESSED_SRGB8_ALPHA8_ASTC_10x10_KHR):
		return vec3.T{10, 10, 1}
	case (GL_COMPRESSED_SRGB8_ALPHA8_ASTC_12x10_KHR):
		return vec3.T{12, 10, 1}
	case (GL_COMPRESSED_SRGB8_ALPHA8_ASTC_12x12_KHR):
		return vec3.T{12, 12, 1}

	}
	return vec3.T{1, 1, 1}
}

func ComputeBlockSize(pixelFormat int, packing uint) uint {
	switch pixelFormat {
	case (GL_COMPRESSED_RGB_S3TC_DXT1_EXT):
	case (GL_COMPRESSED_RGBA_S3TC_DXT1_EXT):
		if packing > 8 {
			return packing
		}
		return 8

	case (GL_COMPRESSED_RGBA_S3TC_DXT3_EXT):
	case (GL_COMPRESSED_RGBA_S3TC_DXT5_EXT):
	case (GL_COMPRESSED_RGB_PVRTC_2BPPV1_IMG):
	case (GL_COMPRESSED_RGBA_PVRTC_2BPPV1_IMG):
	case (GL_COMPRESSED_RGB_PVRTC_4BPPV1_IMG):
	case (GL_COMPRESSED_RGBA_PVRTC_4BPPV1_IMG):
	case (GL_ETC1_RGB8_OES):
		if packing > 16 {
			return packing
		}
		return 16

	case (GL_COMPRESSED_SIGNED_RED_RGTC1_EXT):
	case (GL_COMPRESSED_RED_RGTC1_EXT):
		if packing > 8 {
			return packing
		}
		return 8

	case (GL_COMPRESSED_SIGNED_RED_GREEN_RGTC2_EXT):
	case (GL_COMPRESSED_RED_GREEN_RGTC2_EXT):
		if packing > 16 {
			return packing
		}
		return 16

	case (GL_COMPRESSED_RGB8_ETC2):
	case (GL_COMPRESSED_SRGB8_ETC2):
	case (GL_COMPRESSED_RGB8_PUNCHTHROUGH_ALPHA1_ETC2):
	case (GL_COMPRESSED_SRGB8_PUNCHTHROUGH_ALPHA1_ETC2):
	case (GL_COMPRESSED_R11_EAC):
	case (GL_COMPRESSED_SIGNED_R11_EAC):
		if packing > 8 {
			return packing
		}
		return 8

	case (GL_COMPRESSED_RGBA8_ETC2_EAC):
	case (GL_COMPRESSED_SRGB8_ALPHA8_ETC2_EAC):
	case (GL_COMPRESSED_RG11_EAC):
	case (GL_COMPRESSED_SIGNED_RG11_EAC):
		if packing > 16 {
			return packing
		}
		return 16

	case (GL_COMPRESSED_RGBA_ASTC_4x4_KHR):
	case (GL_COMPRESSED_RGBA_ASTC_5x4_KHR):
	case (GL_COMPRESSED_RGBA_ASTC_5x5_KHR):
	case (GL_COMPRESSED_RGBA_ASTC_6x5_KHR):
	case (GL_COMPRESSED_RGBA_ASTC_6x6_KHR):
	case (GL_COMPRESSED_RGBA_ASTC_8x5_KHR):
	case (GL_COMPRESSED_RGBA_ASTC_8x6_KHR):
	case (GL_COMPRESSED_RGBA_ASTC_8x8_KHR):
	case (GL_COMPRESSED_RGBA_ASTC_10x5_KHR):
	case (GL_COMPRESSED_RGBA_ASTC_10x6_KHR):
	case (GL_COMPRESSED_RGBA_ASTC_10x8_KHR):
	case (GL_COMPRESSED_RGBA_ASTC_10x10_KHR):
	case (GL_COMPRESSED_RGBA_ASTC_12x10_KHR):
	case (GL_COMPRESSED_RGBA_ASTC_12x12_KHR):
	case (GL_COMPRESSED_SRGB8_ALPHA8_ASTC_4x4_KHR):
	case (GL_COMPRESSED_SRGB8_ALPHA8_ASTC_5x4_KHR):
	case (GL_COMPRESSED_SRGB8_ALPHA8_ASTC_5x5_KHR):
	case (GL_COMPRESSED_SRGB8_ALPHA8_ASTC_6x5_KHR):
	case (GL_COMPRESSED_SRGB8_ALPHA8_ASTC_6x6_KHR):
	case (GL_COMPRESSED_SRGB8_ALPHA8_ASTC_8x5_KHR):
	case (GL_COMPRESSED_SRGB8_ALPHA8_ASTC_8x6_KHR):
	case (GL_COMPRESSED_SRGB8_ALPHA8_ASTC_8x8_KHR):
	case (GL_COMPRESSED_SRGB8_ALPHA8_ASTC_10x5_KHR):
	case (GL_COMPRESSED_SRGB8_ALPHA8_ASTC_10x6_KHR):
	case (GL_COMPRESSED_SRGB8_ALPHA8_ASTC_10x8_KHR):
	case (GL_COMPRESSED_SRGB8_ALPHA8_ASTC_10x10_KHR):
	case (GL_COMPRESSED_SRGB8_ALPHA8_ASTC_12x10_KHR):
	case (GL_COMPRESSED_SRGB8_ALPHA8_ASTC_12x12_KHR):
		if packing > 16 {
			return packing
		}
		return 16

	}
	return packing
}

func ComputeRowWidthInBytes(width int, pixelFormat int,
	t int, packing int) uint {
	pixelSize := ComputePixelSizeInBits(pixelFormat, t)
	widthInBits := width * int(pixelSize)
	var packingInBits int = 8
	if packing != 0 {
		packingInBits = packing * 8
	}
	var i int = 0
	if widthInBits%packingInBits > 0 {
		i = 1
	}
	return uint((widthInBits/packingInBits + i) * packing)
}

func ComputeImageSizeInBytes(width int, height int,
	depth int, pixelFormat int,
	ty int, packing int,
	slice_packing int,
	image_packing int) uint {
	if width <= 0 || height <= 0 || depth <= 0 {
		return 0
	}

	blockSize := ComputeBlockSize(pixelFormat, 0)
	if blockSize > 0 {
		footprint := ComputeBlockFootprint(pixelFormat)
		width = (width + int(footprint[0]) - 1) / int(footprint[0])
		height = (height + int(footprint[1]) - 1) / int(footprint[1])

		size := int(blockSize) * width
		size = RoudUpToMultiple(size, packing)
		size *= height
		size = RoudUpToMultiple(size, slice_packing)
		size *= depth
		size = RoudUpToMultiple(size, image_packing)
		return uint(size)
	}

	size := int(img.ComputeRowWidthInBytes(width, pixelFormat, ty, packing))

	size *= height
	size += slice_packing - 1
	size -= size % slice_packing

	size *= depth
	size += image_packing - 1
	size -= size % image_packing

	st := int(img.ComputeBlockSize(pixelFormat, uint(packing)))
	if size > st {
		return uint(size)
	}
	return uint(st)
}

func ComputeNearestPowerOfTwo(s int, bias float32) int {
	if (s & (s - 1)) != 0 {
		p2 := math.Log(float64(s)) / math.Log(2.0)
		rounded_p2 := math.Floor(p2 + float64(bias))
		s = (int)(math.Pow(2.0, rounded_p2))
	}
	return s
}

func RoudUpToMultiple(s int, pack int) int {
	if pack < 2 {
		return s
	}
	s += pack - 1
	s -= s % pack
	return s
}

func (img *Image) IsCompressed() bool {
	switch img.PixelFormat {
	case (GL_COMPRESSED_ALPHA_ARB):
	case (GL_COMPRESSED_INTENSITY_ARB):
	case (GL_COMPRESSED_LUMINANCE_ALPHA_ARB):
	case (GL_COMPRESSED_LUMINANCE_ARB):
	case (GL_COMPRESSED_RGBA_ARB):
	case (GL_COMPRESSED_RGB_ARB):
	case (GL_COMPRESSED_RGB_S3TC_DXT1_EXT):
	case (GL_COMPRESSED_RGBA_S3TC_DXT1_EXT):
	case (GL_COMPRESSED_RGBA_S3TC_DXT3_EXT):
	case (GL_COMPRESSED_RGBA_S3TC_DXT5_EXT):
	case (GL_COMPRESSED_SIGNED_RED_RGTC1_EXT):
	case (GL_COMPRESSED_RED_RGTC1_EXT):
	case (GL_COMPRESSED_SIGNED_RED_GREEN_RGTC2_EXT):
	case (GL_COMPRESSED_RED_GREEN_RGTC2_EXT):
	case (GL_COMPRESSED_RGB_PVRTC_4BPPV1_IMG):
	case (GL_COMPRESSED_RGB_PVRTC_2BPPV1_IMG):
	case (GL_COMPRESSED_RGBA_PVRTC_4BPPV1_IMG):
	case (GL_COMPRESSED_RGBA_PVRTC_2BPPV1_IMG):
	case (GL_ETC1_RGB8_OES):
	case (GL_COMPRESSED_RGB8_ETC2):
	case (GL_COMPRESSED_SRGB8_ETC2):
	case (GL_COMPRESSED_RGB8_PUNCHTHROUGH_ALPHA1_ETC2):
	case (GL_COMPRESSED_SRGB8_PUNCHTHROUGH_ALPHA1_ETC2):
	case (GL_COMPRESSED_RGBA8_ETC2_EAC):
	case (GL_COMPRESSED_SRGB8_ALPHA8_ETC2_EAC):
	case (GL_COMPRESSED_R11_EAC):
	case (GL_COMPRESSED_SIGNED_R11_EAC):
	case (GL_COMPRESSED_RG11_EAC):
	case (GL_COMPRESSED_SIGNED_RG11_EAC):
	case (GL_COMPRESSED_RGBA_ASTC_4x4_KHR):
	case (GL_COMPRESSED_RGBA_ASTC_5x4_KHR):
	case (GL_COMPRESSED_RGBA_ASTC_5x5_KHR):
	case (GL_COMPRESSED_RGBA_ASTC_6x5_KHR):
	case (GL_COMPRESSED_RGBA_ASTC_6x6_KHR):
	case (GL_COMPRESSED_RGBA_ASTC_8x5_KHR):
	case (GL_COMPRESSED_RGBA_ASTC_8x6_KHR):
	case (GL_COMPRESSED_RGBA_ASTC_8x8_KHR):
	case (GL_COMPRESSED_RGBA_ASTC_10x5_KHR):
	case (GL_COMPRESSED_RGBA_ASTC_10x6_KHR):
	case (GL_COMPRESSED_RGBA_ASTC_10x8_KHR):
	case (GL_COMPRESSED_RGBA_ASTC_10x10_KHR):
	case (GL_COMPRESSED_RGBA_ASTC_12x10_KHR):
	case (GL_COMPRESSED_RGBA_ASTC_12x12_KHR):
	case (GL_COMPRESSED_SRGB8_ALPHA8_ASTC_4x4_KHR):
	case (GL_COMPRESSED_SRGB8_ALPHA8_ASTC_5x4_KHR):
	case (GL_COMPRESSED_SRGB8_ALPHA8_ASTC_5x5_KHR):
	case (GL_COMPRESSED_SRGB8_ALPHA8_ASTC_6x5_KHR):
	case (GL_COMPRESSED_SRGB8_ALPHA8_ASTC_6x6_KHR):
	case (GL_COMPRESSED_SRGB8_ALPHA8_ASTC_8x5_KHR):
	case (GL_COMPRESSED_SRGB8_ALPHA8_ASTC_8x6_KHR):
	case (GL_COMPRESSED_SRGB8_ALPHA8_ASTC_8x8_KHR):
	case (GL_COMPRESSED_SRGB8_ALPHA8_ASTC_10x5_KHR):
	case (GL_COMPRESSED_SRGB8_ALPHA8_ASTC_10x6_KHR):
	case (GL_COMPRESSED_SRGB8_ALPHA8_ASTC_10x8_KHR):
	case (GL_COMPRESSED_SRGB8_ALPHA8_ASTC_10x10_KHR):
	case (GL_COMPRESSED_SRGB8_ALPHA8_ASTC_12x10_KHR):
	case (GL_COMPRESSED_SRGB8_ALPHA8_ASTC_12x12_KHR):
		return true
	}
	return false
}
