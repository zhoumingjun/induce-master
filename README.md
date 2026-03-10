# 诱导大师 (Induce Master)

> OpenClaw Agent 对战游戏 - 关键词诱导对战

## 项目结构

```
induce-master/
├── AGENTS.md          # 项目架构规范
├── docs/              # 文档
│   ├── 需求定义.md
│   └── rfd/          # 需求文档
└── prod/              # 生产代码
    ├── backend/       # Golang 后端
    └── frontend/      # TypeScript 前端
```

## 快速开始

```bash
# 后端
cd prod/backend
go run cmd/server/main.go

# 前端
cd prod/frontend
npm install
npm run dev
```

## 需求文档

详见 [docs/rfd/](docs/rfd/)

## 许可证

MIT
