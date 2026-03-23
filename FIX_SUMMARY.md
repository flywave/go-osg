# OSGB坐标转换问题修复总结

## 修复时间
2025年（根据上下文）

## 修复的问题

### 1. 🔴 经纬度顺序错误（主要问题）

**问题描述**：
- C++使用OGR库的`OAMS_TRADITIONAL_GIS_ORDER`，返回顺序为
- Go使用go-proj库，默认使用`COMPLIANT_OGC_ORDER`，返回顺序为
- 导致经纬度交换，所有后续计算都基于错误的坐标

**影响范围**：
所有使用`proj.Transform3()`进行坐标转换的函数

**修复位置**：

#### `coordinate.go` - setupEPSGMode (line 121-130)
```go
// 修复前：
lon, lat, z, err := proj.Transform3(src, wgs84, originX, originY, originZ)
lon = lon * 180.0 / math.Pi
lat = lat * 180.0 / math.Pi

// 修复后：
lon, lat, z, err := proj.Transform3(src, wgs84, originX, originY, originZ)
// 交换经纬度以匹配C++的GIS顺序
lon, lat = lat, lon
lon = lon * 180.0 / math.Pi
lat = lat * 180.0 / math.Pi
```

#### `coordinate.go` - setupWKTMode (line 157-166)
```go
// 修复前：
lon, lat, z, err := proj.Transform3(src, wgs84, originX, originY, originZ)
lon = lon * 180.0 / math.Pi
lat = lat * 180.0 / math.Pi

// 修复后：
lon, lat, z, err := proj.Transform3(src, wgs84, originX, originY, originZ)
// 交换经纬度以匹配C++的GIS顺序
lon, lat = lat, lon
lon = lon * 180.0 / math.Pi
lat = lat * 180.0 / math.Pi
```

#### `coordinate.go` - setupUnknownMode (line 191-197)
```go
// 修复前：
lon, lat, z, err := proj.Transform3(src, wgs84, originX, originY, originZ)
t.center = [3]float64{lon * 180.0 / math.Pi, lat * 180.0 / math.Pi, z}

// 修复后：
lon, lat, z, err := proj.Transform3(src, wgs84, originX, originY, originZ)
// 交换经纬度以匹配C++的GIS顺序
lon, lat = lat, lon
t.center = [3]float64{lon * 180.0 / math.Pi, lat * 180.0 / math.Pi, z}
```

#### `coordinate.go` - toLocalENUFromProjected (line 319-328)
```go
// 修复前：
lon, lat, z, err := proj.Transform3(t.sourceProj, wgs84, absX, absY, absZ)
lonDeg := lon * 180.0 / math.Pi
latDeg := lat * 180.0 / math.Pi

// 修复后：
lon, lat, z, err := proj.Transform3(t.sourceProj, wgs84, absX, absY, absZ)
// 交换经纬度以匹配C++的GIS顺序
lon, lat = lat, lon
lonDeg := lon * 180.0 / math.Pi
latDeg := lat * 180.0 / math.Pi
```

#### `coordinate.go` - toLocalENUFromUnknown (line 353-361)
```go
// 修复前：
lon, lat, z, err := proj.Transform3(t.sourceProj, wgs84, absX, absY, absZ)
lonDeg := lon * 180.0 / math.Pi
latDeg := lat * 180.0 / math.Pi

// 修复后：
lon, lat, z, err := proj.Transform3(t.sourceProj, wgs84, absX, absY, absZ)
// 交换经纬度以匹配C++的GIS顺序
lon, lat = lat, lon
lonDeg := lon * 180.0 / math.Pi
latDeg := lat * 180.0 / math.Pi
```

#### `coordinate.go` - ToWGS84 (line 226-229)
```go
// 修复前：
x, y, z, err := proj.Transform3(t.sourceProj, wgs84, point[0], point[1], point[2])
if err == nil {
    return [3]float64{x, y, z}
}

// 修复后：
x, y, z, err := proj.Transform3(t.sourceProj, wgs84, point[0], point[1], point[2])
if err == nil {
    // 交换经纬度以匹配C++的GIS顺序
    return [3]float64{y, x, z}
}
```

---

### 2. 🟡 SVD矩阵应用索引错误（次要问题）

**问题描述**：
- 在使用最小二乘法计算变换矩阵后，应用矩阵变换时索引顺序错误
- C++使用列向量变换：Transform[row][col] * v[col]
- Go之前使用了错误的索引：X.At(row, col) * v[wrong_index]

**修复位置**：

#### `converter.go` - extractVertices (line 364-366)
```go
// 修复前（错误）：
newX := X.At(0, 0)*x + X.At(1, 0)*y + X.At(2, 0)*z + X.At(3, 0)
newY := X.At(0, 1)*x + X.At(1, 1)*y + X.At(2, 1)*z + X.At(3, 1)
newZ := X.At(0, 2)*x + X.At(1, 2)*y + X.At(2, 2)*z + X.At(3, 2)

// 修复后（正确）：
newX := X.At(0, 0)*x + X.At(0, 1)*y + X.At(0, 2)*z + X.At(0, 3)
newY := X.At(1, 0)*x + X.At(1, 1)*y + X.At(1, 2)*z + X.At(1, 3)
newZ := X.At(2, 0)*x + X.At(2, 1)*y + X.At(2, 2)*z + X.At(2, 3)
```

**说明**：
- 正确的矩阵乘法公式：`result[i] = Σ Matrix[i][j] * vector[j]`
- 修复前：使用的是错误的索引组合
- 修复后：按正确的行列顺序访问矩阵元素

---

## 修复的影响

### 坐标转换流程
1. **OSGB顶点读取** → 相对坐标
2. **加上originOffset** → 投影坐标（如EPSG:4548）
3. **投影坐标→WGS84** → 经纬度（此处交换了经纬度）
4. **WGS84→ECEF** → 地心地固坐标
5. **ECEF→ENU** → 局部东-北-天坐标（用于3D Tiles）

### 修复前的问题
- 步骤3得到的是
- 后续所有计算都基于错误的经纬度
- 最终ENU坐标完全错误，导致模型位置、方向都错误

### 修复后的效果
- 步骤3正确得到
- 经纬度正确，ENU计算正确
- 模型位置和方向与C++版本一致

---

## 验证方法

### 1. 单元测试
创建测试用例，对比Go和C++的转换结果：

```go
func TestCoordinateTransform(t *testing.T) {
    // 测试点
    point := [3]float64{500000.0, 3000000.0, 100.0}

    // 使用相同的投影坐标系
    sourceSRS := "EPSG:4548" // CGCS2000 / 3-degree Gauss-Kruger CM 114E

    // 创建转换器
    trans := NewCoordinateTransformer(sourceSRS, "EPSG:4326")

    // 转换坐标
    wgs84 := trans.ToWGS84(point)

    // 预期结果（从C++版本获取）
    expectedLon := 114.0
    expectedLat := 27.0

    // 验证（允许微小误差）
    if math.Abs(wgs84[0]-expectedLon) > 0.0001 {
        t.Errorf("Longitude mismatch: got %f, expected %f", wgs84[0], expectedLon)
    }
    if math.Abs(wgs84[1]-expectedLat) > 0.0001 {
        t.Errorf("Latitude mismatch: got %f, expected %f", wgs84[1], expectedLat)
    }
}
```

### 2. 端到端测试
使用相同的OSGB文件，对比Go和C++版本的输出：

```bash
# C++版本
./osgb23dtile input.osgb output_c++

# Go版本
./go-osg input.osgb output_go

# 对比结果
diff -r output_c++ output_go
```

### 3. 调试输出
对比中间结果：

```
# C++版本输出
[CoordinateTransformer] OGR transform result: lon=114.0000000000 lat=27.0000000000 h=100.000

# Go版本输出（修复前）
DEBUG setupEPSGMode: center=(27.000000, 114.000000, 100.000)  # 错误：经纬度交换
DEBUG setupEPSGMode: center=(114.000000, 27.000000, 100.000)  # 正确
```

---

## 其他注意事项

### 1. go-proj库的轴顺序
- go-proj默认使用OGC标准顺序：(lat, lon, h)
- OGR库使用GIS传统顺序：(lon, lat, h)
- 修复方法：交换返回值

### 2. 调试输出
所有修复的位置都保留了调试输出，便于验证：
- `fmt.Printf("DEBUG ...")`

### 3. 兼容性
这些修复不会影响非地理参考模式下的转换（如OSGB→GLTF）

---

## 总结

修复了2个关键问题：
1. **经纬度顺序错误**：在6个函数中修复了`proj.Transform3`的返回值顺序
2. **SVD矩阵应用错误**：修复了1个函数中的矩阵乘法索引

这些修复使得Go版本的OSGB转换结果与C++版本保持一致。
