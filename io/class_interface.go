package io

import (
	"github.com/flywave/go-osg/model"
	"github.com/ungerik/go3d/mat4"
	"github.com/ungerik/go3d/quaternion"
	"github.com/ungerik/go3d/vec2"
	"github.com/ungerik/go3d/vec3"
	"github.com/ungerik/go3d/vec4"
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
	case vec2.T:
		return "VEC2F"
	case [2]float64:
		return "VEC2D"
	case vec3.T:
		return "VEC3F"
	case [3]float64:
		return "VEC3D"
	case vec4.T:
		return "VEC4F"
	case [4]float64:
		return "VEC4D"
	case quaternion.T:
		return "QUAT"
	case model.Planef:
		return "PLANE"
	case string:
		return "STRING"
	case mat4.T:
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
	case model.Sphere3f:
		return "BOUNDINGSPHEREF"
	case model.Sphere3d:
		return "BOUNDINGSPHERED"
	}
	return "UNDEFINED"
}

func GetTypeEnum(t string) SerType {
	switch t {
	case "IMAGE":
		return RW_IMAGE
	case "Object":
		return RW_OBJECT
	case "BOOL":
		return RW_BOOL
	case "UCHAR":
		return RW_UCHAR
	case "SHORT":
		return RW_SHORT
	case "INT":
		return RW_INT
	case "UINT":
		return RW_UINT
	case "FLOAT":
		return RW_FLOAT
	case "DOUBLE":
		return RW_DOUBLE
	case "VEC2F":
		return RW_VEC2F
	case "VEC2D":
		return RW_VEC2D
	case "VEC3F":
		return RW_VEC3F
	case "VEC3D":
		return RW_VEC3D
	case "VEC4F":
		return RW_VEC4F
	case "VEC4D":
		return RW_VEC4D
	case "QUAT":
		return RW_QUAT
	case "PLANE":
		return RW_PLANE
	case "STRING":
		return RW_STRING
	case "MATRIXF":
		return RW_MATRIXF
	case "MATRIXD":
		return RW_MATRIXD
	case "VEC2UB":
		return RW_VEC2UB
	case "VEC2B":
		return RW_VEC2B
	case "VEC2S":
		return RW_VEC2S
	case "VEC2US":
		return RW_VEC2US
	case "VEC2I":
		return RW_VEC2I
	case "VEC2UI":
		return RW_VEC2UI
	case "VEC3UB":
		return RW_VEC3UB
	case "VEC3B":
		return RW_VEC3B
	case "VEC3S":
		return RW_VEC3S
	case "VEC3US":
		return RW_VEC3US
	case "VEC3I":
		return RW_VEC3I
	case "VEC3UI":
		return RW_VEC3UI
	case "VEC4UB":
		return RW_VEC4UB
	case "VEC4B":
		return RW_VEC4B
	case "VEC4S":
		return RW_VEC4S
	case "VEC4US":
		return RW_VEC4US
	case "VEC4I":
		return RW_VEC4I
	case "VEC4UI":
		return RW_VEC4UI
	case "BOUNDINGBOXF":
		return RW_BOUNDINGBOXF
	case "BOUNDINGBOXD":
		return RW_BOUNDINGBOXD
	case "BOUNDINGSPHEREF":
		return RW_BOUNDINGSPHEREF
	case "BOUNDINGSPHERED":
		return RW_BOUNDINGSPHERED
	}
	return RW_UNDEFINED
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

func NewClassInterface() ClassInterface {
	cf := ClassInterface{PropertyMap: make(PropertyMapType), TypeMap: make(TypeMapType), WhiteList: make(ObjectPropertyMapType), BlackList: make(ObjectPropertyMapType)}
	cf.TypeMap[RW_UNDEFINED] = "UNDEFINED"
	cf.PropertyMap["UNDEFINED"] = RW_UNDEFINED

	cf.TypeMap[RW_USER] = "USER"
	cf.PropertyMap["USER"] = RW_USER

	cf.TypeMap[RW_OBJECT] = "OBJECT"
	cf.PropertyMap["OBJECT"] = RW_OBJECT

	cf.TypeMap[RW_IMAGE] = "IMAGE"
	cf.PropertyMap["IMAGE"] = RW_IMAGE

	cf.TypeMap[RW_LIST] = "LIST"
	cf.PropertyMap["LIST"] = RW_LIST

	cf.TypeMap[RW_BOOL] = "BOOL"
	cf.PropertyMap["BOOL"] = RW_BOOL

	cf.TypeMap[RW_BOOL] = "BOOL"
	cf.PropertyMap["BOOL"] = RW_BOOL

	cf.TypeMap[RW_CHAR] = "CHAR"
	cf.PropertyMap["CHAR"] = RW_CHAR

	cf.TypeMap[RW_UCHAR] = "UCHAR"
	cf.PropertyMap["UCHAR"] = RW_UCHAR

	cf.TypeMap[RW_SHORT] = "SHORT"
	cf.PropertyMap["SHORT"] = RW_SHORT

	cf.TypeMap[RW_USHORT] = "USHORT"
	cf.PropertyMap["USHORT"] = RW_USHORT

	cf.TypeMap[RW_INT] = "INT"
	cf.PropertyMap["INT"] = RW_INT

	cf.TypeMap[RW_UINT] = "UINT"
	cf.PropertyMap["UINT"] = RW_UINT

	cf.TypeMap[RW_FLOAT] = "FLOAT"
	cf.PropertyMap["FLOAT"] = RW_FLOAT

	cf.TypeMap[RW_DOUBLE] = "DOUBLE"
	cf.PropertyMap["DOUBLE"] = RW_DOUBLE

	cf.TypeMap[RW_VEC2F] = "VEC2F"
	cf.PropertyMap["VEC2F"] = RW_VEC2F

	cf.TypeMap[RW_VEC2D] = "VEC2D"
	cf.PropertyMap["VEC2D"] = RW_VEC2D

	cf.TypeMap[RW_VEC3F] = "VEC3F"
	cf.PropertyMap["VEC3F"] = RW_VEC3F

	cf.TypeMap[RW_VEC3D] = "VEC3D"
	cf.PropertyMap["VEC3D"] = RW_VEC3D

	cf.TypeMap[RW_VEC4F] = "VEC4F"
	cf.PropertyMap["VEC4F"] = RW_VEC4F

	cf.TypeMap[RW_VEC4D] = "VEC4D"
	cf.PropertyMap["VEC4D"] = RW_VEC4D

	cf.TypeMap[RW_QUAT] = "QUAT"
	cf.PropertyMap["QUAT"] = RW_QUAT

	cf.TypeMap[RW_PLANE] = "PLANE"
	cf.PropertyMap["PLANE"] = RW_PLANE

	cf.TypeMap[RW_MATRIXF] = "MATRIXF"
	cf.PropertyMap["MATRIXF"] = RW_MATRIXF

	cf.TypeMap[RW_MATRIXD] = "MATRIXD"
	cf.PropertyMap["MATRIXD"] = RW_MATRIXD

	cf.TypeMap[RW_MATRIX] = "MATRIX"
	cf.PropertyMap["MATRIX"] = RW_MATRIX

	cf.TypeMap[RW_BOUNDINGBOXF] = "BOUNDINGBOXF"
	cf.PropertyMap["BOUNDINGBOXF"] = RW_BOUNDINGBOXF
	cf.TypeMap[RW_BOUNDINGBOXD] = "BOUNDINGBOXD"
	cf.PropertyMap["BOUNDINGBOXD"] = RW_BOUNDINGBOXD

	cf.TypeMap[RW_BOUNDINGSPHEREF] = "BOUNDINGSPHEREF"
	cf.PropertyMap["BOUNDINGSPHEREF"] = RW_BOUNDINGSPHEREF

	cf.TypeMap[RW_BOUNDINGSPHERED] = "BOUNDINGSPHERED"
	cf.PropertyMap["BOUNDINGSPHERED"] = RW_BOUNDINGSPHERED

	cf.TypeMap[RW_GLENUM] = "GLENUM"
	cf.PropertyMap["GLENUM"] = RW_GLENUM

	cf.TypeMap[RW_STRING] = "STRING"
	cf.PropertyMap["STRING"] = RW_STRING

	cf.TypeMap[RW_ENUM] = "ENUM"
	cf.PropertyMap["ENUM"] = RW_ENUM

	cf.TypeMap[RW_VEC2B] = "VEC2B"
	cf.PropertyMap["VEC2B"] = RW_VEC2B

	cf.TypeMap[RW_VEC2UB] = "VEC2UB"
	cf.PropertyMap["VEC2UB"] = RW_VEC2UB

	cf.TypeMap[RW_VEC2S] = "VEC2S"
	cf.PropertyMap["VEC2S"] = RW_VEC2S

	cf.TypeMap[RW_VEC2US] = "VEC2US"
	cf.PropertyMap["VEC2US"] = RW_VEC2US

	cf.TypeMap[RW_VEC2I] = "VEC2I"
	cf.PropertyMap["VEC2I"] = RW_VEC2I

	cf.TypeMap[RW_VEC2UI] = "VEC2UI"
	cf.PropertyMap["VEC2UI"] = RW_VEC2UI

	cf.TypeMap[RW_VEC3B] = "VEC3B"
	cf.PropertyMap["VEC3B"] = RW_VEC3B

	cf.TypeMap[RW_VEC3UB] = "VEC3UB"
	cf.PropertyMap["VEC3UB"] = RW_VEC3UB

	cf.TypeMap[RW_VEC3S] = "VEC3S"
	cf.PropertyMap["VEC3S"] = RW_VEC3S

	cf.TypeMap[RW_VEC3US] = "VEC3US"
	cf.PropertyMap["VEC3US"] = RW_VEC3US

	cf.TypeMap[RW_VEC3I] = "VEC3I"
	cf.PropertyMap["VEC3I"] = RW_VEC3I

	cf.TypeMap[RW_VEC3UI] = "VEC3UI"
	cf.PropertyMap["VEC3UI"] = RW_VEC3UI

	cf.TypeMap[RW_VEC4B] = "VEC4B"
	cf.PropertyMap["VEC4B"] = RW_VEC4B

	cf.TypeMap[RW_VEC4UB] = "VEC4UB"
	cf.PropertyMap["VEC4UB"] = RW_VEC4UB

	cf.TypeMap[RW_VEC4S] = "VEC4S"
	cf.PropertyMap["VEC4S"] = RW_VEC4S

	cf.TypeMap[RW_VEC4US] = "VEC4US"
	cf.PropertyMap["VEC4US"] = RW_VEC4US

	cf.TypeMap[RW_VEC4I] = "VEC4I"
	cf.PropertyMap["VEC4I"] = RW_VEC4I

	cf.TypeMap[RW_VEC4UI] = "VEC4UI"
	cf.PropertyMap["VEC4UI"] = RW_VEC4UI

	cf.TypeMap[RW_LIST] = "LIST"
	cf.PropertyMap["LIST"] = RW_LIST

	cf.TypeMap[RW_VECTOR] = "VECTOR"
	cf.PropertyMap["VECTOR"] = RW_VECTOR

	cf.TypeMap[RW_MAP] = "MAP"
	cf.PropertyMap["MAP"] = RW_MAP
	return cf
}

func (cl *ClassInterface) AreTypesCompatible(lhs SerType, rhs SerType) bool {
	if lhs == rhs {
		return true
	}

	if lhs == RW_MATRIX {
		lhs = RW_MATRIXD
	}
	if rhs == RW_MATRIX {
		rhs = RW_MATRIXD
	}

	if lhs == RW_GLENUM {
		lhs = RW_UINT
	}
	if rhs == RW_GLENUM {
		rhs = RW_UINT
	}

	if lhs == RW_ENUM {
		lhs = RW_INT
	}
	if rhs == RW_ENUM {
		rhs = RW_INT
	}

	if lhs == RW_IMAGE {
		lhs = RW_OBJECT
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
	return RW_UNDEFINED
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

func (cl *ClassInterface) CreateObject(compoundClassName string) *model.Object {
	ow := GetObjectWrapperManager().FindWrap(compoundClassName)
	if ow != nil {
		return ow.CreateInstance()
	}
	return nil
}
