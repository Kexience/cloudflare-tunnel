.PHONY: dev build prepare clean test

# 开发运行
dev:
	go run cmd/web/main.go

# 构建前准备 - 下载 cloudflared 二进制
prepare:
	./scripts/download-cloudflared.sh

# 生产构建
build: prepare
	go build -o dist/server cmd/web/main.go

# 清理构建产物
clean:
	rm -rf dist/
	rm -f pkg/cloudflared/bin/cloudflared-*

# 运行测试
test:
	go test ./...

# 代码检查
vet:
	go vet ./...
