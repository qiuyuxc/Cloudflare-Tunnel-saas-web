# Tunnel Manager

Cloudflare Tunnel 可视化管理面板。通过 Web UI 一键完成隧道选择、域名绑定、DNS 优选、回退源配置。

## 架构

```
┌──────────────┐     ┌──────────────┐     ┌──────────────────┐
│  Vue 3 前端   │────▶│  Go 后端 API  │────▶│ Cloudflare API   │
│  Naive UI    │     │  chi router  │     │ Tunnels / DNS    │
└──────────────┘     └──────────────┘     └──────────────────┘
```

## 功能

- 隧道管理：列出、选择 Cloudflare Tunnel
- 域名绑定：自动配置 Tunnel 路由 + DNS CNAME + SaaS Custom Hostname
- 优选 CNAME：自定义全局优选域名
- 回退源设置：一键配置 fallback origin
- 管理员认证：用户名/密码登录，支持修改密码

## 技术栈

| 层 | 技术 |
|---|---|
| 前端 | Vue 3, TypeScript, Naive UI, Vite, Pinia |
| 后端 | Go, chi, JSON file store |
| 部署 | Docker, GitHub Actions (GHCR) |

## 部署

一键安装，只需 Docker：

```bash
curl -sO https://raw.githubusercontent.com/qiuyuxc/Cloudflare-Tunnel-saas-web/main/install.sh && bash install.sh
```

脚本会自动：
1. 下载 docker-compose.yml（国内自动检测 GHCR 连通性，慢则切换镜像加速）
2. 引导填写 Cloudflare 凭据
3. 拉取预构建镜像并启动

首次启动可自定义管理员密码，不填则自动生成（查看日志获取）。

### 环境变量

| 变量 | 说明 | 必填 |
|------|------|------|
| `CF_API_TOKEN` | Cloudflare API Token | 是 |
| `CF_ACCOUNT_ID` | Cloudflare Account ID | 是 |
| `API_KEY` | API 访问密钥 | 否 |
| `ADMIN_PASSWORD` | 管理员密码 | 否（不填随机生成） |

## 更新

在项目目录下执行：

```bash
./install.sh
```

自动拉取最新镜像并重建容器。

## API

| 方法 | 路径 | 说明 | 鉴权 |
|------|------|------|------|
| GET | `/api/health` | 健康检查 | 无 |
| POST | `/api/admin/login` | 管理员登录 | 无 |
| GET | `/api/admin/status` | 登录状态 | 无 |
| GET | `/api/config` | 获取配置 | 需要 |
| POST | `/api/config/tunnel` | 设置隧道 ID | 需要 |
| POST | `/api/config/service` | 设置转发地址 | 需要 |
| POST | `/api/config/preferred-cname` | 设置优选 CNAME | 需要 |
| GET | `/api/tunnels` | 列出隧道 | 需要 |
| GET | `/api/zones` | 列出 Zone | 需要 |
| POST | `/api/domain/bind` | 绑定域名 | 需要 |
| POST | `/api/domain/fallback` | 设置回退源 | 需要 |

## 项目结构

```
├── backend/
│   ├── main.go           # 入口，路由注册
│   ├── handlers/         # HTTP handlers
│   ├── services/         # Cloudflare API 客户端
│   ├── models/           # 数据类型
│   └── store/            # JSON 文件存储
├── frontend/
│   ├── src/
│   │   ├── views/        # 页面
│   │   ├── components/   # 组件
│   │   ├── stores/       # Pinia stores
│   │   ├── api/          # API 封装
│   │   └── router/       # 路由
│   └── vite.config.ts
├── .github/workflows/    # CI 自动构建
├── install.sh            # 一键部署/更新脚本
├── Dockerfile
└── docker-compose.yml
```
