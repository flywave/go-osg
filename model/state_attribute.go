package model

import "reflect"

const (
	TEXTURE                        = 0
	POLYGONMODE                    = 1
	POLYGONOFFSET                  = 2
	MATERIAL                       = 3
	ALPHAFUNC                      = 4
	ANTIALIAS                      = 5
	COLORTABLE                     = 6
	CULLFACE                       = 7
	FOG                            = 8
	FRONTFACE                      = 9
	LIGHT                          = 10
	POINT                          = 11
	LINEWIDTH                      = 12
	LINESTIPPLE                    = 13
	POLYGONSTIPPLE                 = 14
	SHADEMODEL                     = 15
	TEXENV                         = 16
	TEXENVFILTER                   = 17
	TEXGEN                         = 18
	TEXMAT                         = 19
	LIGHTMODEL                     = 20
	BLENDFUNC                      = 21
	BLENDEQUATION                  = 22
	LOGICOP                        = 23
	STENCIL                        = 24
	COLORMASK                      = 25
	DEPTH                          = 26
	VIEWPORT                       = 27
	SCISSOR                        = 28
	BLENDCOLOR                     = 29
	MULTISAMPLE                    = 30
	CLIPPLANE                      = 31
	COLORMATRIX                    = 32
	VERTEXPROGRAM                  = 33
	FRAGMENTPROGRAM                = 34
	POINTSPRITE                    = 35
	PROGRAM                        = 36
	CLAMPCOLOR                     = 37
	HINT                           = 38
	SAMPLEMASKI                    = 39
	PRIMITIVERESTARTINDEX          = 40
	CLIPCONTROL                    = 41
	VALIDATOR                      = 42
	VIEWMATRIXEXTRACTOR            = 43
	OSGNV_PARAMETER_BLOCK          = 44
	OSGNVEXT_TEXTURE_SHADER        = 45
	OSGNVEXT_VERTEX_PROGRAM        = 46
	OSGNVEXT_REGISTER_COMBINERS    = 47
	OSGNVCG_PROGRAM                = 48
	OSGNVSLANG_PROGRAM             = 49
	OSGNVPARSE_PROGRAM_PARSER      = 50
	UNIFORMBUFFERBINDING           = 51
	TRANSFORMFEEDBACKBUFFERBINDING = 52
	ATOMICCOUNTERBUFFERBINDING     = 53
	PATCH_PARAMETER                = 54
	FRAME_BUFFER_OBJECT            = 55
	VERTEX_ATTRIB_DIVISOR          = 56
	SHADERSTORAGEBUFFERBINDING     = 57
	INDIRECTDRAWBUFFERBINDING      = 58
	VIEWPORTINDEXED                = 59
	DEPTHRANGEINDEXED              = 60
	SCISSORINDEXED                 = 61
	BINDIMAGETEXTURE               = 62
	SAMPLER                        = 63
	CAPABILITY                     = 64
	OFF                            = 0x0
	ON                             = 0x1
	OVERRIDE                       = 0x2
	PROTECTED                      = 0x4
	INHERIT                        = 0x8

	FLAT   = 0x1D00
	SMOOTH = 0x1D01

	SHADEMODEL_T     string = "osg::ShadeModel"
	STATEATRRIBUTE_T string = "osg::StateAttribute"
)

type StateAttributeInterface interface {
	IsStateAttributeInterface() bool
}

type StateAttribute struct {
	Object
	Parents        []*StateSet
	UpdateCallback *Callback
	EventCallback  *Callback
}

func (s *StateAttribute) IsStateAttributeInterface() bool {
	return true
}

func NewStateAttribute() StateAttribute {
	obj := NewObject()
	obj.Type = STATEATRRIBUTE_T
	return StateAttribute{Object: obj}
}

type TypeMemberPair struct {
	First  int
	Second int32
}

func (s *StateAttribute) GetTypeMember() TypeMemberPair {
	return TypeMemberPair{First: TEXTURE, Second: 0}
}

func (s *StateAttribute) IsTextureAttribute() bool {
	return false
}

type ShadeModel struct {
	StateAttribute
	Mode int
}

func NewShadeModel() ShadeModel {
	a := NewStateAttribute()
	a.Type = SHADEMODEL_T
	return ShadeModel{StateAttribute: a}
}

func IsStateAttribute(obj interface{}) bool {
	ss := StateAttribute{}
	baset := reflect.TypeOf(ss)
	t := reflect.TypeOf(obj)
	return t.Implements(baset)
}
