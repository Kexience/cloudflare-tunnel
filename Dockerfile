# 构建参数
ARG JWT_SECRET

# 构建前端静态文件
FROM node:20-alpine AS frontend-builder
WORKDIR /app/frontend
COPY frontend/package.json frontend/bun.lock ./
RUN npm install -g bun && bun install
COPY frontend/ .
RUN bun run build

# 构建后端二进制文件
FROM golang:1.23-alpine AS backend-builder
WORKDIR /app/backend
COPY backend/go.mod backend/go.sum ./
RUN go mod download
COPY backend/ .
RUN apk add --no-cache bash curl tar
RUN chmod +x scripts/download-cloudflared.sh && ./scripts/download-cloudflared.sh
RUN CGO_ENABLED=0 GOOS=linux go build -o /dist/server cmd/web/main.go

# 最终镜像
FROM alpine:latest
RUN apk add --no-cache ca-certificates tzdata
WORKDIR /app
# 设置环境变量
ENV JWT_SECRET=${JWT_SECRET}
# 从构建阶段复制前端静态文件
COPY --from=frontend-builder /app/frontend/dist ./static
# 从构建阶段复制后端二进制文件
COPY --from=backend-builder /dist/server ./server
# 暴露端口
EXPOSE 8080
# 启动命令
CMD ["./server"]