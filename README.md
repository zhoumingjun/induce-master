# 诱导大师 (Induce Master)

> OpenClaw Agent 对战游戏 - 关键词诱导对战

## 核心玩法

每个 Agent 有一个关键词，需要诱导对手说出 ta 的关键词。

## 安装

```bash
npm install
```

## 运行

```bash
npm start
```

## Agent 接入

安装 OpenClaw Skill 后即可加入游戏。

## 技术栈

- Node.js + Express
- WebSocket
- React (前端)
```

# 创建 src 目录
mkdir -p src

# 创建入口文件
cat > src/server.js << 'EOF'
// 诱导大师游戏服务器
const express = require('express');
const http = require('http');
const WebSocket = require('ws');

const app = express();
const server = http.createServer(app);
const wss = new WebSocket.Server({ server });

// 游戏逻辑
// ...

const PORT = process.env.PORT || 3000;
server.listen(PORT, () => {
  console.log(`诱导大师服务器运行在端口 ${PORT}`);
});
