# Go语言笔记
    main 和 init函数均自动调用。
| 维度     | init   | main   |
|:---------|:------:|-------:|
| 调用方式     | 自动调用（初始化包里的变量）  | 自动调用（程序入口）    |
| 定义位置   | 可在任意包中，支持多个定义  | 仅 main 包中，且唯一     |
| 执行顺序   | （复杂）同一文件中的 init 按代码顺序从上到下执行。<br/>同一包不同文件按文件名字符串顺序（字典序）执行。<br/>不同包的 init 按导入依赖关系执行：先被依赖的包先执行，无依赖时按 import 的逆序执行  | init执行完执行     |


## 指针
定义方式：
    
    //定义变量i,指针p1,p1接受i的内存地址
    i := 1
    var p1 *int
    p1 = &i
    
    //定义指针p2,p2储存p1的内存地址，即指针的指针
    p2 := &p1

    fmt.Println(p1) //输出i的内存地址
    fmt.Println(p2) //输出p1的内存地址

    fmt.Println(*p1,"--",**p2)//输出1--1，访问指针和访问指针的指针。

需要注意*访问指针可能会只想空指针。

指针和uintptr转换关系：*T <---> unsafe.Pointer <---> uintptr

    a := "Hello, world!"
    upA := uintptr(unsafe.Pointer(&a))
    upA += 1

    c := (*uint8)(unsafe.Pointer(upA))
    fmt.Println(*c)

## 结构体struct
#### 用于聚合多个不同类型字段的数据结构，Go当中没有类的概念
#### 首字母大写的字段或方法可在包外访问，否则仅限包内

语法：

        type 结构体名 struct {
            field1 int
            Person  // 匿名字段（继承Person字段）
            // ...
        }

实例化方式:

        p := Person{Name: "Alice", Age: 25} 
        p := new(Person)  // *Person 指针类型
        user := struct{Name string; Age int}{Nzame: "Bob"}//匿名结构

标签：
        
        //用作序列化
        type User struct {
            Name string `json:"name"`
            Age  int    `json:"age"`
        }

结构体方法：

        // 值接收者（操作副本）
        func (p Person) GetName() string {
            return p.Name
        }

        // 指针接收者（操作实际对象）
        func (p *Person) SetAge(age int) {
            p.Age = age
        }        

## 常量和枚举

    const a int = 1 //基础定义
    //定义多个
    const (
        h    byte = 3
        i         = "value"
        j, k      = "v", 4
        l, m      = 5, false
    )

    //枚举
    type Weekday int
    const (
        Sunday Weekday = iota // 0
        Monday                // 1
        Tuesday               // 2
    )
    
    //_ 跳值
    const (
        A = iota // 0
        _        // 跳过，iota=1
        B        // 2
    )

    // 位掩码（左移操作）
    const (
        FlagUp = 1 << iota // 1<<0=1
        FlagBroadcast      // 1<<1=2
    )

    // 数学表达式（如数量级定义）
    const (
        KB = 1 << (10 * iota) // 1<<10=1024
        MB                    // 1<<20
    )

    // 多常量同行声明
    const (
        A, B = iota, iota+1 // A=0, B=1
        C, D                // C=1, D=2
    )

    // 插值与重置
    const (
        A = 100        // iota=0（但显式赋值为100）
        B = iota       // iota=1 → B=1
        C              // iota=2 → C=2
    )

    const (
        A = iota       // 0
        B = "test"     // iota=1（但显式赋值）
        C = iota       // iota=2 → C=2
    )

## channel
    先进先出（FIFO）
    ch := make(chan int)      // 无缓冲通道（同步模式）
    chBuf := make(chan int, 3) // 有缓冲通道（异步模式，容量3）

发送阻塞：无缓冲通道需双方就绪，否则阻塞；有缓冲通道在缓冲区满时阻塞。

接收阻塞：无数据时阻塞，接收已关闭通道返回零值（通过v, ok := <-ch判断关闭状态）。

关闭通道：close(ch)，关闭后发送会panic，接收仍可取剩余数据

    // 无缓冲通道导致死锁（未配对）
    ch := make(chan int)
    ch <- 1 // 阻塞直到接收端就绪，若无接收者则deadlock

    // 有缓冲通道允许短暂异步
    ch := make(chan int, 3)
    ch <- 1; ch <- 2 // 不阻塞，缓冲区未满

参数限制仅发送或仅接受

    //仅发送数据
    func <method_name>(<channel_name> chan <- <type>)

    //仅接收数据
    func <method_name>(<channel_name> <-chan <type>)

## sync包
    互斥锁（Mutex）
    var mutex sync.Mutex
    mutex.Lock()   // 加锁
    defer mutex.Unlock() // 解锁
    // 临界区代码

    WaitGroup
    var wg sync.WaitGroup
    wg.Add(3)       // 添加3个任务
    go func() {
        defer wg.Done() // 任务完成，计数器-1
    }()
    wg.Wait()       // 阻塞直到计数器归零

