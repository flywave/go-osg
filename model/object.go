package model

type DataVariance int32

var Type_Mapping map[string]interface{}

func init() {
	Type_Mapping = make(map[string]interface{})
	Type_Mapping["AlphaFunc"] = NewAlphaFunc()
	Type_Mapping["Array"] = NewArray()
	Type_Mapping["CullFace"] = NewCullFace()
	Type_Mapping["ObjectProperty"] = NewObjectProperty()
	Type_Mapping["ObjectMark"] = NewObjectMark()
	Type_Mapping["ObjectGlenum"] = NewObjectGlenum()
	Type_Mapping["Drawable"] = NewDrawable()
	Type_Mapping["Geode"] = NewGeode()
	Type_Mapping["Geometry"] = NewGeometry()
	Type_Mapping["Group"] = NewGroup()
	Type_Mapping["Image"] = NewImage()
	Type_Mapping["Lod"] = NewLod()
	Type_Mapping["PagedLod"] = NewPagedLod()
	Type_Mapping["Node"] = NewNode()
	Type_Mapping["Object"] = NewObject()
	Type_Mapping["PrimitiveSet"] = NewPrimitiveSet()
	Type_Mapping["Shape"] = NewShape()
	Type_Mapping["StateAttribute"] = NewStateAttribute()
	Type_Mapping["StateSet"] = NewStateSet()
	Type_Mapping["TexGen"] = NewTexGen()
	Type_Mapping["TextureCubeMap"] = NewTextureCubeMap()
	Type_Mapping["TextureRectangle"] = NewTextureRectangle()
	Type_Mapping["Texture"] = NewTexture()
	Type_Mapping["Texture1d"] = NewTexture1d()
	Type_Mapping["Texture2d"] = NewTexture2d()
	Type_Mapping["Texture3d"] = NewTexture3d()
	Type_Mapping["Transform"] = NewTransform()
}

const (
	DYNAMIC     DataVariance = 0
	STATIC      DataVariance = 1
	UNSPECIFIED DataVariance = 2
	OBJECT_T    string       = "osg::Object"
)

type Object struct {
	Name         string
	Type         string
	Propertys    map[string]string
	DataVariance DataVariance
	Udc          *UserDataContainer
}

func NewObject() Object {
	return Object{Type: OBJECT_T, DataVariance: UNSPECIFIED, Propertys: make(map[string]string)}
}

type Callback struct {
	Object
	Callback *Callback
}
