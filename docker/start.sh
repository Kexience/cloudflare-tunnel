#!/bin/bash
set -e

# 颜色定义
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
RED='\033[0;31m'
NC='\033[0m'

echo -e "${GREEN}=== Cloudflare Tunnel 管理工具 - Docker 部署 ===${NC}"

# 创建数据目录
mkdir -p data

# 检查 .env 文件
if [ ! -f .env ]; then
    echo -e "${YELLOW}未检测到 .env 文件，正在从模板生成...${NC}"
    
    # 从模板复制
    cp .env.example .env
    
    # 生成随机密钥
    JWT_SECRET=$(openssl rand -base64 32)
    CREDENTIAL_SECRET=$(openssl rand -base64 24)
    
    # 替换模板中的空值
    if [[ "$OSTYPE" == "darwin"* ]]; then
        # macOS
        sed -i '' "s/^JWT_SECRET=$/JWT_SECRET=${JWT_SECRET}/" .env
        sed -i '' "s/^CREDENTIAL_SECRET=$/CREDENTIAL_SECRET=${CREDENTIAL_SECRET}/" .env
    else
        # Linux
        sed -i "s/^JWT_SECRET=$/JWT_SECRET=${JWT_SECRET}/" .env
        sed -i "s/^CREDENTIAL_SECRET=$/CREDENTIAL_SECRET=${CREDENTIAL_SECRET}/" .env
    fi
    
    echo -e "${GREEN}.env 文件已自动生成${NC}"
    echo ""
fi

# 加载环境变量
source .env

# 检查必要变量
if [ -z "$JWT_SECRET" ]; then
    echo -e "${RED}错误: JWT_SECRET 未设置，请编辑 .env 文件${NC}"
    exit 1
fi

if [ -z "$CREDENTIAL_SECRET" ]; then
    echo -e "${RED}错误: CREDENTIAL_SECRET 未设置，请编辑 .env 文件${NC}"
    exit 1
fi

echo -e "${GREEN}启动服务...${NC}"
docker compose up -d

echo ""
echo -e "${GREEN}服务已启动！${NC}"
echo -e "访问地址: http://localhost:${APP_PORT:-8083}"
echo -e "查看日志: docker compose logs -f"