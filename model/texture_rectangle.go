package model

type TextureRectangle struct {
	Texture
	Width  int
	Height int
	Image  *Image
}

func NewTextureRectangle() *TextureRectangle {
	t := NewTexture()
	return &TextureRectangle{Texture: *t}
}
