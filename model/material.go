package model

const (
	MTLT              = "osg::Material"
	AMBIENT           = 0x1200
	DIFFUSE           = 0x1201
	SPECULAR          = 0x1202
	EMISSION          = 0x1600
	AMBIENTANDDIFFUSE = 0x1602
	MTLOFF            = 0x1602 + 1
)

type Material struct {
	StateAttribute
	AmbientFrontAndBack   bool
	DiffuseFrontAndBack   bool
	SpecularFrontAndBack  bool
	EmissionFrontAndBack  bool
	ShininessFrontAndBack bool

	AmbientFront [4]float32
	AmbientBack  [4]float32

	DiffuseFront [4]float32
	DiffuseBack  [4]float32

	SpecularFront [4]float32
	SpecularBack  [4]float32

	EmissionFront [4]float32
	EmissionBack  [4]float32

	ShininessFront float32
	ShininessBack  float32
	Cmod           int
}

func NewMaterial() *Material {
	st := NewStateAttribute()
	st.Type = MTLT
	return &Material{StateAttribute: *st,
		AmbientFrontAndBack:   true,
		DiffuseFrontAndBack:   true,
		SpecularFrontAndBack:  true,
		EmissionFrontAndBack:  true,
		ShininessFrontAndBack: true,
		AmbientFront:          [4]float32{0.2, 0.2, 0.2, 1.0},
		AmbientBack:           [4]float32{0.2, 0.2, 0.2, 1.0},
		DiffuseFront:          [4]float32{0.8, 0.8, 0.8, 1.0},
		DiffuseBack:           [4]float32{0.8, 0.8, 0.8, 1.0},
		SpecularFront:         [4]float32{0.0, 0.0, 0.0, 1.0},
		SpecularBack:          [4]float32{0.0, 0.0, 0.0, 1.0},
		EmissionFront:         [4]float32{0.0, 0.0, 0.0, 1.0},
		EmissionBack:          [4]float32{0.0, 0.0, 0.0, 1.0},
		ShininessFront:        0.0,
		ShininessBack:         0.0,
		Cmod:                  MTLOFF,
	}
}

func (mt *Material) SetAmbient(face int, ambient [4]float32) {
	switch face {
	case GLFRONT:
		mt.AmbientFrontAndBack = false
		mt.AmbientFront = ambient
		break
	case GLBACK:
		mt.AmbientFrontAndBack = false
		mt.AmbientBack = ambient
		break
	case GLFRONTANDBACK:
		mt.AmbientFrontAndBack = true
		mt.AmbientFront = ambient
		mt.AmbientBack = mt.AmbientFront
		break
	}
}

func (mt *Material) GetAmbient(face int) *[4]float32 {
	switch face {
	case GLFRONT:
		return &mt.AmbientFront
	case GLBACK:
		return &mt.AmbientBack
	case GLFRONTANDBACK:
		if !mt.AmbientFrontAndBack {
			return &mt.AmbientFront
		}
	}
	return &mt.AmbientFront
}
func (mt *Material) SetDiffuse(face int, diffuse [4]float32) {
	switch face {
	case GLFRONT:
		mt.DiffuseFrontAndBack = false
		mt.DiffuseFront = diffuse
		break
	case GLBACK:
		mt.DiffuseFrontAndBack = false
		mt.DiffuseBack = diffuse
		break
	case GLFRONTANDBACK:
		mt.DiffuseFrontAndBack = true
		mt.DiffuseFront = diffuse
		mt.DiffuseBack = mt.DiffuseFront
		break
	}
}

func (mt *Material) GetDiffuse(face int) *[4]float32 {
	switch face {
	case GLFRONT:
		return &mt.DiffuseFront
	case GLBACK:
		return &mt.DiffuseBack
	case GLFRONTANDBACK:
		if !mt.DiffuseFrontAndBack {
			return &mt.DiffuseFront
		}
	}
	return &mt.DiffuseFront
}

func (mt *Material) SetSpecular(face int, specular [4]float32) {
	switch face {
	case GLFRONT:
		mt.SpecularFrontAndBack = false
		mt.SpecularFront = specular
		break
	case GLBACK:
		mt.SpecularFrontAndBack = false
		mt.SpecularBack = specular
		break
	case GLFRONTANDBACK:
		mt.SpecularFrontAndBack = true
		mt.SpecularFront = specular
		mt.SpecularBack = mt.SpecularFront
		break
	default:
		break
	}
}

func (mt *Material) GetSpecular(face int) *[4]float32 {
	switch face {
	case GLFRONT:
		return &mt.SpecularFront
	case GLBACK:
		return &mt.SpecularBack
	case GLFRONTANDBACK:
		if !mt.SpecularFrontAndBack {
			return &mt.SpecularFront
		}
	}
	return &mt.SpecularFront
}

func (mt *Material) SetEmission(face int, emission [4]float32) {
	switch face {
	case GLFRONT:
		mt.EmissionFrontAndBack = false
		mt.EmissionFront = emission
		break
	case GLBACK:
		mt.EmissionFrontAndBack = false
		mt.EmissionBack = emission
		break
	case GLFRONTANDBACK:
		mt.EmissionFrontAndBack = true
		mt.EmissionFront = emission
		mt.EmissionBack = mt.EmissionFront
		break
	default:
		break
	}
}

func (mt *Material) GetEmission(face int) *[4]float32 {
	switch face {
	case GLFRONT:
		return &mt.EmissionFront
	case GLBACK:
		return &mt.EmissionBack
	case GLFRONTANDBACK:
		if !mt.EmissionFrontAndBack {
		}
		return &mt.EmissionFront
	}
	return &mt.EmissionFront
}

func (mt *Material) SetShininess(face int, shininess float32) {
	if shininess < 0 {
		shininess = 0.0
	} else if shininess > 128.0 {
		shininess = 128.0
	}

	switch face {
	case GLFRONT:
		mt.ShininessFrontAndBack = false
		mt.ShininessFront = shininess
		break
	case GLBACK:
		mt.ShininessFrontAndBack = false
		mt.ShininessBack = shininess
		break
	case GLFRONTANDBACK:
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
	case GLFRONT:
		return mt.ShininessFront
	case GLBACK:
		return mt.ShininessBack
	case GLFRONTANDBACK:
		if !mt.ShininessFrontAndBack {
		}
		return mt.ShininessFront
	}
	return mt.ShininessFront
}

func (mt *Material) SetTransparency(face int, transparency float32) {
	if face == GLFRONT || face == GLFRONTANDBACK {
		mt.AmbientFront[3] = 1.0 - transparency
		mt.DiffuseFront[3] = 1.0 - transparency
		mt.SpecularFront[3] = 1.0 - transparency
		mt.EmissionFront[3] = 1.0 - transparency
	}

	if face == GLBACK || face == GLFRONTANDBACK {
		mt.AmbientBack[3] = 1.0 - transparency
		mt.DiffuseBack[3] = 1.0 - transparency
		mt.SpecularBack[3] = 1.0 - transparency
		mt.EmissionBack[3] = 1.0 - transparency
	}
}
