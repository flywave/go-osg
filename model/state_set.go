package model

type RenderingHint uint32

type ModeListType map[int]int

type RefAttributePair struct {
	First  *StateAttribute
	Second int
}

type AttributeListType map[string]RefAttributePair

type TextureModeListType []ModeListType

type TextureAttributeListType []AttributeListType

type RefUniformPair struct {
	First  *UniformBase
	Second int
}

type UniformListType map[string]RefUniformPair

type DefinePair struct {
	First  string
	Second int
}

type DefineListType map[string]DefinePair

const (
	DEFAULT_BIN     RenderingHint = 0
	OPAQUE_BIN      RenderingHint = 1
	TRANSPARENT_BIN RenderingHint = 2

	INHERIT_RENDERBIN_DETAILS            = 0
	USE_RENDERBIN_DETAILS                = 1
	OVERRIDE_RENDERBIN_DETAILS           = 2
	PROTECTED_RENDERBIN_DETAILS          = 4
	OVERRIDE_PROTECTED_RENDERBIN_DETAILS = 6

	STATESET_T string = "osg::StateSet"
)

type StateSet struct {
	Object
	Parents              []interface{}
	ModeList             ModeListType
	AttributeList        AttributeListType
	TextureModeList      TextureModeListType
	TextureAttributeList TextureAttributeListType
	UniformList          UniformListType
	DefineList           DefineListType

	RenderingHint RenderingHint
	BinMode       int

	BinNum         int
	BinName        string
	NestRenderBins bool

	UpdateCallback *Callback
	EventCallback  *Callback
}

func NewStateSet() StateSet {
	obj := NewObject()
	obj.Type = STATESET_T
	return StateSet{Object: obj, RenderingHint: DEFAULT_BIN, BinMode: INHERIT_RENDERBIN_DETAILS, NestRenderBins: true, BinNum: 0, BinName: ""}
}
