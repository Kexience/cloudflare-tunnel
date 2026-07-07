# 前端静态文件镜像（用于 Nginx 部署）
FROM node:20-alpine AS builder
WORKDIR /app

# 构建参数：API 基础路径
ARG VITE_API_BASE_URL=/api
ENV VITE_API_BASE_URL=${VITE_API_BASE_URL}

# 安装 bun
RUN npm install -g bun

# 复制依赖文件
COPY package.json ./

# 安装依赖
RUN bun install

# 复制源码
COPY . .

# 构建
RUN bun run build

# Nginx 镜像
FROM nginx:alpine

# 复制构建产物
COPY --from=builder /app/dist /usr/share/nginx/html

# 复制 Nginx 配置
COPY nginx.conf /etc/nginx/conf.d/default.conf

EXPOSE 80
CMD ["nginx", "-g", "daemon off;"]