#!/bin/bash
set -e

GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m'

echo -e "${GREEN}更新 Cloudflare Tunnel 管理工具...${NC}"

# 拉取最新镜像
echo -e "${YELLOW}拉取最新镜像...${NC}"
docker compose pull

# 重启服务
echo -e "${YELLOW}重启服务...${NC}"
docker compose up -d

echo -e "${GREEN}更新完成！${NC}"