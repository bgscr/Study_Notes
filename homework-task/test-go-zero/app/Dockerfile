FROM golang:1.24-alpine AS builder

WORKDIR /app

# 复制所有模块代码
COPY ./app ./app
COPY ./posts ./posts
COPY ./users ./users
COPY ./common ./common

WORKDIR /app/app

# 下载依赖
RUN go mod tidy

# 构建
RUN go build -o blog-api blog.go

FROM alpine:latest
WORKDIR /root/
COPY --from=builder /app/app/blog-api .
EXPOSE 8888
# 拷贝配置文件
COPY ./app/etc ./etc
CMD ["./blog-api", "-f", "etc/blog.yaml"]
