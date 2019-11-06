package model

import "github.com/ungerik/go3d/vec4"

type WrapParameter uint32
type FilterParameter uint32
type InternalFormatMode uint32
type InternalFormatType uint32

const (
	TEXTURE_T string = "osg::Texture"

	MIN_FILTER FilterParameter = 0
	MAG_FILTER FilterParameter = 1

	LINEAR                 = GL_LINEAR
	LINEAR_MIPMAP_LINEAR   = GL_LINEAR_MIPMAP_LINEAR
	LINEAR_MIPMAP_NEAREST  = GL_LINEAR_MIPMAP_NEAREST
	NEAREST                = GL_NEAREST
	NEAREST_MIPMAP_LINEAR  = GL_NEAREST_MIPMAP_LINEAR
	NEAREST_MIPMAP_NEAREST = GL_NEAREST_MIPMAP_NEAREST

	WRAP_S WrapParameter = 0
	WRAP_T WrapParameter = 1
	WRAP_R WrapParameter = 2

	CLAMP           = GL_CLAMP
	CLAMP_TO_EDGE   = GL_CLAMP_TO_EDGE
	CLAMP_TO_BORDER = GL_CLAMP_TO_BORDER
	REPEAT          = GL_REPEAT
	MIRROR          = GL_MIRROR

	NORMALIZED       InternalFormatType = 0x0
	FLOAT            InternalFormatType = 0x1
	SIGNED_INTEGER   InternalFormatType = 0x2
	UNSIGNED_INTEGER InternalFormatType = 0x4
)

type Texture struct {
	StateAttribute
	Wrap_S int
	Wrap_T int
	Wrap_R int

	Min_Filter int
	Mag_Filter int

	MaxAnisotropy               float32
	Minlod                      float32
	Maxlod                      float32
	Lodbias                     float32
	Swizzle                     vec4.T
	UseHardwareMipmapGeneration bool
	UnrefImageDataAfterApply    bool
	ClientStorageHint           bool
	ResizeNonPowerOfTwoHint     bool

	BorderColor vec4.T
	BorderWidth int

	InternalFormatMode InternalFormatMode
	InternalFormatType InternalFormatType
	InternalFormat     int
	SourceFormat       int
	SourceType         int

	UseShadowComparison bool
	ShadowCompareFunc   int
	ShadowTextureMode   int
	ShadowAmbient       float32
}
