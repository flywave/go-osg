package model

type ModeListType map[int]int

type RefAttributePair struct {
	First  interface{}
	Second int
}

type AttributeListType map[int]*RefAttributePair

type TextureModeListType []ModeListType

type TextureAttributeListType []AttributeListType

type RefUniformPair struct {
	First  *UniformBase
	Second int
}

type UniformListType map[int]RefUniformPair

type DefinePair struct {
	First  string
	Second int
}

type DefineListType map[string]DefinePair

const (
	DEFAULT_BIN     = 0
	OPAQUE_BIN      = 1
	TRANSPARENT_BIN = 2

	INHERIT_RENDERBIN_DETAILS            = 0
	USE_RENDERBIN_DETAILS                = 1
	OVERRIDE_RENDERBIN_DETAILS           = 2
	PROTECTED_RENDERBIN_DETAILS          = 4
	OVERRIDE_PROTECTED_RENDERBIN_DETAILS = 6

	STATESET_T string = "osg::StateSet"
)

var textureGlmodeMap map[int]bool

func init() {
	textureGlmodeMap = make(map[int]bool)
	textureGlmodeMap[GL_TEXTURE_1D] = true
	textureGlmodeMap[GL_TEXTURE_2D] = true
	textureGlmodeMap[GL_TEXTURE_3D] = true
	textureGlmodeMap[GL_TEXTURE_BUFFER] = true

	textureGlmodeMap[GL_TEXTURE_CUBE_MAP] = true
	textureGlmodeMap[GL_TEXTURE_RECTANGLE] = true
	textureGlmodeMap[GL_TEXTURE_2D_ARRAY] = true

	textureGlmodeMap[GL_TEXTURE_GEN_Q] = true
	textureGlmodeMap[GL_TEXTURE_GEN_R] = true
	textureGlmodeMap[GL_TEXTURE_GEN_S] = true
	textureGlmodeMap[GL_TEXTURE_GEN_T] = true
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

	RenderingHint int
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

func (ss *StateSet) setMode3(unit int, mode int, val int) {
	l := len(ss.TextureModeList)
	if l <= unit {
		s := l - 1 - unit
		tmp := make(TextureModeListType, s, s)
		ss.TextureModeList = append(ss.TextureModeList, tmp...)
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

func (ss *StateSet) setMode2(mode int, val int) {
	_, ok := textureGlmodeMap[mode]
	if ok {
		ss.SetTextureMode(0, mode, val)
	} else if mode == GL_COLOR_MATERIAL {

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
func (ss *StateSet) SetTextureMode(unit int, mode int, val int) {
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
		sz := l - 1 - unit
		tmp := make([]AttributeListType, sz, sz)
		ss.TextureAttributeList = append(ss.TextureAttributeList, tmp...)
	}
	lst := ss.TextureAttributeList[unit]
	return lst
}

func (ss *StateSet) setAttribute3(lst AttributeListType, attr interface{}, val int) {

	if attr != nil {
		key := attr.(StateAttributeInterface).GetType()
		par, ok := lst[key]
		if ok {
			par.Second = val & (OVERRIDE | PROTECTED)
		} else {
			lst[key] = &RefAttributePair{First: attr, Second: val & (OVERRIDE | PROTECTED)}
		}
	}
}

func (ss *StateSet) setAttribute2(attr interface{}, val int) {
	if attr != nil {
		if attr.(StateAttributeInterface).IsTextureAttribute() {
			ss.SetTextureAttribute(0, attr, val)
		} else {
			ss.setAttribute3(ss.AttributeList, attr, val)
		}
	}
}

func (ss *StateSet) SetTextureAttribute(unit int, attr interface{}, val int) {
	if attr != nil {
		if attr.(StateAttributeInterface).IsTextureAttribute() {
			ss.setAttribute3(ss.createOrGetAttributeList(unit), attr, val)
		} else {
			ss.setAttribute2(attr, val)
		}
	}
}

func (ss *StateSet) SetDefine(k string, first string, value int) {
	dp := DefinePair{First: first, Second: value}
	ss.DefineList[k] = dp
}
