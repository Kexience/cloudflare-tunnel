# AGENTS.md — cloudflared-tunnel-go

## 项目概述

内网穿透工具：Go+Gin 后端管理 Cloudflare Tunnel，Vite+React(TS) 前端提供可视化控制台。

## 目录结构

```
backend/          # Go + Gin API 服务
frontend/         # Vite + React (TypeScript) 管理界面
```

## 技术栈

| 层 | 技术 |
|---|------|
| 后端 | Go, Gin, cloudflared SDK/CLI |
| 前端 | Vite, React 18+, TypeScript, Tailwind CSS |
| 隧道 | Cloudflare Tunnel (argo tunnel) |
| 部署 | 前后端独立部署，前端构建后由 Nginx/CDN 托管 |
| 认证 | JWT |
| 数据库 | SQLite |

## 开发命令

> 项目初始化后，用实际命令替换以下占位符。

```bash
# 后端
cd backend && go run ./cmd/server      # 启动开发服务器
cd backend && go build ./cmd/server    # 编译
cd backend && go test ./...            # 运行全部测试
cd backend && go test ./... -run TestXxx  # 运行单个测试

# 前端
cd frontend && npm install             # 安装依赖
cd frontend && npm run dev             # 启动开发服务器
cd frontend && npm run build           # 构建生产版本
cd frontend && npm run lint            # 代码检查
cd frontend && npx tsc --noEmit        # 类型检查
```

## 后端架构要点

- **入口**: `backend/cmd/server/main.go`
- **路由**: `backend/internal/router/` — Gin 路由定义
- **处理器**: `backend/internal/handler/` — HTTP 处理函数
- **隧道管理**: `backend/internal/tunnel/` — Cloudflare Tunnel 生命周期管理
- **配置**: 环境变量优先，支持 `.env` 文件

### 关键约束

- Tunnel 操作需要 Cloudflare API Token，通过环境变量 `CF_API_TOKEN` 传入
- `cloudflared` 二进制必须在 PATH 中，或通过配置指定路径
- 后端 API 需要处理 tunnel 进程的启动/停止/状态查询
- 使用 Cloudflare SDK 管理隧道配置（创建/删除/查询）
- 使用 `cloudflared` CLI 启动隧道进程
- JWT 认证中间件保护 API 端点
- SQLite 存储隧道配置和用户数据

## 前端架构要点

- **入口**: `frontend/src/main.tsx`
- **API 层**: `frontend/src/api/` — 与后端通信
- **组件**: `frontend/src/components/` — UI 组件
- **页面**: `frontend/src/pages/` — 路由页面

### 关键约束

- 开发模式下需配置 Vite proxy 转发 API 请求到后端
- TypeScript 严格模式开启

## 开发工作流

1. 后端和前端可独立开发，各自启动开发服务器
2. 前端通过 proxy 访问后端 API（`/api/*`）
3. 提交前运行：`go vet ./...` + `npm run lint` + `npx tsc --noEmit`

## 通用约束

1. **禁止使用 `interface{}` 或 `any`**：前后端代码都必须使用明确的类型定义
2. **禁止主动提交**：除非用户明确要求，否则不得执行 git commit
3. **读取 AGENTS.md**：操作前后端项目前，必须先读取对应目录的 AGENTS.md
4. **构建产物输出目录**：所有构建产物（测试、打包、测试打包）必须输出到 `dist/` 目录，严禁在项目根目录生成任何构建产物，避免污染版本库

## 待确认项

- [x] 前端 UI 框架选择（Tailwind / Ant Design / shadcn/ui）
- [x] 前端是否嵌入 Go 二进制（go:embed）还是独立部署
- [x] 认证方案（JWT / session / 无认证）
- [x] 数据库选择（SQLite / 文件配置 / 无状态）
- [x] Cloudflare Tunnel 管理方式（SDK 调用 vs 调用 cloudflared CLI）
