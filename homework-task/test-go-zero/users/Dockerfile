FROM golang:1.24-alpine AS builder

WORKDIR /app

COPY ./users ./users

WORKDIR /app/users
RUN go mod tidy
RUN go build -o users-rpc users.go

FROM alpine:latest
WORKDIR /root/
COPY --from=builder /app/users/users-rpc .
# 拷贝配置文件
COPY ./users/etc ./etc
CMD ["./users-rpc", "-f", "etc/users.yaml"]
