FROM node:18-alpine AS builder

# 设置工作目录
WORKDIR /app

# 复制package.json和package-lock.json（如果有）
COPY package*.json ./

# 安装依赖
RUN npm ci

# 复制所有源文件
COPY . .

# 构建应用
RUN npm run build

# 使用nginx作为生产环境服务器
FROM nginx:alpine

# 复制构建产物到nginx服务目录
COPY --from=builder /app/dist /usr/share/nginx/html

# 复制nginx配置
COPY nginx.conf /etc/nginx/conf.d/default.conf

# 暴露端口
EXPOSE 80

# 启动nginx
CMD ["nginx", "-g", "daemon off;"] 