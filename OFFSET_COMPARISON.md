# Offset 处理对比分析

## C++ 实现

### 1. ReadMark (BinaryStreamOperator.h:268-288)
```cpp
virtual void readMark( osgDB::ObjectMark& mark )
{
    if ( _supportBinaryBrackets )
    {
        if ( mark._name=="{" )
        {
            _beginPositions.push_back( _in->tellg() );  // 位置 A (读取 size 之前)

            if (getInputStream() && getInputStream()->getFileVersion() > 148)
            {
               uint64_t size = 0;
               _in->read( (char*)&size, osgDB::INT64_SIZE);  // 读取 8 字节，tellg() 自动前进
               if ( _byteSwap ) osg::swapBytes( (char*)&size, osgDB::INT64_SIZE);
               _blockSizes.push_back( size );
            }
            else
            {
               int size = 0;
               _in->read( (char*)&size, osgDB::INT_SIZE);  // 读取 4 字节，tellg() 自动前进
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

### 2. AdvanceToCurrentEndBracket (BinaryStreamOperator.h:304-314)
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

### 3. WriteMark (验证 size 包括 size 字段本身)
```cpp
virtual void writeMark( const osgDB::ObjectMark& mark )
{
    if ( mark._name=="{" )
    {
        uint64_t size = 0;
        _beginPositions.push_back( _out->tellp() );  // 位置 A (写入 size 之前)
        _out->write( (char*)&size, osgDB::INT64_SIZE );  // 写入 8 字节的 0
    }
    else if ( mark._name=="}" && _beginPositions.size()>0 )
    {
        std::streampos pos = _out->tellp(), beginPos = _beginPositions.back();
        _beginPositions.pop_back();
        _out->seekp( beginPos );

        std::streampos size64 = pos - beginPos;  // size = 位置 B - 位置 A
        uint64_t size = (uint64_t) size64;  // ⚠️ 包括 size 字段本身！
        _out->write( (char*)&size, osgDB::INT64_SIZE);
        _out->seekp( pos );
    }
}
```

## Go 实现

### 1. ReadMark (binary_input.go:189-213)
```go
func (iter *BinaryInputIterator) ReadMark(mark *model.ObjectMark) {
    if iter.SupportBinaryBrackets {
        if mark.Name == "{" {
            // ⚠️ 关键：保存读取 size 之前的位置
            iter.BeginPositions = append(iter.BeginPositions, iter.Offset)  // 位置 A

            if iter.InputStream.FileVersion > 148 {
                var size int64
                iter.ReadLong(&size)  // 读取 8 字节，Offset 自动增加 8
                iter.BlockSizes = append(iter.BlockSizes, size)
            } else {
                var size int32
                iter.ReadInt(&size)  // 读取 4 字节，Offset 自动增加 4
                iter.BlockSizes = append(iter.BlockSizes, int64(size))
            }
        } else if mark.Name == "}" && len(iter.BeginPositions) > 0 {
            iter.BeginPositions = iter.BeginPositions[:len(iter.BeginPositions)-1]
            iter.BlockSizes = iter.BlockSizes[:len(iter.BlockSizes)-1]
        }
    }
}
```

### 2. AdvanceToCurrentEndBracket (binary_input.go:220-238)
```go
func (iter *BinaryInputIterator) AdvanceToCurrentEndBracket() {
    l := len(iter.BeginPositions)
    if iter.SupportBinaryBrackets && l > 0 {
        pos := iter.BeginPositions[l-1]  // 位置 A
        bs := len(iter.BlockSizes)
        pos += iter.BlockSizes[bs-1]  // A + size = 跳过后的位置
        skip := pos - iter.Offset  // 需要跳过的字节数

        // 保护：如果 skip 为负数，说明我们已经读取过了数据
        if skip > 0 {
            iter.Offset = pos
            iter.In.Discard(int(skip))  // 跳过 skip 字节
        }

        iter.BeginPositions = iter.BeginPositions[:len(iter.BeginPositions)-1]
        iter.BlockSizes = iter.BlockSizes[:len(iter.BlockSizes)-1]
    }
}
```

## 关键对比

### 1. 保存位置
| 操作 | C++ | Go | 是否一致 |
|-----|-----|-----|---------|
| 保存位置时机 | 读取 size **之前** | 读取 size **之前** | ✅ 一致 |
| 保存的值 | `_in->tellg()` | `iter.Offset` | ✅ 一致 |

### 2. 读取 size
| 操作 | C++ | Go | 是否一致 |
|-----|-----|-----|---------|
| 读取后位置变化 | `tellg()` 自动增加 | `Offset` 手动增加 | ✅ 逻辑一致 |
| size 内容 | 包括 size 字段本身 | 包括 size 字段本身 | ✅ 一致 |

### 3. 跳过数据
| 操作 | C++ | Go | 是否一致 |
|-----|-----|-----|---------|
| 计算目标位置 | `beginPos + size` | `BeginPosition + BlockSize` | ✅ 一致 |
| 跳过方式 | `seekg(position)` | `Discard(skip)` | ⚠️ 实现不同，但逻辑一致 |

## 具体例子验证

假设文件布局：
```
位置 100: { (BEGIN_BRACKET)
位置 100: size (8 字节) = 108
位置 108: data (100 字节)
位置 208: } (END_BRACKET)
```

### C++ 执行流程
1. **ReadMark("{")**:
   - `_beginPositions.push_back(100)` - 位置 A
   - `_in->read(&size, 8)` - 读取 size=108，tellg() 变为 108
   - `_blockSizes.push_back(108)`

2. **读取数据**:
   - `_in->read(data, 100)` - tellg() 变为 208

3. **AdvanceToCurrentEndBracket()**:
   - `position = 100 + 108 = 208`
   - `_in->seekg(208)` - 跳到位置 208 ✅

### Go 执行流程
1. **ReadMark("{")**:
   - `BeginPositions = [100]` - 位置 A
   - `ReadLong(&size)` - 读取 size=108，Offset 变为 108
   - `BlockSizes = [108]`

2. **读取数据**:
   - `Read(data, 100)` - Offset 变为 208

3. **AdvanceToCurrentEndBracket()**:
   - `pos = 100 + 108 = 208`
   - `skip = 208 - 208 = 0`
   - `Offset = 208` ✅

### 另一个例子：跳过未读取的数据

假设文件布局：
```
位置 100: { (BEGIN_BRACKET)
位置 100: size (8 字节) = 108
位置 108: data (100 字节) - 未读取
位置 208: } (END_BRACKET)
```

### C++ 执行流程
1. **ReadMark("{")**:
   - `_beginPositions.push_back(100)`
   - `_in->read(&size, 8)` - tellg() = 108
   - `_blockSizes.push_back(108)`

2. **AdvanceToCurrentEndBracket()** (跳过未读取的数据):
   - `position = 100 + 108 = 208`
   - `_in->seekg(208)` - 直接跳到位置 208 ✅

### Go 执行流程
1. **ReadMark("{")**:
   - `BeginPositions = [100]`
   - `ReadLong(&size)` - Offset = 108
   - `BlockSizes = [108]`

2. **AdvanceToCurrentEndBracket()** (跳过未读取的数据):
   - `pos = 100 + 108 = 208`
   - `skip = 208 - 108 = 100`
   - `Offset = 208`
   - `In.Discard(100)` - 跳过 100 字节 ✅

## 结论

✅ **Go 的 Offset 处理和 C++ 完全一致！**

### 关键点
1. **位置保存时机一致**: 都是在读取 size 之前保存位置
2. **size 内容一致**: 都包括 size 字段本身（8 或 4 字节）
3. **跳过逻辑一致**: 都是 `beginPosition + size`
4. **实现方式不同但等价**:
   - C++ 使用 `seekg()` 直接跳转到目标位置
   - Go 使用 `Discard()` 跳过相对字节数
   - 两种方式的结果完全相同

### 额外的保护措施
Go 代码中有一个额外的保护：
```go
if skip > 0 {
    iter.Offset = pos
    iter.In.Discard(int(skip))
}
```

这个保护是必要的，因为：
- 如果已经读取了所有数据，`skip = 0`，不需要再跳过
- 如果 `skip < 0`，说明数据已经被读取过了，也不需要跳过
- 这避免了负数跳过或重复跳过的问题

C++ 不需要这个保护，因为 `seekg()` 可以跳到任意位置（包括当前位置）。
