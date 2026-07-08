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
| 部署 | Docker (multi-stage), docker-compose |

## 部署

```bash
tar xzf tunnel-manager.tar.gz
cd tunnel-manager
./install.sh
```

脚本会自动检测 Docker 环境，引导填写 Cloudflare 凭据，选择镜像源（国内/官方），构建并启动服务。首次启动后查看日志获取初始密码：

```bash
docker compose logs | grep 密
```

> 环境变量说明见压缩包内 `.env.example`。

### 镜像源配置

安装脚本会自动选择镜像源，也可通过环境变量手动指定：

| 变量 | 说明 | 国内镜像 | 默认值 |
|------|------|----------|--------|
| `NPM_REGISTRY` | npm 镜像地址 | `https://registry.npmmirror.com` | `https://registry.npmjs.org` |
| `GOPROXY` | Go 模块代理 | `https://goproxy.cn,direct` | `https://proxy.golang.org,direct` |
| `ALPINE_MIRROR` | Alpine apk 镜像 | `https://mirrors.aliyun.com/alpine/` | `https://dl-cdn.alpinelinux.org/alpine` |

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
├── install.sh            # 一键部署脚本
├── Dockerfile
└── docker-compose.yml
```

