import type { User, RoomDetail, LoginResponse, RoomListResponse, RankingItem } from '../types'

const API_BASE = 'http://localhost:8080/api/v1'

class ApiService {
  private token: string = ''

  setToken(token: string) {
    this.token = token
    localStorage.setItem('token', token)
  }

  getToken(): string {
    if (!this.token) {
      this.token = localStorage.getItem('token') || ''
    }
    return this.token
  }

  clearToken() {
    this.token = ''
    localStorage.removeItem('token')
  }

  private async request<T>(path: string, options: RequestInit = {}): Promise<T> {
    const url = `${API_BASE}${path}`
    const headers: Record<string, string> = {
      'Content-Type': 'application/json',
      ...options.headers as Record<string, string>,
    }

    const token = this.getToken()
    if (token) {
      headers['Authorization'] = `Bearer ${token}`
    }

    const response = await fetch(url, {
      ...options,
      headers,
    })

    if (!response.ok) {
      const error = await response.json().catch(() => ({ error: 'Unknown error' }))
      throw new Error(error.error || 'Request failed')
    }

    return response.json()
  }

  // 认证
  async register(username: string, password: string, displayName: string): Promise<LoginResponse> {
    return this.request('/auth/register', {
      method: 'POST',
      body: JSON.stringify({ username, password, display_name: displayName }),
    })
  }

  async login(username: string, password: string): Promise<LoginResponse> {
    return this.request('/auth/login', {
      method: 'POST',
      body: JSON.stringify({ username, password }),
    })
  }

  async me(): Promise<{ user: User }> {
    return this.request('/users/me')
  }

  // 房间
  async listRooms(): Promise<RoomListResponse> {
    return this.request('/rooms')
  }

  async createRoom(name: string, maxPlayers: number = 4, password?: string): Promise<{ room_id: string; name: string }> {
    return this.request('/rooms', {
      method: 'POST',
      body: JSON.stringify({ name, max_players: maxPlayers, password }),
    })
  }

  async getRoom(roomId: string): Promise<{ room: RoomDetail }> {
    return this.request(`/rooms/${roomId}`)
  }

  async joinRoom(roomId: string, userId: string): Promise<{ success: boolean }> {
    return this.request(`/rooms/${roomId}/join`, {
      method: 'POST',
      body: JSON.stringify({ user_id: userId }),
    })
  }

  async leaveRoom(roomId: string, userId: string): Promise<{ success: boolean }> {
    return this.request(`/rooms/${roomId}/leave`, {
      method: 'POST',
      body: JSON.stringify({ user_id: userId }),
    })
  }

  async setReady(roomId: string, userId: string, ready: boolean): Promise<{ success: boolean; ready: boolean }> {
    return this.request(`/rooms/${roomId}/ready`, {
      method: 'POST',
      body: JSON.stringify({ user_id: userId, ready }),
    })
  }

  async startGame(roomId: string): Promise<{ success: boolean }> {
    return this.request(`/rooms/${roomId}/start`, {
      method: 'POST',
    })
  }

  // 排行榜
  async getRanking(): Promise<{ ranking: RankingItem[] }> {
    return this.request('/users/ranking')
  }
}

export const api = new ApiService()
