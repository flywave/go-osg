package model

const (
	POSITIVE_X = 0
	NEGATIVE_X = 1
	POSITIVE_Y = 2
	NEGATIVE_Y = 3
	POSITIVE_Z = 4
	NEGATIVE_Z = 5
)

type TextureCubeMap struct {
	Texture
	Width  int
	Height int
	Images []*Image
}

func NewTextureCubeMap() TextureCubeMap {
	t := NewTexture()
	return TextureCubeMap{Texture: t, Images: make([]*Image, 6, 6)}
}

func (t *TextureCubeMap) GetImage(f int) *Image {
	return t.Images[f]
}

func (t *TextureCubeMap) SetImage(f int, img *Image) {
	t.Images[f] = img
}
