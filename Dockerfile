# 构建后端
FROM golang:1.20-alpine AS backend-builder

# 设置工作目录
WORKDIR /app

# 设置Go环境变量
ENV CGO_ENABLED=0 \
    GOOS=linux \
    GOARCH=amd64 \
    GO111MODULE=on

# 复制go.mod和go.sum文件
COPY server/go.mod server/go.sum ./

# 下载所有依赖
RUN go mod download

# 复制源代码
COPY server/ ./

# 构建应用
RUN go build -o main .

# 构建前端
FROM node:18-alpine AS frontend-builder

# 设置工作目录
WORKDIR /app

# 复制package.json和package-lock.json
COPY web/package*.json ./

# 安装依赖
RUN npm ci

# 复制所有源文件
COPY web/ ./

# 构建应用
RUN npm run build

# 最终镜像
FROM alpine:latest

WORKDIR /app

# 安装必要的软件
RUN apk --no-cache add ca-certificates tzdata nginx supervisor && \
    mkdir -p /run/nginx && \
    mkdir -p /var/log/supervisor

# 设置时区
RUN cp /usr/share/zoneinfo/Asia/Shanghai /etc/localtime && \
    echo "Asia/Shanghai" > /etc/timezone

# 从builder阶段复制编译好的二进制文件和配置
COPY --from=backend-builder /app/main ./
COPY --from=backend-builder /app/config ./config

# 复制前端构建产物
COPY --from=frontend-builder /app/dist /usr/share/nginx/html

# 复制配置文件
COPY nginx.conf /etc/nginx/http.d/default.conf
COPY supervisord.conf /etc/supervisor/conf.d/supervisord.conf

# 暴露端口
EXPOSE 80 8080

# 健康检查
HEALTHCHECK --interval=30s --timeout=10s --start-period=10s --retries=3 \
    CMD wget -q --spider http://localhost:8080/health || exit 1

# 启动supervisord管理服务
CMD ["/usr/bin/supervisord", "-c", "/etc/supervisor/conf.d/supervisord.conf"] 