# Go 语言数组与切片对比学习笔记

# 一、数组
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

# 二、切片
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

## 切片内存管理机制
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

# Golang 哈希表Map
### 初始化路径
| 场景                        | 处理策略                                                  |
| ------------------------- | ----------------------------------------------------- |
| **字面量 ≤25 项**             | 直接 `make` + 逐条赋值                                      |
| **字面量 >25 项**             | 生成两个切片 → for-loop 赋值                                  |
| **`make(map[K]V, hint)`** | `runtime.makemap` → 计算最小 2^B；预分配 `(2^B + 2^(B-4))` 个桶 |

### 读写流程
| 操作         | 函数             | 关键步骤                                          |
| ---------- | -------------- | --------------------------------------------- |
| **Read**   | `mapaccess1/2` | hash → bucketMask → bucket → tophash → key 匹配 |
| **Write**  | `mapassign`    | 同上定位 → 找空位/溢出桶 → 设置 tophash → 触发扩容            |
| **Delete** | `mapdelete`    | 置 `tophash = emptyOne` & `key,value = nil`    |



### 顶层结构 `hmap`
```go
type hmap struct {
    count     int       // 当前元素数量
    flags     uint8     // 状态标志（迭代、写入等）
    B         uint8     // 桶数量对数（桶数 = 2^B）
    noverflow uint16    // 溢出桶数量
    hash0     uint32    // 哈希种子（防碰撞攻击）
    
    buckets    unsafe.Pointer  // 主桶数组指针
    oldbuckets unsafe.Pointer  // 扩容旧桶指针
    nevacuate  uintptr         // 迁移进度计数器
    
    extra *mapextra  // 溢出桶管理结构
}
```

### `bmap`
```go
type bmap struct {
    tophash [8]uint8      // 哈希值高8位数组
    keys    [8]keyType    // 键数组（内存连续）
    values  [8]valueType  // 值数组（内存连续）
    overflow *bmap        // 溢出桶链表指针
}
```
    键值分离存储减少内存对齐浪费，tophash 快速过滤不匹配键。

### 哈希计算
    通过 hash0 随机化哈希值，防止哈希洪水攻击。
    hash := alg.hash(key, uintptr(h.hash0))  // 使用哈希种子计算

### 桶定位
    桶索引：取哈希值低 B 位
    bucketIndex = hash & (1<<B - 1)
    槽位定位：用哈希高8位（tophash）匹配桶内键。

### 哈希冲突解决
    1. 桶内线性探测
    每个桶存8个键值对，插入时顺序查找空槽位。
    快速失败：若 tophash 为空（emptyOne），直接插入。
    2. 溢出桶链表
    桶满时创建新桶并链接到链表尾部。
    查询代价：最坏时间复杂度 O(n)，但实际因负载因子控制而表现良好。

### 扩容机制
    1. 触发条件
    条件类型	描述
    负载因子 > 6.5	平均每个桶元素超过6.5个
    溢出桶过多	溢出桶数量 ≥ 常规桶数量 ∧ B < 15
    2. 扩容类型
    类型	描述	场景
    增量扩容	容量翻倍（B+1）	负载因子过高
    等量扩容	桶数不变，重排数据减少溢出桶	大量删除导致桶稀疏
    3. 渐进式迁移
    写时迁移：每次插入/修改时迁移1~2个旧桶。
    查询双桶：未迁移完成时需同时检查新旧桶。

### 关键操作分析
    1. 插入流程
    计算哈希值，定位目标桶。
    遍历桶链：
    存在相同键 → 更新值
    找到空槽 → 插入键值
    桶链已满 → 创建溢出桶
    检查扩容条件。
    2. 查询流程    
    val, ok := map[key]
    计算哈希，定位桶链。
    遍历桶链，比较 tophash 和键。
    返回匹配值或 nil。
    3. 删除流程
    标记删除：设置 tophash 为 emptyOne，不释放内存。
    碎片整理：等量扩容时合并稀疏桶。