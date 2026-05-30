# go-osg

Go 语言实现的 OpenSceneGraph 格式读写库。支持 `.osgb`（二进制）、`.osgt`（ASCII）、`.osg2`（扩展格式）的读取，并提供 OSGB 倾斜摄影数据到 3D Tiles 的完整转换能力。

---

## 目录

- [项目结构](#项目结构)
- [快速开始](#快速开始)
  - [读取 OSGB 文件](#读取-osgb-文件)
  - [读取 OSGT（ASCII）文件](#读取-osgtascii-文件)
  - [遍历场景图](#遍历场景图)
  - [转换为 3D Tiles（GLB）](#转换为-3d-tilesglb)
  - [转换为 3D Tiles（B3DM + tileset.json）](#转换为-3d-tilesb3dm--tilesetjson)
- [架构说明](#架构说明)
  - [核心 OSG 读写层](#核心-osg-读写层)
  - [序列化器系统](#序列化器系统)
  - [版本阈值机制](#版本阈值机制)
  - [3D Tiles 转换层](#3d-tiles-转换层)
- [3D Tiles 转换详解](#3d-tiles-转换详解)
  - [完整转换流程](#完整转换流程)
  - [节点遍历逻辑](#节点遍历逻辑)
  - [坐标转换系统](#坐标转换系统)
  - [顶点提取](#顶点提取)
  - [索引提取与三角化](#索引提取与三角化)
  - [纹理材质提取](#纹理材质提取)
  - [网格简化](#网格简化)
  - [Draco 压缩](#draco-压缩)
  - [B3DM / GLB 生成](#b3dm--glb-生成)
  - [tileset.json 生成](#tilesetjson-生成)
  - [PagedLOD 递归加载](#pagedlod-递归加载)
- [转换配置选项](#转换配置选项)
  - [ConverterOptions 详解](#converteroptions-详解)
  - [典型配置场景](#典型配置场景)
- [OSG 二进制格式详解](#osg-二进制格式详解)
  - [文件头部结构](#文件头部结构)
  - [版本阈值表](#版本阈值表)
  - [数组格式差异（版本 112 分界）](#数组格式差异版本-112-分界)
  - [PrimitiveSet 格式差异（版本 96 分界）](#primitiveset-格式差异版本-96-分界)
  - [Image 格式差异（版本 94 分界）](#image-格式差异版本-94-分界)
  - [二进制大括号机制（版本 148 分界）](#二进制大括号机制版本-148-分界)
- [测试](#测试)
  - [运行全部测试](#运行全部测试)
  - [测试分类说明](#测试分类说明)
  - [测试覆盖范围](#测试覆盖范围)
- [构建依赖](#构建依赖)
  - [replace 指令说明](#replace-指令说明)
  - [编译 tiles3d 的前提](#编译-tiles3d-的前提)
- [常见问题](#常见问题)
  - [坐标转换失败（PROJ_DATA 未设置）](#坐标转换失败proj_data-未设置)
  - [缺少测试数据文件](#缺少测试数据文件)
  - [Draco 压缩失败](#draco-压缩失败)
  - [读取特定版本 OSGB 文件出错](#读取特定版本-osgb-文件出错)
- [代码参考（C++ 对照）](#代码参考c-对照)

---

## 项目结构

```
go-osg/
├── *.go                         # 核心 OSG 读写（包 osg）
│   ├── input_stream.go          #   输入流处理
│   ├── binary_input.go          #   二进制输入迭代器
│   ├── binary_output.go         #   二进制输出迭代器
│   ├── ascii_input.go           #   ASCII 输入迭代器
│   ├── ascii_output.go          #   ASCII 输出迭代器
│   ├── read_write.go            #   入口：ReadNode / ReadWrite
│   ├── register.go              #   ObjectWrapper 管理器
│   ├── serializer.go            #   序列化器类型定义
│   ├── geometry.go              #   Geometry 序列化器注册
│   ├── group.go / node.go       #   Group / Node 序列化器注册
│   ├── paged_lod.go / lod.go    #   PagedLOD / LOD 序列化器注册
│   ├── array.go                 #   Array 序列化器注册
│   ├── image.go                 #   Image 序列化器注册
│   ├── object.go                #   Object 序列化器注册
│   ├── state_set.go             #   StateSet 序列化器注册
│   ├── drawable.go              #   Drawable 序列化器注册
│   ├── primitive_set.go         #   PrimitiveSet 序列化器注册
│   ├── texture*.go              #   Texture 系列序列化器注册
│   ├── cull_face.go             #   CullFace 序列化器注册
│   ├── alpha_func.go            #   AlphaFunc 序列化器注册
│   ├── material.go              #   Material 序列化器注册
│   ├── tex_env.go / tex_gen.go  #   TexEnv / TexGen 序列化器注册
│   ├── transform.go             #   Transform 序列化器注册
│   ├── matrix_transform.go      #   MatrixTransform 序列化器注册
│   └── ...
├── model/                       # 数据模型（包 model）
│   ├── node.go / group.go       #   Node, Group, Geode 等
│   ├── geometry.go              #   Geometry（顶点、图元等）
│   ├── array.go                 #   Array（各种顶点数组类型）
│   ├── paged_lod.go / lod.go    #   PagedLOD, LOD
│   ├── state_set.go             #   StateSet 状态集
│   ├── constants.go             #   GL 常量和 OSG 类型常量
│   └── ...
├── tiles3d/                     # OSGB → 3D Tiles（包 tiles3d）
│   ├── converter.go             #   GeometryConverter：几何提取与坐标转换
│   ├── tileset.go               #   Converter / B3DMGenerator / TilesetGenerator
│   ├── coordinate.go            #   CoordinateTransformer：坐标转换
│   ├── metadata.go              #   metadata.xml 解析与 SRS 检测
│   ├── geoid.go                 #   GeoidConverter：大地水准面校正
│   └── options.go               #   配置类型定义
├── cmd/check_vertices/          # 诊断工具：检查 OSGB 二进制顶点布局
├── test_data/                   # 测试数据（cessna.osgb, simpleroom.osgt 等）
├── OpenSceneGraph-master/       # C++ 参考代码
├── cpp_compat_test.go           # C++ 一致性测试（63 个）
├── AGENTS.md                    # OpenCode 指令文件
└── README.md                    # 本文件
```

---

## 快速开始

### 读取 OSGB 文件

```go
package main

import (
    "fmt"
    "github.com/flywave/go-osg"
    "github.com/flywave/go-osg/model"
)

func main() {
    rw := osg.NewReadWrite()
    res := rw.ReadNode("model.osgb", nil)
    if res == nil || res.GetNode() == nil {
        panic("failed to read node")
    }
    node := res.GetNode() // model.NodeInterface

    // 类型断言访问具体类型
    switch n := node.(type) {
    case *model.Group:
        fmt.Printf("Group with %d children\n", len(n.GetChildren()))
    case *model.PagedLod:
        fmt.Printf("PagedLOD: center=%v, ranges=%v\n", n.Center, n.RangeList)
    }
}
```

### 读取 OSGT（ASCII）文件

```go
opts := osg.NewOsgIstreamOptions()
opts.FileType = "Ascii"
res := rw.ReadNode("scene.osgt", opts)
node := res.GetNode()
```

### 遍历场景图

```go
import "github.com/flywave/go-osg/model"

func walkNode(n interface{}, depth int) {
    indent := strings.Repeat("  ", depth)
    switch v := n.(type) {
    case *model.Group:
        fmt.Printf("%sGroup\n", indent)
        for _, c := range v.GetChildren() {
            walkNode(c, depth+1)
        }
    case *model.Geode:
        fmt.Printf("%sGeode\n", indent)
        for _, c := range v.GetChildren() {
            walkNode(c, depth+1)
        }
    case *model.Geometry:
        fmt.Printf("%sGeometry: %d vertices\n", indent,
            len(v.VertexArray.Data.([][3]float32)))
    }
}
```

### 转换为 3D Tiles（GLB）

```go
import (
    "github.com/flywave/go-osg/tiles3d"
)

func convertToGLB(inputPath, outputPath string) error {
    opts := tiles3d.DefaultConverterOptions()
    glbData, err := tiles3d.OSGBToGLB(inputPath, opts)
    if err != nil {
        return err
    }
    return os.WriteFile(outputPath, glbData, 0644)
}
```

### 转换为 3D Tiles（B3DM + tileset.json）

```go
func convertTo3DTiles(inputPath, outputDir string) error {
    opts := tiles3d.DefaultConverterOptions()
    opts.SourceSRS = "EPSG:4548"  // CGCS2000 / 3-degree Gauss-Kruger CM 114E

    result, err := tiles3d.OSGBTo3DTiles(inputPath, outputDir, opts)
    if err != nil {
        return err
    }
    fmt.Printf("tileset.json: %d bytes\n", len(result.JSON))
    fmt.Printf("Bounding box center: (%f, %f, %f)\n",
        result.BoundingBox[0], result.BoundingBox[1], result.BoundingBox[2])
    return nil
}
```

---

## 架构说明

### 核心 OSG 读写层

OSG 文件的读写通过 `ReadWrite` 对象完成，该对象支持 `ReadNode`、`ReadImage`、`ReadObject` 等入口方法。根据文件扩展名（`.osgb`/`.osgt`/`.osg2`）自动选择二进制或 ASCII 输入迭代器。

文件读取流程：

```
ReadNode(path, opts)
  │
  ├─ PrepareReading() → 检查扩展名，设置 FileType
  ├─ OpenReader() → 打开文件
  ├─ ReadInputIterator() → 创建 BinaryInputIterator 或 AsciiInputIterator
  ├─ OsgIstream.Start() → 读取文件头部（类型、版本、属性）
  ├─ OsgIstream.Decompress() → 解压缩（可选）
  └─ OsgIstream.ReadObject() → 递归读取对象树
       │
       ├─ ReadString() → 读取类名（如 "osg::Group"）
       ├─ Read(BEGIN_BRACKET) → 读取块大小
       ├─ Read(UniqueID) → 读取对象 ID
       ├─ ReadObjectFields() → 查找 ObjectWrapper，遍历关联类
       │    ├─ Object.Read() → Name, DataVariance, UserDataContainer
       │    ├─ Node.Read() → InitialBound, StateSet, CullingActive...
       │    └─ Group.Read() → Children（递归调用 ReadObject）
       └─ AdvanceToCurrentEndBracket() → 跳过块剩余数据
```

### 序列化器系统

每个可序列化的 OSG 类通过 `ObjectWrapper` 注册，包装器包含：

- **类名**：如 `"Group"`、`"Geometry"`
- **关联链**：如 `"osg::Object osg::Node osg::Group"`
- **序列化器列表**：有序字段序列化器，每个有其版本范围

序列化器类型：

| 类型 | 含义 | C++ 宏 | 序列化行为（二进制） |
|------|------|--------|-------------------|
| RWUSER | 用户自定义 | ADD_USER_SERIALIZER | 读取 bool 标志位，再读取自定义数据 |
| RWOBJECT | 对象引用 | ADD_OBJECT_SERIALIZER | 读取 bool 标志位，再调用 ReadObject |
| RWVECTOR | 对象向量 | ADD_VECTOR_SERIALIZER | 读取 size，再依次读取各元素 |
| RWSTRING | 字符串 | ADD_STRING_SERIALIZER | 读取 int32 长度 + 字符串数据 |
| RWENUM | 枚举 | ADD_ENUM_SERIALIZER | 读取 int32 值 |
| RWBOOL | 布尔 | ADD_BOOL_SERIALIZER | 读取 1 字节（char）|
| RWINT / RWUINT | 整数 | ADD_INT_SERIALIZER / ADD_UINT_SERIALIZER | 读取 4 字节 |
| RWFLOAT | 浮点 | ADD_FLOAT_SERIALIZER | 读取 4 字节 |
| RWDOUBLE | 双精度 | ADD_DOUBLE_SERIALIZER | 读取 8 字节 |
| RWVEC3F / RWVEC3D | 向量 | ADD_VEC3D_SERIALIZER | 读取 3 个分量 |
| RWMATRIXF | 矩阵 | ADD_MATRIX_SERIALIZER | 读取 4×4 浮点矩阵 |
| RWGLENUM | GL 枚举 | ADD_GLENUM_SERIALIZER | 读取 GLenum（int32）|
| RWIMAGE | 图像 | ADD_IMAGE_SERIALIZER | 读取图像数据 |

### 版本阈值机制

ObjectWrapper 的序列化器通过 FirstVersion/LastVersion 控制生效范围：

```go
// Geometry 在版本 112 处切换：
uv := AddUpdateWrapperVersionProxy(wrap, 112)
wrap.MarkSerializerAsRemoved("VertexData")     // VertexData LastVersion=111
wrap.MarkSerializerAsRemoved("NormalData")     // NormalData  LastVersion=111
// ...
ser11 := NewObjectSerializer("VertexArray", ...) // FirstVersion=112
wrap.AddSerializer(ser11, RWOBJECT)
uv.SetLastVersion()
```

当读入文件版本 ≥ 112 时，VertexData（旧格式）被跳过，VertexArray（新格式）被读取。

### 3D Tiles 转换层

`tiles3d` 包将 OSGB 场景图转换为 3D Tiles 标准格式，核心组件：

```
Converter (tileset.go)
  ├── CoordinateTransformer (coordinate.go)  — 坐标投影变换
  ├── GeoidConverter (geoid.go)               — 大地水准面校正
  ├── GeometryConverter (converter.go)         — 几何提取 / 简化 / 压缩
  ├── TilesetGenerator (tileset.go)            — tileset.json 生成
  └── B3DMGenerator (tileset.go)               — B3DM / GLB 生成
```

---

## 3D Tiles 转换详解

### 完整转换流程

```
OSGBTo3DTiles(inputPath, outputPath, opts)
  │
  ├─ NewConverter(opts)
  │    ├─ NewCoordinateTransformer(sourceSRS, targetSRS)
  │    ├─ NewGeoidConverter(model, dataPath)
  │    ├─ NewGeometryConverter(opts, coordTrans, geoidConv)
  │    ├─ NewTilesetGenerator(opts)
  │    └─ NewB3DMGenerator(opts)
  │
  ├─ Converter.Convert(inputPath, outputPath)
  │    │
  │    ├─ detectMetadata() → 查找 metadata.xml，设置坐标参考
  │    ├─ LoadOSGB(inputPath) → 读取 OSGB 文件
  │    ├─ convertNodeToTile(node, id) → 创建根 Tile
  │    │    ├─ geomConverter.Convert(node) → 提取几何数据
  │    │    ├─ optimizeMesh() → 简化网格（可选）
  │    │    └─ calculateBoundingBox() / calculateGeometricError()
  │    │
  │    ├─ loadPagedLODs(tile, outputPath) → 递归加载子节点
  │    │    ├─ getPagedLODChildren(tile) → 获取子文件列表
  │    │    │    ├─ getPagedLODChildrenFromNode() → 从 PerRangeDataList 提取
  │    │    │    └─ scanDirectoryForSiblings() → 目录扫描获取同级图块
  │    │    └─ LoadOSGB(childPath) → 递归读取每个子文件
  │    │         └─ convertNodeToTile() → 创建子 Tile
  │    │
  │    ├─ extendTileBox(tile) → 合并子节点包围盒
  │    ├─ calcGeometricError(tile) → 递归计算几何误差
  │    ├─ tilesetGen.Generate(tile) → 生成 tileset.json 结构
  │    ├─ tilesetGen.Write() → 写入 tileset.json
  │    └─ writeTiles(tile, outputPath) → 写入所有 B3DM 文件
  │         └─ b3dmGen.Write(tile, path)
  │              ├─ buildGLTF(content) → 构建 GLTF 文档
  │              ├─ draco.EncodeAll(doc, opts) → Draco 压缩（可选）
  │              └─ b3dm.Write() → 编码为 B3DM 二进制格式
  │
  └─ 返回 ConvertResult{JSON, BoundingBox}
```

### 节点遍历逻辑

`GeometryConverter.Convert(node)` 遍历 OSG 节点树：

| 节点类型 | 处理方式 |
|----------|---------|
| `*model.Group` | 递归遍历所有子节点 |
| `*model.Lod` | 递归遍历第一个子节点（LOD 0） |
| `*model.PagedLod` | 递归遍历第一个内联子节点 |
| `*model.MatrixTransform` | 递归遍历所有子节点，对提取的顶点应用矩阵变换 |
| `*model.Geode` | 遍历所有 Drawable，找到 Geometry 后调用 extractGeometry |

### 坐标转换系统

`CoordinateTransformer` 支持多种坐标参考系：

```
SRS 类型检测：
  EPSG:4548     → SRSTypeEPSG
  ENU:114,34    → SRSTypeENU
  PROJCS[...]   → SRSTypeWKT
  unknown       → SRSTypeUnknown

坐标转换流程（投影坐标系 → 局部 ENU）：
  1. 读取 metadata.xml 获取 SRS 和 SRSOrigin
  2. 使用 go-proj 将坐标原点从源 SRS 转换为 WGS84
  3. go-proj 返回 (lat, lon) 在 OGC 顺序，需交换为 GIS 顺序 (lon, lat)
  4. 可选：大地水准面校正（正高 → 椭球高）
  5. WGS84 → ECEF（地心地固坐标系）
  6. ECEF → ENU（东-北-天局部坐标系）
  
对于未知 SRS：
  - 尝试常见投影（EPSG:4548, 4547, 4549, 4490, 4326）
  - 选择成功转换的投影
```

### 顶点提取

`extractVertices` 支持多种数组类型：

| 数组类型 | Go 类型 | 提取方式 |
|-----------|---------|----------|
| Vec3Array (IDVEC3ARRAY) | `[][3]float32` | 扁平化为 `[]float32` |
| FloatArray (IDFLOATARRAY) | `[]float32` | 直接复制 |
| DoubleArray (IDDOUBLEARRAY) | `[]float64` | 转换为 `float32` |
| ShortArray (IDSHORTARRAY) | `[]int16` | 转换为 `float32` |
| UShortArray (IDUSHORTARRAY) | `[]uint16` | 转换为 `float32` |

### 索引提取与三角化

`extractIndices` 支持所有 Primitive 类型：

| Primitive 类型 | 处理方式 |
|----------------|----------|
| DrawArrays | 生成 `[first, first+1, ..., first+count-1]` 连续索引 |
| DrawElementsUByte | 读取 `uint8` 索引并提升为 `uint32` |
| DrawElementsUShort | 读取 `uint16` 索引 |
| DrawElementsUInt | 直接传递 `uint32` 索引 |

四边形三角化：
- **GL_QUADS**：每 4 个顶点生成 2 个三角形 `(a,b,c)` + `(a,c,d)`
- **GL_QUAD_STRIP**：每 4 个顶点生成 2 个三角形 `(a,b,c)` + `(b,d,c)`

### 纹理材质提取

从 Geometry 的 StateSet 中提取：
- **纹理**：从 `TextureAttributeList` 中查找 `osg::Texture`，提取 `Image` 数据（内联数据或外部文件）
- **材质**：从 `AttributeList` 中查找 `osg::Material`，提取漫反射/镜面反射/自发光颜色

### 网格简化

`SimplifyMesh` 使用 meshoptimizer 库：

```go
// 启用简化：
opts.EnableSimplify = true
opts.SimplifyRatio = 0.5       // 保留 50% 的三角形
opts.SimplifyTargetError = 0.01 // 最大误差
```

简化流程：
1. 构建 meshopt.Mesh 对象（顶点、法线、纹理坐标、索引）
2. 调用 `meshopt.ProcessMesh(mesh, settings)`
3. 用简化后的索引替换原始索引

### Draco 压缩

支持 `KHR_draco_mesh_compression` 扩展：

```go
// 启用 Draco 压缩：
opts.EnableDraco = true
opts.DracoPositionBits = 14  // 位置量化位数（默认 11）
opts.DracoNormalBits = 10    // 法向量化位数（默认 10）
opts.DracoTexCoordBits = 12  // 纹理坐标量化位数（默认 12）
```

在 B3DM/GLB 生成时自动应用：
1. `buildGLTF(content)` 构建标准 GLTF 文档（含原始顶点缓冲）
2. `draco.EncodeAll(doc, quantOpts)` 压缩所有图元：
   - 读取访问器数据（位置、法线、纹理坐标、索引）
   - 编码为 Draco 压缩数据
   - 创建新的 Buffer + BufferView
   - 在图元上设置 KHR_draco_mesh_compression 扩展
   - 清空原始访问器的 BufferView 引用
   - 注册扩展到 ExtensionsUsed

### B3DM / GLB 生成

`B3DMGenerator` 支持两种输出格式：

- **B3DM**（Batch 3D Model）：标准 3D Tiles 瓦片格式，B3DM 头 + GLTF 文档
- **GLB**（GLTF Binary）：GLTF 二进制格式

生成流程：
1. `buildGLTF(content)`：构建完整的 GLTF 文档
   - 顶点缓冲区 → BufferView → Accessor（POSITION）
   - 法线缓冲区 → BufferView → Accessor（NORMAL）
   - 纹理坐标缓冲区 → BufferView → Accessor（TEXCOORD_0）
   - 索引缓冲区 → BufferView → Accessor（indices）
   - 材质 & 纹理（如存在）
   - 网格 & 节点 & 场景
2. `draco.EncodeAll()`（可选）
3. `b3dm.Write()` 或 `gltf.Encode()`

### tileset.json 生成

`TilesetGenerator` 生成符合 3D Tiles 1.0 标准的 tileset.json：

```json
{
  "asset": {
    "version": "1.0",
    "generator": "go-osg-3dtiles",
    "gltfUpAxis": "Z"
  },
  "geometricError": 500.0,
  "root": {
    "boundingVolume": { "box": [...] },
    "geometricError": 100.0,
    "refine": "REPLACE",
    "content": { "uri": "Tile_+003_+003_L18_000.osgb.b3dm" },
    "children": [...],
    "transform": [...]  // WGS84 → ECEF 变换矩阵（有地理参考时）
  }
}
```

### PagedLOD 递归加载

`Converter` 支持递归加载 OSGB 的 PagedLOD 层级结构：

1. `getPagedLODChildrenFromNode(tile)`：从 PagedLOD 的 `PerRangeDataList` 中提取子文件名
2. `scanDirectoryForSiblingsBylod(tile)`：扫描同级目录中具有相同 LOD 级别的图块文件
3. 合并两个列表（去重）
4. 对每个子文件调用 `LoadOSGB`，递归创建子 Tile
5. `extendTileBox(tile)`：递归合并子节点包围盒
6. `calcGeometricError(tile)`：递归计算全层级几何误差

可通过 `MaxLOD` 限制最大加载层级：
```go
opts.MaxLOD = 15  // 只加载 LOD 0~15 的图块
```

---

## 转换配置选项

### ConverterOptions 详解

```go
type ConverterOptions struct {
    // 输入输出
    InputPath  string
    OutputPath string

    // 坐标系
    SourceSRS string             // 源坐标系，如 "EPSG:4548"
    TargetSRS string             // 目标坐标系，默认 "EPSG:4326"

    // 大地水准面
    GeoidModel    string         // "none" / "EGM96" / "EGM2008"
    GeoidDataPath string         // 大地水准面数据路径

    // 网格简化（meshoptimizer）
    EnableSimplify      bool     // 启用简化
    SimplifyRatio       float32  // 简化比例 0.0~1.0，1.0 为不简化
    SimplifyTargetError float32  // 简化最大误差，默认 0.01

    // Draco 压缩
    EnableDraco         bool     // 启用 Draco 压缩
    DracoPositionBits   int      // 位置量化位数，默认 11
    DracoNormalBits     int      // 法向量化位数，默认 10
    DracoTexCoordBits   int      // 纹理坐标量化位数，默认 12

    // 纹理
    EnableTexture   bool         // 是否输出纹理，默认 true
    TextureCompress string       // 纹理压缩方式，默认 "jpeg"
    EnableKTX2      bool         // 是否使用 KTX2 格式

    // 渲染
    EnableUnlit bool             // 是否使用 KHR_materials_unlit，默认 true

    // 手动指定中心（无 metadata.xml 时）
    CenterLongitude float64
    CenterLatitude  float64
    CenterHeight    float64

    // 其他
    GeometricError float64       // 根节点几何误差，默认 500.0
    MaxLOD int                   // 最大 LOD 层级，-1 为不限制
}
```

### 典型配置场景

**场景 1：无坐标参考的裸 OSGB → GLB**

```go
opts := tiles3d.DefaultConverterOptions()
opts.SourceSRS = ""
glbData, _ := tiles3d.OSGBToGLB("model.osgb", opts)
```

**场景 2：EPSG:4548 倾斜摄影 → 3D Tiles（含 Draco 压缩）**

```go
opts := tiles3d.DefaultConverterOptions()
opts.SourceSRS = "EPSG:4548"
opts.EnableDraco = true
opts.DracoPositionBits = 14
result, _ := tiles3d.OSGBTo3DTiles("Tile_+002_+000.osgb", "./output", opts)
```

**场景 3：ENU 坐标系 + 网格简化**

```go
opts := tiles3d.DefaultConverterOptions()
opts.SourceSRS = "ENU:114,34"
opts.EnableSimplify = true
opts.SimplifyRatio = 0.3  // 保留 30% 三角形
```

**场景 4：限制 LOD 层级**

```go
opts := tiles3d.DefaultConverterOptions()
opts.SourceSRS = "EPSG:4548"
opts.MaxLOD = 18  // 只加载 LOD 18 及以下
```

---

## OSG 二进制格式详解

### 文件头部结构

```
偏移  大小  字段         说明
0     4    magic_low    0x6C910EA1 (OSG 魔数低 32 位)
4     4    magic_high   0x1AFB4545 (OSG 魔数高 32 位)
8     4    type         1=Scene, 2=Image, 3=Object
12    4    version      文件版本号（如 130）
16    4    attributes   属性位图
                        bit 0 (0x1) = 自定义域
                        bit 1 (0x2) = Schema 数据
                        bit 2 (0x4) = 二进制大括号支持
20    var   compressor  压缩器名称字符串（int32 长度 + 数据）
                          "0" = 无压缩
```

### 版本阈值表

| 版本 | 行为 | C++ 位置 |
|-------|---------|-----------|
| ≥ 112 | Array 作为对象读取（wrapper + classname） | InputStream.cpp:163 |
| > 96  | PrimitiveSet 含 numInstances 字段 | InputStream.cpp:588 |
| > 94  | Image 含 ClassName 属性 | InputStream.cpp:675 |
| > 148 | 二进制大括号块大小用 uint64（8 字节）| BinaryStreamOperator.h:275 |
| ≤ 148 | 二进制大括号块大小用 int32（4 字节）| BinaryStreamOperator.h:283 |
| = 77  | Object/Node 的 UserData 被 UserDataContainer 替换 | Object.cpp, Node.cpp |
| = 70  | PagedLOD 的 FrameNumberOfLastTraversal 被移除 | PagedLOD.cpp |

### 数组格式差异（版本 112 分界）

**版本 < 112（旧格式）**：

```
ArrayID → ArrayType 枚举 → size → BEGIN_BRACKET → 原始分量数据 → END_BRACKET
```

通过 `readArray()` 读取，每种数组类型枚举值（IDVEC3ARRAY=16 等）对应一种数据读取方式。

**版本 ≥ 112（新格式）**：

```
类名 "osg::Vec3Array" → 块大小 → UniqueID → Array 字段序列化 → Vec3Array 向量数据
```

Array 作为完整对象通过 `ReadObject` 读取：
- Binding（枚举，BINDUNDEFINED/BINDOFF/BINDOVERALL/BINDPERVERTEX）
- Normalize（布尔）
- PreserveDataType（布尔）
- 向量数据（IsAVectorSerializer，size × 分量数）

### PrimitiveSet 格式差异（版本 96 分界）

**版本 ≤ 96**：
```
PrimitiveType → Mode → (类型相关数据)
```

**版本 > 96**：
```
PrimitiveType → Mode → numInstances → (类型相关数据)
```

Primitive 类型：
- `ID_DRAW_ARRAYS`：First + Count
- `ID_DRAW_ARRAY_LENGTH`：First + size + 边长列表
- `ID_DRAW_ELEMENTS_UBYTE`：size + uint8 索引列表
- `ID_DRAW_ELEMENTS_USHORT`：size + uint16 索引列表
- `ID_DRAW_ELEMENTS_UINT`：size + uint32 索引列表

≥ 112 版本时，PrimitiveSet 也通过 ReadObject 读取。

### Image 格式差异（版本 94 分界）

**版本 ≤ 94**：直接读取 UniqueID → FileName → WriteHint → 图像数据
**版本 > 94**：先读取 ClassName → 再读取 UniqueID → ...

图像数据模式：
- `IMAGE_INLINE_DATA`：原始像素数据内联
- `IMAGE_INLINE_FILE`：嵌入的文件二进制数据
- `IMAGE_EXTERNAL`：引用外部文件，运行时加载

### 二进制大括号机制（版本 148 分界）

二进制大括号 `{` `}` 用于标记对象块的边界。当 `SupportBinaryBrackets` 标志（attributes bit 2）启用时：

**写入**：
```
{ 的位置写入占位块大小 → 写入数据 → } 时计算大小并回填
```

**读取**：
```
{ 时读取块大小 → 解析对象字段 → } 或 AdvanceToCurrentEndBracket 时跳转到块尾
```

**块大小编码**：
- 版本 > 148：uint64（8 字节）
- 版本 ≤ 148：int32（4 字节）

块大小包含 size 字段本身。例如 int32 块大小 = 100 表示从 size 字段开头到块尾共 100 字节。

---

## 测试

### 运行全部测试

```bash
# 全部 245 个测试
go test -count=1 ./...

# 仅核心 OSG 包
go test -count=1 .

# 仅 3D Tiles 包
go test -count=1 ./tiles3d/...

# 仅 C++ 一致性测试
go test -run "^Test(Wrapper|Threshold|Getter|Default|File|Enums|Factories|Array_New|SerializerTypes)"
```

### 测试分类说明

```bash
# 35 个类型包装器验证（与 C++ 对比）
go test -run "^TestWrapper_"

# 12 个枚举值验证（80+ GL 常量）
go test -run "^TestEnums_"

# 5 个文件集成测试
go test -run "^TestFile_"

# 4 个版本阈值测试
go test -run "^TestThresholds_"

# 3D Tiles 提取器测试
go test -run "^TestExtract" ./tiles3d/...

# 3D Tiles Draco/简化测试
go test -run "^TestCompress|^TestSimplify|^TestOptimize" ./tiles3d/...

# 3D Tiles 生成测试
go test -run "^TestB3DM|^TestTileset|^TestBuildGLTF" ./tiles3d/...
```

### 测试覆盖范围

| 类别 | 数量 | 说明 |
|----------|:----:|------|
| C++ 一致性测试 | 63 | 所有 35 个类型的包装器结构、序列化器版本范围、枚举值 |
| 版本阈值测试 | 4 | 94(Image)、96(PrimitiveSet)、112(Geometry/Array)、148(大括号) |
| 枚举值测试 | 12 | CullFace、AlphaFunc、ShadeModel、TexEnv、TexGen、PrimitiveSet Mode、Transform ReferenceFrame、Material ColorMode、Image Origin、Texture Wrap/Filter |
| 文件集成 | 5 | cessna.osgb、simpleroom.osgt、skydome.osgt、Tile_+003_+003_L18_000.osgb、Tile_+002_+000_L22_000020.osgb |
| 提取器测试 | 20+ | 顶点/法线/纹理坐标/索引提取，全部数据类型覆盖 |
| 转换器遍历 | 10+ | Group/Geode/PagedLod/MatrixTransform 节点树遍历 |
| 索引三角化 | 3 | QUADS、QUAD_STRIP、边界条件 |
| 矩阵变换 | 2 | 单位矩阵、平移矩阵 |
| Draco 压缩 | 5 | 启用/禁用测试、B3DM + GLB 集成 |
| 网格简化 | 6 | 禁用/满比例/空索引、meshopt 集成 |
| B3DM/GLB 生成 | 6 | 常规 + Draco + 写入文件 + nil 保护 |
| tileset.json | 3 | 生成、嵌套层级、写入文件 |
| 坐标转换 | 7 | ECEF、ENU、SRS 检测、元数据解析 |
| 边界框/误差 | 3 | 计算、空输入、几何误差公式 |
| 其他 | 10+ | 包围盒合并、纹理 MIME、文件提取 |

---

## 构建依赖

### replace 指令说明

`go.mod` 使用以下 `replace` 指令将依赖指向本地仓库：

```go
replace github.com/flywave/go-proj     => ../go-proj
replace github.com/flywave/go-draco    => ../go-draco
replace github.com/flywave/go-meshopt  => ../go-meshopt
replace github.com/flywave/gltf        => ../gltf
replace github.com/flywave/go-geoid    => ../go-geoid
replace github.com/flywave/go-3dtile   => ../go-3dtile
```

### 编译 tiles3d 的前提

`tiles3d` 包依赖上述 6 个外部仓库，编译前需确保它们存在于上级目录：

```
../
├── go-proj/         # PROJ 坐标投影转换
├── go-draco/        # Draco 网格压缩
├── go-meshopt/      # 网格简化
├── gltf/            # GLTF 文档构建
├── go-geoid/        # 大地水准面校正
├── go-3dtile/       # B3DM 格式读写
└── go-osg/          # 本仓库
```

如需单独编译核心 OSG 读写包（不包含 tiles3d），`model/` 和根包不依赖上述仓库。

### 环境变量

`tiles3d` 的坐标转换功能需要 PROJ 数据文件：

```bash
export PROJ_DATA=/path/to/proj-data
export PROJ_LIB=/path/to/proj-data
```

---

## 常见问题

### 坐标转换失败（PROJ_DATA 未设置）

```
proj_create: Cannot find proj.db
```

**解决**：设置 PROJ_DATA 环境变量指向 proj-data 目录：
```bash
export PROJ_DATA=/usr/share/proj
```

### 缺少测试数据文件

部分测试引用 `test_data/OSGB/Data/*` 和 `test_data/0131/Data/*` 等不存在的数据文件。这些测试会自动跳过。可用的测试数据：

```
test_data/
├── cessna.osgb               # 飞机模型（含纹理）
├── simpleroom.osgt           # 简化房间（ASCII 格式）
├── skydome.osgt              # 天空穹（ASCII 格式）
├── Tile_+003_+003_L18_000.osgb  # 3D 图块
└── images/skymap.jpg         # 天空纹理
```

### Draco 压缩失败

当网格过于简单（如单个三角形）时，Draco 可能报 "All triangles are degenerate"。这不是代码错误，而是 Draco 库对最小三角形的限制。在生产数据的大网格上不会出现此问题。

### 读取特定版本 OSGB 文件出错

Go 代码已根据 C++ 参考实现了所有版本阈值（94/96/112/148）。如果特定版本文件读取失败，请检查：

1. 文件版本号（二进制头部偏移 8）
2. 版本对应的阈值行为（见版本阈值表）
3. 是否有自定义域（attributes bit 0）

---

## 代码参考（C++ 对照）

C++ 参考代码位于 `OpenSceneGraph-master/` 目录，关键对照文件：

| C++ 文件 | Go 文件 | 对照内容 |
|----------|---------|---------|
| `src/osgDB/InputStream.cpp` | `input_stream.go`, `binary_input.go` | 输入流、ReadArray、ReadObject、版本阈值 |
| `src/osgPlugins/osg/BinaryStreamOperator.h` | `binary_input.go` | 二进制大括号处理、块大小 |
| `src/osgWrappers/serializers/osg/Array.cpp` | `array.go` | Array 包装器注册（Binding/Normalize/PreserveDataType） |
| `src/osgWrappers/serializers/osg/Geometry.cpp` | `geometry.go` | Geometry 版本 112 分界 |
| `src/osgWrappers/serializers/osg/Group.cpp` | `group.go` | Group 子节点读写 |
| `src/osgWrappers/serializers/osg/PagedLOD.cpp` | `paged_lod.go` | PagedLOD 序列化 |
| `src/osgWrappers/serializers/osg/LOD.cpp` | `lod.go` | LOD 层级字段 |
| `src/osgWrappers/serializers/osg/Node.cpp` | `node.go` | Node 通用字段 |
| `src/osgWrappers/serializers/osg/Object.cpp` | `object.go` | Object 基类字段 |
| `src/osgWrappers/serializers/osg/Image.cpp` | `image.go` | Image 序列化 |
| `src/osgWrappers/serializers/osg/Drawable.cpp` | `drawable.go` | Drawable 版本阈值 |
| `src/osgWrappers/serializers/osg/PrimitiveSet.cpp` | `primitive_set.go` | PrimitiveSet 序列化 |
| `src/osgWrappers/serializers/osg/StateSet.cpp` | `state_set.go` | StateSet 模式/属性/纹理列表 |
| `src/osgWrappers/serializers/osg/CullFace.cpp` | `cull_face.go` | CullFace 裁剪面模式 |
| `src/osgWrappers/serializers/osg/AlphaFunc.cpp` | `alpha_func.go` | AlphaFunc 比较函数 |
| `src/osgWrappers/serializers/osg/ShadeModel.cpp` | `shade_model.go` | ShadeModel 着色模式 |
| `src/osgWrappers/serializers/osg/Material.cpp` | `material.go` | Material 材质属性 |
| `src/osgWrappers/serializers/osg/TexEnv.cpp` | `tex_env.go` | TexEnv 纹理环境 |
| `src/osgWrappers/serializers/osg/TexGen.cpp` | `tex_gen.go` | TexGen 纹理坐标生成 |
| `src/osgWrappers/serializers/osg/Texture*.cpp` | `texture*.go` | Texture 系列包装器 |
| `src/osgWrappers/serializers/osg/Transform.cpp` | `transform.go` | Transform 参考系 |
| `src/osgWrappers/serializers/osg/MatrixTransform.cpp` | `matrix_transform.go` | MatrixTransform 矩阵 |
| `src/osgWrappers/serializers/osg/PositionAttitudeTransform.cpp` | `position_attitude_transform.go` | PAT 位置/姿态/缩放 |
| `src/osgWrappers/serializers/osg/BufferData.cpp` | `buffer_data.go` | BufferData 基类 |
