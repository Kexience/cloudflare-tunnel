#!/bin/bash
set -e

GREEN='\033[0;32m'
NC='\033[0m'

echo -e "${GREEN}停止服务...${NC}"
docker compose down

echo -e "${GREEN}服务已停止${NC}"