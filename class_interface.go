package osg

import (
	"github.com/flywave/go-osg/model"
)

func GetTypeString(iter interface{}) string {
	switch iter.(type) {
	case *model.Image:
		return "IMAGE"
	case *model.Object:
		return "OBJECT"
	case bool:
		return "BOOL"
	case uint8:
		return "UCHAR"
	case int16:
		return "SHORT"
	case int:
		return "INT"
	case uint:
		return "UINT"
	case float32:
		return "FLOAT"
	case float64:
		return "DOUBLE"
	case [2]float32:
		return "VEC2F"
	case [2]float64:
		return "VEC2D"
	case [3]float32:
		return "VEC3F"
	case [3]float64:
		return "VEC3D"
	case [4]float32:
		return "VEC4F"
	case [4]float64:
		return "VEC4D"
	case *model.Quaternion:
		return "QUAT"
	case *model.Planef:
		return "PLANE"
	case string:
		return "STRING"
	case [4][4]float32:
		return "MATRIXF"
	case [4][4]float64:
		return "MATRIXD"
	case [2]byte:
		return "VEC2UB"
	case [2]int8:
		return "VEC2B"
	case [2]int16:
		return "VEC2S"
	case [2]uint16:
		return "VEC2US"
	case [2]int32:
		return "VEC2I"
	case [2]uint32:
		return "VEC2UI"
	case [3]byte:
		return "VEC3UB"
	case [3]int8:
		return "VEC3B"
	case [3]int16:
		return "VEC3S"
	case [3]uint16:
		return "VEC3US"
	case [3]int32:
		return "VEC3I"
	case [3]uint32:
		return "VEC3UI"
	case [4]byte:
		return "VEC4UB"
	case [4]int8:
		return "VEC4B"
	case [4]int16:
		return "VEC4S"
	case [4]uint16:
		return "VEC4US"
	case [4]int32:
		return "VEC4I"
	case [4]uint32:
		return "VEC4UI"
	case [2][3]float32:
		return "BOUNDINGBOXF"
	case [2][3]float64:
		return "BOUNDINGBOXD"
	case *model.Sphere3f:
		return "BOUNDINGSPHEREF"
	case *model.Sphere3d:
		return "BOUNDINGSPHERED"
	}
	return "UNDEFINED"
}

func GetTypeEnum(t string) SerType {
	switch t {
	case "IMAGE":
		return RWIMAGE
	case "Object":
		return RWOBJECT
	case "BOOL":
		return RWBOOL
	case "UCHAR":
		return RWUCHAR
	case "SHORT":
		return RWSHORT
	case "INT":
		return RWINT
	case "UINT":
		return RWUINT
	case "FLOAT":
		return RWFLOAT
	case "DOUBLE":
		return RWDOUBLE
	case "VEC2F":
		return RWVEC2F
	case "VEC2D":
		return RWVEC2D
	case "VEC3F":
		return RWVEC3F
	case "VEC3D":
		return RWVEC3D
	case "VEC4F":
		return RWVEC4F
	case "VEC4D":
		return RWVEC4D
	case "QUAT":
		return RWQUAT
	case "PLANE":
		return RWPLANE
	case "STRING":
		return RWSTRING
	case "MATRIXF":
		return RWMATRIXF
	case "MATRIXD":
		return RWMATRIXD
	case "VEC2UB":
		return RWVEC2UB
	case "VEC2B":
		return RWVEC2B
	case "VEC2S":
		return RWVEC2S
	case "VEC2US":
		return RWVEC2US
	case "VEC2I":
		return RWVEC2I
	case "VEC2UI":
		return RWVEC2UI
	case "VEC3UB":
		return RWVEC3UB
	case "VEC3B":
		return RWVEC3B
	case "VEC3S":
		return RWVEC3S
	case "VEC3US":
		return RWVEC3US
	case "VEC3I":
		return RWVEC3I
	case "VEC3UI":
		return RWVEC3UI
	case "VEC4UB":
		return RWVEC4UB
	case "VEC4B":
		return RWVEC4B
	case "VEC4S":
		return RWVEC4S
	case "VEC4US":
		return RWVEC4US
	case "VEC4I":
		return RWVEC4I
	case "VEC4UI":
		return RWVEC4UI
	case "BOUNDINGBOXF":
		return RWBOUNDINGBOXF
	case "BOUNDINGBOXD":
		return RWBOUNDINGBOXD
	case "BOUNDINGSPHEREF":
		return RWBOUNDINGSPHEREF
	case "BOUNDINGSPHERED":
		return RWBOUNDINGSPHERED
	}
	return RWUNDEFINED
}

type PropertyMapType map[string]SerType
type TypeMapType map[SerType]string

type ObjectPropertyMapType map[string]PropertyMapType

type ClassInterface struct {
	PropertyMap PropertyMapType
	TypeMap     TypeMapType

	WhiteList ObjectPropertyMapType
	BlackList ObjectPropertyMapType
}

func NewClassInterface() *ClassInterface {
	cf := &ClassInterface{PropertyMap: make(PropertyMapType), TypeMap: make(TypeMapType), WhiteList: make(ObjectPropertyMapType), BlackList: make(ObjectPropertyMapType)}
	cf.TypeMap[RWUNDEFINED] = "UNDEFINED"
	cf.PropertyMap["UNDEFINED"] = RWUNDEFINED

	cf.TypeMap[RWUSER] = "USER"
	cf.PropertyMap["USER"] = RWUSER

	cf.TypeMap[RWOBJECT] = "OBJECT"
	cf.PropertyMap["OBJECT"] = RWOBJECT

	cf.TypeMap[RWIMAGE] = "IMAGE"
	cf.PropertyMap["IMAGE"] = RWIMAGE

	cf.TypeMap[RWLIST] = "LIST"
	cf.PropertyMap["LIST"] = RWLIST

	cf.TypeMap[RWBOOL] = "BOOL"
	cf.PropertyMap["BOOL"] = RWBOOL

	cf.TypeMap[RWBOOL] = "BOOL"
	cf.PropertyMap["BOOL"] = RWBOOL

	cf.TypeMap[RWCHAR] = "CHAR"
	cf.PropertyMap["CHAR"] = RWCHAR

	cf.TypeMap[RWUCHAR] = "UCHAR"
	cf.PropertyMap["UCHAR"] = RWUCHAR

	cf.TypeMap[RWSHORT] = "SHORT"
	cf.PropertyMap["SHORT"] = RWSHORT

	cf.TypeMap[RWUSHORT] = "USHORT"
	cf.PropertyMap["USHORT"] = RWUSHORT

	cf.TypeMap[RWINT] = "INT"
	cf.PropertyMap["INT"] = RWINT

	cf.TypeMap[RWUINT] = "UINT"
	cf.PropertyMap["UINT"] = RWUINT

	cf.TypeMap[RWFLOAT] = "FLOAT"
	cf.PropertyMap["FLOAT"] = RWFLOAT

	cf.TypeMap[RWDOUBLE] = "DOUBLE"
	cf.PropertyMap["DOUBLE"] = RWDOUBLE

	cf.TypeMap[RWVEC2F] = "VEC2F"
	cf.PropertyMap["VEC2F"] = RWVEC2F

	cf.TypeMap[RWVEC2D] = "VEC2D"
	cf.PropertyMap["VEC2D"] = RWVEC2D

	cf.TypeMap[RWVEC3F] = "VEC3F"
	cf.PropertyMap["VEC3F"] = RWVEC3F

	cf.TypeMap[RWVEC3D] = "VEC3D"
	cf.PropertyMap["VEC3D"] = RWVEC3D

	cf.TypeMap[RWVEC4F] = "VEC4F"
	cf.PropertyMap["VEC4F"] = RWVEC4F

	cf.TypeMap[RWVEC4D] = "VEC4D"
	cf.PropertyMap["VEC4D"] = RWVEC4D

	cf.TypeMap[RWQUAT] = "QUAT"
	cf.PropertyMap["QUAT"] = RWQUAT

	cf.TypeMap[RWPLANE] = "PLANE"
	cf.PropertyMap["PLANE"] = RWPLANE

	cf.TypeMap[RWMATRIXF] = "MATRIXF"
	cf.PropertyMap["MATRIXF"] = RWMATRIXF

	cf.TypeMap[RWMATRIXD] = "MATRIXD"
	cf.PropertyMap["MATRIXD"] = RWMATRIXD

	cf.TypeMap[RWMATRIX] = "MATRIX"
	cf.PropertyMap["MATRIX"] = RWMATRIX

	cf.TypeMap[RWBOUNDINGBOXF] = "BOUNDINGBOXF"
	cf.PropertyMap["BOUNDINGBOXF"] = RWBOUNDINGBOXF
	cf.TypeMap[RWBOUNDINGBOXD] = "BOUNDINGBOXD"
	cf.PropertyMap["BOUNDINGBOXD"] = RWBOUNDINGBOXD

	cf.TypeMap[RWBOUNDINGSPHEREF] = "BOUNDINGSPHEREF"
	cf.PropertyMap["BOUNDINGSPHEREF"] = RWBOUNDINGSPHEREF

	cf.TypeMap[RWBOUNDINGSPHERED] = "BOUNDINGSPHERED"
	cf.PropertyMap["BOUNDINGSPHERED"] = RWBOUNDINGSPHERED

	cf.TypeMap[RWGLENUM] = "GLENUM"
	cf.PropertyMap["GLENUM"] = RWGLENUM

	cf.TypeMap[RWSTRING] = "STRING"
	cf.PropertyMap["STRING"] = RWSTRING

	cf.TypeMap[RWENUM] = "ENUM"
	cf.PropertyMap["ENUM"] = RWENUM

	cf.TypeMap[RWVEC2B] = "VEC2B"
	cf.PropertyMap["VEC2B"] = RWVEC2B

	cf.TypeMap[RWVEC2UB] = "VEC2UB"
	cf.PropertyMap["VEC2UB"] = RWVEC2UB

	cf.TypeMap[RWVEC2S] = "VEC2S"
	cf.PropertyMap["VEC2S"] = RWVEC2S

	cf.TypeMap[RWVEC2US] = "VEC2US"
	cf.PropertyMap["VEC2US"] = RWVEC2US

	cf.TypeMap[RWVEC2I] = "VEC2I"
	cf.PropertyMap["VEC2I"] = RWVEC2I

	cf.TypeMap[RWVEC2UI] = "VEC2UI"
	cf.PropertyMap["VEC2UI"] = RWVEC2UI

	cf.TypeMap[RWVEC3B] = "VEC3B"
	cf.PropertyMap["VEC3B"] = RWVEC3B

	cf.TypeMap[RWVEC3UB] = "VEC3UB"
	cf.PropertyMap["VEC3UB"] = RWVEC3UB

	cf.TypeMap[RWVEC3S] = "VEC3S"
	cf.PropertyMap["VEC3S"] = RWVEC3S

	cf.TypeMap[RWVEC3US] = "VEC3US"
	cf.PropertyMap["VEC3US"] = RWVEC3US

	cf.TypeMap[RWVEC3I] = "VEC3I"
	cf.PropertyMap["VEC3I"] = RWVEC3I

	cf.TypeMap[RWVEC3UI] = "VEC3UI"
	cf.PropertyMap["VEC3UI"] = RWVEC3UI

	cf.TypeMap[RWVEC4B] = "VEC4B"
	cf.PropertyMap["VEC4B"] = RWVEC4B

	cf.TypeMap[RWVEC4UB] = "VEC4UB"
	cf.PropertyMap["VEC4UB"] = RWVEC4UB

	cf.TypeMap[RWVEC4S] = "VEC4S"
	cf.PropertyMap["VEC4S"] = RWVEC4S

	cf.TypeMap[RWVEC4US] = "VEC4US"
	cf.PropertyMap["VEC4US"] = RWVEC4US

	cf.TypeMap[RWVEC4I] = "VEC4I"
	cf.PropertyMap["VEC4I"] = RWVEC4I

	cf.TypeMap[RWVEC4UI] = "VEC4UI"
	cf.PropertyMap["VEC4UI"] = RWVEC4UI

	cf.TypeMap[RWLIST] = "LIST"
	cf.PropertyMap["LIST"] = RWLIST

	cf.TypeMap[RWVECTOR] = "VECTOR"
	cf.PropertyMap["VECTOR"] = RWVECTOR

	cf.TypeMap[RWMAP] = "MAP"
	cf.PropertyMap["MAP"] = RWMAP
	return cf
}

func (cl *ClassInterface) AreTypesCompatible(lhs SerType, rhs SerType) bool {
	if lhs == rhs {
		return true
	}

	if lhs == RWMATRIX {
		lhs = RWMATRIXD
	}
	if rhs == RWMATRIX {
		rhs = RWMATRIXD
	}

	if lhs == RWGLENUM {
		lhs = RWUINT
	}
	if rhs == RWGLENUM {
		rhs = RWUINT
	}

	if lhs == RWENUM {
		lhs = RWINT
	}
	if rhs == RWENUM {
		rhs = RWINT
	}

	if lhs == RWIMAGE {
		lhs = RWOBJECT
	}

	return lhs == rhs
}

func (cl *ClassInterface) GetTypeName(t SerType) string {
	str, ok := cl.TypeMap[t]
	if ok {
		return str
	}
	return ""
}

func (cl *ClassInterface) GetType(name string) SerType {
	t, ok := cl.PropertyMap[name]
	if ok {
		return t
	}
	return RWUNDEFINED
}

func (cl *ClassInterface) GetObjectWrapper(obj *model.Object) *ObjectWrapper {
	return GetObjectWrapperManager().FindWrap(obj.Name)
}

func (cl *ClassInterface) GetSerializer(obj *model.Object, propertyName string, t *SerType) *BaseSerializer {
	ow := cl.GetObjectWrapper(obj)
	if ow != nil {
		ow.GetSerializerAndType(propertyName, t)
	}
	return nil
}

func (cl *ClassInterface) CreateObject(compoundClassName string) interface{} {
	ow := GetObjectWrapperManager().FindWrap(compoundClassName)
	if ow != nil {
		return ow.CreateInstance()
	}
	return nil
}
