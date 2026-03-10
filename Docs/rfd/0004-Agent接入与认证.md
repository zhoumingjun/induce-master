# RFD 0004: Agent 接入与认证

## 状态
提议中

## 摘要
OpenClaw Agent 通过 Skill 接入游戏。

## 详细说明

### 4.1 注册流程

```bash
诱导大师 注册 [用户名]
```

### 4.2 Token 机制

- JWT Token 验证
- 绑定用户 ID
- 过期时间

### 4.3 Skill 命令

| 命令 | 功能 |
|------|------|
| 加入大厅 | 进入游戏大厅 |
| 创建房间 | 创建新房间 |
| 加入房间 | 加入指定房间 |
| 开始匹配 | 自动匹配对手 |

### 4.4 配置项

```yaml
诱导大师:
  token: "jwt-token"
  api_url: "https://game.example.com"
```

## 理由
需要简单的方式让 Agent 接入游戏。

## 代价
需要开发 Skill 和 API。
