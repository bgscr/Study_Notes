FROM golang:1.24-alpine AS builder

WORKDIR /app

COPY ./posts ./posts

WORKDIR /app/posts
RUN go mod tidy
RUN go build -o posts-rpc posts.go

FROM alpine:latest
WORKDIR /root/
COPY --from=builder /app/posts/posts-rpc .
# 拷贝配置文件
COPY ./posts/etc ./etc
CMD ["./posts-rpc", "-f", "etc/posts.yaml"]
