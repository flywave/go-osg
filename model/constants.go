package model

type Sizes uint32
type TypeIds uint32
type Glenum uint32
type ArrayTable uint32
type TextureInternalFormatMode uint32
type TextureShadowCompareFunc uint32
type TextureShadowTextureMode uint32
type PrimitiveTableEnum uint32

const (
	INDENT_VALUE uint32 = 2
	BOOL_SIZE    Sizes  = 1
	CHAR_SIZE    Sizes  = 1
	SHORT_SIZE   Sizes  = 2
	INT_SIZE     Sizes  = 4
	LONG_SIZE    Sizes  = 4
	INT64_SIZE   Sizes  = 8
	FLOAT_SIZE   Sizes  = 4
	DOUBLE_SIZE  Sizes  = 8
	GLENUM_SIZE  Sizes  = 4

	ID_BYTE_ARRAY   TypeIds = 0
	ID_UBYTE_ARRAY  TypeIds = 1
	ID_SHORT_ARRAY  TypeIds = 2
	ID_USHORT_ARRAY TypeIds = 3
	ID_INT_ARRAY    TypeIds = 4
	ID_UINT_ARRAY   TypeIds = 5
	ID_FLOAT_ARRAY  TypeIds = 6
	ID_DOUBLE_ARRAY TypeIds = 7
	ID_VEC2B_ARRAY  TypeIds = 8
	ID_VEC3B_ARRAY  TypeIds = 9
	ID_VEC4B_ARRAY  TypeIds = 10
	ID_VEC4UB_ARRAY TypeIds = 11
	ID_VEC2S_ARRAY  TypeIds = 12
	ID_VEC3S_ARRAY  TypeIds = 13
	ID_VEC4S_ARRAY  TypeIds = 14
	ID_VEC2_ARRAY   TypeIds = 15
	ID_VEC3_ARRAY   TypeIds = 16
	ID_VEC4_ARRAY   TypeIds = 17
	ID_VEC2D_ARRAY  TypeIds = 18
	ID_VEC3D_ARRAY  TypeIds = 19
	ID_VEC4D_ARRAY  TypeIds = 20
	ID_VEC2UB_ARRAY TypeIds = 21
	ID_VEC3UB_ARRAY TypeIds = 22
	ID_VEC2US_ARRAY TypeIds = 23
	ID_VEC3US_ARRAY TypeIds = 24
	ID_VEC4US_ARRAY TypeIds = 25

	ID_VEC2I_ARRAY  TypeIds = 26
	ID_VEC3I_ARRAY  TypeIds = 27
	ID_VEC4I_ARRAY  TypeIds = 28
	ID_VEC2UI_ARRAY TypeIds = 29
	ID_VEC3UI_ARRAY TypeIds = 30
	ID_VEC4UI_ARRAY TypeIds = 31

	ID_UINT64_ARRAY TypeIds = 32
	ID_INT64_ARRAY  TypeIds = 33

	ID_DRAWARRAYS            TypeIds = 50
	ID_DRAWARRAY_LENGTH      TypeIds = 51
	ID_DRAWELEMENTS_UBYTE    TypeIds = 52
	ID_DRAWELEMENTS_USHORT   TypeIds = 53
	ID_DRAWELEMENTS_UINT     TypeIds = 54
	GL_ALPHA_TEST            Glenum  = 0x0BC0
	GL_BLEND                 Glenum  = 0x0BE2
	GL_COLOR_LOGIC_OP        Glenum  = 0x0BF2
	GL_COLOR_MATERIAL        Glenum  = 0x0B57
	GL_CULL_FACE             Glenum  = 0x0B44
	GL_DEPTH_TEST            Glenum  = 0x0B71
	GL_FOG                   Glenum  = 0x0B60
	GL_FRAGMENT_PROGRAM_ARB  Glenum  = 0x8804
	GL_LINE_STIPPLE          Glenum  = 0x0B24
	GL_POINT_SMOOTH          Glenum  = 0x0B10
	GL_POINT_SPRITE_ARB      Glenum  = 0x8861
	GL_POLYGON_OFFSET_FILL   Glenum  = 0x8037
	GL_POLYGON_OFFSET_LINE   Glenum  = 0x2A02
	GL_POLYGON_OFFSET_POINT  Glenum  = 0x2A01
	GL_POLYGON_STIPPLE       Glenum  = 0x0B42
	GL_SCISSOR_TEST          Glenum  = 0x0C11
	GL_STENCIL_TEST          Glenum  = 0x0B90
	GL_STENCIL_TEST_TWO_SIDE Glenum  = 0x8910
	GL_VERTEX_PROGRAM_ARB    Glenum  = 0x8620

	GL_COLOR_SUM      Glenum = 0x8458
	GL_LIGHTING       Glenum = 0x0B50
	GL_NORMALIZE      Glenum = 0x0BA1
	GL_RESCALE_NORMAL Glenum = 0x803A

	GL_TEXTURE_1D        Glenum = 0x0DE0
	GL_TEXTURE_2D        Glenum = 0x0DE1
	GL_TEXTURE_3D        Glenum = 0x806F
	GL_TEXTURE_BUFFER    Glenum = 0x8C2A
	GL_TEXTURE_CUBE_MAP  Glenum = 0x8513
	GL_TEXTURE_RECTANGLE Glenum = 0x84F5
	GL_TEXTURE_GEN_Q     Glenum = 0x0C63
	GL_TEXTURE_GEN_R     Glenum = 0x0C62
	GL_TEXTURE_GEN_S     Glenum = 0x0C60
	GL_TEXTURE_GEN_T     Glenum = 0x0C61
	GL_TEXTURE_2D_ARRAY  Glenum = 0x8C1A

	GL_CLIP_PLANE0 Glenum = 0x3000
	GL_CLIP_PLANE1 Glenum = 0x3001
	GL_CLIP_PLANE2 Glenum = 0x3002
	GL_CLIP_PLANE3 Glenum = 0x3003
	GL_CLIP_PLANE4 Glenum = 0x3004
	GL_CLIP_PLANE5 Glenum = 0x3005

	GL_LIGHT0 Glenum = 0x4000
	GL_LIGHT1 Glenum = 0x4001
	GL_LIGHT2 Glenum = 0x4002
	GL_LIGHT3 Glenum = 0x4003
	GL_LIGHT4 Glenum = 0x4004
	GL_LIGHT5 Glenum = 0x4005
	GL_LIGHT6 Glenum = 0x4006
	GL_LIGHT7 Glenum = 0x4007

	GL_VERTEX_PROGRAM_POINT_SIZE Glenum = 0x8642
	GL_VERTEX_PROGRAM_TWO_SIDE   Glenum = 0x8643

	// Functions
	GL_NEVER    Glenum = 0x0200
	GL_LESS     Glenum = 0x0201
	GL_EQUAL    Glenum = 0x0202
	GL_LEQUAL   Glenum = 0x0203
	GL_GREATER  Glenum = 0x0204
	GL_NOTEQUAL Glenum = 0x0205
	GL_GEQUAL   Glenum = 0x0206
	GL_ALWAYS   Glenum = 0x0207

	// Texture environment states
	GL_REPLACE     Glenum = 0x1E01
	GL_MODULATE    Glenum = 0x2100
	GL_ADD         Glenum = 0x0104
	GL_ADD_SIGNED  Glenum = 0x8574
	GL_INTERPOLATE Glenum = 0x8575
	GL_SUBTRACT    Glenum = 0x84E7
	GL_DOT3_RGB    Glenum = 0x86AE
	GL_DOT3_RGBA   Glenum = 0x86AF

	GL_CONSTANT      Glenum = 0x8576
	GL_PRIMARY_COLOR Glenum = 0x8577
	GL_PREVIOUS      Glenum = 0x8578
	GL_TEXTURE       Glenum = 0x1702
	GL_TEXTURE0      Glenum = 0x84C0
	GL_TEXTURE1      Glenum = 0x84C1
	GL_TEXTURE2      Glenum = 0x84C2
	GL_TEXTURE3      Glenum = 0x84C3
	GL_TEXTURE4      Glenum = 0x84C4
	GL_TEXTURE5      Glenum = 0x84C5
	GL_TEXTURE6      Glenum = 0x84C6
	GL_TEXTURE7      Glenum = 0x84C7

	GL_COMBINE_ARB        Glenum = 0x8570
	GL_COMBINE_RGB_ARB    Glenum = 0x8571
	GL_COMBINE_ALPHA_ARB  Glenum = 0x8572
	GL_SOURCE0_RGB_ARB    Glenum = 0x8580
	GL_SOURCE1_RGB_ARB    Glenum = 0x8581
	GL_SOURCE2_RGB_ARB    Glenum = 0x8582
	GL_SOURCE0_ALPHA_ARB  Glenum = 0x8588
	GL_SOURCE1_ALPHA_ARB  Glenum = 0x8589
	GL_SOURCE2_ALPHA_ARB  Glenum = 0x858A
	GL_OPERAND0_RGB_ARB   Glenum = 0x8590
	GL_OPERAND1_RGB_ARB   Glenum = 0x8591
	GL_OPERAND2_RGB_ARB   Glenum = 0x8592
	GL_OPERAND0_ALPHA_ARB Glenum = 0x8598
	GL_OPERAND1_ALPHA_ARB Glenum = 0x8599
	GL_OPERAND2_ALPHA_ARB Glenum = 0x859A
	GL_RGB_SCALE_ARB      Glenum = 0x8573
	GL_ADD_SIGNED_ARB     Glenum = 0x8574
	GL_INTERPOLATE_ARB    Glenum = 0x8575
	GL_SUBTRACT_ARB       Glenum = 0x84E7
	GL_CONSTANT_ARB       Glenum = 0x8576
	GL_PRIMARY_COLOR_ARB  Glenum = 0x8577
	GL_PREVIOUS_ARB       Glenum = 0x8578

	GL_DOT3_RGB_ARB  Glenum = 0x86AE
	GL_DOT3_RGBA_ARB Glenum = 0x86AF

	// Texture clamp modes
	GL_CLAMP               Glenum = 0x2900
	GL_CLAMP_TO_EDGE       Glenum = 0x812F
	GL_CLAMP_TO_BORDER     Glenum = 0x812D
	GL_REPEAT              Glenum = 0x2901
	GL_MIRROR              Glenum = 0x8370
	GL_CLAMP_TO_BORDER_ARB Glenum = 0x812D
	GL_MIRRORED_REPEAT_IBM Glenum = 0x8370

	// Texture filter modes
	GL_LINEAR                 Glenum = 0x2601
	GL_LINEAR_MIPMAP_LINEAR   Glenum = 0x2703
	GL_LINEAR_MIPMAP_NEAREST  Glenum = 0x2701
	GL_NEAREST                Glenum = 0x2600
	GL_NEAREST_MIPMAP_LINEAR  Glenum = 0x2702
	GL_NEAREST_MIPMAP_NEAREST Glenum = 0x2700

	// Texture formats
	GL_INTENSITY                                 Glenum = 0x8049
	GL_LUMINANCE                                 Glenum = 0x1909
	GL_ALPHA                                     Glenum = 0x1906
	GL_LUMINANCE_ALPHA                           Glenum = 0x190A
	GL_RGB                                       Glenum = 0x1907
	GL_RGBA                                      Glenum = 0x1908
	GL_COMPRESSED_ALPHA_ARB                      Glenum = 0x84E9
	GL_COMPRESSED_LUMINANCE_ARB                  Glenum = 0x84EA
	GL_COMPRESSED_INTENSITY_ARB                  Glenum = 0x84EC
	GL_COMPRESSED_LUMINANCE_ALPHA_ARB            Glenum = 0x84EB
	GL_COMPRESSED_RGB_ARB                        Glenum = 0x84ED
	GL_COMPRESSED_RGBA_ARB                       Glenum = 0x84EE
	GL_COMPRESSED_RGB_S3TC_DXT1_EXT              Glenum = 0x83F0
	GL_COMPRESSED_RGBA_S3TC_DXT1_EXT             Glenum = 0x83F1
	GL_COMPRESSED_RGBA_S3TC_DXT3_EXT             Glenum = 0x83F2
	GL_COMPRESSED_RGBA_S3TC_DXT5_EXT             Glenum = 0x83F3
	GL_COMPRESSED_RGB_PVRTC_4BPPV1_IMG           Glenum = 0x8C00
	GL_COMPRESSED_RGB_PVRTC_2BPPV1_IMG           Glenum = 0x8C01
	GL_COMPRESSED_RGBA_PVRTC_4BPPV1_IMG          Glenum = 0x8C02
	GL_COMPRESSED_RGBA_PVRTC_2BPPV1_IMG          Glenum = 0x8C03
	GL_ETC1_RGB8_OES                             Glenum = 0x8D64
	GL_COMPRESSED_RGB8_ETC2                      Glenum = 0x9274
	GL_COMPRESSED_SRGB8_ETC2                     Glenum = 0x9275
	GL_COMPRESSED_RGB8_PUNCHTHROUGH_ALPHA1_ETC2  Glenum = 0x9276
	GL_COMPRESSED_SRGB8_PUNCHTHROUGH_ALPHA1_ETC2 Glenum = 0x9277
	GL_COMPRESSED_RGBA8_ETC2_EAC                 Glenum = 0x9278
	GL_COMPRESSED_SRGB8_ALPHA8_ETC2_EAC          Glenum = 0x9279
	GL_COMPRESSED_R11_EAC                        Glenum = 0x9270
	GL_COMPRESSED_SIGNED_R11_EAC                 Glenum = 0x9271
	GL_COMPRESSED_RG11_EAC                       Glenum = 0x9272
	GL_COMPRESSED_SIGNED_RG11_EAC                Glenum = 0x9273

	// Texture source types
	GL_BYTE           Glenum = 0x1400
	GL_SHORT          Glenum = 0x1402
	GL_INT            Glenum = 0x1404
	GL_FLOAT          Glenum = 0x1406
	GL_DOUBLE         Glenum = 0x140A
	GL_UNSIGNED_BYTE  Glenum = 0x1401
	GL_UNSIGNED_SHORT Glenum = 0x1403
	GL_UNSIGNED_INT   Glenum = 0x1405

	// Blend values
	GL_DST_ALPHA                Glenum = 0x0304
	GL_DST_COLOR                Glenum = 0x0306
	GL_ONE                      Glenum = 1
	GL_ONE_MINUS_DST_ALPHA      Glenum = 0x0305
	GL_ONE_MINUS_DST_COLOR      Glenum = 0x0307
	GL_ONE_MINUS_SRC_ALPHA      Glenum = 0x0303
	GL_ONE_MINUS_SRC_COLOR      Glenum = 0x0301
	GL_SRC_ALPHA                Glenum = 0x0302
	GL_SRC_ALPHA_SATURATE       Glenum = 0x0308
	GL_SRC_COLOR                Glenum = 0x0300
	GL_CONSTANT_COLOR           Glenum = 0x8001
	GL_ONE_MINUS_CONSTANT_COLOR Glenum = 0x8002
	GL_CONSTANT_ALPHA           Glenum = 0x8003
	GL_ONE_MINUS_CONSTANT_ALPHA Glenum = 0x8004
	GL_ZERO                     Glenum = 0

	// Fog coordinate sources
	GL_COORDINATE Glenum = 0x8451
	GL_DEPTH      Glenum = 0x8452

	GL_GENERATE_MIPMAP_SGIS      Glenum = 0x8191
	GL_GENERATE_MIPMAP_HINT_SGIS Glenum = 0x8192

	GL_TEXTURE_COMPRESSION_HINT_ARB Glenum = 0x84EF

	// Hint targets
	GL_FOG_HINT                        Glenum = 0x0C54
	GL_GENERATE_MIPMAP_HINT            Glenum = 0x8192
	GL_LINE_SMOOTH_HINT                Glenum = 0x0C52
	GL_PERSPECTIVE_CORRECTION_HINT     Glenum = 0x0C50
	GL_POINT_SMOOTH_HINT               Glenum = 0x0C51
	GL_POLYGON_SMOOTH_HINT             Glenum = 0x0C53
	GL_TEXTURE_COMPRESSION_HINT        Glenum = 0x84EF
	GL_FRAGMENT_SHADER_DERIVATIVE_HINT Glenum = 0x8B8B

	GL_FOG_COORDINATE Glenum = 0x8451
	GL_FRAGMENT_DEPTH Glenum = 0x8452

	// Polygon modes
	GL_POINT Glenum = 0x1B00
	GL_LINE  Glenum = 0x1B01
	GL_FILL  Glenum = 0x1B00

	// Misc
	GL_BACK           Glenum = 0x0405
	GL_FRONT          Glenum = 0x0404
	GL_FRONT_AND_BACK Glenum = 0x0408
	GL_FIXED_ONLY     Glenum = 0x891D
	GL_FASTEST        Glenum = 0x1101
	GL_NICEST         Glenum = 0x1101
	GL_DONT_CARE      Glenum = 0x1100

	GL_OBJECT_LINEAR Glenum = 0x2401
	GL_EYE_LINEAR    Glenum = 0x2400
	GL_SPHERE_MAP    Glenum = 0x2402

	GL_NORMAL_MAP Glenum = 0x8511

	GL_REFLECTION_MAP Glenum = 0x8512

	GL_RED   Glenum = 0x1903
	GL_GREEN Glenum = 0x1904
	GL_BLUE  Glenum = 0x1905

	ArrayType ArrayTable = 0

	ByteArrayType  ArrayTable = 1
	ShortArrayType ArrayTable = 2
	IntArrayType   ArrayTable = 3

	UByteArrayType  ArrayTable = 4
	UShortArrayType ArrayTable = 5
	UIntArrayType   ArrayTable = 6

	FloatArrayType  ArrayTable = 7
	DoubleArrayType ArrayTable = 8

	Vec2bArrayType ArrayTable = 9
	Vec3bArrayType ArrayTable = 10
	Vec4bArrayType ArrayTable = 11

	Vec2sArrayType ArrayTable = 12
	Vec3sArrayType ArrayTable = 13
	Vec4sArrayType ArrayTable = 14

	Vec2iArrayType ArrayTable = 15
	Vec3iArrayType ArrayTable = 16
	Vec4iArrayType ArrayTable = 17

	Vec2ubArrayType ArrayTable = 18
	Vec3ubArrayType ArrayTable = 19
	Vec4ubArrayType ArrayTable = 20

	Vec2usArrayType ArrayTable = 21
	Vec3usArrayType ArrayTable = 22
	Vec4usArrayType ArrayTable = 23

	Vec2uiArrayType ArrayTable = 24
	Vec3uiArrayType ArrayTable = 25
	Vec4uiArrayType ArrayTable = 26

	Vec2ArrayType ArrayTable = 27
	Vec3ArrayType ArrayTable = 28
	Vec4ArrayType ArrayTable = 29

	Vec2dArrayType ArrayTable = 30
	Vec3dArrayType ArrayTable = 31
	Vec4dArrayType ArrayTable = 32

	MatrixArrayType  ArrayTable = 33
	MatrixdArrayType ArrayTable = 34

	QuatArrayType ArrayTable = 35

	UInt64ArrayType ArrayTable = 36
	Int64ArrayType  ArrayTable = 37

	LastArrayType ArrayTable = 37

	DRAWARRAYS         PrimitiveTableEnum = 50
	DRAWARRAYSLENGTH   PrimitiveTableEnum = 51
	DRAWElEMENTSUBYTE  PrimitiveTableEnum = 52
	DRAWElEMENTSUSHORT PrimitiveTableEnum = 53
	DRAWElEMENTSUINT   PrimitiveTableEnum = 54

	GL_POINTS         PrimitiveTableEnum = 0x0000
	GL_LINES          PrimitiveTableEnum = 0x0001
	GL_LINE_STRIP     PrimitiveTableEnum = 0x0003
	GL_LINE_LOOP      PrimitiveTableEnum = 0x0002
	GL_TRIANGLES      PrimitiveTableEnum = 0x0004
	GL_TRIANGLE_STRIP PrimitiveTableEnum = 0x0005
	GL_TRIANGLE_FAN   PrimitiveTableEnum = 0x0006
	GL_QUADS          PrimitiveTableEnum = 0x0007
	GL_QUAD_STRIP     PrimitiveTableEnum = 0x0008
	GL_POLYGON        PrimitiveTableEnum = 0x0009

	GL_LINES_ADJACENCY              PrimitiveTableEnum = 0x000A
	GL_LINES_ADJACENCY_EXT          PrimitiveTableEnum = 0x000A
	GL_LINE_STRIP_ADJACENCY_EXT     PrimitiveTableEnum = 0x000B
	GL_LINE_STRIP_ADJACENCY         PrimitiveTableEnum = 0x000B
	GL_TRIANGLES_ADJACENCY          PrimitiveTableEnum = 0x000C
	GL_TRIANGLES_ADJACENCY_EXT      PrimitiveTableEnum = 0x000C
	GL_TRIANGLE_STRIP_ADJACENCY     PrimitiveTableEnum = 0x000D
	GL_TRIANGLE_STRIP_ADJACENCY_EXT PrimitiveTableEnum = 0x000D

	GL_PATCHES PrimitiveTableEnum = 0x000E

	USE_IMAGE_DATA_FORMAT      TextureInternalFormatMode = 0
	USE_USER_DEFINED_FORMAT    TextureInternalFormatMode = 1
	USE_ARB_COMPRESSION        TextureInternalFormatMode = 2
	USE_S3TC_DXT1_COMPRESSION  TextureInternalFormatMode = 3
	USE_S3TC_DXT3_COMPRESSION  TextureInternalFormatMode = 4
	USE_S3TC_DXT5_COMPRESSION  TextureInternalFormatMode = 5
	USE_PVRTC_2BPP_COMPRESSION TextureInternalFormatMode = 6
	USE_PVRTC_4BPP_COMPRESSION TextureInternalFormatMode = 7
	USE_ETC_COMPRESSION        TextureInternalFormatMode = 8
	USE_ETC2_COMPRESSION       TextureInternalFormatMode = 9
	USE_RGTC1_COMPRESSION      TextureInternalFormatMode = 10
	USE_RGTC2_COMPRESSION      TextureInternalFormatMode = 11
	USE_S3TC_DXT1c_COMPRESSION TextureInternalFormatMode = 12
	USE_S3TC_DXT1a_COMPRESSION TextureInternalFormatMode = 13

	NEVER    TextureShadowCompareFunc = 0x0200
	LESS     TextureShadowCompareFunc = 0x0201
	EQUAL    TextureShadowCompareFunc = 0x0202
	LEQUAL   TextureShadowCompareFunc = 0x0203
	GREATER  TextureShadowCompareFunc = 0x0204
	NOTEQUAL TextureShadowCompareFunc = 0x0205
	GEQUAL   TextureShadowCompareFunc = 0x0206
	ALWAYS   TextureShadowCompareFunc = 0x0207

	LUMINANCE TextureShadowTextureMode = 0x1909
	INTENSITY TextureShadowTextureMode = 0x8049
	ALPHA     TextureShadowTextureMode = 0x1906
	NONE      TextureShadowTextureMode = 0x0000
)
