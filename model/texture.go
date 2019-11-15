package model

type FilterParameter uint32
type InternalFormatMode uint32
type InternalFormatType uint32

const (
	TEXTURE_T   string = "osg::Texture"
	TEXTURE1D_T string = "osg::Texture1D"
	TEXTURE2D_T string = "osg::Texture2D"
	TEXTURE3D_T string = "osg::Texture3D"

	MIN_FILTER FilterParameter = 0
	MAG_FILTER FilterParameter = 1

	LINEAR                 = GL_LINEAR
	LINEAR_MIPMAP_LINEAR   = GL_LINEAR_MIPMAP_LINEAR
	LINEAR_MIPMAP_NEAREST  = GL_LINEAR_MIPMAP_NEAREST
	NEAREST                = GL_NEAREST
	NEAREST_MIPMAP_LINEAR  = GL_NEAREST_MIPMAP_LINEAR
	NEAREST_MIPMAP_NEAREST = GL_NEAREST_MIPMAP_NEAREST

	WRAP_S = 0
	WRAP_T = 1
	WRAP_R = 2

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
	Swizzle                     [4]int
	UseHardwareMipmapGeneration bool
	UnRefImageDataAfterApply    bool
	ClientStorageHint           bool
	ResizeNonPowerOfTwoHint     bool

	BorderColor [4]float64
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

	Image         *Image
	TextureWidth  uint64
	TextureHeight uint64
	TextureDepth  uint64
	TextureTarget int
}

func (tex *Texture) IsTextureAttribute() bool {
	return true
}

func IsCompressedInternalFormat(internalFormat int) bool {
	switch internalFormat {
	case GL_COMPRESSED_ALPHA_ARB:
	case GL_COMPRESSED_INTENSITY_ARB:
	case GL_COMPRESSED_LUMINANCE_ALPHA_ARB:
	case GL_COMPRESSED_LUMINANCE_ARB:
	case GL_COMPRESSED_RGBA_ARB:
	case GL_COMPRESSED_RGB_ARB:
	case GL_COMPRESSED_RGB_S3TC_DXT1_EXT:
	case (GL_COMPRESSED_RGBA_S3TC_DXT1_EXT):
	case (GL_COMPRESSED_RGBA_S3TC_DXT3_EXT):
	case (GL_COMPRESSED_RGBA_S3TC_DXT5_EXT):
	case (GL_COMPRESSED_SIGNED_RED_RGTC1_EXT):
	case (GL_COMPRESSED_RED_RGTC1_EXT):
	case (GL_COMPRESSED_SIGNED_RED_GREEN_RGTC2_EXT):
	case (GL_COMPRESSED_RED_GREEN_RGTC2_EXT):
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
	case (GL_COMPRESSED_RGB_PVRTC_4BPPV1_IMG):
	case (GL_COMPRESSED_RGB_PVRTC_2BPPV1_IMG):
	case (GL_COMPRESSED_RGBA_PVRTC_4BPPV1_IMG):
	case (GL_COMPRESSED_RGBA_PVRTC_2BPPV1_IMG):
		return true
	}
	return false
}

func (tex *Texture) IsCompressedInternalFormat() bool {
	return IsCompressedInternalFormat(tex.InternalFormat)
}

func (tex *Texture) GetCompressedSize(internalFormat int, width int, height int,
	depth int) (int, int) {
	var blockSize int = 0
	var size int = 0
	if tex.InternalFormat == GL_COMPRESSED_RGB_S3TC_DXT1_EXT ||
		tex.InternalFormat == GL_COMPRESSED_RGBA_S3TC_DXT1_EXT {
		blockSize = 8
	} else if tex.InternalFormat == GL_COMPRESSED_RGBA_S3TC_DXT3_EXT ||
		tex.InternalFormat == GL_COMPRESSED_RGBA_S3TC_DXT5_EXT {
		blockSize = 16
	} else if tex.InternalFormat == GL_ETC1_RGB8_OES {
		blockSize = 8
	} else if tex.InternalFormat == GL_COMPRESSED_RGB8_ETC2 ||
		tex.InternalFormat == GL_COMPRESSED_SRGB8_ETC2 {
		blockSize = 8
	} else if tex.InternalFormat == GL_COMPRESSED_RGB8_PUNCHTHROUGH_ALPHA1_ETC2 ||
		tex.InternalFormat == GL_COMPRESSED_SRGB8_PUNCHTHROUGH_ALPHA1_ETC2 {
		blockSize = 8
	} else if tex.InternalFormat == GL_COMPRESSED_RGBA8_ETC2_EAC ||
		tex.InternalFormat == GL_COMPRESSED_SRGB8_ALPHA8_ETC2_EAC {
		blockSize = 16
	} else if tex.InternalFormat == GL_COMPRESSED_R11_EAC ||
		tex.InternalFormat == GL_COMPRESSED_SIGNED_R11_EAC {
		blockSize = 8
	} else if tex.InternalFormat == GL_COMPRESSED_RG11_EAC ||
		tex.InternalFormat == GL_COMPRESSED_SIGNED_RG11_EAC {
		blockSize = 16
	} else if tex.InternalFormat == GL_COMPRESSED_RED_RGTC1_EXT ||
		tex.InternalFormat == GL_COMPRESSED_SIGNED_RED_RGTC1_EXT {
		blockSize = 8
	} else if tex.InternalFormat == GL_COMPRESSED_RED_GREEN_RGTC2_EXT ||
		tex.InternalFormat == GL_COMPRESSED_SIGNED_RED_GREEN_RGTC2_EXT {
		blockSize = 16
	} else if tex.InternalFormat == GL_COMPRESSED_RGBA_PVRTC_2BPPV1_IMG ||
		tex.InternalFormat == GL_COMPRESSED_RGB_PVRTC_2BPPV1_IMG {
		blockSize = 8 * 4
		widthBlocks := width / 8
		heightBlocks := height / 4
		bpp := 2

		if widthBlocks < 2 {
			widthBlocks = 2
		}
		if heightBlocks < 2 {
			heightBlocks = 2
		}
		size = widthBlocks * heightBlocks * blockSize * bpp / 8
	} else if tex.InternalFormat == GL_COMPRESSED_RGBA_PVRTC_4BPPV1_IMG ||
		tex.InternalFormat == GL_COMPRESSED_RGB_PVRTC_4BPPV1_IMG {
		blockSize = 4 * 4
		widthBlocks := width / 4
		heightBlocks := height / 4
		bpp := 4

		if widthBlocks < 2 {
			widthBlocks = 2
		}
		if heightBlocks < 2 {
			heightBlocks = 2
		}
		size = widthBlocks * heightBlocks * blockSize * bpp / 8
	} else {
		blockSize = 0
	}

	size = width + 3/4*height + 3/4*depth*blockSize
	return blockSize, size
}

func (tex *Texture) SetWrap(which int, wrap int) {
	switch which {
	case WRAP_S:
		tex.Wrap_S = wrap
		break
	case WRAP_T:
		tex.Wrap_T = wrap
		break
	case WRAP_R:
		tex.Wrap_R = wrap
		break
	default:
		break
	}
}

func (tex *Texture) GetWrap(wrap int) int {
	if tex.Wrap_S == wrap {
		return WRAP_S
	}
	if tex.Wrap_T == wrap {
		return WRAP_T
	}
	if tex.Wrap_R == wrap {
		return WRAP_R
	}
	return WRAP_S
}

func (tex *Texture) SetFilter(which FilterParameter, filter int) {
	switch which {
	case MIN_FILTER:
		tex.Min_Filter = filter
		break
	case MAG_FILTER:
		tex.Mag_Filter = filter
		break
	default:
		break
	}
}

func (tex *Texture) GetFilter(which FilterParameter) int {
	switch which {
	case MIN_FILTER:
		return tex.Min_Filter
	case MAG_FILTER:
		return tex.Mag_Filter
	default:
		return tex.Min_Filter
	}
}

func NewTexture() Texture {
	st := NewStateAttribute()
	st.Type = TEXTURE_T
	return Texture{
		StateAttribute:              st,
		MaxAnisotropy:               1,
		UseHardwareMipmapGeneration: true,
		UnRefImageDataAfterApply:    false,
		ClientStorageHint:           false,
		ResizeNonPowerOfTwoHint:     true,
		InternalFormat:              0,
		UseShadowComparison:         false,
		ShadowAmbient:               0,
		BorderWidth:                 1,
	}
}

func NewTexture1d() Texture {
	t := NewTexture()
	t.Type = TEXTURE1D_T
	t.TextureHeight = 1
	t.TextureDepth = 1
	return t
}

func NewTexture2d() Texture {
	t := NewTexture()
	t.Type = TEXTURE2D_T
	t.TextureDepth = 1
	return t
}

func NewTexture3d() Texture {
	t := NewTexture()
	t.Type = TEXTURE3D_T
	return t
}
