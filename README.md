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
- 密码重置：忘记密码时通过 CLI 命令重置

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

脚本会自动检测 Docker 环境，引导填写 Cloudflare 凭据，构建并启动服务。首次启动后查看日志获取初始密码：

```bash
docker compose logs | grep 密
```

> 环境变量说明见压缩包内 `.env.example`。

### 重置密码

忘记密码时，在容器内执行：

```bash
# 随机生成新密码
docker compose exec tunnel-manager ./tunnel-manager --reset-password

# 设置为指定密码
docker compose exec tunnel-manager ./tunnel-manager --set-password=新密码
```

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
