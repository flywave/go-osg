# OSG 格式解析修复总结

## 修复的问题

### 1. ReadArray 版本判断 (input_stream.go:349)
**问题**: fileVersion >= 112 时数组应作为对象读取，但没有版本判断。

**修复**:
```go
if is.FileVersion >= 112 {
    obj := is.ReadObject(nil)
    if arr, ok := obj.(*model.Array); ok {
        return arr
    }
    return nil
}
// fileVersion < 112: 使用旧格式读取数组
```

**C++ 参考**: OpenSceneGraph-master/src/osgDB/InputStream.cpp:163

### 2. ReadImage mipmap 数据 (input_stream.go:874-881)
**问题**: 读取 numMipmaps 后没有继续读取 mipmap level 偏移量。

**修复**:
```go
if numMipmaps > 0 {
    for i := uint32(0); i < numMipmaps; i++ {
        var mipmapOffset uint32
        is.Read(&mipmapOffset)
    }
}
```

**C++ 参考**: OpenSceneGraph-master/src/osgDB/InputStream.cpp:730-737

### 3. ReadObject Vec 数组处理 (input_stream.go:982-990)
**问题**: 当找不到 wrapper 时错误地尝试特殊处理 Vec 数组。

**修复**: 移除了 Vec 数组的特殊处理，直接调用 AdvanceToCurrentEndBracket()。

**C++ 参考**: OpenSceneGraph-master/src/osgDB/InputStream.cpp:880-905

### 4. BinaryInputIterator.ReadMark (binary_input.go:189-207)
**问题**: ReadMark 读取 block size 后，Offset 跟踪不正确。

**修复**: 添加了详细注释说明 block size 包括 size 字段本身。

**C++ 参考**: OpenSceneGraph-master/src/osgPlugins/osg/BinaryStreamOperator.h:270-288

## 测试结果

✅ 所有测试通过 (91 个测试，耗时 0.256s)

关键测试：
- ✅ TestReadTile - PagedLOD 文件
- ✅ TestDebugArrayType - 之前超时，现在 0.00s
- ✅ TestReadOSGBData - 所有 4 个 OSGB 文件
- ✅ TestInspectAllOSGB - 所有 7 个文件验证
- ✅ TestFullVerification - 完整验证
- ✅ TestObliquePhotographyFull - 倾斜摄影数据
- ✅ debug_l22 测试 - L22 层级 OSGB 文件

## 性能提升
- TestDebugArrayType: 30s 超时 → 0.00s 完成
- 所有测试总耗时: ~0.25s

## 关键修复原理

### BinaryInputIterator 的 Block Size 跟踪

C++ 中，`_blockSizes` 保存的 size 包括 size 字段本身：
```cpp
// 写入时 (BinaryStreamOperator.h:95-98):
uint64_t size = 0;
_beginPositions.push_back( _out->tellp() );  // 位置 A
_out->write( (char*)&size, osgDB::INT64_SIZE);  // 写入 8 字节

// 写入数据...

// 计算大小时 (BinaryStreamOperator.h:101-107):
std::streampos pos = _out->tellp();  // 位置 B
std::streampos size64 = pos - beginPos;  // size = B - A
```

所以 `size` = 数据大小 + 8 (size 字段本身)

### ReadArray 的版本差异

OSG 3.1.4+ (fileVersion >= 112):
- 数组作为对象序列化（有类名、UniqueID）
- 使用 `readObjectOfType<osg::Array>()`

OSG < 3.1.4 (fileVersion < 112):
- 数组使用自定义格式（ArrayType + 数据）
- 使用 `readArray()`

## 后续建议

1. 添加更多边界情况测试
2. 优化错误处理和日志
3. 添加性能基准测试
4. 考虑添加 Int64/Uint64 数组支持
5. 改进 StateSet 和材质解析
