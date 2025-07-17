# CPU profiling
    func main() {
        f, _ := os.OpenFile("cpu.profile", os.O_CREATE|os.O_RDWR, 0644)
        defer f.Close()
        pprof.StartCPUProfile(f)
        defer pprof.StopCPUProfile()

        n := 10
        for i := 1; i <= 5; i++ {
            fmt.Printf("fib(%d)=%d\n", n, fib(n))
            n += 3 * i
        }
    }

    运行后生成cpu.profile
    执行 go tool pprof cpu.profile

    输入top 10命令显示好事最高的10个函数

    输入 list命令后跟一个表示方法名，例如 list main.fib


# Memory profiling
    记录堆内存：
        f, _ := os.OpenFile("mem.profile", os.O_CREATE|os.O_RDWR, 0644)
        defer f.Close()

        // 执行业务代码...

        pprof.Lookup("heap").WriteTo(f, 0)
    分析命令：go tool pprof mem.profile

# HTTP服务集成
    import _ "net/http/pprof"
    func NewProfileHttpServer(addr string) {
        go func() {
            log.Fatalln(http.ListenAndServe(addr, nil))
        }()
    }

    NewProfileHttpServer(":9999")
    
    访问http://localhost:9999/debug/pprof查看实时数据

### 命令远程采样
    go tool pprof -http :8080 localhost:9999/debug/pprof/profile?seconds=120

    也可以使用go tool pprof命令生成文件下载到本地分析，然后执行：
    go tool pprof -http :8888 pprof.samples.cpu.001.pb.gz
