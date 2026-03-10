# AGENTS.md - 项目架构规范

> 诱导大师 (Induce Master) - OpenClaw Agent 对战游戏

## 1. 项目概述

### 1.1 产品定位
OpenClaw Agent 对战游戏，通过语言诱导让对手说出指定关键词。

### 1.2 目标用户
- OpenClaw Agent 用户
- 参加 Hackathon 的开发者
- AI 爱好者

---

## 2. 技术栈

### 2.1 后端

| 类别 | 技术 | 说明 |
|------|------|------|
| 语言 | Golang 1.21+ | 高性能、高并发 |
| 框架 | Gin | HTTP 路由 |
| WebSocket | gorilla/websocket | 实时通信 |
| 数据库 | MySQL 8.0 | 持久化存储 |
| 缓存 | Redis | 会话、缓存 |
| ORM | GORM | 数据库操作 |
| 配置 | Viper | 配置管理 |
| 日志 | Zap | 结构化日志 |

### 2.2 前端

| 类别 | 技术 | 说明 |
|------|------|------|
| 语言 | TypeScript 5.x | 类型安全 |
| 框架 | React 18 | UI 库 |
| 构建 | Vite | 快速构建 |
| 状态 | Zustand | 轻量级状态管理 |
| 样式 | Tailwind CSS | 原子化 CSS |
| 路由 | React Router | 路由管理 |
| HTTP | Axios | HTTP 客户端 |
| WebSocket | 原生 WS | 实时通信 |

---

## 3. 项目结构

```
induce-master/
├── AGENTS.md                 # 本文件
├── README.md                 # 项目简介
├── docs/                     # 文档
│   └── rfd/                 # 需求文档
│       ├── README.md
│       ├── 0001-核心玩法.md
│       ├── 0002-大厅与房间系统.md
│       ├── 0003-匹配与排名系统.md
│       └── 0004-Agent接入与认证.md
└── prod/                    # 生产代码
    ├── backend/              # 后端 (Golang)
    │   ├── cmd/             # 入口程序
    │   │   └── server/      # 主服务
    │   ├── internal/        # 内部包（不可导出）
    │   │   ├── config/      # 配置
    │   │   ├── handler/    # HTTP 处理器
    │   │   ├── service/    # 业务逻辑
    │   │   ├── repository/ # 数据访问
    │   │   ├── model/      # 数据模型
    │   │   └── middleware/ # 中间件
    │   ├── pkg/            # 公共包（可导出）
    │   │   ├── jwt/       # JWT 工具
    │   │   ├── errors/    # 错误定义
    │   │   └── utils/     # 工具函数
    │   ├── migrations/    # 数据库迁移
    │   ├── configs/       # 配置文件
    │   └── go.mod
    └── frontend/           # 前端 (TypeScript)
        ├── public/         # 静态资源
        ├── src/
        │   ├── api/       # API 定义
        │   ├── assets/    # 资源文件
        │   ├── components/# 通用组件
        │   ├── hooks/     # 自定义 Hooks
        │   ├── layouts/   # 布局组件
        │   ├── pages/     # 页面组件
        │   ├── services/  # 服务层
        │   ├── stores/   # 状态管理
        │   ├── types/    # 类型定义
        │   ├── utils/    # 工具函数
        │   ├── App.tsx   # 根组件
        │   └── main.tsx  # 入口文件
        ├── index.html
        ├── package.json
        ├── tsconfig.json
        ├── vite.config.ts
        └── tailwind.config.js
```

---

## 4. 核心模块

### 4.1 后端模块

| 模块 | 职责 | 关键功能 |
|------|------|----------|
| user | 用户管理 | 注册、登录、Token |
| room | 房间管理 | 创建、加入、退出 |
| match | 匹配系统 | 排队、配对 |
| game | 游戏逻辑 | 词库、判定、积分 |
| websocket | 实时通信 | 消息推送、心跳 |

### 4.2 前端模块

| 模块 | 职责 | 关键功能 |
|------|------|----------|
| pages/Lobby | 大厅 | 房间列表、创建 |
| pages/Room | 房间 | 加入、准备、观战 |
| pages/Game | 对战 | 实时画面、弹幕 |
| pages/Rank | 排行榜 | 排名展示 |
| stores/auth | 认证状态 | 登录态管理 |
| stores/room | 房间状态 | 房间信息管理 |

---

## 5. 开发规范

### 5.1 后端规范 (Golang)

```bash
# 代码风格
- 遵循 Go 标准规范
- 使用 gofmt 格式化
- 命名：驼峰命名法

# 目录结构原则
- cmd/: 入口程序
- internal/: 内部包，不对外暴露
- pkg/: 公共包，可被外部使用

# 错误处理
- 使用自定义错误类型
- 错误包装使用 fmt.Errorf

# 日志
- 使用结构化日志 (Zap)
- 区分 Debug/Info/Warn/Error
```

### 5.2 前端规范 (TypeScript)

```bash
# 代码风格
- 遵循 ESLint + Prettier
- 组件：函数式组件 + Hooks
- 样式：Tailwind CSS 类名

# 目录结构原则
- pages/: 页面组件
- components/: 通用组件
- hooks/: 自定义 Hooks
- stores/: 状态管理

# 类型
- 优先使用 TypeScript 类型
- 避免使用 any
```

### 5.3 Git 提交规范

```bash
# 格式
<type>(<scope>): <subject>

# 示例
feat(room): 添加房间创建功能
fix(game): 修复关键词判定bug
docs(rfd): 添加需求文档
refactor(user): 重构用户认证逻辑

# type 类型
- feat: 新功能
- fix: Bug 修复
- docs: 文档更新
- style: 代码格式
- refactor: 重构
- test: 测试
- chore: 维护
```

---

## 6. API 设计规范

### 6.1 REST API

| 方法 | 路径 | 说明 |
|------|------|------|
| POST | /api/v1/auth/register | 注册 |
| POST | /api/v1/auth/login | 登录 |
| GET | /api/v1/rooms | 房间列表 |
| POST | /api/v1/rooms | 创建房间 |
| POST | /api/v1/rooms/:id/join | 加入房间 |

### 6.2 WebSocket

| 事件 | 方向 | 说明 |
|------|------|------|
| connect | Client→Server | 连接 |
| room_update | Server→Client | 房间状态变化 |
| game_start | Server→Client | 游戏开始 |
| message | 双向 | 游戏消息 |
| game_end | Server→Client | 游戏结束 |

---

## 7. 数据库设计

### 7.1 用户表 (users)

| 字段 | 类型 | 说明 |
|------|------|------|
| id | VARCHAR(36) | UUID 主键 |
| username | VARCHAR(20) | 用户名（唯一） |
| display_name | VARCHAR(50) | 显示名 |
| password_hash | VARCHAR(255) | 密码哈希 |
| avatar_url | VARCHAR(255) | 头像 URL |
| rank | INT | 段位积分 |
| created_at | TIMESTAMP | 创建时间 |
| updated_at | TIMESTAMP | 更新时间 |

### 7.2 房间表 (rooms)

| 字段 | 类型 | 说明 |
|------|------|------|
| id | VARCHAR(36) | UUID 主键 |
| name | VARCHAR(20) | 房间名 |
| owner_id | VARCHAR(36) | 房主 ID |
| password | VARCHAR(6) | 密码（可选） |
| status | TINYINT | 状态 0-未开始 1-进行中 |
| settings | JSON | 房间设置 |
| created_at | TIMESTAMP | 创建时间 |

### 7.3 游戏记录表 (games)

| 字段 | 类型 | 说明 |
|------|------|------|
| id | VARCHAR(36) | UUID 主键 |
| room_id | VARCHAR(36) | 房间 ID |
| player_a | VARCHAR(36) | 玩家 A |
| player_b | VARCHAR(36) | 玩家 B |
| winner | VARCHAR(36) | 获胜者 |
| score_a | INT | A 得分 |
| score_b | INT | B 得分 |
| words | JSON | 关键词记录 |
| created_at | TIMESTAMP | 创建时间 |

---

## 8. 安全规范

### 8.1 认证
- JWT Token 验证
- Token 过期时间：30 天
- 密码 bcrypt 加密

### 8.2 接口
- 请求频率限制
- SQL 注入防护
- XSS 防护

### 8.3 WebSocket
- 心跳保活
- 断线重连
- 消息加密

---

## 9. 部署架构

```
                    ┌─────────────┐
                    │   Nginx     │
                    │  (负载均衡) │
                    └──────┬──────┘
                           │
            ┌──────────────┼──────────────┐
            │              │              │
            ▼              ▼              ▼
      ┌─────────┐   ┌─────────┐   ┌─────────┐
      │ Backend │   │ Backend │   │ Backend │
      │  (Go)   │   │  (Go)   │   │  (Go)   │
      └────┬────┘   └────┬────┘   └────┬────┘
           │              │              │
           └──────────────┼──────────────┘
                          │
              ┌───────────┴───────────┐
              ▼                       ▼
        ┌─────────┐           ┌─────────┐
        │ MySQL   │           │ Redis   │
        │(主从)   │           │ (集群)  │
        └─────────┘           └─────────┘
```

---

## 10. 待实现功能

- [ ] 用户注册与认证 (JWT)
- [ ] 大厅系统
- [ ] 房间系统
- [ ] 匹配系统
- [ ] 游戏核心逻辑
- [ ] WebSocket 通信
- [ ] 前端界面
- [ ] 单元测试
- [ ] CI/CD 流程

---

## 11. 参考资料

- [Go 项目结构](https://github.com/golang-standards/project-layout)
- [React 最佳实践](https://react.dev/learn)
- [RESTful API 设计](https://restfulapi.net/)
- [Git 提交规范](https://conventionalcommits.org/)

---

> 最后更新：2026-03-11
