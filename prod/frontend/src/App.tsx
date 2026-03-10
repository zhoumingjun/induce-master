import { useState, useEffect, useCallback } from 'react'
import { authAPI, roomAPI, rankingAPI } from './services/api'

type Page = 'login' | 'lobby' | 'room' | 'game' | 'rank'

interface Room {
  id: string
  name: string
  owner_id: string
  status: number
  player_count: number
}

interface Player {
  user_id: string
  username: string
  ready: boolean
  score?: number
}

function App() {
  const [currentPage, setCurrentPage] = useState<Page>('lobby')
  const [isLoggedIn, setIsLoggedIn] = useState(false)
  const [username, setUsername] = useState('')
  const [password, setPassword] = useState('')
  const [loading, setLoading] = useState(false)
  const [error, setError] = useState('')

  // 检查登录状态
  useEffect(() => {
    if (authAPI.isLoggedIn()) {
      setIsLoggedIn(true)
      setCurrentPage('lobby')
    } else {
      setCurrentPage('login')
    }
  }, [])

  const handleLogin = async () => {
    if (!username || !password) {
      setError('请输入用户名和密码')
      return
    }
    setLoading(true)
    setError('')
    try {
      await authAPI.login(username, password)
      setIsLoggedIn(true)
      setCurrentPage('lobby')
    } catch (err: any) {
      setError(err.message || '登录失败')
    } finally {
      setLoading(false)
    }
  }

  const handleRegister = async () => {
    if (!username || !password) {
      setError('请输入用户名和密码')
      return
    }
    setLoading(true)
    setError('')
    try {
      await authAPI.register(username, password)
      setIsLoggedIn(true)
      setCurrentPage('lobby')
    } catch (err: any) {
      setError(err.message || '注册失败')
    } finally {
      setLoading(false)
    }
  }

  const handleLogout = () => {
    authAPI.logout()
    setIsLoggedIn(false)
    setCurrentPage('login')
    setUsername('')
    setPassword('')
  }

  return (
    <div className="min-h-screen bg-gray-900 text-white">
      {/* 头部 */}
      <header className="bg-gray-800 p-4 flex justify-between items-center">
        <h1 className="text-2xl font-bold">🎮 诱导大师</h1>
        <div>
          {isLoggedIn ? (
            <div className="flex items-center gap-4">
              <span className="text-gray-400">{authAPI.getUsername()}</span>
              <button 
                onClick={handleLogout}
                className="px-4 py-2 bg-red-600 rounded hover:bg-red-700"
              >
                退出
              </button>
            </div>
          ) : (
            <span className="text-gray-400">未登录</span>
          )}
        </div>
      </header>

      {/* 主内容 */}
      <main className="p-4">
        {currentPage === 'login' && (
          <LoginPage 
            username={username}
            password={password}
            loading={loading}
            error={error}
            onUsernameChange={setUsername}
            onPasswordChange={setPassword}
            onLogin={handleLogin}
            onRegister={handleRegister}
          />
        )}
        {currentPage === 'lobby' && (
          <LobbyPage onNavigate={setCurrentPage} />
        )}
        {currentPage === 'room' && (
          <RoomPage onNavigate={setCurrentPage} />
        )}
        {currentPage === 'game' && (
          <GamePage onNavigate={setCurrentPage} />
        )}
        {currentPage === 'rank' && (
          <RankingPage onNavigate={setCurrentPage} />
        )}
      </main>
    </div>
  )
}

// 登录页面
function LoginPage({ 
  username, 
  password, 
  loading, 
  error,
  onUsernameChange, 
  onPasswordChange, 
  onLogin, 
  onRegister 
}: {
  username: string
  password: string
  loading: boolean
  error: string
  onUsernameChange: (v: string) => void
  onPasswordChange: (v: string) => void
  onLogin: () => void
  onRegister: () => void
}) {
  return (
    <div className="max-w-md mx-auto mt-20">
      <h2 className="text-3xl font-bold text-center mb-8">欢迎来到诱导大师</h2>
      
      {error && (
        <div className="bg-red-600 text-white p-3 rounded mb-4">{error}</div>
      )}
      
      <div className="space-y-4">
        <div>
          <label className="block text-gray-400 mb-2">用户名</label>
          <input
            type="text"
            value={username}
            onChange={(e) => onUsernameChange(e.target.value)}
            className="w-full px-4 py-2 bg-gray-800 rounded border border-gray-700 focus:border-blue-500 outline-none"
            placeholder="请输入用户名"
          />
        </div>
        
        <div>
          <label className="block text-gray-400 mb-2">密码</label>
          <input
            type="password"
            value={password}
            onChange={(e) => onPasswordChange(e.target.value)}
            className="w-full px-4 py-2 bg-gray-800 rounded border border-gray-700 focus:border-blue-500 outline-none"
            placeholder="请输入密码"
            onKeyDown={(e) => e.key === 'Enter' && onLogin()}
          />
        </div>
        
        <div className="flex gap-4">
          <button
            onClick={onLogin}
            disabled={loading}
            className="flex-1 py-3 bg-blue-600 rounded hover:bg-blue-700 disabled:opacity-50"
          >
            {loading ? '登录中...' : '登录'}
          </button>
          <button
            onClick={onRegister}
            disabled={loading}
            className="flex-1 py-3 bg-green-600 rounded hover:bg-green-700 disabled:opacity-50"
          >
            {loading ? '注册中...' : '注册'}
          </button>
        </div>
      </div>
    </div>
  )
}

// 大厅页面
function LobbyPage({ onNavigate }: { onNavigate: (p: Page) => void }) {
  const [rooms, setRooms] = useState<Room[]>([])
  const [loading, setLoading] = useState(false)
  const [error, setError] = useState('')

  const loadRooms = useCallback(async () => {
    setLoading(true)
    setError('')
    try {
      const data = await roomAPI.list()
      setRooms(data)
    } catch (err: any) {
      setError(err.message || '加载房间失败')
    } finally {
      setLoading(false)
    }
  }, [])

  useEffect(() => {
    loadRooms()
    const interval = setInterval(loadRooms, 5000) // 每5秒刷新
    return () => clearInterval(interval)
  }, [loadRooms])

  const handleCreateRoom = async () => {
    const name = prompt('请输入房间名称', '房间' + Date.now())
    if (!name) return
    
    try {
      await roomAPI.create(name)
      onNavigate('room')
    } catch (err: any) {
      alert(err.message || '创建房间失败')
    }
  }

  const handleJoinRoom = async (roomId: string) => {
    try {
      await roomAPI.join(roomId)
      onNavigate('room')
    } catch (err: any) {
      alert(err.message || '加入房间失败')
    }
  }

  return (
    <div>
      <div className="flex gap-4 mb-6">
        <button 
          onClick={handleCreateRoom}
          className="px-6 py-3 bg-green-600 rounded-lg hover:bg-green-700"
        >
          创建房间
        </button>
        <button className="px-6 py-3 bg-blue-600 rounded-lg hover:bg-blue-700">
          快速匹配
        </button>
        <button 
          onClick={() => onNavigate('rank')}
          className="px-6 py-3 bg-purple-600 rounded-lg hover:bg-purple-700"
        >
          排行榜
        </button>
        <button 
          onClick={loadRooms}
          className="px-6 py-3 bg-gray-600 rounded-lg hover:bg-gray-700"
        >
          刷新
        </button>
      </div>

      {error && (
        <div className="bg-red-600 text-white p-3 rounded mb-4">{error}</div>
      )}

      <h2 className="text-xl mb-4">房间列表 {loading && <span className="text-gray-400 text-sm">(刷新中...)</span>}</h2>
      
      {rooms.length === 0 ? (
        <div className="text-gray-400 text-center py-8">暂无房间</div>
      ) : (
        <div className="space-y-2">
          {rooms.map(room => (
            <div key={room.id} className="bg-gray-800 p-4 rounded-lg flex justify-between items-center">
              <div>
                <h3 className="font-bold">{room.name}</h3>
                <p className="text-gray-400">{room.player_count}/4 人</p>
              </div>
              <button 
                onClick={() => handleJoinRoom(room.id)}
                disabled={room.player_count >= 4}
                className="px-4 py-2 bg-blue-600 rounded disabled:opacity-50 hover:bg-blue-700"
              >
                加入
              </button>
            </div>
          ))}
        </div>
      )}
    </div>
  )
}

// 房间页面
function RoomPage({ onNavigate }: { onNavigate: (p: Page) => void }) {
  const [roomInfo, setRoomInfo] = useState<any>(null)
  const [loading, setLoading] = useState(true)
  const [myReady, setMyReady] = useState(false)

  // TODO: 从房间列表传入或从 URL 获取
  const roomId = 'room-1'

  const loadRoom = useCallback(async () => {
    try {
      const data = await roomAPI.get(roomId)
      setRoomInfo(data)
    } catch (err: any) {
      console.error('加载房间失败:', err)
    } finally {
      setLoading(false)
    }
  }, [roomId])

  useEffect(() => {
    loadRoom()
    const interval = setInterval(loadRoom, 3000)
    return () => clearInterval(interval)
  }, [loadRoom])

  const handleReady = async () => {
    try {
      await roomAPI.ready(roomId, !myReady)
      setMyReady(!myReady)
    } catch (err: any) {
      alert(err.message || '操作失败')
    }
  }

  const handleStart = async () => {
    try {
      await roomAPI.start(roomId)
      onNavigate('game')
    } catch (err: any) {
      alert(err.message || '开始游戏失败')
    }
  }

  const handleLeave = async () => {
    try {
      await roomAPI.leave(roomId)
      onNavigate('lobby')
    } catch (err: any) {
      alert(err.message || '离开房间失败')
    }
  }

  if (loading) {
    return <div className="text-center py-8">加载中...</div>
  }

  const players = roomInfo?.players || []

  return (
    <div>
      <div className="flex justify-between items-center mb-4">
        <h2 className="text-xl">房间 - {roomInfo?.name || '房间'}</h2>
        <button 
          onClick={handleLeave}
          className="px-4 py-2 bg-red-600 rounded hover:bg-red-700"
        >
          离开
        </button>
      </div>

      <div className="bg-gray-800 p-4 rounded-lg mb-4">
        <h3 className="font-bold mb-2">玩家列表</h3>
        {players.length === 0 ? (
          <p className="text-gray-400">等待玩家加入...</p>
        ) : (
          players.map((p: Player, i: number) => (
            <p key={i} className="mb-1">
              {p.username}：{p.ready ? '已准备 ✅' : '等待准备...'}
            </p>
          ))
        )}
      </div>

      <div className="flex gap-4">
        <button 
          onClick={handleReady}
          className={`px-6 py-3 rounded-lg ${myReady ? 'bg-yellow-600' : 'bg-green-600'}`}
        >
          {myReady ? '取消准备' : '准备'}
        </button>
        <button 
          onClick={handleStart}
          className="px-6 py-3 bg-blue-600 rounded-lg hover:bg-blue-700"
        >
          开始游戏
        </button>
      </div>
    </div>
  )
}

// 游戏页面
function GamePage({ onNavigate }: { onNavigate: (p: Page) => void }) {
  const [messages, setMessages] = useState<any[]>([])
  const [input, setInput] = useState('')

  // TODO: 完善游戏逻辑

  const handleSend = () => {
    if (!input.trim()) return
    
    const newMsg = {
      round: 1,
      user_id: authAPI.getUserID(),
      username: authAPI.getUsername(),
      content: input,
      time: new Date(),
    }
    
    setMessages([...messages, newMsg])
    setInput('')
  }

  return (
    <div>
      <div className="text-center mb-4">
        <h2 className="text-xl">对战进行中</h2>
        <p className="text-gray-400">你的关键词：苹果 | 对手关键词：香蕉</p>
      </div>
      
      <div className="bg-gray-800 p-4 rounded-lg mb-4 h-64 overflow-y-auto">
        {messages.length === 0 ? (
          <p className="text-gray-400">游戏即将开始...</p>
        ) : (
          messages.map((msg, i) => (
            <div key={i} className="mb-2">
              [{msg.username}] {msg.content}
            </div>
          ))
        )}
      </div>

      <div className="flex gap-2">
        <input 
          type="text" 
          placeholder="输入消息..."
          value={input}
          onChange={(e) => setInput(e.target.value)}
          onKeyDown={(e) => e.key === 'Enter' && handleSend()}
          className="flex-1 px-4 py-2 bg-gray-700 rounded"
        />
        <button 
          onClick={handleSend}
          className="px-6 py-2 bg-green-600 rounded hover:bg-green-700"
        >
          发送
        </button>
      </div>

      <button 
        onClick={() => onNavigate('lobby')}
        className="mt-4 px-6 py-2 bg-gray-600 rounded"
      >
        返回大厅
      </button>
    </div>
  )
}

// 排行榜页面
function RankingPage({ onNavigate }: { onNavigate: (p: Page) => void }) {
  const [ranking, setRanking] = useState<any[]>([])
  const [loading, setLoading] = useState(true)

  useEffect(() => {
    rankingAPI.get()
      .then(data => setRanking(data))
      .catch(console.error)
      .finally(() => setLoading(false))
  }, [])

  if (loading) {
    return <div className="text-center py-8">加载中...</div>
  }

  return (
    <div>
      <h2 className="text-xl mb-4">排行榜</h2>
      
      {ranking.length === 0 ? (
        <p className="text-gray-400 text-center py-8">暂无数据</p>
      ) : (
        <div className="space-y-2">
          {ranking.map((r, i) => (
            <div key={i} className="bg-gray-800 p-4 rounded-lg flex justify-between">
              <span className="font-bold">#{i + 1} {r.username}</span>
              <span>{r.score} 分</span>
            </div>
          ))}
        </div>
      )}
      
      <button 
        onClick={() => onNavigate('lobby')}
        className="mt-4 px-6 py-2 bg-gray-600 rounded hover:bg-gray-700"
      >
        返回大厅
      </button>
    </div>
  )
}

export default App
