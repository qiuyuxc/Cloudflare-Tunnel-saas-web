#!/bin/bash
set -e

echo "=============================="
echo "  Tunnel Manager 一键部署"
echo "=============================="

# Check dependencies
command -v docker >/dev/null 2>&1 || { echo "❌ 未安装 Docker，请先安装"; exit 1; }
if docker compose version >/dev/null 2>&1; then
  COMPOSE="docker compose"
elif command -v docker-compose >/dev/null 2>&1; then
  COMPOSE="docker-compose"
else
  echo "❌ 未安装 Docker Compose"; exit 1
fi

# Get env vars
if [ ! -f .env ]; then
  echo ""
  echo "📝 首次运行，配置环境变量："
  read -p "CF_API_TOKEN: " CF_TOKEN
  read -p "CF_ACCOUNT_ID: " CF_AID
  read -p "API_KEY (可选，直接回车跳过): " API_KEY

  cat > .env <<EOF
CF_API_TOKEN=${CF_TOKEN}
CF_ACCOUNT_ID=${CF_AID}
API_KEY=${API_KEY}
EOF
  echo "✅ .env 已生成"
fi

# Create data directory for persistence
mkdir -p data

# Mirror source selection
echo ""
echo "🌐 选择镜像源："
echo "  1) 国内镜像 (npmmirror / goproxy.cn / 阿里云 Alpine)"
echo "  2) 官方源 (npmjs.org / proxy.golang.org / Alpine CDN)"
read -p "请选择 [1/2，默认1]: " MIRROR_CHOICE
MIRROR_CHOICE=${MIRROR_CHOICE:-1}

if [ "$MIRROR_CHOICE" = "1" ]; then
  export NPM_REGISTRY=https://registry.npmmirror.com
  export GOPROXY=https://goproxy.cn,direct
  export ALPINE_MIRROR=https://mirrors.aliyun.com/alpine/
  echo "✅ 使用国内镜像源"
else
  echo "✅ 使用官方源"
fi

# Build and start
echo ""
echo "🔨 构建镜像..."
$COMPOSE build

echo ""
echo "🚀 启动服务..."
$COMPOSE up -d

echo ""
echo "=============================="
echo "  ✅ 部署完成！"
echo "  访问: http://$(hostname -I | awk '{print $1}'):8080"
echo ""
echo "  ⚠️  首次启动会自动生成随机密码"
echo "  请查看日志获取初始密码:"
echo "    $COMPOSE logs | grep 密"
echo ""
echo "  常用命令:"
echo "    查看日志: $COMPOSE logs -f"
echo "    停止服务: $COMPOSE down"
echo "    重启服务: $COMPOSE restart"
echo "=============================="
