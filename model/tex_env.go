package model

type TexEnv struct {
	StateAttribute
	Mode  int32
	Color [4]float32
}

func NewTexEnv() *TexEnv {
	sa := NewStateAttribute()
	return &TexEnv{StateAttribute: *sa, Mode: GLMODULATE}
}
