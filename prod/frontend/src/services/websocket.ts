type MessageType = 
  | 'connect' 
  | 'room_update' 
  | 'game_start' 
  | 'message' 
  | 'game_end' 
  | 'error'
  | 'ping'
  | 'pong'

interface WSMessage {
  type: MessageType
  payload: unknown
}

type MessageHandler = (payload: unknown) => void

class WebSocketService {
  private ws: WebSocket | null = null
  private userId: string = ''
  private token: string = ''
  private handlers: Map<MessageType, MessageHandler[]> = new Map()
  private reconnectTimer: number | null = null
  private pingInterval: number | null = null

  connect(userId: string, token: string) {
    this.userId = userId
    this.token = token
    
    const wsUrl = `ws://localhost:8080/ws?user_id=${userId}&token=${token}`
    this.ws = new WebSocket(wsUrl)

    this.ws.onopen = () => {
      console.log('WebSocket connected')
      this.startPing()
    }

    this.ws.onmessage = (event) => {
      try {
        const message: WSMessage = JSON.parse(event.data)
        this.handleMessage(message)
      } catch (e) {
        console.error('Failed to parse message:', e)
      }
    }

    this.ws.onclose = () => {
      console.log('WebSocket disconnected')
      this.stopPing()
      this.reconnect()
    }

    this.ws.onerror = (error) => {
      console.error('WebSocket error:', error)
    }
  }

  disconnect() {
    this.stopPing()
    if (this.reconnectTimer) {
      clearTimeout(this.reconnectTimer)
      this.reconnectTimer = null
    }
    if (this.ws) {
      this.ws.close()
      this.ws = null
    }
  }

  private reconnect() {
    if (this.reconnectTimer) return
    
    this.reconnectTimer = window.setTimeout(() => {
      this.reconnectTimer = null
      if (this.userId && this.token) {
        this.connect(this.userId, this.token)
      }
    }, 3000)
  }

  private startPing() {
    this.pingInterval = window.setInterval(() => {
      this.send('ping', {})
    }, 30000)
  }

  private stopPing() {
    if (this.pingInterval) {
      clearInterval(this.pingInterval)
      this.pingInterval = null
    }
  }

  private handleMessage(message: WSMessage) {
    const handlers = this.handlers.get(message.type) || []
    handlers.forEach(handler => handler(message.payload))
  }

  on(type: MessageType, handler: MessageHandler) {
    if (!this.handlers.has(type)) {
      this.handlers.set(type, [])
    }
    this.handlers.get(type)!.push(handler)
  }

  off(type: MessageType, handler: MessageHandler) {
    const handlers = this.handlers.get(type)
    if (handlers) {
      const index = handlers.indexOf(handler)
      if (index > -1) {
        handlers.splice(index, 1)
      }
    }
  }

  send(type: string, payload: unknown) {
    if (this.ws && this.ws.readyState === WebSocket.OPEN) {
      this.ws.send(JSON.stringify({ type, payload }))
    }
  }

  // 便捷方法
  joinRoom(roomId: string) {
    this.send('join_room', { room_id: roomId })
  }

  leaveRoom(roomId: string) {
    this.send('leave_room', { room_id: roomId })
  }

  sendMessage(content: string) {
    this.send('message', { content })
  }

  sendGameAction(action: string, data: Record<string, unknown> = {}) {
    this.send('game_action', { action, ...data })
  }
}

export const ws = new WebSocketService()

// 事件类型导出
export const WSEvents = {
  ROOM_UPDATE: 'room_update' as MessageType,
  GAME_START: 'game_start' as MessageType,
  GAME_MESSAGE: 'message' as MessageType,
  GAME_END: 'game_end' as MessageType,
  ERROR: 'error' as MessageType,
}
