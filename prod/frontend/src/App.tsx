import { useEffect } from 'react'
import { useStore } from './stores'
import { api } from './api'

function App() {
  const { 
    currentPage, 
    setPage, 
    isLoggedIn, 
    user,
    logout,
    fetchRooms,
    rooms,
    currentRoom,
    gameStart,
    gameMessages,
    gameEnd,
    joinRoom,
    createRoom,
    leaveRoom,
    setReady,
    startGame,
    sendMessage,
    clearGame,
  } = useStore()

  // 初始加载 - 检查 URL token
  useEffect(() => {
    const params = new URLSearchParams(window.location.search)
    const token = params.get('token')
    if (token && !isLoggedIn) {
      // 从 URL token 恢复登录状态
      api.setToken(token)
      api.me().then(res => {
        useStore.setState({ 
          token, 
          user: res.user, 
          isLoggedIn: true 
        })
        // 初始化 WebSocket
        useStore.getState().initWebSocket()
      }).catch(() => {})
    }
  }, [])

  // 初始加载
  useEffect(() => {
    if (isLoggedIn) {
      fetchRooms()
    }
  }, [isLoggedIn])

  // 未登录跳转登录页
  useEffect(() => {
    if (!isLoggedIn && currentPage !== 'login' && currentPage !== 'register') {
      setPage('login')
    }
  }, [isLoggedIn])

  if (!isLoggedIn) {
    return currentPage === 'register' ? (
      <Register onSwitch={() => setPage('login')} />
    ) : (
      <Login onSwitch={() => setPage('register')} />
    )
  }

  return (
    <div className="min-h-screen bg-gray-900 text-white">
      {/* 头部 */}
      <header className="bg-gray-800 p-4 flex justify-between items-center">
        <h1 className="text-2xl font-bold">🎮 诱导大师</h1>
        <div className="flex items-center gap-4">
          <span>{user?.display_name}</span>
          <span className="text-yellow-400">🏆 {user?.rank || 0}</span>
          <button 
            onClick={logout}
            className="px-4 py-2 bg-red-600 rounded hover:bg-red-700"
          >
            退出
          </button>
        </div>
      </header>

      {/* 主内容 */}
      <main className="p-4">
        {currentPage === 'lobby' && (
          <Lobby 
            rooms={rooms}
            onJoinRoom={joinRoom}
            onCreateRoom={createRoom}
            onNavigate={setPage}
          />
        )}
        {currentPage === 'room' && currentRoom && (
          <Room 
            room={currentRoom}
            userId={user?.id || ''}
            onLeave={leaveRoom}
            onReady={setReady}
            onStart={startGame}
          />
        )}
        {currentPage === 'game' && gameStart && (
          <Game 
            game={gameStart}
            messages={gameMessages}
            gameEnd={gameEnd}
            userId={user?.id || ''}
            onSendMessage={sendMessage}
            onBack={() => {
              clearGame()
              setPage('lobby')
              fetchRooms()
            }}
          />
        )}
        {currentPage === 'rank' && (
          <Ranking onNavigate={setPage} />
        )}
      </main>
    </div>
  )
}

// 登录组件
function Login({ onSwitch }: { onSwitch: () => void }) {
  const login = useStore(s => s.login)
  const [username, setUsername] = useState('')
  const [password, setPassword] = useState('')
  const [error, setError] = useState('')
  const [loading, setLoading] = useState(false)

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault()
    setError('')
    setLoading(true)
    try {
      await login(username, password)
    } catch (e: any) {
      setError(e.message || '登录失败')
    }
    setLoading(false)
  }

  return (
    <div className="min-h-screen bg-gray-900 flex items-center justify-center">
      <form onSubmit={handleSubmit} className="bg-gray-800 p-8 rounded-lg w-96">
        <h2 className="text-2xl font-bold mb-6 text-center">登录</h2>
        {error && <p className="text-red-500 mb-4">{error}</p>}
        <input
          type="text"
          placeholder="用户名"
          value={username}
          onChange={e => setUsername(e.target.value)}
          className="w-full px-4 py-2 mb-4 bg-gray-700 rounded"
          required
        />
        <input
          type="password"
          placeholder="密码"
          value={password}
          onChange={e => setPassword(e.target.value)}
          className="w-full px-4 py-2 mb-6 bg-gray-700 rounded"
          required
        />
        <button 
          type="submit"
          disabled={loading}
          className="w-full py-2 bg-blue-600 rounded hover:bg-blue-700 disabled:opacity-50"
        >
          {loading ? '登录中...' : '登录'}
        </button>
        <p className="mt-4 text-center text-gray-400">
          没有账号？<button type="button" onClick={onSwitch} className="text-blue-400">注册</button>
        </p>
      </form>
    </div>
  )
}

// 注册组件
function Register({ onSwitch }: { onSwitch: () => void }) {
  const register = useStore(s => s.register)
  const [username, setUsername] = useState('')
  const [password, setPassword] = useState('')
  const [displayName, setDisplayName] = useState('')
  const [error, setError] = useState('')
  const [loading, setLoading] = useState(false)

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault()
    setError('')
    setLoading(true)
    try {
      await register(username, password, displayName)
    } catch (e: any) {
      setError(e.message || '注册失败')
    }
    setLoading(false)
  }

  return (
    <div className="min-h-screen bg-gray-900 flex items-center justify-center">
      <form onSubmit={handleSubmit} className="bg-gray-800 p-8 rounded-lg w-96">
        <h2 className="text-2xl font-bold mb-6 text-center">注册</h2>
        {error && <p className="text-red-500 mb-4">{error}</p>}
        <input
          type="text"
          placeholder="用户名"
          value={username}
          onChange={e => setUsername(e.target.value)}
          className="w-full px-4 py-2 mb-4 bg-gray-700 rounded"
          required
        />
        <input
          type="text"
          placeholder="显示名称"
          value={displayName}
          onChange={e => setDisplayName(e.target.value)}
          className="w-full px-4 py-2 mb-4 bg-gray-700 rounded"
          required
        />
        <input
          type="password"
          placeholder="密码"
          value={password}
          onChange={e => setPassword(e.target.value)}
          className="w-full px-4 py-2 mb-6 bg-gray-700 rounded"
          required
        />
        <button 
          type="submit"
          disabled={loading}
          className="w-full py-2 bg-green-600 rounded hover:bg-green-700 disabled:opacity-50"
        >
          {loading ? '注册中...' : '注册'}
        </button>
        <p className="mt-4 text-center text-gray-400">
          已有账号？<button type="button" onClick={onSwitch} className="text-blue-400">登录</button>
        </p>
      </form>
    </div>
  )
}

// 大厅组件
function Lobby({ 
  rooms, 
  onJoinRoom, 
  onCreateRoom,
  onNavigate,
}: { 
  rooms: any[]
  onJoinRoom: (id: string) => Promise<void>
  onCreateRoom: (name: string) => Promise<string>
  onNavigate: (p: any) => void
}) {
  const [showCreate, setShowCreate] = useState(false)
  const [roomName, setRoomName] = useState('')

  const handleCreate = async () => {
    if (!roomName.trim()) return
    const roomId = await onCreateRoom(roomName)
    await onJoinRoom(roomId)
  }

  return (
    <div>
      <div className="flex gap-4 mb-4">
        <button 
          onClick={() => setShowCreate(true)}
          className="px-6 py-3 bg-green-600 rounded-lg hover:bg-green-700"
        >
          创建房间
        </button>
        <button 
          onClick={() => onNavigate('rank')}
          className="px-6 py-3 bg-purple-600 rounded-lg hover:bg-purple-700"
        >
          排行榜
        </button>
      </div>

      {/* 创建房间弹窗 */}
      {showCreate && (
        <div className="fixed inset-0 bg-black/50 flex items-center justify-center">
          <div className="bg-gray-800 p-6 rounded-lg">
            <h3 className="text-xl mb-4">创建房间</h3>
            <input
              type="text"
              placeholder="房间名"
              value={roomName}
              onChange={e => setRoomName(e.target.value)}
              className="w-full px-4 py-2 mb-4 bg-gray-700 rounded"
            />
            <div className="flex gap-2">
              <button 
                onClick={handleCreate}
                className="px-4 py-2 bg-green-600 rounded"
              >
                创建
              </button>
              <button 
                onClick={() => setShowCreate(false)}
                className="px-4 py-2 bg-gray-600 rounded"
              >
                取消
              </button>
            </div>
          </div>
        </div>
      )}

      <h2 className="text-xl mb-4">房间列表</h2>
      {rooms.length === 0 ? (
        <p className="text-gray-400">暂无房间</p>
      ) : (
        <div className="space-y-2">
          {rooms.map(room => (
            <div key={room.id} className="bg-gray-800 p-4 rounded-lg flex justify-between items-center">
              <div>
                <h3 className="font-bold">{room.name}</h3>
                <p className="text-gray-400">
                  {room.player_count}/{room.max_players || 4} 人
                  {room.status === 1 ? ' 🔴 进行中' : ' 🟢 等待中'}
                </p>
              </div>
              <button 
                onClick={() => onJoinRoom(room.id)}
                disabled={room.status === 1 || room.player_count >= (room.max_players || 4)}
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

// 房间组件
function Room({ 
  room, 
  userId,
  onLeave, 
  onReady, 
  onStart,
}: { 
  room: any
  userId: string
  onLeave: () => Promise<void>
  onReady: (ready: boolean) => Promise<void>
  onStart: () => Promise<void>
}) {
  const [myReady, setMyReady] = useState(false)
  const isOwner = room.owner_id === userId

  const allReady = room.players?.every((p: any) => p.ready) && room.players?.length >= 2

  const handleReady = async () => {
    const newReady = !myReady
    setMyReady(newReady)
    await onReady(newReady)
  }

  return (
    <div>
      <div className="flex justify-between items-center mb-4">
        <h2 className="text-xl">房间: {room.name}</h2>
        <button 
          onClick={onLeave}
          className="px-4 py-2 bg-red-600 rounded hover:bg-red-700"
        >
          退出房间
        </button>
      </div>

      <div className="bg-gray-800 p-4 rounded-lg mb-4">
        <h3 className="font-bold mb-2">玩家列表</h3>
        {room.players?.length === 0 ? (
          <p className="text-gray-400">等待玩家加入...</p>
        ) : (
          room.players?.map((p: any) => (
            <div key={p.id} className="flex justify-between py-2">
              <span>
                {p.username || p.user_id} {p.user_id === userId && '(你)'} {p.user_id === room.owner_id && '👑'}
              </span>
              <span className={p.ready ? 'text-green-400' : 'text-gray-400'}>
                {p.ready ? '✅ 已准备' : '⏳ 未准备'}
              </span>
            </div>
          ))
        )}
      </div>

      <div className="flex gap-4">
        <button 
          onClick={handleReady}
          className={`px-6 py-3 rounded-lg ${myReady ? 'bg-yellow-600' : 'bg-blue-600'} hover:opacity-90`}
        >
          {myReady ? '取消准备' : '准备'}
        </button>
        {isOwner && (
          <button 
            onClick={onStart}
            disabled={!allReady}
            className="px-6 py-3 bg-green-600 rounded-lg disabled:opacity-50 hover:bg-green-700"
          >
            开始游戏
          </button>
        )}
      </div>

      {!allReady && room.players?.length >= 2 && (
        <p className="mt-4 text-yellow-400">等待所有玩家准备...</p>
      )}
      {room.players?.length < 2 && (
        <p className="mt-4 text-gray-400">等待更多玩家加入...</p>
      )}
    </div>
  )
}

// 游戏组件
function Game({ 
  game,
  messages,
  gameEnd,
  userId,
  onSendMessage,
  onBack,
}: { 
  game: any
  messages: any[]
  gameEnd: any | null
  userId: string
  onSendMessage: (content: string) => void
  onBack: () => void
}) {
  const [input, setInput] = useState('')
  const messagesEndRef = useRef<HTMLDivElement>(null)

  useEffect(() => {
    messagesEndRef.current?.scrollIntoView({ behavior: 'smooth' })
  }, [messages])

  const handleSend = (e: React.FormEvent) => {
    e.preventDefault()
    if (!input.trim()) return
    onSendMessage(input)
    setInput('')
  }

  const isWinner = gameEnd?.winner_id === userId

  return (
    <div>
      {/* 游戏信息 */}
      <div className="bg-gray-800 p-4 rounded-lg mb-4 flex justify-between items-center">
        <div>
          <p className="text-lg">你的关键词: <span className="text-yellow-400 font-bold">{game.your_word}</span></p>
          <p className="text-gray-400">对手关键词: {game.opponent?.user_id ? '***' : '等待中'}</p>
        </div>
        <div className="text-right">
          <p>第 {game.round} / {game.max_rounds || 10} 轮</p>
          <p className="text-gray-400">限时 {game.time_limit || 60}秒</p>
        </div>
      </div>

      {/* 游戏结束弹窗 */}
      {gameEnd && (
        <div className="fixed inset-0 bg-black/50 flex items-center justify-center">
          <div className="bg-gray-800 p-8 rounded-lg text-center">
            <h2 className="text-3xl mb-4">
              {isWinner ? '🎉 你赢了！' : '😢 你输了'}
            </h2>
            <p className="mb-4">
              你的词: {gameEnd.words?.[userId]} | 对手词: {gameEnd.words?.[game.opponent?.user_id]}
            </p>
            <button 
              onClick={onBack}
              className="px-6 py-3 bg-blue-600 rounded-lg hover:bg-blue-700"
            >
              返回大厅
            </button>
          </div>
        </div>
      )}

      {/* 消息区域 */}
      <div className="bg-gray-800 p-4 rounded-lg mb-4 h-96 overflow-y-auto">
        {messages.length === 0 ? (
          <p className="text-gray-400 text-center">游戏开始！诱导对手说出关键词吧！</p>
        ) : (
          messages.map((msg, i) => (
            <div 
              key={i} 
              className={`mb-2 ${msg.user_id === userId ? 'text-right' : 'text-left'}`}
            >
              <span className={msg.user_id === userId ? 'text-green-400' : 'text-blue-400'}>
                [{msg.user_id === userId ? '你' : '对手'}]
              </span>
              {' '}
              {msg.content}
              {msg.is_keyword && (
                <span className="text-red-500 ml-2">🚨 说出关键词！</span>
              )}
            </div>
          ))
        )}
        <div ref={messagesEndRef} />
      </div>

      {/* 输入框 */}
      <form onSubmit={handleSend} className="flex gap-2">
        <input
          type="text"
          value={input}
          onChange={e => setInput(e.target.value)}
          placeholder="输入消息..."
          disabled={!!gameEnd}
          className="flex-1 px-4 py-2 bg-gray-700 rounded disabled:opacity-50"
        />
        <button 
          type="submit"
          disabled={!!gameEnd}
          className="px-6 py-2 bg-green-600 rounded disabled:opacity-50 hover:bg-green-700"
        >
          发送
        </button>
      </form>
    </div>
  )
}

// 排行榜组件
function Ranking({ onNavigate }: { onNavigate: (p: any) => void }) {
  const [ranking, setRanking] = useState<any[]>([])
  const [loading, setLoading] = useState(true)

  useEffect(() => {
    api.getRanking()
      .then(res => setRanking(res.ranking || []))
      .catch(console.error)
      .finally(() => setLoading(false))
  }, [])

  return (
    <div>
      <h2 className="text-xl mb-4">排行榜</h2>
      {loading ? (
        <p className="text-gray-400">加载中...</p>
      ) : ranking.length === 0 ? (
        <p className="text-gray-400">暂无数据</p>
      ) : (
        <div className="space-y-2">
          {ranking.map((r, i) => (
            <div key={r.user_id} className="bg-gray-800 p-4 rounded-lg flex justify-between items-center">
              <span className="font-bold">
                #{i + 1} {r.username}
              </span>
              <span className="text-yellow-400">{r.rank} 分</span>
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

// 引入 useState 和 useRef
import { useState, useRef } from 'react'

export default App
