# 环境搭建
    安装goctl: go install github.com/zeromicro/go-zero/tools/goctl@latest
    安装protobuf:在https://github.com/protocolbuffers/protobuf 选择最新release版本下载，并将解压后的bin\protoc.exe放到gopath的bin下
    安装protoc-gen-go ，如果没有安装请先安装 go install google.golang.org/protobuf/cmd/protoc-gen-go@latest
    安装protoc-gen-go-grpc  ，如果没有安装请先安装 go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest

    goctl 1.3.3版本开始支持以下命令一键安装剩余package:goctl env check -i -f

# 生成代码
### 生成api go代码，和配置
    这个语句在windows的cmd或者powershell，似乎只能单个文件执行，例如user.api,而不能使用*.api
    goctl api go -api user.api -dir ./  --style=goZero

### 生成Dockerfile
    goctl docker -go user.go
### 生成K8S 部署文件
    goctl kube deploy -name redis -namespace adhoc -image redis:6-alpine -o redis.yaml -port 6379

### 生成rpc代码
    goctl rpc protoc user-rpc/pb/user.proto --go_out=./ --go-grpc_out=./  --zrpc_out=./ --style=goZero 

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
