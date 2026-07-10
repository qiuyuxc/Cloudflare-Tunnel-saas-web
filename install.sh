#!/bin/bash
set -e

SCRIPT_DIR="$(cd "$(dirname "$0")" && pwd)"
cd "$SCRIPT_DIR"

echo "=============================="
echo "  Tunnel Manager 安装/更新"
echo "=============================="

command -v docker >/dev/null 2>&1 || { echo "❌ 未安装 Docker，请先安装"; exit 1; }
if docker compose version >/dev/null 2>&1; then
  COMPOSE="docker compose"
elif command -v docker-compose >/dev/null 2>&1; then
  COMPOSE="docker-compose"
else
  echo "❌ 未安装 Docker Compose"; exit 1
fi

# Download docker-compose.yml if missing (first run via curl | bash)
if [ ! -f docker-compose.yml ]; then
  echo ""
  echo "📥 首次使用，下载配置文件..."
  echo "🌐 选择下载源："
  echo "  1) 国内加速 (ghfast.top)"
  echo "  2) GitHub 直连"
  read -p "请选择 [1/2，默认1]: " DL_MIRROR
  DL_MIRROR=${DL_MIRROR:-1}
  if [ "$DL_MIRROR" = "1" ]; then
    BASE="https://ghfast.top/https://raw.githubusercontent.com/qiuyuxc/Cloudflare-Tunnel-saas-web/main"
  else
    BASE="https://raw.githubusercontent.com/qiuyuxc/Cloudflare-Tunnel-saas-web/main"
  fi
  curl -sO "$BASE/docker-compose.yml"
  # Auto-detect: if ghcr.io is slow (>3s), switch to NJU mirror
  GHCR_TIME=$(curl -s -o /dev/null -w "%{time_total}" --connect-timeout 3 --max-time 5 "https://ghcr.io/v2/" 2>/dev/null || echo "99")
  if [ "$(echo "$GHCR_TIME > 3" | bc 2>/dev/null || echo 1)" = "1" ]; then
    sed -i 's|ghcr.io|ghcr.nju.edu.cn|g' docker-compose.yml
    echo "✅ GHCR 响应慢 (${GHCR_TIME}s)，已切换镜像加速"
  fi
  echo "✅ docker-compose.yml 已下载"
fi

# First-time config
if [ ! -f .env ]; then
  echo ""
  echo "📝 首次运行，配置环境变量："
  read -p "CF_API_TOKEN: " CF_TOKEN
  read -p "CF_ACCOUNT_ID: " CF_AID
  read -p "API_KEY (可选，直接回车跳过): " API_KEY
  read -p "ADMIN_PASSWORD (可选，直接回车随机生成): " ADMIN_PASS

  cat > .env <<EOF
CF_API_TOKEN=${CF_TOKEN}
CF_ACCOUNT_ID=${CF_AID}
API_KEY=${API_KEY}
ADMIN_PASSWORD=${ADMIN_PASS}
EOF
  echo "✅ .env 已生成"
else
  echo ""
  echo "📦 检测到已有配置，执行更新..."
  source .env
fi

mkdir -p data

# Pull and start
echo ""
echo "⬇️  拉取镜像..."
if ! $COMPOSE pull 2>/dev/null; then
  echo "⚠️  镜像拉取失败，尝试本地构建..."
  $COMPOSE build
fi

echo ""
echo "🚀 启动服务..."
$COMPOSE up -d --force-recreate

echo ""
echo "=============================="
echo "  ✅ 部署完成！"
echo "  访问: http://$(hostname -I | awk '{print $1}'):8080"
echo ""

if [ -z "$ADMIN_PASS" ]; then
  echo "  ⚠️  未设置密码，已自动生成"
  echo "  查看日志获取初始密码:"
  echo "    $COMPOSE logs | grep 密码"
else
  echo "  🔑 管理员密码: $ADMIN_PASS"
fi

echo ""
echo "  常用命令:"
echo "    查看日志: $COMPOSE logs -f"
echo "    停止服务: $COMPOSE down"
echo "    重启服务: $COMPOSE restart"
echo "    更新服务: $0"
echo "=============================="
