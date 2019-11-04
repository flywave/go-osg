package model

type RenderingHint uint32
type RenderBinMode uint32

type ModeListType map[Glmode]GlmodeValue

type RefAttributePair struct {
	First  *StateAttribute
	Second OverrideValue
}

type AttributeListType map[TypeMemberPair]RefAttributePair

type TextureModeListType []ModeListType

type TextureAttributeListType []AttributeListType

type RefUniformPair struct {
	First  *UniformBase
	Second OverrideValue
}

type UniformListType map[string]RefUniformPair

type DefinePair struct {
	First  string
	Second OverrideValue
}

type DefineListType map[string]DefinePair

const (
	DEFAULT_BIN     RenderingHint = 0
	OPAQUE_BIN      RenderingHint = 1
	TRANSPARENT_BIN RenderingHint = 2

	INHERIT_RENDERBIN_DETAILS            RenderBinMode = 0
	USE_RENDERBIN_DETAILS                RenderBinMode = 1
	OVERRIDE_RENDERBIN_DETAILS           RenderBinMode = 2
	PROTECTED_RENDERBIN_DETAILS          RenderBinMode = 4
	OVERRIDE_PROTECTED_RENDERBIN_DETAILS RenderBinMode = 6

	StateSetType string = "osg::StateSet"
)

type StateSet struct {
	Object
	Parents              []*Node
	ModeList             ModeListType
	AttributeList        AttributeListType
	TextureModeList      TextureModeListType
	TextureAttributeList TextureAttributeListType
	UniformList          UniformListType
	DefineList           DefineListType

	RenderingHint RenderingHint
	BinMode       RenderBinMode

	BinNum         int
	BinName        string
	NestRenderBins bool

	UpdateCallback *Callback
	EventCallback  *Callback
}

func NewStateSet() StateSet {
	obj := NewObject()
	obj.Type = StateSetType
	return StateSet{Object: obj, RenderingHint: DEFAULT_BIN, BinMode: INHERIT_RENDERBIN_DETAILS, NestRenderBins: true, BinNum: 0, BinName: ""}
}
