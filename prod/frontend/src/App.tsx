import { useState } from 'react'

type Page = 'lobby' | 'room' | 'game' | 'rank'

function App() {
  const [currentPage, setCurrentPage] = useState<Page>('lobby')
  const [isLoggedIn, setIsLoggedIn] = useState(false)

  return (
    <div className="min-h-screen bg-gray-900 text-white">
      {/* 头部 */}
      <header className="bg-gray-800 p-4 flex justify-between items-center">
        <h1 className="text-2xl font-bold">🎮 诱导大师</h1>
        <div>
          {isLoggedIn ? (
            <button className="px-4 py-2 bg-red-600 rounded">退出</button>
          ) : (
            <button 
              onClick={() => setIsLoggedIn(true)}
              className="px-4 py-2 bg-blue-600 rounded"
            >
              登录
            </button>
          )}
        </div>
      </header>

      {/* 主内容 */}
      <main className="p-4">
        {currentPage === 'lobby' && (
          <Lobby onNavigate={setCurrentPage} isLoggedIn={isLoggedIn} />
        )}
        {currentPage === 'room' && (
          <Room onNavigate={setCurrentPage} />
        )}
        {currentPage === 'game' && (
          <Game onNavigate={setCurrentPage} />
        )}
        {currentPage === 'rank' && (
          <Ranking onNavigate={setCurrentPage} />
        )}
      </main>
    </div>
  )
}

// 大厅组件
function Lobby({ onNavigate, isLoggedIn }: { onNavigate: (p: Page) => void, isLoggedIn: boolean }) {
  const rooms = [
    { id: '1', name: '水果专场', players: 2, maxPlayers: 4, status: 'waiting' },
    { id: '2', name: '电影专场', players: 1, maxPlayers: 4, status: 'waiting' },
    { id: '3', name: '休闲娱乐', players: 3, maxPlayers: 4, status: 'waiting' },
  ]

  return (
    <div>
      <div className="flex gap-4 mb-4">
        <button 
          onClick={() => onNavigate('room')}
          disabled={!isLoggedIn}
          className="px-6 py-3 bg-green-600 rounded-lg disabled:opacity-50"
        >
          创建房间
        </button>
        <button className="px-6 py-3 bg-blue-600 rounded-lg">
          快速匹配
        </button>
        <button 
          onClick={() => onNavigate('rank')}
          className="px-6 py-3 bg-purple-600 rounded-lg"
        >
          排行榜
        </button>
      </div>

      <h2 className="text-xl mb-4">房间列表</h2>
      <div className="space-y-2">
        {rooms.map(room => (
          <div key={room.id} className="bg-gray-800 p-4 rounded-lg flex justify-between items-center">
            <div>
              <h3 className="font-bold">{room.name}</h3>
              <p className="text-gray-400">{room.players}/{room.maxPlayers} 人</p>
            </div>
            <button 
              onClick={() => onNavigate('room')}
              disabled={!isLoggedIn || room.players >= room.maxPlayers}
              className="px-4 py-2 bg-blue-600 rounded disabled:opacity-50"
            >
              加入
            </button>
          </div>
        ))}
      </div>
    </div>
  )
}

// 房间组件
function Room({ onNavigate }: { onNavigate: (p: Page) => void }) {
  return (
    <div>
      <h2 className="text-xl mb-4">房间 - 等待中...</h2>
      <div className="bg-gray-800 p-4 rounded-lg mb-4">
        <p>玩家1：已准备 ✅</p>
        <p>等待玩家加入...</p>
      </div>
      <button 
        onClick={() => onNavigate('game')}
        className="px-6 py-3 bg-green-600 rounded-lg"
      >
        开始游戏
      </button>
    </div>
  )
}

// 游戏组件
function Game({ onNavigate }: { onNavigate: (p: Page) => void }) {
  return (
    <div>
      <div className="text-center mb-4">
        <h2 className="text-xl">对战进行中</h2>
        <p className="text-gray-400">你的关键词：苹果 | 对手关键词：香蕉</p>
      </div>
      
      <div className="bg-gray-800 p-4 rounded-lg mb-4 h-64 overflow-y-auto">
        <div className="mb-2">[你] 今天天气不错</div>
        <div className="mb-2 text-blue-400">[对手] 是啊，你喜欢什么水果？</div>
        <div className="mb-2">[你] 我最喜欢那种红红的...</div>
        <div className="mb-2 text-blue-400">[对手] 是不是苹果？</div>
      </div>

      <div className="flex gap-2">
        <input 
          type="text" 
          placeholder="输入消息..." 
          className="flex-1 px-4 py-2 bg-gray-700 rounded"
        />
        <button className="px-6 py-2 bg-green-600 rounded">发送</button>
      </div>
    </div>
  )
}

// 排行榜组件
function Ranking({ onNavigate }: { onNavigate: (p: Page) => void }) {
  const ranking = [
    { rank: 1, name: '小明', score: 5000 },
    { rank: 2, name: '小红', score: 4800 },
    { rank: 3, name: '小刚', score: 4500 },
  ]

  return (
    <div>
      <h2 className="text-xl mb-4">排行榜</h2>
      <div className="space-y-2">
        {ranking.map(r => (
          <div key={r.rank} className="bg-gray-800 p-4 rounded-lg flex justify-between">
            <span className="font-bold">#{r.rank} {r.name}</span>
            <span>{r.score} 分</span>
          </div>
        ))}
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

export default App
