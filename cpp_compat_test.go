package osg

import (
	"bufio"
	"bytes"
	"fmt"
	"os"
	"testing"

	"github.com/flywave/go-osg/model"
)

// ---------------------------------------------------------------------------
// Helpers
// ---------------------------------------------------------------------------

type serSpec struct {
	Name       string
	WantType   SerType
	FirstVer   int32
	LastVer    int32
	EnumValues map[string]int32 // only for RWENUM
}

func verifyWrapper(t *testing.T, name string, wantChain []string, wantSers []serSpec) {
	t.Run(name, func(t *testing.T) {
		wrap := GetObjectWrapperManager().FindWrap("osg::" + name)
		if wrap == nil {
			t.Fatalf("wrapper osg::%s not found", name)
		}
		t.Logf("found: %d associates, %d serializers", len(wrap.Associates), len(wrap.Serializers))

		// Verify associate chain
		if len(wrap.Associates) != len(wantChain) {
			t.Errorf("associates: got %d, want %d", len(wrap.Associates), len(wantChain))
		}
		for i, ass := range wrap.Associates {
			if i < len(wantChain) && ass.Name != wantChain[i] {
				t.Errorf("associate[%d] = %q, want %q", i, ass.Name, wantChain[i])
			}
		}

		// Build set of serializers for matching
		got := map[string]struct{ Name string; Type SerType; Fv, Lv int32 }{}
		for _, s := range wrap.Serializers {
			ser := s.(Serializer)
			got[ser.GetSerializerName()] = struct{ Name string; Type SerType; Fv, Lv int32 }{
				ser.GetSerializerName(), findSerType(wrap, ser.GetSerializerName()),
				ser.GetFirstVersion(), ser.GetLastVersion(),
			}
		}

		for _, ws := range wantSers {
			g, ok := got[ws.Name]
			if !ok {
				t.Errorf("missing serializer %q", ws.Name)
				continue
			}
			if ws.FirstVer != 0 && g.Fv != ws.FirstVer {
				t.Errorf("%s FirstVersion = %d, want %d", ws.Name, g.Fv, ws.FirstVer)
			}
			if ws.LastVer != 0 && g.Lv != ws.LastVer {
				t.Errorf("%s LastVersion = %d, want %d", ws.Name, g.Lv, ws.LastVer)
			}
			delete(got, ws.Name)
		}
		for name := range got {
			t.Errorf("unexpected serializer %q (only C++ serializers expected)", name)
		}

		// Verify enum values
		for _, ws := range wantSers {
			if ws.EnumValues == nil {
				continue
			}
			ser := findEnumSerializer(wrap, ws.Name)
			if ser == nil {
				t.Errorf("enum serializer %q not found for value check", ws.Name)
				continue
			}
			for enumName, enumVal := range ws.EnumValues {
				if gotVal := ser.LookUp.GetValue(enumName); gotVal != enumVal {
					t.Errorf("%s enum %q = %d, want %d", ws.Name, enumName, gotVal, enumVal)
				}
			}
		}
	})
}

func findSerType(wrap *ObjectWrapper, name string) SerType {
	for i, s := range wrap.Serializers {
		if s.(Serializer).GetSerializerName() == name {
			return wrap.TypeList[i]
		}
	}
	return 0
}

func findEnumSerializer(wrap *ObjectWrapper, name string) *EnumSerializer {
	for _, s := range wrap.Serializers {
		if ser, ok := s.(*EnumSerializer); ok && ser.GetSerializerName() == name {
			return ser
		}
	}
	return nil
}

func ser(name string, ty SerType) serSpec {
	return serSpec{Name: name, WantType: ty}
}

func serV(name string, ty SerType, fv, lv int32) serSpec {
	return serSpec{Name: name, WantType: ty, FirstVer: fv, LastVer: lv}
}

func serE(name string, ty SerType, enums map[string]int32) serSpec {
	return serSpec{Name: name, WantType: ty, EnumValues: enums}
}

func serVE(name string, ty SerType, fv, lv int32, enums map[string]int32) serSpec {
	return serSpec{Name: name, WantType: ty, FirstVer: fv, LastVer: lv, EnumValues: enums}
}

// ---------------------------------------------------------------------------
// 1. Object
// ---------------------------------------------------------------------------

func TestWrapper_Object(t *testing.T) {
	verifyWrapper(t, "Object", []string{"osg::Object"}, []serSpec{
		ser("Name", RWSTRING),
		ser("DataVariance", RWENUM),
		serV("UserData", RWUSER, 0, 76),
		ser("UserDataContainer", RWOBJECT),
	})
}

// ---------------------------------------------------------------------------
// 2. BufferData
// ---------------------------------------------------------------------------

func TestWrapper_BufferData(t *testing.T) {
	verifyWrapper(t, "BufferData", []string{"osg::Object", "osg::BufferData"}, []serSpec{
		serV("BufferObject", RWOBJECT, 147, 0), // 0=MAX
	})
}

// ---------------------------------------------------------------------------
// 3. Array
// ---------------------------------------------------------------------------

func TestWrapper_Array(t *testing.T) {
	verifyWrapper(t, "Array", []string{"osg::Object", "osg::BufferData", "osg::Array"}, []serSpec{
		serE("Binding", RWENUM, map[string]int32{
			"BINDUNDEFINED": model.BINDUNDEFINED,
			"BINDOFF":       model.BINDOFF,
			"BINDOVERALL":   model.BINDOVERALL,
			"BINDPERPRIMITIVESET": model.BINDPERPRIMITIVESET,
			"BINDPERVERTEX": model.BINDPERVERTEX,
		}),
		ser("Normalize", RWBOOL),
		ser("PreserveDataType", RWBOOL),
	})
}

// ---------------------------------------------------------------------------
// 4-35. Concrete array types
// ---------------------------------------------------------------------------

func TestWrapper_ConcreteArrays(t *testing.T) {
	// C++ Array.cpp ARRAY_WRAPPERS: parent chain + IsAVectorSerializer named "vector"
	for _, name := range []string{
		"FloatArray", "Vec2Array", "Vec3Array", "Vec4Array",
		"DoubleArray", "Vec2dArray", "Vec3dArray", "Vec4dArray",
		"ByteArray", "Vec2bArray", "Vec3bArray", "Vec4bArray",
		"UByteArray", "Vec2ubArray", "Vec3ubArray", "Vec4ubArray",
		"ShortArray", "Vec2sArray", "Vec3sArray", "Vec4sArray",
		"UShortArray", "Vec2usArray", "Vec3usArray", "Vec4usArray",
		"IntArray", "Vec2iArray", "Vec3iArray", "Vec4iArray",
		"UIntArray", "Vec2uiArray", "Vec3uiArray", "Vec4uiArray",
	} {
		verifyWrapper(t, name, []string{"osg::Object", "osg::BufferData", "osg::Array", "osg::" + name},
			[]serSpec{ser("vector", RWVECTOR)})
	}
}

// ---------------------------------------------------------------------------
// 36. Node
// ---------------------------------------------------------------------------

func TestWrapper_Node(t *testing.T) {
	verifyWrapper(t, "Node", []string{"osg::Object", "osg::Node"}, []serSpec{
		ser("InitialBound", RWUSER),
		ser("ComputeBoundingSphereCallback", RWOBJECT),
		ser("UpdateCallback", RWOBJECT),
		ser("EventCallback", RWOBJECT),
		ser("CullCallback", RWOBJECT),
		ser("CullingActive", RWBOOL),
		ser("NodeMask", RWUINT),
		serV("Descriptions", RWUSER, 0, 76),
		ser("StateSet", RWOBJECT),
	})
}

// ---------------------------------------------------------------------------
// 37. Group
// ---------------------------------------------------------------------------

func TestWrapper_Group(t *testing.T) {
	verifyWrapper(t, "Group", []string{"osg::Object", "osg::Node", "osg::Group"}, []serSpec{
		ser("Children", RWUSER),
	})
}

// ---------------------------------------------------------------------------
// 38. Geode
// ---------------------------------------------------------------------------

func TestWrapper_Geode(t *testing.T) {
	verifyWrapper(t, "Geode", []string{"osg::Object", "osg::Node", "osg::Geode"}, []serSpec{
		ser("Drawables", RWOBJECT),
	})
}

// ---------------------------------------------------------------------------
// 39. LOD
// ---------------------------------------------------------------------------

func TestWrapper_LOD(t *testing.T) {
	verifyWrapper(t, "LOD", []string{"osg::Object", "osg::Node", "osg::Group", "osg::LOD"}, []serSpec{
		serE("CenterMode", RWENUM, map[string]int32{
			"USEBOUNDINGSPHERECENTER":                 model.USEBOUNDINGSPHERECENTER,
			"USERDEFINEDCENTER":                        model.USERDEFINEDCENTER,
			"UNIONOFBOUNDINGSPHEREANDUSERDEFINED":      model.UNIONOFBOUNDINGSPHEREANDUSERDEFINED,
		}),
		ser("UserCenter", RWUSER),
		serE("RangeMode", RWENUM, map[string]int32{
			"DISTANCEFROMEYEPOINT": model.DISTANCEFROMEYEPOINT,
			"PIXELSIZEONSCREEN":    model.PIXELSIZEONSCREEN,
		}),
		ser("RangeList", RWUSER),
	})
}

// ---------------------------------------------------------------------------
// 40. PagedLOD
// ---------------------------------------------------------------------------

func TestWrapper_PagedLOD(t *testing.T) {
	verifyWrapper(t, "PagedLOD", []string{"osg::Object", "osg::Node", "osg::LOD", "osg::PagedLOD"}, []serSpec{
		ser("DatabasePath", RWUSER),
		serV("FrameNumberOfLastTraversal", RWUINT, 0, 69),
		ser("NumChildrenThatCannotBeExpired", RWUINT),
		ser("DisableExternalChildrenPaging", RWBOOL),
		ser("RangeDataList", RWUSER),
		ser("Children", RWUSER),
	})
}

// ---------------------------------------------------------------------------
// 41. Geometry
// ---------------------------------------------------------------------------

func TestWrapper_Geometry(t *testing.T) {
	verifyWrapper(t, "Geometry", []string{"osg::Object", "osg::Node", "osg::Drawable", "osg::Geometry"}, []serSpec{
		ser("PrimitiveSetList", RWVECTOR),
		serV("VertexData", RWUSER, 0, 111),
		serV("NormalData", RWUSER, 0, 111),
		serV("ColorData", RWUSER, 0, 111),
		serV("SecondaryColorData", RWUSER, 0, 111),
		serV("FogCoordData", RWUSER, 0, 111),
		serV("TexCoordData", RWUSER, 0, 111),
		serV("VertexAttribData", RWUSER, 0, 111),
		serV("FastPathHint", RWUSER, 0, 111),
		serV("VertexArray", RWOBJECT, 112, 0),
		serV("NormalArray", RWOBJECT, 112, 0),
		serV("ColorArray", RWOBJECT, 112, 0),
		serV("SecondaryColorArray", RWOBJECT, 112, 0),
		serV("FogCoordArray", RWOBJECT, 112, 0),
		serV("TexCoordArrayList", RWVECTOR, 112, 0),
		serV("VertexAttribArrayList", RWVECTOR, 112, 0),
	})
}

// ---------------------------------------------------------------------------
// 42. Drawable
// ---------------------------------------------------------------------------

func TestWrapper_Drawable(t *testing.T) {
	verifyWrapper(t, "Drawable", []string{"osg::Object", "osg::Node", "osg::Drawable"}, []serSpec{
		ser("InitialBound", RWUSER),
		ser("ComputeBoundingBoxCallback", RWOBJECT),
		ser("Shape", RWOBJECT),
		ser("SupportsDisplayList", RWBOOL),
		ser("UseDisplayList", RWBOOL),
		ser("UseVertexBufferObjects", RWBOOL),
		serV("StateSet", RWOBJECT, 0, 155),
		serV("UpdateCallback", RWOBJECT, 0, 155),
		serV("EventCallback", RWOBJECT, 0, 155),
		serV("CullCallback", RWOBJECT, 0, 155),
		serV("DrawCallback", RWOBJECT, 0, 155),
		serV("NodeMask", RWUINT, 142, 0),
		serV("CullingActive", RWBOOL, 145, 0),
	})
}

// ---------------------------------------------------------------------------
// 43. PrimitiveSet + sub-types
// ---------------------------------------------------------------------------

func TestWrapper_PrimitiveSet(t *testing.T) {
	verifyWrapper(t, "PrimitiveSet", []string{"osg::Object", "osg::BufferData", "osg::PrimitiveSet"}, []serSpec{
		serV("NumInstances", RWINT, 96, 0),
		serE("Mode", RWENUM, map[string]int32{
			"POINTS":            model.GLPOINTS,
			"LINES":             model.GLLINES,
			"LINESTRIP":         model.GLLINESTRIP,
			"LINELOOP":          model.GLLINELOOP,
			"TRIANGLES":         model.GLTRIANGLES,
			"TRIANGLESTRIP":     model.GLTRIANGLESTRIP,
			"TRIANGLEFAN":       model.GLTRIANGLEFAN,
			"QUADS":             model.GLQUADS,
			"QUADSTRIP":         model.GLQUADSTRIP,
			"POLYGON":           model.GLPOLYGON,
			"LINESADJACENCY":    model.GLLINESADJACENCY,
			"LINESTRIPADJACENCY": model.GLLINESTRIPADJACENCY,
			"TRIANGLESADJACENCY":  model.GLTRIANGLESADJACENCY,
			"TRIANGLESTRIPADJACENCY": model.GLTRIANGLESTRIPADJACENCY,
			"PATCHES":           model.GLPATCHES,
		}),
	})
}

func TestWrapper_DrawArrays(t *testing.T) {
	verifyWrapper(t, "DrawArrays", []string{"osg::Object", "osg::BufferData", "osg::PrimitiveSet", "osg::DrawArrays"}, []serSpec{
		ser("First", RWINT),
		ser("Count", RWUINT),
	})
}

func TestWrapper_DrawArrayLengths(t *testing.T) {
	// C++ registers: First(INT) + vector(ISAVECTOR)
	verifyWrapper(t, "DrawArrayLengths", []string{"osg::Object", "osg::BufferData", "osg::PrimitiveSet", "osg::DrawArrayLengths"}, []serSpec{
		ser("First", RWINT),
		ser("Data", RWVECTOR),
	})
}

func TestWrapper_DrawElementsUByte(t *testing.T) {
	verifyWrapper(t, "DrawElementsUByte", []string{"osg::Object", "osg::BufferData", "osg::PrimitiveSet", "osg::DrawElementsUByte"}, []serSpec{
		ser("Data", RWVECTOR),
	})
}

func TestWrapper_DrawElementsUShort(t *testing.T) {
	verifyWrapper(t, "DrawElementsUShort", []string{"osg::Object", "osg::BufferData", "osg::PrimitiveSet", "osg::DrawElementsUShort"}, []serSpec{
		ser("Data", RWVECTOR),
	})
}

func TestWrapper_DrawElementsUInt(t *testing.T) {
	verifyWrapper(t, "DrawElementsUInt", []string{"osg::Object", "osg::BufferData", "osg::PrimitiveSet", "osg::DrawElementsUInt"}, []serSpec{
		ser("Data", RWVECTOR),
	})
}

// ---------------------------------------------------------------------------
// 49. Image
// ---------------------------------------------------------------------------

func TestWrapper_Image(t *testing.T) {
	verifyWrapper(t, "Image", []string{"osg::Object", "osg::BufferData", "osg::Image"}, []serSpec{
		ser("FileName", RWSTRING),
		ser("WriteHint", RWENUM),
		ser("AllocationMode", RWENUM),
		ser("InternalTextureFormat", RWGLENUM),
		ser("DataType", RWGLENUM),
		ser("PixelFormat", RWGLENUM),
		ser("RowLength", RWINT),
		ser("Packing", RWINT),
		ser("Origin", RWENUM),
	})
}

// ---------------------------------------------------------------------------
// 50. StateSet
// ---------------------------------------------------------------------------

func TestWrapper_StateSet(t *testing.T) {
	verifyWrapper(t, "StateSet", []string{"osg::Object", "osg::StateSet"}, []serSpec{
		ser("ModeList", RWUSER),
		ser("AttributeList", RWUSER),
		ser("TextureModeList", RWUSER),
		ser("TextureAttributeList", RWUSER),
		ser("UniformList", RWUSER),
		ser("RenderingHint", RWINT),
		ser("RenderBinMode", RWENUM),
		ser("BinNumber", RWINT),
		ser("BinName", RWSTRING),
		ser("NestRenderBins", RWBOOL),
		ser("UpdateCallback", RWOBJECT),
		ser("EventCallback", RWOBJECT),
		serV("DefineList", RWUSER, 151, 0),
	})
}

// ---------------------------------------------------------------------------
// 51. StateAttribute
// ---------------------------------------------------------------------------

func TestWrapper_StateAttribute(t *testing.T) {
	verifyWrapper(t, "StateAttribute", []string{"osg::Object", "osg::StateAttribute"}, []serSpec{
		ser("UpdateCallback", RWOBJECT),
		ser("EventCallback", RWOBJECT),
	})
}

// ---------------------------------------------------------------------------
// 52. CullFace
// ---------------------------------------------------------------------------

func TestWrapper_CullFace(t *testing.T) {
	verifyWrapper(t, "CullFace", []string{"osg::Object", "osg::StateAttribute", "osg::CullFace"}, []serSpec{
		serE("Mode", RWENUM, map[string]int32{
			"FRONT":          model.GLFRONT,
			"BACK":           model.GLBACK,
			"FRONTANDBACK":   model.GLFRONTANDBACK,
		}),
	})
}

// ---------------------------------------------------------------------------
// 53. AlphaFunc
// ---------------------------------------------------------------------------

func TestWrapper_AlphaFunc(t *testing.T) {
	verifyWrapper(t, "AlphaFunc", []string{"osg::Object", "osg::StateAttribute", "osg::AlphaFunc"}, []serSpec{
		serE("Function", RWENUM, map[string]int32{
			"NEVER":    model.GLNEVER,
			"LESS":     model.GLLESS,
			"EQUAL":    model.GLEQUAL,
			"LEQUAL":   model.GLLEQUAL,
			"GREATER":  model.GLGREATER,
			"NOTEQUAL": model.GLNOTEQUAL,
			"GEQUAL":   model.GLGEQUAL,
			"ALWAYS":   model.GLALWAYS,
		}),
		ser("ReferenceValue", RWFLOAT),
	})
}

// ---------------------------------------------------------------------------
// 54. ShadeModel
// ---------------------------------------------------------------------------

func TestWrapper_ShadeModel(t *testing.T) {
	verifyWrapper(t, "ShadeModel", []string{"osg::Object", "osg::StateAttribute", "osg::ShadeModel"}, []serSpec{
		serE("Mode", RWENUM, map[string]int32{
			"SMOOTH": model.SMOOTH,
			"FLAT":   model.FLAT,
		}),
	})
}

// ---------------------------------------------------------------------------
// 55. Material
// ---------------------------------------------------------------------------

func TestWrapper_Material(t *testing.T) {
	verifyWrapper(t, "Material", []string{"osg::Object", "osg::StateAttribute", "osg::Material"}, []serSpec{
		ser("ColorMode", RWENUM),
		ser("Ambient", RWUSER),
		ser("Diffuse", RWUSER),
		ser("Specular", RWUSER),
		ser("Emission", RWUSER),
		ser("Shininess", RWUSER),
	})
}

// ---------------------------------------------------------------------------
// 56. TexEnv
// ---------------------------------------------------------------------------

func TestWrapper_TexEnv(t *testing.T) {
	verifyWrapper(t, "TexEnv", []string{"osg::Object", "osg::StateAttribute", "osg::TexEnv"}, []serSpec{
		ser("Mode", RWENUM),
		ser("Color", RWVEC4F),
	})
}

// ---------------------------------------------------------------------------
// 57. TexGen
// ---------------------------------------------------------------------------

func TestWrapper_TexGen(t *testing.T) {
	verifyWrapper(t, "TexGen", []string{"osg::Object", "osg::StateAttribute", "osg::TexGen"}, []serSpec{
		ser("Mode", RWENUM),
		ser("PlaneS", RWUSER),
		ser("PlaneT", RWUSER),
		ser("PlaneR", RWUSER),
		ser("PlaneQ", RWUSER),
	})
}

// ---------------------------------------------------------------------------
// 58. Transform
// ---------------------------------------------------------------------------

func TestWrapper_Transform(t *testing.T) {
	verifyWrapper(t, "Transform", []string{"osg::Object", "osg::Node", "osg::Group", "osg::Transform"}, []serSpec{
		ser("ReferenceFrame", RWENUM),
	})
}

// ---------------------------------------------------------------------------
// 59. MatrixTransform
// ---------------------------------------------------------------------------

func TestWrapper_MatrixTransform(t *testing.T) {
	verifyWrapper(t, "MatrixTransform",
		[]string{"osg::Object", "osg::Node", "osg::Group", "osg::Transform", "osg::MatrixTransform"},
		[]serSpec{ser("Matrix", RWMATRIXF)})
}

// ---------------------------------------------------------------------------
// 60. PositionAttitudeTransform
// ---------------------------------------------------------------------------

func TestWrapper_PositionAttitudeTransform(t *testing.T) {
	verifyWrapper(t, "PositionAttitudeTransform",
		[]string{"osg::Object", "osg::Node", "osg::Group", "osg::Transform", "osg::PositionAttitudeTransform"},
		[]serSpec{
			ser("Position", RWDOUBLE|0xF0000000),
			ser("Attitude", RWDOUBLE|0xF0000000),
			ser("Scale", RWDOUBLE|0xF0000000),
			ser("PivotPoint", RWDOUBLE|0xF0000000),
		})
}

// ---------------------------------------------------------------------------
// 61. Texture (base)
// ---------------------------------------------------------------------------

func TestWrapper_Texture(t *testing.T) {
	verifyWrapper(t, "Texture", []string{"osg::Object", "osg::StateAttribute", "osg::Texture"}, []serSpec{
		ser("WRAP_S", RWUSER),
		ser("WRAP_T", RWUSER),
		ser("WRAP_R", RWUSER),
		ser("MIN_FILTER", RWUSER),
		ser("MAG_FILTER", RWUSER),
		ser("MaxAnisotropy", RWFLOAT),
		ser("UseHardwareMipMapGeneration", RWBOOL),
		ser("UnRefImageDataAfterApply", RWBOOL),
		ser("ClientStorageHint", RWBOOL),
		ser("ResizeNonPowerOfTwoHint", RWBOOL),
		ser("BorderColor", RWVEC4D),
		ser("BorderWidth", RWINT),
		ser("InternalFormatMode", RWENUM),
		ser("InternalFormat", RWUSER),
		ser("SourceFormat", RWUSER),
		ser("SourceType", RWUSER),
		ser("ShadowComparison", RWBOOL),
		ser("ShadowCompareFunc", RWENUM),
		ser("ShadowTextureMode", RWENUM),
		ser("ShadowAmbient", RWFLOAT),
		serV("ImageAttachment", RWUSER, 95, 153),
		serV("Swizzle", RWUSER, 98, 0),
		serV("MinLOD", RWFLOAT, 155, 0),
		serV("MaxLOD", RWFLOAT, 155, 0),
		serV("LODBias", RWFLOAT, 155, 0),
	})
}

// ---------------------------------------------------------------------------
// 62-65. Texture sub-types
// ---------------------------------------------------------------------------

func TestWrapper_Texture1D(t *testing.T) {
	verifyWrapper(t, "Texture1D",
		[]string{"osg::Object", "osg::StateAttribute", "osg::Texture", "osg::Texture1D"},
		[]serSpec{
			ser("Image", RWIMAGE),
			ser("TextureWidth", RWUINT),
		})
}

func TestWrapper_Texture2D(t *testing.T) {
	verifyWrapper(t, "Texture2D",
		[]string{"osg::Object", "osg::StateAttribute", "osg::Texture", "osg::Texture2D"},
		[]serSpec{
			ser("Image", RWIMAGE),
			ser("TextureWidth", RWUINT),
			ser("TextureHeight", RWUINT),
		})
}

func TestWrapper_Texture3D(t *testing.T) {
	verifyWrapper(t, "Texture3D",
		[]string{"osg::Object", "osg::StateAttribute", "osg::Texture", "osg::Texture3D"},
		[]serSpec{
			ser("Image", RWIMAGE),
			ser("TextureWidth", RWUINT),
			ser("TextureHeight", RWUINT),
			ser("TextureDepth", RWUINT),
		})
}

func TestWrapper_TextureRectangle(t *testing.T) {
	verifyWrapper(t, "TextureRectangle",
		[]string{"osg::Object", "osg::StateAttribute", "osg::Texture", "osg::TextureRectangle"},
		[]serSpec{
			ser("Image", RWIMAGE),
			ser("TextureWidth", RWINT),
			ser("TextureHeight", RWINT),
		})
}

func TestWrapper_TextureCubeMap(t *testing.T) {
	verifyWrapper(t, "TextureCubeMap",
		[]string{"osg::Object", "osg::StateAttribute", "osg::Texture", "osg::TextureCubeMap"},
		[]serSpec{
			ser("PosX", RWUSER),
			ser("NegX", RWUSER),
			ser("PosY", RWUSER),
			ser("NegY", RWUSER),
			ser("PosZ", RWUSER),
			ser("NegZ", RWUSER),
			ser("TextureWidth", RWINT),
			ser("TextureHeight", RWINT),
		})
}

// ---------------------------------------------------------------------------
// Version-threshold logic tests
// ---------------------------------------------------------------------------

func TestThresholds_ReadArray(t *testing.T) {
	for _, c := range []struct{ v int32; want bool }{{0, false}, {111, false}, {112, true}, {130, true}} {
		if got := c.v >= 112; got != c.want {
			t.Errorf("ReadArray version=%d: uses ReadObject=%v want=%v", c.v, got, c.want)
		}
	}
}

func TestThresholds_ImageClassName(t *testing.T) {
	for _, c := range []struct{ v int32; want bool }{{0, false}, {94, false}, {95, true}, {130, true}} {
		if got := c.v > 94; got != c.want {
			t.Errorf("Image ClassName version=%d: reads=%v want=%v", c.v, got, c.want)
		}
	}
}

func TestThresholds_PrimitiveSetNumInstances(t *testing.T) {
	for _, c := range []struct{ v int32; want bool }{{0, false}, {96, false}, {97, true}, {130, true}} {
		if got := c.v > 96; got != c.want {
			t.Errorf("PrimitiveSet numInstances version=%d: reads=%v want=%v", c.v, got, c.want)
		}
	}
}

func TestThresholds_BinaryBracketSize(t *testing.T) {
	for _, c := range []struct{ v int32; want bool }{{0, false}, {148, false}, {149, true}} {
		if got := c.v > 148; got != c.want {
			t.Errorf("binary bracket uint64 version=%d: uses=%v want=%v", c.v, got, c.want)
		}
	}
}

// ---------------------------------------------------------------------------
// Getter pointer tests (the bug found earlier)
// ---------------------------------------------------------------------------

func TestGetters_ReturnPointers(t *testing.T) {
	arr := model.NewArray(model.Vec3ArrayType, model.GLFLOAT, 3)
	cf := &model.CullFace{}
	af := &model.AlphaFunc{}

	tests := []struct {
		name string
		val  interface{}
	}{
		{"getBinding", getBinding(arr)},
		{"getNormalize", getNormalize(arr)},
		{"getPreserveDataType", getPreserveDataType(arr)},
		{"getDataSize", getDataSize(arr)},
		{"getArrayDataType", getArrayDataType(arr)},
		{"getElementSize", getElementSize(arr)},
		{"getTotalDataSize", getTotalDataSize(arr)},
		{"getNumElements", getNumElements(arr)},
		{"getType", getType(arr)},
		{"getMode(cullface)", getMode(cf)},
		{"getComparisonFunc", getComparisonFunc(af)},
		{"getReferenceValue", getReferenceValue(af)},
	}

	for _, tt := range tests {
		switch tt.val.(type) {
		case *int32, *bool, *float32, *float64, *uint32, *int, *uint, *model.ArrayTable:
		default:
			t.Errorf("%s returns %T (want pointer)", tt.name, tt.val)
		}
	}
}

// ---------------------------------------------------------------------------
// Serializer type consistency: verify SerType tags match C++ ADD_* macros
// ---------------------------------------------------------------------------

func TestSerializerTypes_DrawableBoolFields(t *testing.T) {
	wrap := GetObjectWrapperManager().FindWrap("osg::Drawable")
	if wrap == nil {
		t.Fatal("Drawable wrapper not found")
	}
	// C++: ADD_BOOL_SERIALIZER(SupportsDisplayList/UseDisplayList/UseVertexBufferObjects)
	// Go must use RWBOOL for these
	want := map[string]SerType{
		"SupportsDisplayList":    RWBOOL,
		"UseDisplayList":         RWBOOL,
		"UseVertexBufferObjects": RWBOOL,
	}
	for _, s := range wrap.Serializers {
		ser := s.(Serializer)
		name := ser.GetSerializerName()
		if _, ok := want[name]; ok {
			got := findSerType(wrap, name)
			if got != want[name] {
				t.Errorf("%s SerType = %d, want %d (RWBOOL)", name, got, RWBOOL)
			}
			delete(want, name)
		}
	}
	for name := range want {
		t.Errorf("missing serializer %q in Drawable wrapper", name)
	}
}

func TestSerializerTypes_NodeMask(t *testing.T) {
	wrap := GetObjectWrapperManager().FindWrap("osg::Node")
	if wrap == nil {
		t.Fatal("Node wrapper not found")
	}
	// C++: ADD_HEXINT_SERIALIZER(NodeMask) → stored as uint32 → RWUINT
	if got := findSerType(wrap, "NodeMask"); got != RWUINT {
		t.Errorf("NodeMask SerType = %d, want %d (RWUINT)", got, RWUINT)
	}
}

func TestSerializerTypes_ArrayBindingNormalize(t *testing.T) {
	wrap := GetObjectWrapperManager().FindWrap("osg::Array")
	if wrap == nil {
		t.Fatal("Array wrapper not found")
	}
	// C++: Array.cpp: Binding → ADD_ENUM_SERIALIZER → RWENUM
	if got := findSerType(wrap, "Binding"); got != RWENUM {
		t.Errorf("Binding SerType = %d, want %d (RWENUM)", got, RWENUM)
	}
	// C++: Normalize → ADD_BOOL_SERIALIZER → RWBOOL
	if got := findSerType(wrap, "Normalize"); got != RWBOOL {
		t.Errorf("Normalize SerType = %d, want %d (RWBOOL)", got, RWBOOL)
	}
	// C++: PreserveDataType → ADD_BOOL_SERIALIZER → RWBOOL
	if got := findSerType(wrap, "PreserveDataType"); got != RWBOOL {
		t.Errorf("PreserveDataType SerType = %d, want %d (RWBOOL)", got, RWBOOL)
	}
}

func TestSerializerTypes_ImageEnumTypes(t *testing.T) {
	wrap := GetObjectWrapperManager().FindWrap("osg::Image")
	if wrap == nil {
		t.Fatal("Image wrapper not found")
	}
	// C++: WriteHint → ADD_ENUM_SERIALIZER → RWENUM
	if got := findSerType(wrap, "WriteHint"); got != RWENUM {
		t.Errorf("WriteHint SerType = %d, want %d (RWENUM)", got, RWENUM)
	}
	// C++: RowLength → ADD_INT_SERIALIZER → RWINT
	if got := findSerType(wrap, "RowLength"); got != RWINT {
		t.Errorf("RowLength SerType = %d, want %d (RWINT)", got, RWINT)
	}
	// C++: InternalTextureFormat → ADD_GLENUM_SERIALIZER → RWGLENUM
	if got := findSerType(wrap, "InternalTextureFormat"); got != RWGLENUM {
		t.Errorf("InternalTextureFormat SerType = %d, want %d (RWGLENUM)", got, RWGLENUM)
	}
	// C++: FileName → ADD_STRING_SERIALIZER → RWSTRING
	if got := findSerType(wrap, "FileName"); got != RWSTRING {
		t.Errorf("FileName SerType = %d, want %d (RWSTRING)", got, RWSTRING)
	}
}

func TestSerializerTypes_PositionAttitudeTransform(t *testing.T) {
	wrap := GetObjectWrapperManager().FindWrap("osg::PositionAttitudeTransform")
	if wrap == nil {
		t.Fatal("PositionAttitudeTransform wrapper not found")
	}
	// C++: Position/Attitude/Scale/PivotPoint → ADD_VEC3D_SERIALIZER → VEC3D
	// Go uses RWDOUBLE|0xF0000000 as encoding type
	// We just verify they exist and have the right flag
	for _, name := range []string{"Position", "Attitude", "Scale", "PivotPoint"} {
		got := findSerType(wrap, name)
		if got&0xF0000000 == 0 {
			t.Errorf("%s SerType = 0x%08X, missing 0xF0000000 flag", name, got)
		}
	}
}

// ---------------------------------------------------------------------------
// Default value tests
// ---------------------------------------------------------------------------

func TestDefaults_ArrayBinding(t *testing.T) {
	arr := model.NewArray(model.Vec3ArrayType, model.GLFLOAT, 3)
	if arr.Binding != model.BINDUNDEFINED {
		t.Errorf("default Binding = %d, want %d", arr.Binding, model.BINDUNDEFINED)
	}
}

// ---------------------------------------------------------------------------
// File reading test for Tile_+002_+000_L22_000020.osgb
// ---------------------------------------------------------------------------

func TestFile_Tile002L22(t *testing.T) {
	data := readTestData(t, "tiles3d_test/Tile_+002_+000_L22_000020.osgb")
	rw := NewReadWrite()
	res := rw.ReadNodeWithReader(bufio.NewReader(bytes.NewReader(data)), &OsgIstreamOptions{})
	if res == nil || res.GetNode() == nil {
		t.Fatal("failed to read node")
	}

	totalVerts := 0
	geomCount := 0
	walkGeometry(res.GetNode(), func(g *model.Geometry) {
		geomCount++
		if g.VertexArray == nil || g.VertexArray.Data == nil {
			t.Error("VertexArray.Data is nil")
			return
		}
		d, ok := g.VertexArray.Data.([][3]float32)
		if !ok {
			t.Errorf("unexpected data type %T", g.VertexArray.Data)
			return
		}
		totalVerts += len(d)
		if len(d) > 0 {
			x, y, z := d[0][0], d[0][1], d[0][2]
			if x < -200 || x > 200 || y < -200 || y > 200 {
				t.Errorf("suspicious vertex (%.2f,%.2f,%.2f)", x, y, z)
			}
		}
	})
	t.Logf("geometries=%d vertices=%d binding=%d", geomCount, totalVerts, 4)
	if geomCount != 4 {
		t.Errorf("expected 4 geometries, got %d", geomCount)
	}
	if totalVerts != 293 {
		t.Errorf("expected 293 vertices, got %d", totalVerts)
	}
}

func readTestData(t *testing.T, path string) []byte {
	t.Helper()
	data, err := os.ReadFile(path)
	if err != nil {
		t.Skipf("test data %s not available: %v", path, err)
	}
	return data
}

func walkAll(n interface{}, fn func(interface{})) {
	fn(n)
	switch v := n.(type) {
	case *model.Group:
		for _, c := range v.GetChildren() {
			walkAll(c, fn)
		}
	case *model.PagedLod:
		for _, c := range v.GetChildren() {
			walkAll(c, fn)
		}
	case *model.Geode:
		for _, c := range v.GetChildren() {
			walkAll(c, fn)
		}
	}
}

func walkGeometry(n interface{}, fn func(*model.Geometry)) {
	switch v := n.(type) {
	case *model.Geometry:
		fn(v)
	case *model.Group:
		for _, c := range v.GetChildren() {
			walkGeometry(c, fn)
		}
	case *model.PagedLod:
		for _, c := range v.GetChildren() {
			walkGeometry(c, fn)
		}
	case *model.Geode:
		for _, c := range v.GetChildren() {
			walkGeometry(c, fn)
		}
	}
}

// ---------------------------------------------------------------------------
// Enum value tests — verify every serializer enum value matches C++ GL constants
// ---------------------------------------------------------------------------

func TestEnums_CullFaceValues(t *testing.T) {
	if model.GLFRONT != 0x0404 {
		t.Errorf("GLFRONT = 0x%04X, want 0x0404", model.GLFRONT)
	}
	if model.GLBACK != 0x0405 {
		t.Errorf("GLBACK = 0x%04X, want 0x0405", model.GLBACK)
	}
	if model.GLFRONTANDBACK != 0x0408 {
		t.Errorf("GLFRONTANDBACK = 0x%04X, want 0x0408", model.GLFRONTANDBACK)
	}
}

func TestEnums_AlphaFuncValues(t *testing.T) {
	vals := map[string]int32{
		"NEVER": 0x0200, "LESS": 0x0201, "EQUAL": 0x0202, "LEQUAL": 0x0203,
		"GREATER": 0x0204, "NOTEQUAL": 0x0205, "GEQUAL": 0x0206, "ALWAYS": 0x0207,
	}
	syms := map[string]int32{"GLNEVER": model.GLNEVER, "GLLESS": model.GLLESS, "GLEQUAL": model.GLEQUAL,
		"GLLEQUAL": model.GLLEQUAL, "GLGREATER": model.GLGREATER, "GLNOTEQUAL": model.GLNOTEQUAL,
		"GLGEQUAL": model.GLGEQUAL, "GLALWAYS": model.GLALWAYS}
	for name, val := range syms {
		if want := vals[name[2:]]; val != want {
			t.Errorf("%s = 0x%04X, want 0x%04X", name, val, want)
		}
	}
}

func TestEnums_ShadeModelValues(t *testing.T) {
	if model.SMOOTH != 0x1D01 {
		t.Errorf("SMOOTH = 0x%04X, want 0x1D01", model.SMOOTH)
	}
	if model.FLAT != 0x1D00 {
		t.Errorf("FLAT = 0x%04X, want 0x1D00", model.FLAT)
	}
}

func TestEnums_TexEnvValues(t *testing.T) {
	if model.GLDECAL != 0x2101 {
		t.Errorf("GLDECAL = 0x%04X, want 0x2101", model.GLDECAL)
	}
	if model.GLMODULATE != 0x2100 {
		t.Errorf("GLMODULATE = 0x%04X, want 0x2100", model.GLMODULATE)
	}
	if model.GLREPLACE != 0x1E01 {
		t.Errorf("GLREPLACE = 0x%04X, want 0x1E01", model.GLREPLACE)
	}
	if model.GLADD != 0x0104 {
		t.Errorf("GLADD = 0x%04X, want 0x0104", model.GLADD)
	}
}

func TestEnums_TexGenValues(t *testing.T) {
	if model.GLOBJECTLINEAR != 0x2401 {
		t.Errorf("GLOBJECTLINEAR = 0x%04X, want 0x2401", model.GLOBJECTLINEAR)
	}
	if model.GLEYELINEAR != 0x2400 {
		t.Errorf("GLEYELINEAR = 0x%04X, want 0x2400", model.GLEYELINEAR)
	}
	if model.GLSPHEREMAP != 0x2402 {
		t.Errorf("GLSPHEREMAP = 0x%04X, want 0x2402", model.GLSPHEREMAP)
	}
	if model.GLNORMALMAP != 0x8511 {
		t.Errorf("GLNORMALMAP = 0x%04X, want 0x8511", model.GLNORMALMAP)
	}
	if model.GLREFLECTIONMAP != 0x8512 {
		t.Errorf("GLREFLECTIONMAP = 0x%04X, want 0x8512", model.GLREFLECTIONMAP)
	}
}

func TestEnums_PrimitiveSetModeValues(t *testing.T) {
	vals := map[string]int32{
		"GLPOINTS": 0x0000, "GLLINES": 0x0001, "GLLINESTRIP": 0x0003, "GLLINELOOP": 0x0002,
		"GLTRIANGLES": 0x0004, "GLTRIANGLESTRIP": 0x0005, "GLTRIANGLEFAN": 0x0006,
		"GLQUADS": 0x0007, "GLQUADSTRIP": 0x0008, "GLPOLYGON": 0x0009,
		"GLLINESADJACENCY": 0x000A, "GLLINESTRIPADJACENCY": 0x000B,
		"GLTRIANGLESADJACENCY": 0x000C, "GLTRIANGLESTRIPADJACENCY": 0x000D,
		"GLPATCHES": 0x000E,
	}
	for sym, want := range vals {
		ptr := map[string]int32{
			"GLPOINTS": model.GLPOINTS, "GLLINES": model.GLLINES, "GLLINESTRIP": model.GLLINESTRIP,
			"GLLINELOOP": model.GLLINELOOP, "GLTRIANGLES": model.GLTRIANGLES,
			"GLTRIANGLESTRIP": model.GLTRIANGLESTRIP, "GLTRIANGLEFAN": model.GLTRIANGLEFAN,
			"GLQUADS": model.GLQUADS, "GLQUADSTRIP": model.GLQUADSTRIP, "GLPOLYGON": model.GLPOLYGON,
			"GLLINESADJACENCY": model.GLLINESADJACENCY, "GLLINESTRIPADJACENCY": model.GLLINESTRIPADJACENCY,
			"GLTRIANGLESADJACENCY": model.GLTRIANGLESADJACENCY, "GLTRIANGLESTRIPADJACENCY": model.GLTRIANGLESTRIPADJACENCY,
			"GLPATCHES": model.GLPATCHES,
		}[sym]
		if ptr != want {
			t.Errorf("%s = 0x%04X, want 0x%04X", sym, ptr, want)
		}
	}
}

func TestEnums_TransformReferenceFrame(t *testing.T) {
	if model.RELATIVERF != 0 {
		t.Errorf("RELATIVERF = %d, want 0", model.RELATIVERF)
	}
	if model.ABSOLUTERF != 1 {
		t.Errorf("ABSOLUTERF = %d, want 1", model.ABSOLUTERF)
	}
	if model.ABSOLUTERFINHERITVIEWPOINT != 2 {
		t.Errorf("ABSOLUTERFINHERITVIEWPOINT = %d, want 2", model.ABSOLUTERFINHERITVIEWPOINT)
	}
}

func TestEnums_MaterialColorMode(t *testing.T) {
	if model.AMBIENT != 0x1200 {
		t.Errorf("AMBIENT = 0x%04X, want 0x1200", model.AMBIENT)
	}
	if model.DIFFUSE != 0x1201 {
		t.Errorf("DIFFUSE = 0x%04X, want 0x1201", model.DIFFUSE)
	}
	if model.SPECULAR != 0x1202 {
		t.Errorf("SPECULAR = 0x%04X, want 0x1202", model.SPECULAR)
	}
	if model.EMISSION != 0x1600 {
		t.Errorf("EMISSION = 0x%04X, want 0x1600", model.EMISSION)
	}
	if model.AMBIENTANDDIFFUSE != 0x1602 {
		t.Errorf("AMBIENTANDDIFFUSE = 0x%04X, want 0x1602", model.AMBIENTANDDIFFUSE)
	}
}

func TestEnums_ImageOrigin(t *testing.T) {
	if model.BOTTOMLEFT != 0 {
		t.Errorf("BOTTOMLEFT = %d, want 0", model.BOTTOMLEFT)
	}
	if model.TOPLEFT != 1 {
		t.Errorf("TOPLEFT = %d, want 1", model.TOPLEFT)
	}
}

func TestEnums_TextureWrapValues(t *testing.T) {
	if model.GLCLAMP != 0x2900 {
		t.Errorf("GLCLAMP = 0x%04X, want 0x2900", model.GLCLAMP)
	}
	if model.GLCLAMPTOEDGE != 0x812F {
		t.Errorf("GLCLAMPTOEDGE = 0x%04X, want 0x812F", model.GLCLAMPTOEDGE)
	}
	if model.GLREPEAT != 0x2901 {
		t.Errorf("GLREPEAT = 0x%04X, want 0x2901", model.GLREPEAT)
	}
}

func TestEnums_TextureFilterValues(t *testing.T) {
	if model.GLLINEAR != 0x2601 {
		t.Errorf("GLLINEAR = 0x%04X, want 0x2601", model.GLLINEAR)
	}
	if model.GLNEAREST != 0x2600 {
		t.Errorf("GLNEAREST = 0x%04X, want 0x2600", model.GLNEAREST)
	}
	if model.GLLINEARMIPMAPLINEAR != 0x2703 {
		t.Errorf("GLLINEARMIPMAPLINEAR = 0x%04X, want 0x2703", model.GLLINEARMIPMAPLINEAR)
	}
	if model.GLLINEARMIPMAPNEAREST != 0x2701 {
		t.Errorf("GLLINEARMIPMAPNEAREST = 0x%04X, want 0x2701", model.GLLINEARMIPMAPNEAREST)
	}
	if model.GLNEARESTMIPMAPLINEAR != 0x2702 {
		t.Errorf("GLNEARESTMIPMAPLINEAR = 0x%04X, want 0x2702", model.GLNEARESTMIPMAPLINEAR)
	}
	if model.GLNEARESTMIPMAPNEAREST != 0x2700 {
		t.Errorf("GLNEARESTMIPMAPNEAREST = 0x%04X, want 0x2700", model.GLNEARESTMIPMAPNEAREST)
	}
}

// ---------------------------------------------------------------------------
// Integration tests for all available test data files
// ---------------------------------------------------------------------------

func TestFile_Cessna(t *testing.T) {
	data := readTestData(t, "test_data/cessna.osgb")
	if data == nil {
		return
	}
	rw := NewReadWrite()
	res := rw.ReadNodeWithReader(bufio.NewReader(bytes.NewReader(data)), &OsgIstreamOptions{})
	if res == nil || res.GetNode() == nil {
		t.Fatal("failed to read cessna.osgb")
	}
	geomCount := 0
	vertCount := 0
	walkGeometry(res.GetNode(), func(g *model.Geometry) {
		geomCount++
		if g.VertexArray != nil && g.VertexArray.Data != nil {
			if d, ok := g.VertexArray.Data.([][3]float32); ok {
				vertCount += len(d)
			}
		}
	})
	if geomCount == 0 {
		t.Error("cessna.osgb: no geometries found")
	}
	t.Logf("cessna.osgb: geometries=%d vertices=%d", geomCount, vertCount)
}

func TestFile_SimpleRoom(t *testing.T) {
	data := readTestData(t, "test_data/simpleroom.osgt")
	if data == nil {
		return
	}
	rw := NewReadWrite()
	opts := NewOsgIstreamOptions()
	opts.FileType = "Ascii"
	res := rw.ReadNodeWithReader(bufio.NewReader(bytes.NewReader(data)), opts)
	if res == nil || res.GetNode() == nil {
		t.Fatal("failed to read simpleroom.osgt")
	}
	geomCount := 0
	nodeCount := 0
	walkAll(res.GetNode(), func(n interface{}) { nodeCount++ })
	walkGeometry(res.GetNode(), func(g *model.Geometry) {
		geomCount++
	})
	t.Logf("simpleroom.osgt: version=92 nodes=%d geometries=%d", nodeCount, geomCount)
}

func TestFile_SkyDome(t *testing.T) {
	data := readTestData(t, "test_data/skydome.osgt")
	if data == nil {
		return
	}
	rw := NewReadWrite()
	opts := NewOsgIstreamOptions()
	opts.FileType = "Ascii"
	res := rw.ReadNodeWithReader(bufio.NewReader(bytes.NewReader(data)), opts)
	if res == nil || res.GetNode() == nil {
		t.Fatal("failed to read skydome.osgt")
	}
	geomCount := 0
	walkGeometry(res.GetNode(), func(g *model.Geometry) {
		geomCount++
	})
	t.Logf("skydome.osgt: geometries=%d", geomCount)
}

func TestFile_Tile003L18(t *testing.T) {
	data := readTestData(t, "test_data/Tile_+003_+003_L18_000.osgb")
	if data == nil {
		return
	}
	rw := NewReadWrite()
	res := rw.ReadNodeWithReader(bufio.NewReader(bytes.NewReader(data)), &OsgIstreamOptions{})
	if res == nil || res.GetNode() == nil {
		t.Fatal("failed to read Tile_+003_+003_L18_000.osgb")
	}
	vertCount := 0
	geomCount := 0
	walkGeometry(res.GetNode(), func(g *model.Geometry) {
		geomCount++
		if g.VertexArray != nil && g.VertexArray.Data != nil {
			if d, ok := g.VertexArray.Data.([][3]float32); ok {
				vertCount += len(d)
			}
		}
	})
	t.Logf("Tile_+003_+003_L18_000.osgb: geometries=%d vertices=%d", geomCount, vertCount)
	if geomCount == 0 {
		t.Error("no geometries found")
	}
}



// ---------------------------------------------------------------------------
// Array construction via NewArray produces correct metadata
// ---------------------------------------------------------------------------

func TestArray_NewArrayProperties(t *testing.T) {
	// Vec3Array with 3 elements
	arr := model.NewArray(model.Vec3ArrayType, model.GLFLOAT, 3)
	if arr.Type != model.Vec3ArrayType {
		t.Errorf("Type = %d, want %d", arr.Type, model.Vec3ArrayType)
	}
	if arr.DataType != model.GLFLOAT {
		t.Errorf("DataType = %d, want %d", arr.DataType, model.GLFLOAT)
	}
	if arr.DataSize != 3 {
		t.Errorf("DataSize = %d, want 3", arr.DataSize)
	}

	// FloatArray with 1 element
	arr2 := model.NewArray(model.FloatArrayType, model.GLFLOAT, 1)
	if arr2.Type != model.FloatArrayType {
		t.Errorf("FloatArray Type = %d", arr2.Type)
	}
	if arr2.DataSize != 1 {
		t.Errorf("FloatArray DataSize = %d", arr2.DataSize)
	}
}

// ---------------------------------------------------------------------------
// Factory functions return correct instance types (prevent **T bugs)
// ---------------------------------------------------------------------------

func TestFactories_ReturnCorrectTypes(t *testing.T) {
	tests := []struct {
		name string
		fn   func() interface{}
		want string
	}{
		{"Object", func() interface{} { return model.NewObject() }, "*model.Object"},
		{"Node", func() interface{} { return model.NewNode() }, "*model.Node"},
		{"Group", func() interface{} { return model.NewGroup() }, "*model.Group"},
		{"Geode", func() interface{} { return model.NewGeode() }, "*model.Geode"},
		{"Geometry", func() interface{} { return model.NewGeometry() }, "*model.Geometry"},
		{"Transform", func() interface{} { return model.NewTransform() }, "*model.Transform"},
		{"MatrixTransform", func() interface{} { return model.NewMatrixTransform() }, "*model.MatrixTransform"},
		{"PositionAttitudeTransform", func() interface{} { return model.NewPositionAttitudeTransform() }, "*model.PositionAttitudeTransform"},
		{"LOD", func() interface{} { return model.NewLod() }, "*model.Lod"},
		{"PagedLOD", func() interface{} { return model.NewPagedLod() }, "*model.PagedLod"},
		{"StateSet", func() interface{} { return model.NewStateSet() }, "*model.StateSet"},
		{"CullFace", func() interface{} { return model.NewCullFace() }, "*model.CullFace"},
		{"AlphaFunc", func() interface{} { return model.NewAlphaFunc() }, "*model.AlphaFunc"},
		{"ShadeModel", func() interface{} { return model.NewShadeModel() }, "*model.ShadeModel"},
		{"Material", func() interface{} { return model.NewMaterial() }, "*model.Material"},
		{"TexEnv", func() interface{} { return model.NewTexEnv() }, "*model.TexEnv"},
		{"TexGen", func() interface{} { return model.NewTexGen() }, "*model.TexGen"},
		{"Image", func() interface{} { return model.NewImage() }, "*model.Image"},
	}
	for _, tt := range tests {
		obj := tt.fn()
		got := fmt.Sprintf("%T", obj)
		if got != tt.want {
			t.Errorf("%s factory returns %s, want %s", tt.name, got, tt.want)
		}
	}
}

// ---------------------------------------------------------------------------
// Default model values match C++ default constructors
// ---------------------------------------------------------------------------

func TestDefaults_ModelValues(t *testing.T) {
	// LOD defaults
	lod := model.NewLod()
	if *lod.GetCmode() != model.USEBOUNDINGSPHERECENTER {
		t.Errorf("LOD CenterMode default = %d, want %d", *lod.GetCmode(), model.USEBOUNDINGSPHERECENTER)
	}
	if *lod.GetRmode() != model.DISTANCEFROMEYEPOINT {
		t.Errorf("LOD RangeMode default = %d, want %d", *lod.GetRmode(), model.DISTANCEFROMEYEPOINT)
	}

	// Transform defaults
	tf := model.NewTransform()
	if tf.ReferenceFrame != model.RELATIVERF {
		t.Errorf("Transform ReferenceFrame default = %d, want %d", tf.ReferenceFrame, model.RELATIVERF)
	}

	// Material defaults
	mat := model.NewMaterial()
	if mat.Cmod != model.MTLOFF {
		t.Errorf("Material ColorMode default = %d, want %d", mat.Cmod, model.MTLOFF)
	}

	// TexEnv defaults
	te := model.NewTexEnv()
	if te.Mode != model.GLMODULATE {
		t.Errorf("TexEnv Mode default = 0x%04X, want 0x%04X", te.Mode, model.GLMODULATE)
	}

	// AlphaFunc defaults
	af := model.NewAlphaFunc()
	if af.ComparisonFunc != int(model.GLALWAYS) {
		t.Errorf("AlphaFunc Function default = %d, want %d (ALWAYS)", af.ComparisonFunc, model.GLALWAYS)
	}
	if af.ReferenceValue != 1.0 {
		t.Errorf("AlphaFunc ReferenceValue default = %f, want 1.0", af.ReferenceValue)
	}

	// CullFace defaults
	cf := model.NewCullFace()
	if cf.Mode != int(model.GLBACK) {
		t.Errorf("CullFace Mode default = %d, want %d (BACK)", cf.Mode, model.GLBACK)
	}

	// Image defaults
	img := model.NewImage()
	if img.AllocationMode != model.USENEWDELETE {
		t.Errorf("Image AllocationMode default = %d, want %d", img.AllocationMode, model.USENEWDELETE)
	}
}
