# 第一阶段：构建 Go 二进制
FROM golang:alpine AS builder

# 必备工具（tzdata 保证时区正常）
RUN apk update --no-cache && apk add --no-cache tzdata

WORKDIR /build

# 先复制 go.mod 和 go.sum（利用缓存）
COPY go.mod go.sum ./
RUN go mod download

# 再复制所有源码
COPY . .

# 构建指定服务
ARG SERVICE_PATH
ARG OUTPUT_NAME
RUN go build -ldflags="-s -w" -o /app/${OUTPUT_NAME} ${SERVICE_PATH}

# 第二阶段：运行
FROM alpine:3.20
WORKDIR /app
COPY --from=builder /app/${OUTPUT_NAME} .
EXPOSE 8888
CMD ["./blog"] # 默认启动 blog，可在 docker-compose 中覆盖
