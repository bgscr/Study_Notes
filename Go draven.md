# Go 语言数组与切片对比学习笔记

## 一、数组
### 1.1 底层内存模型
- **连续内存块**：元素地址计算公式 `addr = base + index * size`
- **编译期确定**：长度和类型在编译时静态解析，无法动态修改
- **栈/堆分配**：小数组可能分配在栈空间，大数组可能逃逸到堆

```go
// 栈分配的数组案例
func stackArray() {
    small := [4]int{1,2,3,4}  // 可能保持在栈
}

// 堆逃逸案例
func heapArray() *[1e6]int {
    var big [1e6]int          // 可能逃逸到堆
    return &big
}
```

## 二、切片
 底层引用数组的连续片段（指针 + 长度 + 容量）
    
```go
type SliceHeader struct {
    Data uintptr  // 指向底层数组
    Len  int      // 可访问元素个数
    Cap  int      // 最大可扩展容量
}
```
引用语义：赋值传递SliceHeader结构体（24字节），共享底层数组

```go
// 1. 基于数组创建
arr := [5]int{1,2,3,4,5}
s1 := arr[1:3]     // Len=2, Cap=4（到数组末尾）

// 2. make函数创建
s2 := make([]int, 3, 5)  // 类型，长度，容量

// 3. 字面量初始化
s3 := []int{1,2,3}       // 自动创建底层数组


s := []int{0,1,2,3,4}
a := s[1:3]        // [1,2]   Len=2 Cap=4
b := s[2:3:4]      // [2]     Len=1 Cap=2（三参数切片限制容量）



s := make([]int, 2, 3)
s = append(s, 1)     // 无需扩容
s = append(s, 2)     // 触发扩容：
                     // 新容量 = 原容量*2（当<1024时）
                     // 否则扩容25%


src := []int{1,2,3}
dst := make([]int, 2)
n := copy(dst, src)  // n=2, dst=[1,2]

```

## 内存管理机制
    扩容策略（源码级别）
    预估容量 = max(当前Len + 新增元素数, 2*当前Cap)
    根据元素大小选择内存块：
    < 1024字节：每次翻倍
    ≥ 1024字节：增长因子降为1.25
    内存对齐调整最终容量

```go
// 扩容示例
s := []int{1,2,3}
s = append(s, 4,5)  // 新Cap=6（原Cap=3 → 3*2=6）
```

    nil切片：var s []int（Data=nil, Len=0, Cap=0）
    空切片：s := make([]int, 0)（Data指向空数组）
    内存泄露陷阱
    
```go
var large [1e6]int
small := large[:3]  // 整个大数组无法被GC回收
```

### 进阶使用技巧
切片复用模式
```go
// 重用底层数组
func process(data []byte) []byte {
    // 处理数据...
    return data[:0]  // 清空切片保留容量
}

buf := make([]byte, 0, 1024)
for {
    buf = process(buf)  // 复用缓冲区
}
```

避免意外的数据共享
```go
original := []int{1,2,3}
// 错误方式：共享底层数组
shared := original[:2]  
shared[0] = 9          // 修改会影响original

// 安全复制
copied := make([]int, len(original))
copy(copied, original)
copied[0] = 9          // 不影响original
```

避免意外的数据共享
```go
original := []int{1,2,3}
// 错误方式：共享底层数组
shared := original[:2]  
shared[0] = 9          // 修改会影响original

// 安全复制
copied := make([]int, len(original))
copy(copied, original)
copied[0] = 9          // 不影响original
```

批量追加优化
```go
// 低效方式
for _, v := range data {
    slice = append(slice, v)
}

// 高效方式
slice = append(slice, data...)
```
