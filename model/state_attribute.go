package model

type Type uint32
type Values uint32
type Glmode int
type GlmodeValue uint32
type OverrideValue uint32

const (
	TEXTURE                        Type = 0
	POLYGONMODE                    Type = 1
	POLYGONOFFSET                  Type = 2
	MATERIAL                       Type = 3
	ALPHAFUNC                      Type = 4
	ANTIALIAS                      Type = 5
	COLORTABLE                     Type = 6
	CULLFACE                       Type = 7
	FOG                            Type = 8
	FRONTFACE                      Type = 9
	LIGHT                          Type = 10
	POINT                          Type = 11
	LINEWIDTH                      Type = 12
	LINESTIPPLE                    Type = 13
	POLYGONSTIPPLE                 Type = 14
	SHADEMODEL                     Type = 15
	TEXENV                         Type = 16
	TEXENVFILTER                   Type = 17
	TEXGEN                         Type = 18
	TEXMAT                         Type = 19
	LIGHTMODEL                     Type = 20
	BLENDFUNC                      Type = 21
	BLENDEQUATION                  Type = 22
	LOGICOP                        Type = 23
	STENCIL                        Type = 24
	COLORMASK                      Type = 25
	DEPTH                          Type = 26
	VIEWPORT                       Type = 27
	SCISSOR                        Type = 28
	BLENDCOLOR                     Type = 29
	MULTISAMPLE                    Type = 30
	CLIPPLANE                      Type = 31
	COLORMATRIX                    Type = 32
	VERTEXPROGRAM                  Type = 33
	FRAGMENTPROGRAM                Type = 34
	POINTSPRITE                    Type = 35
	PROGRAM                        Type = 36
	CLAMPCOLOR                     Type = 37
	HINT                           Type = 38
	SAMPLEMASKI                    Type = 39
	PRIMITIVERESTARTINDEX          Type = 40
	CLIPCONTROL                    Type = 41
	VALIDATOR                      Type = 42
	VIEWMATRIXEXTRACTOR            Type = 43
	OSGNV_PARAMETER_BLOCK          Type = 44
	OSGNVEXT_TEXTURE_SHADER        Type = 45
	OSGNVEXT_VERTEX_PROGRAM        Type = 46
	OSGNVEXT_REGISTER_COMBINERS    Type = 47
	OSGNVCG_PROGRAM                Type = 48
	OSGNVSLANG_PROGRAM             Type = 49
	OSGNVPARSE_PROGRAM_PARSER      Type = 50
	UNIFORMBUFFERBINDING           Type = 51
	TRANSFORMFEEDBACKBUFFERBINDING Type = 52
	ATOMICCOUNTERBUFFERBINDING     Type = 53
	PATCH_PARAMETER                Type = 54
	FRAME_BUFFER_OBJECT            Type = 55
	VERTEX_ATTRIB_DIVISOR          Type = 56
	SHADERSTORAGEBUFFERBINDING     Type = 57
	INDIRECTDRAWBUFFERBINDING      Type = 58
	VIEWPORTINDEXED                Type = 59
	DEPTHRANGEINDEXED              Type = 60
	SCISSORINDEXED                 Type = 61
	BINDIMAGETEXTURE               Type = 62
	SAMPLER                        Type = 63
	CAPABILITY                     Type = 64

	OFF       Values = 0x0
	ON        Values = 0x1
	OVERRIDE  Values = 0x2
	PROTECTED Values = 0x4
	INHERIT   Values = 0x8

	FLAT   uint = 0x1D00
	SMOOTH uint = 0x1D01

	SHADEMODEL_T     string = "osg::ShadeModel"
	STATEATRRIBUTE_T string = "osg::StateAttribute"
)

type StateAttribute struct {
	Object
	Parents        []*StateSet
	UpdateCallback *Callback
	EventCallback  *Callback
}

func NewStateAttribute() StateAttribute {
	obj := NewObject()
	obj.Type = STATEATRRIBUTE_T
	return StateAttribute{Object: obj}
}

type TypeMemberPair struct {
	First  Type
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
	Mode uint
}

func NewShadeModel() ShadeModel {
	a := NewStateAttribute()
	a.Type = SHADEMODEL_T
	return ShadeModel{StateAttribute: a}
}
