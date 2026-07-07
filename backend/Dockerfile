# 后端服务镜像
FROM golang:1.26-alpine AS builder
WORKDIR /app

# 安装依赖工具
RUN apk add --no-cache bash curl tar

# 下载依赖
COPY go.mod go.sum ./
RUN go mod download

# 复制源码
COPY . .

# 下载 cloudflared 二进制
RUN chmod +x scripts/download-cloudflared.sh && ./scripts/download-cloudflared.sh

# 构建
RUN CGO_ENABLED=0 GOOS=linux go build -o /dist/server cmd/web/main.go

# 最终镜像
FROM alpine:latest
RUN apk add --no-cache ca-certificates tzdata
WORKDIR /app

# 复制构建产物
COPY --from=builder /dist/server ./server

# 暴露端口
EXPOSE 8083

# 启动命令
CMD ["./server"]