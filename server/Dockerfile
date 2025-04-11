FROM golang:1.20-alpine AS builder

# 设置工作目录
WORKDIR /app

# 设置Go环境变量
ENV CGO_ENABLED=0 \
    GOOS=linux \
    GOARCH=amd64 \
    GO111MODULE=on

# 复制go.mod和go.sum文件
COPY go.mod go.sum ./

# 下载所有依赖
RUN go mod download

# 复制源代码
COPY . .

# 构建应用
RUN go build -o main .

# 使用scratch作为最终镜像
FROM alpine:latest

WORKDIR /app

# 从builder阶段复制编译好的二进制文件
COPY --from=builder /app/main .
COPY --from=builder /app/config ./config

# 设置时区
RUN apk --no-cache add tzdata && \
    cp /usr/share/zoneinfo/Asia/Shanghai /etc/localtime && \
    echo "Asia/Shanghai" > /etc/timezone && \
    apk del tzdata

# 暴露端口
EXPOSE 8080

# 设置入口点
ENTRYPOINT ["/app/main"] 