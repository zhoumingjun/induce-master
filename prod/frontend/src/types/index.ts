// 类型定义

export interface User {
  id: string
  username: string
  display_name: string
  avatar_url: string
  rank: number
}

export interface Room {
  id: string
  name: string
  owner_id: string
  status: number
  player_count: number
  max_players?: number
}

export interface RoomDetail extends Room {
  players: RoomPlayer[]
}

export interface RoomPlayer {
  id: string
  user_id: string
  username?: string
  ready: boolean
  score: number
}

export interface GameMessage {
  round: number
  user_id: string
  username: string
  content: string
  time: string
  is_keyword?: boolean
}

export interface GameStart {
  game_id: string
  round: number
  time_limit: number
  your_word: string
  opponent: {
    user_id: string
    username?: string
  }
}

export interface GameEnd {
  game_id: string
  winner_id: string
  scores: Record<string, number>
  words: Record<string, string>
}

// API 响应
export interface ApiResponse<T> {
  data?: T
  error?: string
}

export interface LoginResponse {
  token: string
  user: User
}

export interface RoomListResponse {
  rooms: Room[]
}

export interface RankingItem {
  user_id: string
  username: string
  rank: number
}
