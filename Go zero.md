# 环境搭建
    安装goctl: go install github.com/zeromicro/go-zero/tools/goctl@latest
    安装protobuf:在https://github.com/protocolbuffers/protobuf 选择最新release版本下载，并将解压后的bin\protoc.exe放到gopath的bin下
    安装protoc-gen-go ，如果没有安装请先安装 go install google.golang.org/protobuf/cmd/protoc-gen-go@latest
    安装protoc-gen-go-grpc  ，如果没有安装请先安装 go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest

    goctl 1.3.3版本开始支持以下命令一键安装剩余package:goctl env check -i -f

# 生成代码
### 生成api go代码，和配置
    这个语句在windows的cmd或者powershell，似乎只能单个文件执行，例如user.api,而不能使用*.api
    goctl api go -api user.api -dir ../  --style=goZero

    etc：静态配置文件目录
    demo.go：程序启动入口文件
    internal：单个服务内部文件，其可见范围仅限当前服务
    config：静态配置文件对应的结构体声明目录
    handler：handler 目录，可选，一般 http 服务会有这一层做路由管理，handler 为固定后缀
    logic：业务目录，所有业务编码文件都存放在这个目录下面，logic 为固定后缀
    svc：依赖注入目录，所有 logic 层需要用到的依赖都要在这里进行显式注入
    types：结构体存放目录

### 生成Dockerfile
    goctl docker -go user.go

    #构建镜像并创建启动容器，跟踪容器的日志并保持运行状态，直到你手动 Ctrl+C 停止，或者容器退出
    docker-compose up --build 

    #会构建镜像并启动容器，但命令本身会立刻退出，容器会继续后台运行。
    docker-compose up -d --build

    # 只构建镜像，不启动容器
    docker-compose build

    # 单独启动容器
    docker-compose start blog-api posts.rpc users.rpc

    # 停止容器
    docker-compose stop blog-api
    
    # 构建并创建容器，但不启动
    docker-compose up --build --no-start
### 生成K8S 部署文件
    goctl kube deploy -name redis -namespace adhoc -image redis:6-alpine -o redis.yaml -port 6379

### 生成rpc代码
    goctl rpc protoc user-rpc/pb/user.proto --go_out=./ --go-grpc_out=./  --zrpc_out=./ --style=goZero 
    
    如果proto文件中使用了import "google/api/annotations.proto"; 那么需要将googleapis的google和protobuf的src/google放置到GOPATH的src地址下，比执行：
    goctl rpc protoc ./posts.proto --go_out=./ --go-grpc_out=./ --zrpc_out=./ --style=goZero -I . -I "$env:GOPATH\src"


### 生成 Model代码
    goctl model mysql datasource -url="${username}:${passwd}@tcp(${host}:${port})/${dbname}" -table="${tables}"  -dir="${modeldir}" -cache=true --style=goZero

# api服务之api文件
```go
syntax = "v1"

info(
	title: "用户中心服务"
	desc: "用户中心服务"
	author: "Mikael"
	email: "13247629622@163.com"
	version: "v1"
)
@server(
    //路由前缀
	prefix: usercenter/v1
    //代表当前 service 代码块下的路由生成代码时都会被放到 login 目录下
	group: user
    // 定义一个鉴权控制的中间件，多个中间件以英文逗号,分割，如 Middleware1,Middleware2,中间件按声明顺序执行
    middleware: AuthInterceptor
    // 定义一个超时时长为 3 秒的超时配置，这里可填写为 time.Duration 的字符串形式，详情可参考 
    // https://pkg.go.dev/time#Duration.String
    timeout: 3s
)

service usercenter {
	
	@doc "register"
	@handler register
	post /user/register (RegisterReq) returns (RegisterResp)
	
	@doc "login"
	@handler login
	post /user/login (LoginReq) returns (LoginResp)
}


```


# MapReduce
    MapReduce：主函数，接收generate（生成数据）、mapper（处理数据）、reducer（聚合数据）三个参数。
    Map/MapVoid：仅执行generate和mapper，无reducer。
    Finish/FinishVoid：处理固定数量的任务，支持并行和错误中断 

    | 函数 | 流程阶段 | 返回值 | 核心用途 |
    | -------           | -------                       | -------       |
    | MapReduce         | generate → mapper → reducer   | 有            | 完整的数据并行处理与聚合
    | Map/MapVoid       | generate → mapper             | 有/无         | 仅生成并处理数据，无需聚合
    | Finish/FinishVoid | 并行执行固定任务（无generate）   | 有/无         | 并发执行独立任务，无数据依赖