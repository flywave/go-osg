package model

type Planef struct {
	Normal   [3]float32
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

type Quaternion struct {
	Value [4]float64
}
