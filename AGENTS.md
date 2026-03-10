# AGENTS.md - 项目架构规范

## 项目概述

**诱导大师 (Induce Master)** - OpenClaw Agent 对战游戏

## 技术栈

### 后端
- **语言**: Golang
- **框架**: Gin / Fiber
- **数据库**: MySQL + Redis
- **通信**: WebSocket

### 前端
- **语言**: TypeScript
- **框架**: React
- **构建工具**: Vite
- **样式**: Tailwind CSS

## 项目结构

```
induce-master/
├── docs/                    # 文档
│   ├── AGENTS.md           # 本文件
│   └── 需求定义.md
└── prod/                   # 生产代码
    ├── backend/            # 后端 (Golang)
    │   ├── cmd/           # 入口
    │   │   └── server/   
    │   ├── internal/      # 业务逻辑
    │   └── pkg/          # 公共包
    └── frontend/          # 前端 (TypeScript)
        └── src/
```

## 核心模块

### 后端模块

| 模块 | 说明 |
|------|------|
| game | 游戏核心逻辑 |
| room | 房间管理 |
| user | 用户认证 |
| match | 匹配系统 |
| websocket | 实时通信 |

### 前端模块

| 模块 | 说明 |
|------|------|
| pages | 页面组件 |
| components | 通用组件 |
| hooks | 自定义 Hooks |
| services | API 服务 |
| stores | 状态管理 |

## 开发规范

1. 后端代码遵循 Go 标准规范
2. 前端使用函数式组件 + Hooks
3. 提交信息使用中文
4. 需要测试覆盖

## 待实现功能

- [ ] 用户注册与认证 (JWT)
- [ ] 大厅系统
- [ ] 房间系统
- [ ] 匹配系统
- [ ] 游戏核心逻辑
- [ ] WebSocket 通信
- [ ] 前端界面
