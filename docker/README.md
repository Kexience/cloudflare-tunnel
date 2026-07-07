# Cloudflare Tunnel 管理工具 - Docker 部署

## 快速开始

```bash
# 1. 进入 docker 目录
cd docker

# 2. 一键启动（首次运行会自动生成 .env 文件）
./start.sh
```

启动后访问 http://localhost:8083

## 常用命令

```bash
# 启动服务
./start.sh

# 停止服务
./stop.sh

# 更新到最新版本
./update.sh

# 查看日志
docker compose logs -f

# 重启服务
docker compose restart
```

## 配置说明

首次运行会自动生成 `.env` 文件，包含以下配置：

| 变量 | 说明 |
|------|------|
| `JWT_SECRET` | JWT 认证密钥（自动生成） |
| `CREDENTIAL_SECRET` | 凭证加密密钥（自动生成） |
| `GITHUB_USER` | GitHub 用户名（需手动修改） |

## 数据持久化

SQLite 数据库存储在 `./data` 目录，备份时只需复制此目录。

## 手动部署

如果不想使用脚本，可以手动操作：

```bash
cd docker

# 复制环境变量模板
cp ../.env.example .env

# 编辑 .env，填写 GITHUB_USER
vim .env

# 创建数据目录
mkdir -p data

# 启动
docker compose up -d
```