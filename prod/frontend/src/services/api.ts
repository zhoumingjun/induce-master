const API_BASE = 'http://localhost:8080/api/v1'

// 通用请求方法
async function request(endpoint: string, options: RequestInit = {}): Promise<any> {
  const token = localStorage.getItem('token')
  
  const config: RequestInit = {
    headers: {
      'Content-Type': 'application/json',
      ...(token && { 'Authorization': `Bearer ${token}` }),
      ...options.headers,
    },
    ...options,
  }

  const response = await fetch(`${API_BASE}${endpoint}`, config)
  
  if (!response.ok) {
    const error = await response.json().catch(() => ({ error: '请求失败' }))
    throw new Error(error.error || '请求失败')
  }

  return response.json()
}

// 认证 API
export const authAPI = {
  // 注册
  async register(username: string, password: string): Promise<any> {
    const data = await request('/auth/register', {
      method: 'POST',
      body: JSON.stringify({ username, password }),
    })
    if (data.token) {
      localStorage.setItem('token', data.token)
      localStorage.setItem('user_id', data.user_id)
      localStorage.setItem('username', data.username)
    }
    return data
  },

  // 登录
  async login(username: string, password: string): Promise<any> {
    const data = await request('/auth/login', {
      method: 'POST',
      body: JSON.stringify({ username, password }),
    })
    if (data.token) {
      localStorage.setItem('token', data.token)
      localStorage.setItem('user_id', data.user_id)
      localStorage.setItem('username', data.username)
    }
    return data
  },

  // 获取当前用户
  async me(): Promise<any> {
    return request('/users/me')
  },

  // 登出
  logout(): void {
    localStorage.removeItem('token')
    localStorage.removeItem('user_id')
    localStorage.removeItem('username')
  },

  // 获取登录状态
  isLoggedIn(): boolean {
    return !!localStorage.getItem('token')
  },

  getUserID(): string {
    return localStorage.getItem('user_id') || ''
  },

  getUsername(): string {
    return localStorage.getItem('username') || ''
  },
}

// 房间 API
export const roomAPI = {
  // 获取房间列表
  async list(): Promise<any[]> {
    const data = await request('/rooms')
    return data.rooms || []
  },

  // 创建房间
  async create(name: string, maxPlayers: number = 4): Promise<any> {
    return request('/rooms', {
      method: 'POST',
      body: JSON.stringify({
        name,
        owner_id: localStorage.getItem('user_id'),
        max_players: maxPlayers,
      }),
    })
  },

  // 获取房间详情
  async get(roomId: string): Promise<any> {
    return request(`/rooms/${roomId}`)
  },

  // 加入房间
  async join(roomId: string): Promise<any> {
    return request(`/rooms/${roomId}/join`, {
      method: 'POST',
      body: JSON.stringify({
        user_id: localStorage.getItem('user_id'),
      }),
    })
  },

  // 离开房间
  async leave(roomId: string): Promise<any> {
    return request(`/rooms/${roomId}/leave`, {
      method: 'POST',
      body: JSON.stringify({
        user_id: localStorage.getItem('user_id'),
      }),
    })
  },

  // 准备
  async ready(roomId: string, ready: boolean): Promise<any> {
    return request(`/rooms/${roomId}/ready`, {
      method: 'POST',
      body: JSON.stringify({
        user_id: localStorage.getItem('user_id'),
        ready,
      }),
    })
  },

  // 开始游戏
  async start(roomId: string): Promise<any> {
    return request(`/rooms/${roomId}/start`, {
      method: 'POST',
    })
  },
}

// 排行榜 API
export const rankingAPI = {
  async get(): Promise<any[]> {
    const data = await request('/users/ranking')
    return data.ranking || []
  },
}

// WebSocket 连接
export function createWebSocket(): WebSocket {
  const userId = localStorage.getItem('user_id')
  const token = localStorage.getItem('token')
  
  if (!userId || !token) {
    throw new Error('未登录')
  }

  const ws = new WebSocket(`ws://localhost:8080/ws?user_id=${userId}&token=${token}`)

  ws.onmessage = (event) => {
    const message = JSON.parse(event.data)
    return message
  }

  return ws
}

export default {
  auth: authAPI,
  room: roomAPI,
  ranking: rankingAPI,
}
