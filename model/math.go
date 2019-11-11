package model

import "github.com/ungerik/go3d/vec3"

type Planef struct {
	Normal   vec3.T
	Distance float32
}

type Planed struct {
	Normal   [3]float64
	Distance float64
}

type Sphere3f struct {
	Center [3]float32
	Radius float32
}
type Sphere3d struct {
	Center [3]float64
	Radius float64
}
