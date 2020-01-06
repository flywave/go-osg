package model

const (
	TEXTURET   string = "osg::Texture"
	TEXTURE1DT string = "osg::Texture1D"
	TEXTURE2DT string = "osg::Texture2D"
	TEXTURE3DT string = "osg::Texture3D"

	MINFILTER uint32 = 0
	MAGFILTER uint32 = 1

	LINEAR               = GLLINEAR
	LINEARMIPMAPLINEAR   = GLLINEARMIPMAPLINEAR
	LINEARMIPMAPNEAREST  = GLLINEARMIPMAPNEAREST
	NEAREST              = GLNEAREST
	NEARESTMIPMAPLINEAR  = GLNEARESTMIPMAPLINEAR
	NEARESTMIPMAPNEAREST = GLNEARESTMIPMAPNEAREST

	WRAPS = 0
	WRAPT = 1
	WRAPR = 2

	CLAMP         = GLCLAMP
	CLAMPTOEDGE   = GLCLAMPTOEDGE
	CLAMPTOBORDER = GLCLAMPTOBORDER
	REPEAT        = GLREPEAT
	MIRROR        = GLMIRROR

	NORMALIZED      uint32 = 0x0
	FLOAT           uint32 = 0x1
	SIGNEDINTEGER   uint32 = 0x2
	UNSIGNEDINTEGER uint32 = 0x4
)

type Texture struct {
	StateAttribute
	WrapS int
	WrapT int
	WrapR int

	MinFilter int
	MagFilter int

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

	InternalFormatMode uint32
	InternalFormatType uint32
	InternalFormat     int
	SourceFormat       int
	SourceType         int

	UseShadowComparison bool
	ShadowCompareFunc   int
	ShadowTextureMode   int
	ShadowAmbient       float32

	Image         *Image
	TextureWidth  uint32
	TextureHeight uint32
	TextureDepth  uint32
	TextureTarget int
}

func (tex *Texture) IsTextureAttribute() bool {
	return true
}

func IsCompressedInternalFormat(internalFormat int) bool {
	switch internalFormat {
	case GLCOMPRESSEDALPHAARB:
	case GLCOMPRESSEDINTENSITYARB:
	case GLCOMPRESSEDLUMINANCEALPHAARB:
	case GLCOMPRESSEDLUMINANCEARB:
	case GLCOMPRESSEDRGBAARB:
	case GLCOMPRESSEDRGBARB:
	case GLCOMPRESSEDRGBS3TCDXT1EXT:
	case (GLCOMPRESSEDRGBAS3TCDXT1EXT):
	case (GLCOMPRESSEDRGBAS3TCDXT3EXT):
	case (GLCOMPRESSEDRGBAS3TCDXT5EXT):
	case (GLCOMPRESSEDSIGNEDREDRGTC1EXT):
	case (GLCOMPRESSEDREDRGTC1EXT):
	case (GLCOMPRESSEDSIGNEDREDGREENRGTC2EXT):
	case (GLCOMPRESSEDREDGREENRGTC2EXT):
	case (GLETC1RGB8OES):
	case (GLCOMPRESSEDRGB8ETC2):
	case (GLCOMPRESSEDSRGB8ETC2):
	case (GLCOMPRESSEDRGB8PUNCHTHROUGHALPHA1ETC2):
	case (GLCOMPRESSEDSRGB8PUNCHTHROUGHALPHA1ETC2):
	case (GLCOMPRESSEDRGBA8ETC2EAC):
	case (GLCOMPRESSEDSRGB8ALPHA8ETC2EAC):
	case (GLCOMPRESSEDR11EAC):
	case (GLCOMPRESSEDSIGNEDR11EAC):
	case (GLCOMPRESSEDRG11EAC):
	case (GLCOMPRESSEDSIGNEDRG11EAC):
	case (GLCOMPRESSEDRGBPVRTC4BPPV1IMG):
	case (GLCOMPRESSEDRGBPVRTC2BPPV1IMG):
	case (GLCOMPRESSEDRGBAPVRTC4BPPV1IMG):
	case (GLCOMPRESSEDRGBAPVRTC2BPPV1IMG):
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
	if tex.InternalFormat == GLCOMPRESSEDRGBS3TCDXT1EXT ||
		tex.InternalFormat == GLCOMPRESSEDRGBAS3TCDXT1EXT {
		blockSize = 8
	} else if tex.InternalFormat == GLCOMPRESSEDRGBAS3TCDXT3EXT ||
		tex.InternalFormat == GLCOMPRESSEDRGBAS3TCDXT5EXT {
		blockSize = 16
	} else if tex.InternalFormat == GLETC1RGB8OES {
		blockSize = 8
	} else if tex.InternalFormat == GLCOMPRESSEDRGB8ETC2 ||
		tex.InternalFormat == GLCOMPRESSEDSRGB8ETC2 {
		blockSize = 8
	} else if tex.InternalFormat == GLCOMPRESSEDRGB8PUNCHTHROUGHALPHA1ETC2 ||
		tex.InternalFormat == GLCOMPRESSEDSRGB8PUNCHTHROUGHALPHA1ETC2 {
		blockSize = 8
	} else if tex.InternalFormat == GLCOMPRESSEDRGBA8ETC2EAC ||
		tex.InternalFormat == GLCOMPRESSEDSRGB8ALPHA8ETC2EAC {
		blockSize = 16
	} else if tex.InternalFormat == GLCOMPRESSEDR11EAC ||
		tex.InternalFormat == GLCOMPRESSEDSIGNEDR11EAC {
		blockSize = 8
	} else if tex.InternalFormat == GLCOMPRESSEDRG11EAC ||
		tex.InternalFormat == GLCOMPRESSEDSIGNEDRG11EAC {
		blockSize = 16
	} else if tex.InternalFormat == GLCOMPRESSEDREDRGTC1EXT ||
		tex.InternalFormat == GLCOMPRESSEDSIGNEDREDRGTC1EXT {
		blockSize = 8
	} else if tex.InternalFormat == GLCOMPRESSEDREDGREENRGTC2EXT ||
		tex.InternalFormat == GLCOMPRESSEDSIGNEDREDGREENRGTC2EXT {
		blockSize = 16
	} else if tex.InternalFormat == GLCOMPRESSEDRGBAPVRTC2BPPV1IMG ||
		tex.InternalFormat == GLCOMPRESSEDRGBPVRTC2BPPV1IMG {
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
	} else if tex.InternalFormat == GLCOMPRESSEDRGBAPVRTC4BPPV1IMG ||
		tex.InternalFormat == GLCOMPRESSEDRGBPVRTC4BPPV1IMG {
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
	case WRAPS:
		tex.WrapS = wrap
		break
	case WRAPT:
		tex.WrapT = wrap
		break
	case WRAPR:
		tex.WrapR = wrap
		break
	default:
		break
	}
}

func (tex *Texture) GetWrap(wrap int) int {
	if tex.WrapS == wrap {
		return WRAPS
	}
	if tex.WrapT == wrap {
		return WRAPT
	}
	if tex.WrapR == wrap {
		return WRAPR
	}
	return WRAPS
}

func (tex *Texture) SetFilter(which uint32, filter int) {
	switch which {
	case MINFILTER:
		tex.MinFilter = filter
		break
	case MAGFILTER:
		tex.MagFilter = filter
		break
	default:
		break
	}
}

func (tex *Texture) GetFilter(which uint32) int {
	switch which {
	case MINFILTER:
		return tex.MinFilter
	case MAGFILTER:
		return tex.MagFilter
	default:
		return tex.MinFilter
	}
}

func NewTexture() *Texture {
	st := NewStateAttribute()
	st.Type = TEXTURET
	return &Texture{
		StateAttribute:              *st,
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

func NewTexture1d() *Texture {
	t := NewTexture()
	t.Type = TEXTURE1DT
	t.TextureHeight = 1
	t.TextureDepth = 1
	return t
}

func NewTexture2d() *Texture {
	t := NewTexture()
	t.Type = TEXTURE2DT
	t.TextureDepth = 1
	return t
}

func NewTexture3d() *Texture {
	t := NewTexture()
	t.Type = TEXTURE3DT
	return t
}
