import { create } from 'zustand'
import { persist } from 'zustand/middleware'
import type { User, Room, RoomDetail, GameMessage, GameStart, GameEnd } from '../types'
import { api } from '../api'
import { ws } from '../services/websocket'

interface AppState {
  // 用户
  user: User | null
  token: string
  isLoggedIn: boolean

  // 房间
  rooms: Room[]
  currentRoom: RoomDetail | null

  // 游戏
  gameStart: GameStart | null
  gameMessages: GameMessage[]
  gameEnd: GameEnd | null

  // UI
  currentPage: 'lobby' | 'room' | 'game' | 'rank' | 'login' | 'register'

  // Actions
  setPage: (page: AppState['currentPage']) => void
  
  // Auth
  login: (username: string, password: string) => Promise<void>
  register: (username: string, password: string, displayName: string) => Promise<void>
  logout: () => void

  // Rooms
  fetchRooms: () => Promise<void>
  createRoom: (name: string, maxPlayers?: number) => Promise<string>
  joinRoom: (roomId: string) => Promise<void>
  leaveRoom: () => Promise<void>
  setReady: (ready: boolean) => Promise<void>
  startGame: () => Promise<void>

  // WebSocket
  initWebSocket: () => void

  // Game
  sendMessage: (content: string) => void
  clearGame: () => void
}

export const useStore = create<AppState>()(
  persist(
    (set, get) => ({
      // 初始状态
      user: null,
      token: '',
      isLoggedIn: false,
      rooms: [],
      currentRoom: null,
      gameStart: null,
      gameMessages: [],
      gameEnd: null,
      currentPage: 'lobby',

      setPage: (page) => set({ currentPage: page }),

      // Auth
      login: async (username, password) => {
        const response = await api.login(username, password)
        api.setToken(response.token)
        set({ user: response.user, token: response.token, isLoggedIn: true })
        get().initWebSocket()
      },

      register: async (username, password, displayName) => {
        const response = await api.register(username, password, displayName)
        api.setToken(response.token)
        set({ user: response.user, token: response.token, isLoggedIn: true })
        get().initWebSocket()
      },

      logout: () => {
        api.clearToken()
        ws.disconnect()
        set({ user: null, token: '', isLoggedIn: false, currentRoom: null })
      },

      // Rooms
      fetchRooms: async () => {
        try {
          const response = await api.listRooms()
          set({ rooms: response.rooms })
        } catch (e) {
          console.error('Failed to fetch rooms:', e)
        }
      },

      createRoom: async (name, maxPlayers = 4) => {
        const response = await api.createRoom(name, maxPlayers)
        await get().fetchRooms()
        return response.room_id
      },

      joinRoom: async (roomId) => {
        const user = get().user
        if (!user) return
        
        await api.joinRoom(roomId, user.id)
        const response = await api.getRoom(roomId)
        set({ currentRoom: response.room, currentPage: 'room' })
      },

      leaveRoom: async () => {
        const { currentRoom, user } = get()
        if (!currentRoom || !user) return
        
        await api.leaveRoom(currentRoom.id, user.id)
        set({ currentRoom: null, currentPage: 'lobby' })
        get().fetchRooms()
      },

      setReady: async (ready) => {
        const { currentRoom, user } = get()
        if (!currentRoom || !user) return
        
        await api.setReady(currentRoom.id, user.id, ready)
      },

      startGame: async () => {
        const { currentRoom } = get()
        if (!currentRoom) return
        
        await api.startGame(currentRoom.id)
      },

      // WebSocket
      initWebSocket: () => {
        const { user, token } = get()
        if (!user || !token) return

        ws.connect(user.id, token)

        // 房间更新
        ws.on('room_update', (payload) => {
          const data = payload as { room: RoomDetail }
          set({ currentRoom: data.room })
        })

        // 游戏开始
        ws.on('game_start', (payload) => {
          const data = payload as GameStart
          set({ 
            gameStart: data, 
            gameMessages: [], 
            gameEnd: null,
            currentPage: 'game' 
          })
        })

        // 游戏消息
        ws.on('message', (payload) => {
          const data = payload as GameMessage
          set(state => ({ 
            gameMessages: [...state.gameMessages, data] 
          }))
        })

        // 游戏结束
        ws.on('game_end', (payload) => {
          const data = payload as GameEnd
          set({ gameEnd: data })
        })
      },

      // Game
      sendMessage: (content) => {
        ws.sendMessage(content)
      },

      clearGame: () => {
        set({ 
          gameStart: null, 
          gameMessages: [], 
          gameEnd: null 
        })
      },
    }),
    {
      name: 'induce-master-storage',
      partialize: (state) => ({ 
        token: state.token, 
        user: state.user,
        isLoggedIn: state.isLoggedIn 
      }),
    }
  )
)
