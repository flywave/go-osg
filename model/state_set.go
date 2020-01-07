package model

type ModeListType map[int32]int32

type RefAttributePair struct {
	First  StateAttributeInterface
	Second int32
}

type AttributeListType map[int32]*RefAttributePair

type TextureModeListType []ModeListType

type TextureAttributeListType []AttributeListType

type RefUniformPair struct {
	First  *UniformBase
	Second int32
}

type UniformListType map[int32]RefUniformPair

type DefinePair struct {
	First  string
	Second int32
}

type DefineListType map[string]DefinePair

const (
	DEFAULTBIN     = 0
	OPAQUEBIN      = 1
	TRANSPARENTBIN = 2

	INHERITRENDERBINDETAILS           = 0
	USERENDERBINDETAILS               = 1
	OVERRIDERENDERBINDETAILS          = 2
	PROTECTEDRENDERBINDETAILS         = 4
	OVERRIDEPROTECTEDRENDERBINDETAILS = 6

	STATESETT string = "osg::StateSet"
)

var textureGlmodeMap map[int32]bool

func init() {
	textureGlmodeMap = make(map[int32]bool)
	textureGlmodeMap[GLTEXTURE1D] = true
	textureGlmodeMap[GLTEXTURE2D] = true
	textureGlmodeMap[GLTEXTURE3D] = true
	textureGlmodeMap[GLTEXTUREBUFFER] = true

	textureGlmodeMap[GLTEXTURECUBEMAP] = true
	textureGlmodeMap[GLTEXTURERECTANGLE] = true
	textureGlmodeMap[GLTEXTURE2DARRAY] = true

	textureGlmodeMap[GLTEXTUREGENQ] = true
	textureGlmodeMap[GLTEXTUREGENR] = true
	textureGlmodeMap[GLTEXTUREGENS] = true
	textureGlmodeMap[GLTEXTUREGENT] = true
}

type StateSet struct {
	Object
	Parents              []interface{}
	ModeList             ModeListType
	AttributeList        AttributeListType
	TextureModeList      TextureModeListType
	TextureAttributeList TextureAttributeListType
	UniformList          UniformListType
	DefineList           DefineListType

	RenderingHint int32
	BinMode       int32

	BinNum         int32
	BinName        string
	NestRenderBins bool

	UpdateCallback *Callback
	EventCallback  *Callback
}

func NewStateSet() *StateSet {
	obj := NewObject()
	obj.Type = STATESETT
	return &StateSet{Object: *obj, ModeList: make(ModeListType), AttributeList: make(AttributeListType), RenderingHint: DEFAULTBIN, BinMode: INHERITRENDERBINDETAILS, NestRenderBins: true, BinNum: 0, BinName: ""}
}

func (ss *StateSet) setMode3(unit int, mode int32, val int32) {
	l := len(ss.TextureModeList)
	if l <= unit {
		s := 1 + unit - l
		for i := 0; i < s; i++ {
			ss.TextureModeList = append(ss.TextureModeList, make(ModeListType))
		}
	}
	list := ss.TextureModeList[unit]
	if (val & INHERIT) > 0 {
		_, ok := list[mode]
		if ok {
			delete(list, mode)
		}
	} else {
		list[mode] = val
	}
}

func (ss *StateSet) setMode2(mode int32, val int32) {
	_, ok := textureGlmodeMap[mode]
	if ok {
		ss.SetTextureMode(0, mode, val)
	} else if mode == GLCOLORMATERIAL {

	} else {
		if (val & INHERIT) > 0 {
			_, ok := ss.ModeList[mode]
			if ok {
				delete(ss.ModeList, mode)
			}
		} else {
			ss.ModeList[mode] = val
		}
	}

}
func (ss *StateSet) SetTextureMode(unit int, mode int32, val int32) {
	_, ok := textureGlmodeMap[mode]
	if ok {
		ss.setMode3(unit, mode, val)
	} else {
		ss.setMode2(mode, val)
	}
}

func (ss *StateSet) IsTextureAttribute() bool {
	return false
}

func (ss *StateSet) createOrGetAttributeList(unit int) AttributeListType {
	l := len(ss.TextureAttributeList)
	if unit >= l {
		sz := 1 + unit - l
		for i := 0; i < sz; i++ {
			ss.TextureAttributeList = append(ss.TextureAttributeList, make(AttributeListType))
		}
	}
	return ss.TextureAttributeList[unit]
}

func (ss *StateSet) setAttribute3(lst AttributeListType, attr interface{}, val int32) {
	if attr != nil {
		key := attr.(StateAttributeInterface).GetType()
		par, ok := lst[key]
		if ok {
			par.Second = val & (OVERRIDE | PROTECTED)
		} else {
			ref := &RefAttributePair{First: attr.(StateAttributeInterface), Second: val & (OVERRIDE | PROTECTED)}
			lst[key] = ref
		}
	}
}

func (ss *StateSet) setAttribute2(attr interface{}, val int32) {
	if attr != nil {
		if attr.(StateAttributeInterface).IsTextureAttribute() {
			ss.SetTextureAttribute(0, attr, val)
		} else {
			ss.setAttribute3(ss.AttributeList, attr, val)
		}
	}
}

func (ss *StateSet) SetTextureAttribute(unit int, attr interface{}, val int32) {
	if attr != nil {
		if attr.(StateAttributeInterface).IsTextureAttribute() {
			ss.setAttribute3(ss.createOrGetAttributeList(unit), attr, val)
		} else {
			ss.setAttribute2(attr, val)
		}
	}
}

func (ss *StateSet) SetDefine(k string, first string, value int32) {
	dp := DefinePair{First: first, Second: value}
	ss.DefineList[k] = dp
}
