package model

import "github.com/ungerik/go3d/vec4"

const (
	MTL_T               = "osg::Material"
	AMBIENT             = 0x1200
	DIFFUSE             = 0x1201
	SPECULAR            = 0x1202
	EMISSION            = 0x1600
	AMBIENT_AND_DIFFUSE = 0x1602
	MTLOFF              = 0x1602 + 1
)

type Material struct {
	StateAttribute
	AmbientFrontAndBack   bool
	DiffuseFrontAndBack   bool
	SpecularFrontAndBack  bool
	EmissionFrontAndBack  bool
	ShininessFrontAndBack bool

	AmbientFront vec4.T
	AmbientBack  vec4.T

	DiffuseFront vec4.T
	DiffuseBack  vec4.T

	SpecularFront vec4.T
	SpecularBack  vec4.T

	EmissionFront vec4.T
	EmissionBack  vec4.T

	ShininessFront float32
	ShininessBack  float32
	Cmod           int
}

func NewMaterial() Material {
	st := NewStateAttribute()
	st.Type = MTL_T
	return Material{StateAttribute: st,
		AmbientFrontAndBack:   true,
		DiffuseFrontAndBack:   true,
		SpecularFrontAndBack:  true,
		EmissionFrontAndBack:  true,
		ShininessFrontAndBack: true,
		AmbientFront:          vec4.T{0.2, 0.2, 0.2, 1.0},
		AmbientBack:           vec4.T{0.2, 0.2, 0.2, 1.0},
		DiffuseFront:          vec4.T{0.8, 0.8, 0.8, 1.0},
		DiffuseBack:           vec4.T{0.8, 0.8, 0.8, 1.0},
		SpecularFront:         vec4.T{0.0, 0.0, 0.0, 1.0},
		SpecularBack:          vec4.T{0.0, 0.0, 0.0, 1.0},
		EmissionFront:         vec4.T{0.0, 0.0, 0.0, 1.0},
		EmissionBack:          vec4.T{0.0, 0.0, 0.0, 1.0},
		ShininessFront:        0.0,
		ShininessBack:         0.0,
		Cmod:                  MTLOFF,
	}
}

func (mt *Material) SetAmbient(face int, ambient vec4.T) {
	switch face {
	case GL_FRONT:
		mt.AmbientFrontAndBack = false
		mt.AmbientFront = ambient
		break
	case GL_BACK:
		mt.AmbientFrontAndBack = false
		mt.AmbientBack = ambient
		break
	case GL_FRONT_AND_BACK:
		mt.AmbientFrontAndBack = true
		mt.AmbientFront = ambient
		mt.AmbientBack = mt.AmbientFront
		break
	}
}

func (mt *Material) GetAmbient(face int) vec4.T {
	switch face {
	case GL_FRONT:
		return mt.AmbientFront
	case GL_BACK:
		return mt.AmbientBack
	case GL_FRONT_AND_BACK:
		if !mt.AmbientFrontAndBack {
			return mt.AmbientFront
		}
	}
	return mt.AmbientFront
}
func (mt *Material) SetDiffuse(face int, diffuse vec4.T) {
	switch face {
	case GL_FRONT:
		mt.DiffuseFrontAndBack = false
		mt.DiffuseFront = diffuse
		break
	case GL_BACK:
		mt.DiffuseFrontAndBack = false
		mt.DiffuseBack = diffuse
		break
	case GL_FRONT_AND_BACK:
		mt.DiffuseFrontAndBack = true
		mt.DiffuseFront = diffuse
		mt.DiffuseBack = mt.DiffuseFront
		break
	}
}

func (mt *Material) GetDiffuse(face int) vec4.T {
	switch face {
	case GL_FRONT:
		return mt.DiffuseFront
	case GL_BACK:
		return mt.DiffuseBack
	case GL_FRONT_AND_BACK:
		if !mt.DiffuseFrontAndBack {
			return mt.DiffuseFront
		}
	}
	return mt.DiffuseFront
}

func (mt *Material) SetSpecular(face int, specular vec4.T) {
	switch face {
	case GL_FRONT:
		mt.SpecularFrontAndBack = false
		mt.SpecularFront = specular
		break
	case GL_BACK:
		mt.SpecularFrontAndBack = false
		mt.SpecularBack = specular
		break
	case GL_FRONT_AND_BACK:
		mt.SpecularFrontAndBack = true
		mt.SpecularFront = specular
		mt.SpecularBack = mt.SpecularFront
		break
	default:
		break
	}
}

func (mt *Material) GetSpecular(face int) vec4.T {
	switch face {
	case GL_FRONT:
		return mt.SpecularFront
	case GL_BACK:
		return mt.SpecularBack
	case GL_FRONT_AND_BACK:
		if !mt.SpecularFrontAndBack {
			return mt.SpecularFront
		}
	}
	return mt.SpecularFront
}

func (mt *Material) SetEmission(face int, emission vec4.T) {
	switch face {
	case GL_FRONT:
		mt.EmissionFrontAndBack = false
		mt.EmissionFront = emission
		break
	case GL_BACK:
		mt.EmissionFrontAndBack = false
		mt.EmissionBack = emission
		break
	case GL_FRONT_AND_BACK:
		mt.EmissionFrontAndBack = true
		mt.EmissionFront = emission
		mt.EmissionBack = mt.EmissionFront
		break
	default:
		break
	}
}

func (mt *Material) GetEmission(face int) vec4.T {
	switch face {
	case GL_FRONT:
		return mt.EmissionFront
	case GL_BACK:
		return mt.EmissionBack
	case GL_FRONT_AND_BACK:
		if !mt.EmissionFrontAndBack {
		}
		return mt.EmissionFront
	}
	return mt.EmissionFront
}

func (mt *Material) SetShininess(face int, shininess float32) {
	if shininess < 0 {
		shininess = 0.0
	} else if shininess > 128.0 {
		shininess = 128.0
	}

	switch face {
	case GL_FRONT:
		mt.ShininessFrontAndBack = false
		mt.ShininessFront = shininess
		break
	case GL_BACK:
		mt.ShininessFrontAndBack = false
		mt.ShininessBack = shininess
		break
	case GL_FRONT_AND_BACK:
		mt.ShininessFrontAndBack = true
		mt.ShininessFront = shininess
		mt.ShininessBack = shininess
		break
	default:
		break
	}
}

func (mt *Material) GetShininess(face int) float32 {
	switch face {
	case GL_FRONT:
		return mt.ShininessFront
	case GL_BACK:
		return mt.ShininessBack
	case GL_FRONT_AND_BACK:
		if !mt.ShininessFrontAndBack {
		}
		return mt.ShininessFront
	}
	return mt.ShininessFront
}

func (mt *Material) SetTransparency(face int, transparency float32) {
	if face == GL_FRONT || face == GL_FRONT_AND_BACK {
		mt.AmbientFront[3] = 1.0 - transparency
		mt.DiffuseFront[3] = 1.0 - transparency
		mt.SpecularFront[3] = 1.0 - transparency
		mt.EmissionFront[3] = 1.0 - transparency
	}

	if face == GL_BACK || face == GL_FRONT_AND_BACK {
		mt.AmbientBack[3] = 1.0 - transparency
		mt.DiffuseBack[3] = 1.0 - transparency
		mt.SpecularBack[3] = 1.0 - transparency
		mt.EmissionBack[3] = 1.0 - transparency
	}
}
