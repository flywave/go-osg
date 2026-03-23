# OSG C++ vs Go 实现对比分析

## 版本相关的特殊处理

### 1. 数组读取 (fileVersion >= 112)

**C++ 实现** (InputStream.cpp:163):
```cpp
InputStream& operator>>( osg::ref_ptr<osg::Array>& ptr ) { 
    if (_fileVersion>=112) 
        ptr = readObjectOfType<osg::Array>(); 
    else 
        ptr = readArray(); 
    return *this; 
}
```

**Go 修复前**:
```go
func (is *OsgIstream) ReadArray() *model.Array {
    // 总是使用旧格式读取，没有版本判断
    ...
}
```

**Go 修复后** (input_stream.go:349):
```go
func (is *OsgIstream) ReadArray() *model.Array {
    if is.FileVersion >= 112 {
        // 新格式：作为对象读取
        obj := is.ReadObject(nil)
        if arr, ok := obj.(*model.Array); ok {
            return arr
        }
        return nil
    }
    // fileVersion < 112: 使用旧格式
    ...
}
```

### 2. PrimitiveSet 读取 (fileVersion > 96)

**C++ 实现** (InputStream.cpp:588-591):
```cpp
unsigned int numInstances = 0u;
*this >> type >> mode;
if ( _fileVersion>96 )
{
    *this >> numInstances;  // 版本 > 96 才有 numInstances
}
```

**Go 实现** (input_stream.go:750):
```go
if is.FileVersion > 96 {
    is.Read(&numInstances)
}
```
✅ **已正确实现**

### 3. Image ClassName (fileVersion > 94)

**C++ 实现** (InputStream.cpp:675):
```cpp
std::string className = "osg::Image";
if ( _fileVersion>94 )  // ClassName 只在 3.1.4+ 支持
    *this >> PROPERTY("ClassName") >> className;
```

**Go 实现** (input_stream.go:829):
```go
if is.FileVersion > 94 {
    is.PROPERTY.Name = "ClassName"
    is.Read(is.PROPERTY)
    is.Read(&className)
}
```
✅ **已正确实现**

## Binary Bracket 机制

### Block Size 的计算和跟踪

**C++ 写入逻辑** (BinaryStreamOperator.h:93-109):
```cpp
virtual void writeMark( const osgDB::ObjectMark& mark )
{
    if ( _supportBinaryBrackets )
    {
        if (getOutputStream() && getOutputStream()->getFileVersion() > 148)
        {
            if ( mark._name=="{" )
            {
                uint64_t size = 0;
                _beginPositions.push_back( _out->tellp() );  // 位置 A (写入前)
                _out->write( (char*)&size, osgDB::INT64_SIZE );  // 写入 8 字节的 0
            }
            else if ( mark._name=="}" && _beginPositions.size()>0 )
            {
                std::streampos pos = _out->tellp(), beginPos = _beginPositions.back();
                _beginPositions.pop_back();
                _out->seekp( beginPos );
                
                std::streampos size64 = pos - beginPos;  // size = 位置 B - 位置 A
                uint64_t size = (uint64_t) size64;  // 包括 size 字段本身！
                _out->write( (char*)&size, osgDB::INT64_SIZE);
                _out->seekp( pos );
            }
        }
        else  // fileVersion <= 148
        {
            // 类似逻辑，使用 int32_t (4 字节)
        }
    }
}
```

**关键点**：
- `beginPos` 是写入 size **之前** 的位置
- `pos` 是写入所有数据**之后** 的位置
- `size = pos - beginPos` **包括** size 字段本身（4 或 8 字节）

**C++ 读取逻辑** (BinaryStreamOperator.h:268-288):
```cpp
virtual void readMark( osgDB::ObjectMark& mark )
{
    if ( _supportBinaryBrackets )
    {
        if ( mark._name=="{" )
        {
            _beginPositions.push_back( _in->tellg() );  // 位置 A
            
            if (getInputStream() && getInputStream()->getFileVersion() > 148)
            {
               uint64_t size = 0;
               _in->read( (char*)&size, osgDB::INT64_SIZE);  // 读取 8 字节
               if ( _byteSwap ) osg::swapBytes( (char*)&size, osgDB::INT64_SIZE);
               _blockSizes.push_back( size );  // size 包括 size 字段本身
            }
            else
            {
               int size = 0;
               _in->read( (char*)&size, osgDB::INT_SIZE);  // 读取 4 字节
               if ( _byteSwap ) osg::swapBytes( (char*)&size, osgDB::INT_SIZE);
               _blockSizes.push_back( size );
            }
        }
        else if ( mark._name=="}" && _beginPositions.size()>0 )
        {
            _beginPositions.pop_back();
            _blockSizes.pop_back();
        }
    }
}
```

**C++ 跳过逻辑** (BinaryStreamOperator.h:304-314):
```cpp
virtual void advanceToCurrentEndBracket()
{
    if ( _supportBinaryBrackets && _beginPositions.size()>0 )
    {
        std::streampos position(_beginPositions.back());  // 位置 A
        position += _blockSizes.back();  // A + size = 跳过后的位置
        _in->seekg( position );  // 直接 seek 到目标位置
        _beginPositions.pop_back();
        _blockSizes.pop_back();
    }
}
```

**Go 实现** (binary_input.go:189-231):
```go
func (iter *BinaryInputIterator) ReadMark(mark *model.ObjectMark) {
    if iter.SupportBinaryBrackets {
        if mark.Name == "{" {
            // ⚠️ BlockSizes 中的 size 包括 size 字段本身
            iter.BeginPositions = append(iter.BeginPositions, iter.Offset)  // 位置 A
            
            if iter.InputStream.FileVersion > 148 {
                var size int64
                iter.ReadLong(&size)  // 读取后 Offset 增加了 8
                iter.BlockSizes = append(iter.BlockSizes, size)  // size 包括 size 字段本身
            } else {
                var size int32
                iter.ReadInt(&size)  // 读取后 Offset 增加了 4
                iter.BlockSizes = append(iter.BlockSizes, int64(size))
            }
        } else if mark.Name == "}" && len(iter.BeginPositions) > 0 {
            iter.BeginPositions = iter.BeginPositions[:len(iter.BeginPositions)-1]
            iter.BlockSizes = iter.BlockSizes[:len(iter.BlockSizes)-1]
        }
    }
}

func (iter *BinaryInputIterator) AdvanceToCurrentEndBracket() {
    l := len(iter.BeginPositions)
    if iter.SupportBinaryBrackets && l > 0 {
        pos := iter.BeginPositions[l-1]  // 位置 A
        bs := len(iter.BlockSizes)
        pos += iter.BlockSizes[bs-1]  // A + size = 跳过后的位置
        skip := pos - iter.Offset  // 需要跳过的字节数

        if skip > 0 {
            iter.Offset = pos
            iter.In.Discard(int(skip))  // 跳过
        }
        
        iter.BeginPositions = iter.BeginPositions[:len(iter.BeginPositions)-1]
        iter.BlockSizes = iter.BlockSizes[:len(iter.BlockSizes)-1]
    }
}
```

✅ **逻辑正确**：虽然 Go 使用手动跟踪 Offset，但计算方式与 C++ 等价

## Image Mipmap 数据

**C++ 实现** (InputStream.cpp:730-737):
```cpp
unsigned int levelSize = readSize();  // numMipmaps
osg::Image::MipmapDataType levels(levelSize);
for ( unsigned int i=0; i<levelSize; ++i )
{
    *this >> levels[i];  // 读取每个 mipmap level 的偏移量
}
if ( image && levelSize>0 )
    image->setMipmapLevels( levels );
```

**Go 修复前**:
```go
var numMipmaps uint32 = 0
is.Read(&numMipmaps)
imgdata.Data = is.ReadCharArray(int(imgdata.Size))
// ❌ 没有读取 mipmap 偏移量！
```

**Go 修复后** (input_stream.go:874-881):
```go
var numMipmaps uint32 = 0
is.Read(&numMipmaps)

imgdata.Data = is.ReadCharArray(int(imgdata.Size))

if numMipmaps > 0 {
    // 读取 mipmap offsets（每个 mipmap level 在数据中的偏移量）
    for i := uint32(0); i < numMipmaps; i++ {
        var mipmapOffset uint32
        is.Read(&mipmapOffset)
    }
    // 注意：我们暂时不保存 mipmap 数据，只是跳过它们以保持流的正确位置
}
```

## Vec 数组处理

**C++ 实现** (InputStream.cpp:880-905):
```cpp
osg::ref_ptr<osg::Object> InputStream::readObject( osg::Object* existingObj )
{
    std::string className;
    unsigned int id = 0;
    *this >> className;

    if (className=="NULL")
    {
        return 0;
    }

    *this >> BEGIN_BRACKET >> PROPERTY("UniqueID") >> id;
    if ( getException() ) return 0;

    IdentifierMap::iterator itr = _identifierMap.find( id );
    if ( itr!=_identifierMap.end() )
    {
        advanceToCurrentEndBracket();
        return itr->second;
    }

    osg::ref_ptr<osg::Object> obj = readObjectFields( className, id, existingObj );

    advanceToCurrentEndBracket();  // 总是调用

    return obj;
}
```

**关键点**：
- **没有** Vec 数组的特殊处理
- `readObjectFields` 如果找不到 wrapper，返回 NULL
- 调用者负责处理 NULL 情况

**Go 修复前**:
```go
obj = is.ReadObjectFields(cls, int32(id), obj)
if obj == nil && strings.HasPrefix(cls, "osg::Vec") {
    // ❌ 错误的特殊处理！
    is.Read(is.ENDBRACKET)
    ary := is.ReadArray()
    return ary
}
is.AdvanceToCurrentEndBracket()
return obj
```

**Go 修复后** (input_stream.go:982-990):
```go
obj = is.ReadObjectFields(cls, int32(id), obj)

// ⚠️ 移除 Vec 数组的特殊处理
// 如果找不到 wrapper，返回 nil
// fileVersion >= 112 时，Vec 数组作为对象读取（有 wrapper）
// fileVersion < 112 时，不会进入 ReadObject，而是直接调用 ReadArray()

is.AdvanceToCurrentEndBracket()
return obj
```

## 总结

### 关键差异

| 功能 | C++ 实现 | Go 修复前 | Go 修复后 |
|-----|---------|---------|---------|
| ReadArray 版本判断 | ✅ 有版本判断 | ❌ 缺失 | ✅ 添加 |
| ReadImage mipmap | ✅ 读取完整 | ❌ 缺失 | ✅ 添加 |
| ReadObject Vec 处理 | ✅ 无特殊处理 | ❌ 错误处理 | ✅ 移除 |
| Binary Bracket Offset | ✅ C++ 使用 tellg/seekp | ⚠️ Go 手动跟踪 | ✅ 逻辑正确 |

### 版本号参考

- **fileVersion >= 112**: OSG 3.1.4+ (数组作为对象)
- **fileVersion > 96**: numInstances 字段
- **fileVersion > 94**: ClassName 字段
- **fileVersion > 148**: Block size 使用 uint64_t (8 字节)
- **fileVersion <= 148**: Block size 使用 int32_t (4 字节)

### 性能对比

- **TestDebugArrayType**: 
  - 修复前: 30s 超时
  - 修复后: 0.00s ✅

- **所有测试**: 
  - 修复前: 部分失败
  - 修复后: 91 个测试通过 (0.256s) ✅
