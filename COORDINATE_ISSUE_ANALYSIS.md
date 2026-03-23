# OSGB转换坐标问题分析

## 问题概述
Go语言版本的OSGB转换过程中，读取的坐标存在极大问题，而C++版本可以正确转换。

## 核心代码路径

### C++版本（正确）
- 主转换逻辑：`/Users/xuning/Code/3dtiles/src/osgb23dtile.cpp:154-246`
- 坐标转换器：`/Users/xuning/Code/3dtiles/src/coordinate_transformer.cpp`

### Go版本（有问题）
- 主转换逻辑：`/Users/xuning/Work/go-osg/tiles3d/converter.go:199-379`
- 坐标转换器：`/Users/xuning/Work/go-osg/tiles3d/coordinate.go`

## 关键差异分析

### 1. 坐标提取顺序

#### C++版本
```cpp
for (int VertexIndex = 0; VertexIndex < vertexArr->size(); VertexIndex++)
{
    osg::Vec3d Vertex = vertexArr->at(VertexIndex);
    glm::dvec3 vertex = glm::dvec3(Vertex.x(), Vertex.y(), Vertex.z());
    Min = glm::min(vertex, Min);
    Max = glm::max(vertex, Max);
}
```
- 直接读取osg::Vec3Array的x, y, z

#### Go版本
```go
vertices := c.extractVertices(geom.VertexArray)
// ... 数据类型转换 ...
x := float64(vertices[i])
y := float64(vertices[i+1])
z := float64(vertices[i+2])
```
- 通过extractVertices提取，数据类型转换
- **可能的问题**：在extractVertices中，数据被提取为[]float32，然后再转换

### 2. 投影坐标转换

#### C++版本
```cpp
// 使用OGR库，设置轴映射策略
outRs.SetAxisMappingStrategy(OAMS_TRADITIONAL_GIS_ORDER);
inRs.SetAxisMappingStrategy(OAMS_TRADITIONAL_GIS_ORDER);
// 转换结果：lon, lat, h
ogr_transform_->Transform(1, &result.x, &result.y, &result.z);
// Geoid校正：注意参数顺序
result.z = ApplyGeoidCorrection(result.y, result.x, result.z); // lat, lon, h
```

#### Go版本
```go
// 使用go-proj库，未设置轴映射
lon, lat, z, err := proj.Transform3(t.sourceProj, wgs84, absX, absY, absZ)
// Geoid校正：注意参数顺序
z = t.geoidConv.ConvertOrthometricToEllipsoidal(latDeg, lonDeg, z) // lat, lon, h
```

**潜在问题**：
- go-proj库可能使用不同的轴顺序（COMPLIANT_OGC_ORDER vs TRADITIONAL_GIS_ORDER）
- OGR库明确设置了轴映射策略为TRADITIONAL_GIS_ORDER（lon, lat）
- 如果go-proj使用默认的OGC顺序（lat, lon），会导致lon和lat交换

### 3. ENU计算公式

#### C++版本
```cpp
double latRad = lat * math.Pi / 180.0;
double lonRad = lon * math.Pi / 180.0;

sinLat = sin(latRad);
cosLat = cos(latRad);
sinLon = sin(lonRad);
cosLon = cos(lonRad);

originECEF = t.ToECEFFromLatLon(latRad, lonRad, height);

dx = ecefPoint[0] - originECEF[0];
dy = ecefPoint[1] - originECEF[1];
dz = ecefPoint[2] - originECEF[2];

xEnu = -sinLon*dx + cosLon*dy;
yEnu = -sinLat*cosLon*dx - sinLat*sinLon*dy + cosLat*dz;
zEnu = cosLat*cosLon*dx + cosLat*sinLon*dy + sinLat*dz;
```
- 使用弧度进行计算

#### Go版本
```go
latRad := lat * math.Pi / 180.0
lonRad := lon * math.Pi / 180.0

sinLat := math.Sin(latRad)
cosLat := math.Cos(latRad)
sinLon := math.Sin(lonRad)
cosLon := math.Cos(lonRad)

originECEF := t.ToECEFFromLatLon(latRad, lonRad, height)

dx := ecefPoint[0] - originECEF[0]
dy := ecefPoint[1] - originECEF[1]
dz := ecefPoint[2] - originECEF[2]

xEnu := -sinLon*dx + cosLon*dy
yEnu := -sinLat*cosLon*dx - sinLat*sinLon*dy + cosLat*dz
zEnu := cosLat*cosLon*dx + cosLat*sinLon*dy + sinLat*dz
```
- 同样使用弧度，公式看起来一致

### 4. 最小二乘法求解

#### C++版本
```cpp
Eigen::MatrixXd A, B;
A.resize(8, 4);
B.resize(8, 4);

// A矩阵：原始点（列主序，但Eigen按行存储）
for (int row = 0; row < 8; row++)
{
    A.row(row) << OriginalPoints[row].x, OriginalPoints[row].y, OriginalPoints[row].z, 1;
}
// B矩阵：校正点
for (int row = 0; row < 8; row++)
{
    B.row(row) << CorrectedPoints[row].x, CorrectedPoints[row].y, CorrectedPoints[row].z, 1;
}

Eigen::BDCSVD<Eigen::MatrixXd> SVD(A, Eigen::ComputeThinU | Eigen::ComputeThinV);
Eigen::MatrixXd X = SVD.solve(B);

// 应用变换
glm::dvec4 v = Transform * glm::dvec4(Vertex.x(), Vertex.y(), Vertex.z(), 1);
```
- Eigen使用列主序存储
- X[i,j]表示第i行第j列

#### Go版本
```go
A := mat.NewDense(8, 4, nil)
B := mat.NewDense(8, 4, nil)

for i := 0; i < 8; i++ {
    A.Set(i, 0, originalCorners[i][0])
    A.Set(i, 1, originalCorners[i][1])
    A.Set(i, 2, originalCorners[i][2])
    A.Set(i, 3, 1.0)
    B.Set(i, 0, correctedCorners[i][0])
    B.Set(i, 1, correctedCorners[i][1])
    B.Set(i, 2, correctedCorners[i][2])
    B.Set(i, 3, 1.0)
}

var svd mat.SVD
svd.Factorize(A, mat.SVDFull)
var X mat.Dense
svd.SolveTo(&X, B, 0)

// 应用变换 - 注意索引顺序
newX := X.At(0, 0)*x + X.At(1, 0)*y + X.At(2, 0)*z + X.At(3, 0)
newY := X.At(0, 1)*x + X.At(1, 1)*y + X.At(2, 1)*z + X.At(3, 1)
newZ := X.At(0, 2)*x + X.At(1, 2)*y + X.At(2, 2)*z + X.At(3, 2)
```
- gonum使用列主序存储
- X.At(i, j)获取第i行第j列

**可能的问题**：
- 在应用变换时，Go版本的索引顺序与C++不同
- C++：X(row, col) * point(col)的转置方式
- Go：X(row, col) * point(row)的方式，可能导致结果错误

## 最可能的问题根源

### 主要问题：投影坐标系轴顺序不一致

**C++版本**：
- 使用OGR库，明确设置`SetAxisMappingStrategy(OAMS_TRADITIONAL_GIS_ORDER)`
- 结果顺序：(lon, lat, h)
- TRADITIONAL_GIS_ORDER：x=lon, y=lat

**Go版本**：
- 使用go-proj库，未设置轴映射策略
- go-proj的默认行为可能是COMPLIANT_OGC_ORDER
- 结果顺序：(lat, lon, h)
- COMPLIANT_OGC_ORDER：x=lat, y=lon

**导致的结果**：
- C++将(lon, lat, h)传入WGS84→ECEF
- Go将(lat, lon, h)传入WGS84→ECEF（如果go-proj使用OGC顺序）
- 经纬度交换导致ENU计算时使用错误的经纬度
- 最终ENU坐标完全错误

### 次要问题：SVD矩阵应用方式不同

**C++版本**：
```cpp
glm::dvec4 v = Transform * glm::dvec4(Vertex.x(), Vertex.y(), Vertex.z(), 1);
```
- 列向量：[x, y, z, 1]^T
- Transform[row][col]
- 结果[row] = Σ Transform[row][col] * v[col]

**Go版本**：
```go
newX := X.At(0, 0)*x + X.At(1, 0)*y + X.At(2, 0)*z + X.At(3, 0)
```
- 如果X是转置的，这个计算可能是错误的
- 正确的应该是：newX := X.At(0, 0)*x + X.At(0, 1)*y + X.At(0, 2)*z + X.At(0, 3)

## 建议的修复方案

### 1. 修复投影坐标系轴顺序

在Go版本中，明确设置go-proj的轴映射策略：

```go
// 检查go-proj是否支持设置轴顺序
// 如果不支持，需要在转换后交换lon和lat
lon, lat, z, err := proj.Transform3(t.sourceProj, wgs84, absX, absY, absZ)
// 如果go-proj返回(lat, lon)，需要交换
temp := lon
lon = lat
lat = temp
```

### 2. 验证并修复SVD矩阵应用

验证Go版本的SVD求解结果，确保矩阵乘法的索引顺序正确：

```go
// 应该使用：
newX := X.At(0, 0)*x + X.At(0, 1)*y + X.At(0, 2)*z + X.At(0, 3)
newY := X.At(1, 0)*x + X.At(1, 1)*y + X.At(1, 2)*z + X.At(1, 3)
newZ := X.At(2, 0)*x + X.At(2, 1)*y + X.At(2, 2)*z + X.At(2, 3)
```

### 3. 添加调试输出

在关键步骤添加调试输出，对比C++和Go的中间结果：

```go
fmt.Printf("DEBUG: Projected point: (%f, %f, %f)\n", absX, absY, absZ)
fmt.Printf("DEBUG: WGS84 (rad): lon=%f, lat=%f, h=%f\n", lon, lat, z)
fmt.Printf("DEBUG: WGS84 (deg): lon=%f, lat=%f, h=%f\n", lonDeg, latDeg, z)
fmt.Printf("DEBUG: ECEF: (%f, %f, %f)\n", ecef[0], ecef[1], ecef[2])
fmt.Printf("DEBUG: ENU: (%f, %f, %f)\n", enu[0], enu[1], enu[2])
```

## 验证方法

1. 使用相同的OSGB文件，分别用C++和Go版本转换
2. 对比以下中间结果：
   - 投影坐标
   - WGS84坐标（经度和纬度）
   - ECEF坐标
   - ENU坐标
   - 最终顶点坐标
3. 找出第一个出现差异的步骤

## 需要进一步调查的点

1. go-proj库的默认轴顺序是什么？
2. go-proj是否支持设置轴映射策略？
3. gonum的SVD求解是否与Eigen的行为一致？
4. OSG的顶点数组在Go中的解析是否正确？
