# Go语言笔记
    main 和 init函数均自动调用。
| 维度     | init   | main   |
|:---------|:------:|-------:|
| 调用方式     | 自动调用（初始化包里的变量）  | 自动调用（程序入口）    |
| 定义位置   | 可在任意包中，支持多个定义  | 仅 main 包中，且唯一     |
| 执行顺序   | （复杂）同一文件中的 init 按代码顺序从上到下执行。<br/>同一包不同文件按文件名字符串顺序（字典序）执行。<br/>不同包的 init 按导入依赖关系执行：先被依赖的包先执行，无依赖时按 import 的逆序执行  | init执行完执行     |


    指针
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
