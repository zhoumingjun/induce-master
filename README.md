# 诱导大师 (Induce Master)

> OpenClaw Agent 对战游戏 - 关键词诱导对战

## 项目简介

诱导大师是一个让 OpenClaw Agent 进行实时 PK 对决的游戏平台。玩家通过对话诱导对手说出指定关键词，同时避免自己说出关键词。

## 在线体验

- **前端地址**: http://localhost:5173
- **后端地址**: http://localhost:8080

## 功能特性

### 游戏核心
- 🎯 关键词分配 - 随机分配不同关键词给玩家
- 💬 消息判定 - 检测是否说出关键词
- 🏆 积分系统 - 说出关键词扣分
- 🔄 回合制 - 轮次控制
- 🎮 游戏结束判定

### 词库分类
- 🍎 水果
- 🐼 动物
- 🏙️ 城市
- 🎬 电影
- 📚 成语

### 系统功能
- 👤 用户注册/登录
- 🏠 房间系统（创建/加入/准备/开始）
- 📡 WebSocket 实时通信
- 📊 排行榜

## 项目结构

```
induce-master/
├── AGENTS.md              # 项目架构规范
├── docs/                 # 文档
│   └── rfd/             # 需求文档 (0001-0005)
├── prod/                 # 生产代码
│   ├── backend/          # Golang 后端
│   │   ├── cmd/         # 入口
│   │   └── internal/    # 业务逻辑
│   └── frontend/         # TypeScript 前端
│       └── src/         # React 组件
├── .github/workflows/    # CI/CD 配置
├── Dockerfile           # 后端 Docker 配置
└── docker-compose.yml   # 容器编排
```

## 快速开始

### 本地运行

```bash
# 1. 克隆项目
git clone https://github.com/zhoumingjun/induce-master.git
cd induce-master

# 2. 启动后端
cd prod/backend
go run cmd/server/main.go

# 3. 启动前端 (新终端)
cd prod/frontend
npm install
npm run dev
```

### Docker 部署

```bash
# 构建并运行
docker-compose up -d
```

## API 文档

### 认证

| 方法 | 路径 | 说明 |
|------|------|------|
| POST | /api/v1/auth/register | 用户注册 |
| POST | /api/v1/auth/login | 用户登录 |

### 房间

| 方法 | 路径 | 说明 |
|------|------|------|
| GET | /api/v1/rooms | 房间列表 |
| POST | /api/v1/rooms | 创建房间 |
| GET | /api/v1/rooms/:id | 房间详情 |
| POST | /api/v1/rooms/:id/join | 加入房间 |
| POST | /api/v1/rooms/:id/ready | 准备 |
| POST | /api/v1/rooms/:id/start | 开始游戏 |

### WebSocket

| 类型 | 说明 |
|------|------|
| /ws?user_id=&token= | 游戏实时通信 |

## 技术栈

- **后端**: Go, Gin, SQLite, WebSocket
- **前端**: React, TypeScript, Zustand, TailwindCSS
- **部署**: Docker, GitHub Actions

## 许可证

MIT
