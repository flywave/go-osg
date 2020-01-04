package model

const (
	POSITIVEX = 0
	NEGATIVEX = 1
	POSITIVEY = 2
	NEGATIVEY = 3
	POSITIVEZ = 4
	NEGATIVEZ = 5
)

type TextureCubeMap struct {
	Texture
	Width  int
	Height int
	Images []*Image
}

func NewTextureCubeMap() *TextureCubeMap {
	t := NewTexture()
	return &TextureCubeMap{Texture: *t, Images: make([]*Image, 6, 6)}
}

func (t *TextureCubeMap) GetImage(f int) *Image {
	return t.Images[f]
}

func (t *TextureCubeMap) SetImage(f int, img *Image) {
	t.Images[f] = img
}
