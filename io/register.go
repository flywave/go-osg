package io

import (
	"strings"

	"github.com/flywave/go-osg/model"
)

var manager *objectWrapperManager

func init() {
	manager = newObjectWrapperManager()
}

type ObjectWrapperAssociate struct {
	FirstVersion int
	LastVersion  int
	Name         string
}

type CreateInstanceFuncType func() interface{}

type ObjectWrapper struct {
	CreateInstanceFunc                   CreateInstanceFuncType
	Domain                               string
	Name                                 string
	Associates                           []*ObjectWrapperAssociate
	Serializers                          []interface{}
	BackupSerializers                    []interface{}
	TypeList                             []SerType
	Version                              int
	IsAssociatesRevisionsInheritanceDone bool
}

func NewObjectWrapper(name string, fn CreateInstanceFuncType, associates string) ObjectWrapper {
	ow := ObjectWrapper{Name: name, CreateInstanceFunc: fn, Version: 0, IsAssociatesRevisionsInheritanceDone: false}
	ow.SplitAssociates(associates, " ")
	return ow
}

func NewObjectWrapper2(name string, domain string, fn CreateInstanceFuncType, associates string) ObjectWrapper {
	ow := ObjectWrapper{Name: name, Domain: domain, CreateInstanceFunc: fn, Version: 0, IsAssociatesRevisionsInheritanceDone: false}
	ow.SplitAssociates(associates, " ")
	return ow
}

func (wp *ObjectWrapper) SplitAssociates(str string, separator string) {
	list := strings.Split(separator, str)
	if separator == "" {
		separator = " "
	}

	for _, l := range list {
		if l != separator {
			owa := ObjectWrapperAssociate{Name: l}
			wp.Associates = append(wp.Associates, &owa)
		}
	}
}

func (wp *ObjectWrapper) CreateInstance() interface{} {
	return wp.CreateInstanceFunc()
}

func (wp *ObjectWrapper) AddSerializer(s interface{}, t SerType) {
	s.(*BaseSerializer).FirstVersion = wp.Version
	wp.Serializers = append(wp.Serializers, s)
	wp.TypeList = append(wp.TypeList, t)
}

func (wp *ObjectWrapper) MarkSerializerAsAdded(name string) {
	for _, as := range wp.Associates {
		if as.Name == name {
			as.FirstVersion = wp.Version
			return
		}
	}
}

func (wp *ObjectWrapper) MarkSerializerAsRemoved(name string) {
	for _, s := range wp.Serializers {
		ser := s.(*BaseSerializer)
		if ser.GetSerializerName() == name {
			ser.LastVersion = wp.Version - 1
		}
	}
}

func (wp *ObjectWrapper) GetSerializer(name string) interface{} {
	for _, s := range wp.Serializers {
		ser := s.(*BaseSerializer)
		if ser.GetSerializerName() == name {
			return s
		}
	}

	for _, as := range wp.Associates {
		w := GetObjectWrapperManager().FindWrap(as.Name)
		if w == nil {
			continue
		}
		for _, s := range w.Serializers {
			ser := s.(*BaseSerializer)
			if ser.GetSerializerName() == name {
				return s
			}
		}
	}
	return nil
}

func (wp *ObjectWrapper) GetSerializerAndType(name string, ty *SerType) interface{} {
	for i, s := range wp.Serializers {
		ser := s.(*BaseSerializer)
		if ser.GetSerializerName() == name {
			*ty = wp.TypeList[i]
			return s
		}
	}

	for _, as := range wp.Associates {
		w := GetObjectWrapperManager().FindWrap(as.Name)
		if w == nil {
			continue
		}
		for _, s := range w.Serializers {
			ser := s.(*BaseSerializer)
			if ser.GetSerializerName() == name {
				*ty = w.TypeList[0]
				return s
			}
		}
	}
	*ty = RW_UNDEFINED
	return nil
}

func (wp *ObjectWrapper) Read(is *OsgIstream, obj *model.Object) {
	inputVersion := is.GetFileVersion(wp.Domain)
	for _, s := range wp.Serializers {
		ser := s.(*BaseSerializer)
		if ser.FirstVersion <= inputVersion &&
			inputVersion <= ser.LastVersion && ser.SupportsGetSet() {
			s := Serializer(ser)
			s.Read(is, obj)
		}
	}
}

func (wp *ObjectWrapper) Write(os *OsgOstream, obj *model.Object) {
	inputVersion := os.GetFileVersion(wp.Domain)
	for _, s := range wp.Serializers {
		ser := s.(*BaseSerializer)
		if ser.FirstVersion <= inputVersion &&
			inputVersion <= ser.LastVersion && ser.SupportsGetSet() {
			s := Serializer(ser)
			s.Write(os, obj)
		}
	}
}

func (wp *ObjectWrapper) ReadSchema(properties []string, types []SerType) bool {
	if len(wp.BackupSerializers) != 0 {
		wp.BackupSerializers = wp.Serializers
	}
	wp.Serializers = wp.Serializers[0:0]
	size := len(properties)
	serializersSize := len(wp.BackupSerializers)

	for i := 0; i < size; i++ {
		if serializersSize < i {
			break
		}
		prop := properties[i]
		ser := wp.BackupSerializers[i].(*BaseSerializer)
		if prop == ser.GetSerializerName() {
			wp.Serializers = append(wp.Serializers, wp.BackupSerializers[i])
		} else {
			for _, s := range wp.Serializers {
				ser := s.(*BaseSerializer)
				if prop != ser.GetSerializerName() {
					continue
				}
				wp.Serializers = append(wp.Serializers, ser)
			}
		}
	}
	return size == len(wp.Serializers)
}

func (wp *ObjectWrapper) WriteSchema(properties []string, types []SerType) {
	ssize := len(wp.Serializers)
	tsize := len(wp.TypeList)
	i := 0
	for {
		if ssize-1 <= i || tsize-1 <= i {
			break
		}
		s := wp.Serializers[i]
		ser := s.(*BaseSerializer)
		t := wp.TypeList[i]
		if ser.SupportsGetSet() {
			properties = append(properties, ser.GetSerializerName())
			types = append(types, t)
		}
		i++
	}
}

type AddPropFuncType func(obj *ObjectWrapper)
type AddPropCustomFuncType func(str string, obj *ObjectWrapper)

func NewRegisterCustomWrapperProxy(inst_func CreateInstanceFuncType, domain string, name string, associates string) {
	wrap := NewObjectWrapper2(name, domain, inst_func, associates)
	wrap.CreateInstanceFunc = inst_func
	wrap.Name = name
	wrap.Domain = domain
	wrap.SplitAssociates(associates, " ")
	ptr := &wrap
	GetObjectWrapperManager().AddWrap(ptr)
}

type objectWrapperManager struct {
	Wraps       map[string]*ObjectWrapper
	Compressors map[string]*CompressorStream
	GlobalMap   map[string]*IntLookup
}

func newObjectWrapperManager() *objectWrapperManager {
	if manager != nil {
		return manager
	}
	obj := objectWrapperManager{Wraps: make(map[string]*ObjectWrapper), Compressors: make(map[string]*CompressorStream), GlobalMap: make(map[string]*IntLookup)}
	lk := NewIntLookup()
	obj.GlobalMap["GL"] = &lk
	lk.Add("GL_ALPHA_TEST", model.GL_ALPHA_TEST)
	lk.Add("GL_BLEND", model.GL_BLEND)
	lk.Add("GL_COLOR_LOGIC_OP", model.GL_COLOR_LOGIC_OP)
	lk.Add("GL_COLOR_MATERIAL", model.GL_COLOR_MATERIAL)
	lk.Add("GL_CULL_FACE", model.GL_CULL_FACE)
	lk.Add("GL_DEPTH_TEST", model.GL_DEPTH_TEST)
	lk.Add("GL_FOG", model.GL_FOG)
	lk.Add("GL_FRAGMENT_PROGRAM_ARB", model.GL_FRAGMENT_PROGRAM_ARB)
	lk.Add("GL_LINE_STIPPLE", model.GL_LINE_STIPPLE)
	lk.Add("GL_POINT_SMOOTH", model.GL_POINT_SMOOTH)
	lk.Add("GL_POINT_SPRITE_ARB", model.GL_POINT_SPRITE_ARB)
	lk.Add("GL_POLYGON_OFFSET_FILL", model.GL_POLYGON_OFFSET_FILL)
	lk.Add("GL_POLYGON_OFFSET_LINE", model.GL_POLYGON_OFFSET_LINE)
	lk.Add("GL_POLYGON_OFFSET_POINT", model.GL_POLYGON_OFFSET_POINT)
	lk.Add("GL_POLYGON_STIPPLE", model.GL_POLYGON_STIPPLE)
	lk.Add("GL_SCISSOR_TEST", model.GL_SCISSOR_TEST)
	lk.Add("GL_STENCIL_TEST", model.GL_STENCIL_TEST)
	lk.Add("GL_STENCIL_TEST_TWO_SIDE", model.GL_STENCIL_TEST_TWO_SIDE)
	lk.Add("GL_VERTEX_PROGRAM_ARB", model.GL_VERTEX_PROGRAM_ARB)

	lk.Add("GL_COLOR_SUM", model.GL_COLOR_SUM)
	lk.Add("GL_LIGHTING", model.GL_LIGHTING)
	lk.Add("GL_NORMALIZE", model.GL_NORMALIZE)
	lk.Add("GL_RESCALE_NORMAL", model.GL_RESCALE_NORMAL)

	lk.Add("GL_TEXTURE_1D", model.GL_TEXTURE_1D)
	lk.Add("GL_TEXTURE_2D", model.GL_TEXTURE_2D)
	lk.Add("GL_TEXTURE_3D", model.GL_TEXTURE_3D)
	lk.Add("GL_TEXTURE_CUBE_MAP", model.GL_TEXTURE_CUBE_MAP)
	lk.Add("GL_TEXTURE_RECTANGLE", model.GL_TEXTURE_RECTANGLE)
	lk.Add("GL_TEXTURE_GEN_Q", model.GL_TEXTURE_GEN_Q)
	lk.Add("GL_TEXTURE_GEN_R", model.GL_TEXTURE_GEN_R)
	lk.Add("GL_TEXTURE_GEN_S", model.GL_TEXTURE_GEN_S)
	lk.Add("GL_TEXTURE_GEN_T", model.GL_TEXTURE_GEN_T)

	lk.Add("GL_CLIP_PLANE0", model.GL_CLIP_PLANE0)
	lk.Add("GL_CLIP_PLANE1", model.GL_CLIP_PLANE1)
	lk.Add("GL_CLIP_PLANE2", model.GL_CLIP_PLANE2)
	lk.Add("GL_CLIP_PLANE3", model.GL_CLIP_PLANE3)
	lk.Add("GL_CLIP_PLANE4", model.GL_CLIP_PLANE4)
	lk.Add("GL_CLIP_PLANE5", model.GL_CLIP_PLANE5)

	lk.Add("GL_LIGHT0", model.GL_LIGHT0)
	lk.Add("GL_LIGHT1", model.GL_LIGHT1)
	lk.Add("GL_LIGHT2", model.GL_LIGHT2)
	lk.Add("GL_LIGHT3", model.GL_LIGHT3)
	lk.Add("GL_LIGHT4", model.GL_LIGHT4)
	lk.Add("GL_LIGHT5", model.GL_LIGHT5)
	lk.Add("GL_LIGHT6", model.GL_LIGHT6)
	lk.Add("GL_LIGHT7", model.GL_LIGHT7)

	lk.Add("GL_VERTEX_PROGRAM_POINT_SIZE", model.GL_VERTEX_PROGRAM_POINT_SIZE)
	lk.Add("GL_VERTEX_PROGRAM_TWO_SIDE", model.GL_VERTEX_PROGRAM_TWO_SIDE)

	// Functions
	lk.Add("NEVER", model.GL_NEVER)
	lk.Add("LESS", model.GL_LESS)
	lk.Add("EQUAL", model.GL_EQUAL)
	lk.Add("LEQUAL", model.GL_LEQUAL)
	lk.Add("GREATER", model.GL_GREATER)
	lk.Add("NOTEQUAL", model.GL_NOTEQUAL)
	lk.Add("GEQUAL", model.GL_GEQUAL)
	lk.Add("ALWAYS", model.GL_ALWAYS)

	// Texture environment states
	lk.Add("REPLACE", model.GL_REPLACE)
	lk.Add("MODULATE", model.GL_MODULATE)
	lk.Add("Add", model.GL_ADD)
	lk.Add("Add_SIGNED", model.GL_ADD_SIGNED_ARB)
	lk.Add("INTERPOLATE", model.GL_INTERPOLATE_ARB)
	lk.Add("SUBTRACT", model.GL_SUBTRACT_ARB)
	lk.Add("DOT3_RGB", model.GL_DOT3_RGB_ARB)
	lk.Add("DOT3_RGBA", model.GL_DOT3_RGBA_ARB)

	lk.Add("CONSTANT", model.GL_CONSTANT_ARB)
	lk.Add("PRIMARY_COLOR", model.GL_PRIMARY_COLOR_ARB)
	lk.Add("PREVIOUS", model.GL_PREVIOUS_ARB)
	lk.Add("TEXTURE", model.GL_TEXTURE)
	lk.Add("TEXTURE0", model.GL_TEXTURE0)
	lk.Add("TEXTURE1", model.GL_TEXTURE0+1)
	lk.Add("TEXTURE2", model.GL_TEXTURE0+2)
	lk.Add("TEXTURE3", model.GL_TEXTURE0+3)
	lk.Add("TEXTURE4", model.GL_TEXTURE0+4)
	lk.Add("TEXTURE5", model.GL_TEXTURE0+5)
	lk.Add("TEXTURE6", model.GL_TEXTURE0+6)
	lk.Add("TEXTURE7", model.GL_TEXTURE0+7)

	// Texture clamp modes
	lk.Add("CLAMP", model.GL_CLAMP)
	lk.Add("CLAMP_TO_EDGE", model.GL_CLAMP_TO_EDGE)
	lk.Add("CLAMP_TO_BORDER", model.GL_CLAMP_TO_BORDER_ARB)
	lk.Add("REPEAT", model.GL_REPEAT)
	lk.Add("MIRROR", model.GL_MIRRORED_REPEAT_IBM)

	// Texture filter modes
	lk.Add("LINEAR", model.GL_LINEAR)
	lk.Add("LINEAR_MIPMAP_LINEAR", model.GL_LINEAR_MIPMAP_LINEAR)
	lk.Add("LINEAR_MIPMAP_NEAREST", model.GL_LINEAR_MIPMAP_NEAREST)
	lk.Add("NEAREST", model.GL_NEAREST)
	lk.Add("NEAREST_MIPMAP_LINEAR", model.GL_NEAREST_MIPMAP_LINEAR)
	lk.Add("NEAREST_MIPMAP_NEAREST", model.GL_NEAREST_MIPMAP_NEAREST)

	// Texture formats
	lk.Add("GL_INTENSITY", model.GL_INTENSITY)
	lk.Add("GL_LUMINANCE", model.GL_LUMINANCE)
	lk.Add("GL_ALPHA", model.GL_ALPHA)
	lk.Add("GL_LUMINANCE_ALPHA", model.GL_LUMINANCE_ALPHA)
	lk.Add("GL_RGB", model.GL_RGB)
	lk.Add("GL_RGBA", model.GL_RGBA)
	lk.Add("GL_COMPRESSED_ALPHA_ARB", model.GL_COMPRESSED_ALPHA_ARB)
	lk.Add("GL_COMPRESSED_LUMINANCE_ARB", model.GL_COMPRESSED_LUMINANCE_ARB)
	lk.Add("GL_COMPRESSED_INTENSITY_ARB", model.GL_COMPRESSED_INTENSITY_ARB)
	lk.Add("GL_COMPRESSED_LUMINANCE_ALPHA_ARB",
		model.GL_COMPRESSED_LUMINANCE_ALPHA_ARB)
	lk.Add("GL_COMPRESSED_RGB_ARB", model.GL_COMPRESSED_RGB_ARB)
	lk.Add("GL_COMPRESSED_RGBA_ARB", model.GL_COMPRESSED_RGBA_ARB)
	lk.Add("GL_COMPRESSED_RGB_S3TC_DXT1_EXT",
		model.GL_COMPRESSED_RGB_S3TC_DXT1_EXT)
	lk.Add("GL_COMPRESSED_RGBA_S3TC_DXT1_EXT",
		model.GL_COMPRESSED_RGBA_S3TC_DXT1_EXT)
	lk.Add("GL_COMPRESSED_RGBA_S3TC_DXT3_EXT",
		model.GL_COMPRESSED_RGBA_S3TC_DXT3_EXT)
	lk.Add("GL_COMPRESSED_RGBA_S3TC_DXT5_EXT",
		model.GL_COMPRESSED_RGBA_S3TC_DXT5_EXT)
	lk.Add("GL_COMPRESSED_RGB_PVRTC_4BPPV1_IMG",
		model.GL_COMPRESSED_RGB_PVRTC_4BPPV1_IMG)
	lk.Add("GL_COMPRESSED_RGB_PVRTC_2BPPV1_IMG",
		model.GL_COMPRESSED_RGB_PVRTC_2BPPV1_IMG)
	lk.Add("GL_COMPRESSED_RGBA_PVRTC_4BPPV1_IMG",
		model.GL_COMPRESSED_RGBA_PVRTC_4BPPV1_IMG)
	lk.Add("GL_COMPRESSED_RGBA_PVRTC_2BPPV1_IMG",
		model.GL_COMPRESSED_RGBA_PVRTC_2BPPV1_IMG)
	lk.Add("GL_ETC1_RGB8_OES", model.GL_ETC1_RGB8_OES)
	lk.Add("GL_COMPRESSED_RGB8_ETC2", model.GL_COMPRESSED_RGB8_ETC2)
	lk.Add("GL_COMPRESSED_SRGB8_ETC2", model.GL_COMPRESSED_SRGB8_ETC2)
	lk.Add("GL_COMPRESSED_RGB8_PUNCHTHROUGH_ALPHA1_ETC2",
		model.GL_COMPRESSED_RGB8_PUNCHTHROUGH_ALPHA1_ETC2)
	lk.Add("GL_COMPRESSED_SRGB8_PUNCHTHROUGH_ALPHA1_ETC2",
		model.GL_COMPRESSED_SRGB8_PUNCHTHROUGH_ALPHA1_ETC2)
	lk.Add("GL_COMPRESSED_RGBA8_ETC2_EAC", model.GL_COMPRESSED_RGBA8_ETC2_EAC)
	lk.Add("GL_COMPRESSED_SRGB8_ALPHA8_ETC2_EAC",
		model.GL_COMPRESSED_SRGB8_ALPHA8_ETC2_EAC)
	lk.Add("GL_COMPRESSED_R11_EAC", model.GL_COMPRESSED_R11_EAC)
	lk.Add("GL_COMPRESSED_SIGNED_R11_EAC", model.GL_COMPRESSED_SIGNED_R11_EAC)
	lk.Add("GL_COMPRESSED_RG11_EAC", model.GL_COMPRESSED_RG11_EAC)
	lk.Add("GL_COMPRESSED_SIGNED_RG11_EAC", model.GL_COMPRESSED_SIGNED_RG11_EAC)

	// Texture source types
	lk.Add("GL_BYTE", model.GL_BYTE)
	lk.Add("GL_SHORT", model.GL_SHORT)
	lk.Add("GL_INT", model.GL_INT)
	lk.Add("GL_FLOAT", model.GL_FLOAT)
	lk.Add("GL_DOUBLE", model.GL_DOUBLE)
	lk.Add("GL_UNSIGNED_BYTE", model.GL_UNSIGNED_BYTE)
	lk.Add("GL_UNSIGNED_SHORT", model.GL_UNSIGNED_SHORT)
	lk.Add("GL_UNSIGNED_INT", model.GL_UNSIGNED_INT)

	// Blend values
	lk.Add("DST_ALPHA", model.GL_DST_ALPHA)
	lk.Add("DST_COLOR", model.GL_DST_COLOR)
	lk.Add("ONE", model.GL_ONE)
	lk.Add("ONE_MINUS_DST_ALPHA", model.GL_ONE_MINUS_DST_ALPHA)
	lk.Add("ONE_MINUS_DST_COLOR", model.GL_ONE_MINUS_DST_COLOR)
	lk.Add("ONE_MINUS_SRC_ALPHA", model.GL_ONE_MINUS_SRC_ALPHA)
	lk.Add("ONE_MINUS_SRC_COLOR", model.GL_ONE_MINUS_SRC_COLOR)
	lk.Add("SRC_ALPHA", model.GL_SRC_ALPHA)
	lk.Add("SRC_ALPHA_SATURATE", model.GL_SRC_ALPHA_SATURATE)
	lk.Add("SRC_COLOR", model.GL_SRC_COLOR)
	lk.Add("CONSTANT_COLOR", model.GL_CONSTANT_COLOR)
	lk.Add("ONE_MINUS_CONSTANT_COLOR", model.GL_ONE_MINUS_CONSTANT_COLOR)
	lk.Add("CONSTANT_ALPHA", model.GL_CONSTANT_ALPHA)
	lk.Add("ONE_MINUS_CONSTANT_ALPHA", model.GL_ONE_MINUS_CONSTANT_ALPHA)
	lk.Add("ZERO", model.GL_ZERO)

	// Fog coordinate sources
	lk.Add("COORDINATE", model.GL_FOG_COORDINATE)
	lk.Add("DEPTH", model.GL_FRAGMENT_DEPTH)

	// Hint targets
	lk.Add("FOG_HINT", model.GL_FOG_HINT)
	lk.Add("GENERATE_MIPMAP_HINT", model.GL_GENERATE_MIPMAP_HINT_SGIS)
	lk.Add("LINE_SMOOTH_HINT", model.GL_LINE_SMOOTH_HINT)
	lk.Add("PERSPECTIVE_CORRECTION_HINT", model.GL_PERSPECTIVE_CORRECTION_HINT)
	lk.Add("POINT_SMOOTH_HINT", model.GL_POINT_SMOOTH_HINT)
	lk.Add("POLYGON_SMOOTH_HINT", model.GL_POLYGON_SMOOTH_HINT)
	lk.Add("TEXTURE_COMPRESSION_HINT", model.GL_TEXTURE_COMPRESSION_HINT_ARB)
	lk.Add("FRAGMENT_SHADER_DERIVATIVE_HINT",
		model.GL_FRAGMENT_SHADER_DERIVATIVE_HINT)

	// Polygon modes
	lk.Add("POINT", model.GL_POINT)
	lk.Add("LINE", model.GL_LINE)
	lk.Add("FILL", model.GL_FILL)

	// Misc
	lk.Add("BACK", model.GL_BACK)
	lk.Add("FRONT", model.GL_FRONT)
	lk.Add("FRONT_AND_BACK", model.GL_FRONT_AND_BACK)
	lk.Add("FIXED_ONLY", model.GL_FIXED_ONLY)
	lk.Add("FASTEST", model.GL_FASTEST)
	lk.Add("NICEST", model.GL_NICEST)
	lk.Add("DONT_CARE", model.GL_DONT_CARE)

	arrayTable := NewIntLookup()
	obj.GlobalMap["ArrayType"] = &arrayTable

	arrayTable.Add("ByteArray", model.ID_BYTE_ARRAY)
	arrayTable.Add("UByteArray", model.ID_UBYTE_ARRAY)
	arrayTable.Add("ShortArray", model.ID_SHORT_ARRAY)
	arrayTable.Add("UShortArray", model.ID_USHORT_ARRAY)
	arrayTable.Add("IntArray", model.ID_INT_ARRAY)
	arrayTable.Add("UIntArray", model.ID_UINT_ARRAY)
	arrayTable.Add("FloatArray", model.ID_FLOAT_ARRAY)
	arrayTable.Add("DoubleArray", model.ID_DOUBLE_ARRAY)

	arrayTable.Add("Vec2bArray", model.ID_VEC2B_ARRAY)
	arrayTable.Add("Vec3bArray", model.ID_VEC3B_ARRAY)
	arrayTable.Add("Vec4bArray", model.ID_VEC4B_ARRAY)
	arrayTable.Add("Vec2ubArray", model.ID_VEC2UB_ARRAY)
	arrayTable.Add("Vec3ubArray", model.ID_VEC3UB_ARRAY)
	arrayTable.Add("Vec4ubArray", model.ID_VEC4UB_ARRAY)
	arrayTable.Add("Vec2sArray", model.ID_VEC2S_ARRAY)
	arrayTable.Add("Vec3sArray", model.ID_VEC3S_ARRAY)
	arrayTable.Add("Vec4sArray", model.ID_VEC4S_ARRAY)
	arrayTable.Add("Vec2usArray", model.ID_VEC2US_ARRAY)
	arrayTable.Add("Vec3usArray", model.ID_VEC3US_ARRAY)
	arrayTable.Add("Vec4usArray", model.ID_VEC4US_ARRAY)
	arrayTable.Add("Vec2fArray", model.ID_VEC2_ARRAY)
	arrayTable.Add("Vec3fArray", model.ID_VEC3_ARRAY)
	arrayTable.Add("Vec4fArray", model.ID_VEC4_ARRAY)
	arrayTable.Add("Vec2dArray", model.ID_VEC2D_ARRAY)
	arrayTable.Add("Vec3dArray", model.ID_VEC3D_ARRAY)
	arrayTable.Add("Vec4dArray", model.ID_VEC4D_ARRAY)

	arrayTable.Add("Vec2iArray", model.ID_VEC2I_ARRAY)
	arrayTable.Add("Vec3iArray", model.ID_VEC3I_ARRAY)
	arrayTable.Add("Vec4iArray", model.ID_VEC4I_ARRAY)
	arrayTable.Add("Vec2uiArray", model.ID_VEC2UI_ARRAY)
	arrayTable.Add("Vec3uiArray", model.ID_VEC3UI_ARRAY)
	arrayTable.Add("Vec4uiArray", model.ID_VEC4UI_ARRAY)

	primitiveTable := NewIntLookup()
	obj.GlobalMap["PrimitiveType"] = &primitiveTable

	primitiveTable.Add("DrawArrays", model.ID_DRAWARRAYS)
	primitiveTable.Add("DrawArraysLength", model.ID_DRAWARRAY_LENGTH)
	primitiveTable.Add("DrawElementsUByte", model.ID_DRAWELEMENTS_UBYTE)
	primitiveTable.Add("DrawElementsUShort", model.ID_DRAWELEMENTS_USHORT)
	primitiveTable.Add("DrawElementsUInt", model.ID_DRAWELEMENTS_UINT)

	primitiveTable.Add("GL_POINTS", model.GL_POINTS)
	primitiveTable.Add("GL_LINES", model.GL_LINES)
	primitiveTable.Add("GL_LINE_STRIP", model.GL_LINE_STRIP)
	primitiveTable.Add("GL_LINE_LOOP", model.GL_LINE_LOOP)
	primitiveTable.Add("GL_TRIANGLES", model.GL_TRIANGLES)
	primitiveTable.Add("GL_TRIANGLE_STRIP", model.GL_TRIANGLE_STRIP)
	primitiveTable.Add("GL_TRIANGLE_FAN", model.GL_TRIANGLE_FAN)
	primitiveTable.Add("GL_QUADS", model.GL_QUADS)
	primitiveTable.Add("GL_QUAD_STRIP", model.GL_QUAD_STRIP)
	primitiveTable.Add("GL_POLYGON", model.GL_POLYGON)

	primitiveTable.Add2("GL_LINES_ADJACENCY_EXT", "GL_LINES_ADJACENCY",
		model.GL_LINES_ADJACENCY)
	primitiveTable.Add2("GL_LINE_STRIP_ADJACENCY_EXT", "GL_LINE_STRIP_ADJACENCY",
		model.GL_LINE_STRIP_ADJACENCY)
	primitiveTable.Add2("GL_TRIANGLES_ADJACENCY_EXT", "GL_TRIANGLES_ADJACENCY",
		model.GL_TRIANGLES_ADJACENCY)
	primitiveTable.Add2("GL_TRIANGLE_STRIP_ADJACENCY_EXT",
		"GL_TRIANGLE_STRIP_ADJACENCY",
		model.GL_TRIANGLE_STRIP_ADJACENCY)

	primitiveTable.Add("GL_PATCHES", model.GL_PATCHES)
	return &obj
}

func GetObjectWrapperManager() *objectWrapperManager {
	return manager
}

func (man *objectWrapperManager) AddWrap(wrap *ObjectWrapper) {
	if wrap == nil {
		return
	}
	manager.Wraps[strings.ToLower(wrap.Name)] = wrap
}

func (man *objectWrapperManager) RemoveWrap(wrap *ObjectWrapper) {
	nm := strings.ToLower(wrap.Name)
	delete(manager.Wraps, nm)
}

func (man *objectWrapperManager) FindWrap(str string) *ObjectWrapper {
	nm := strings.ToLower(str)
	nm = "flywave::" + nm
	w, ok := manager.Wraps[nm]
	if ok {
		return w
	}
	return nil
}

func (man *objectWrapperManager) FindLookup(group string) *IntLookup {
	lk, ok := man.GlobalMap[group]
	if !ok {
		lk, ok = man.GlobalMap["GL"]
	}
	return lk
}

func (man *objectWrapperManager) AddCompressor(st *CompressorStream) {
	man.Compressors[st.Name] = st
}

func (man *objectWrapperManager) RemoveCompressor(st *CompressorStream) {
	if st == nil {
		return
	}
	_, ok := man.Compressors[st.Name]
	if ok {
		delete(man.Compressors, st.Name)
	}
}

func (man *objectWrapperManager) FindCompressor(st string) *CompressorStream {
	lk, ok := man.Compressors[st]
	if ok {
		return lk
	}
	return nil
}

type UpdateWrapperVersionProxy struct {
	Wrap        *ObjectWrapper
	LastVersion int
}

func AddUpdateWrapperVersionProxy(w *ObjectWrapper, v int) {
	w.Version = v
}
