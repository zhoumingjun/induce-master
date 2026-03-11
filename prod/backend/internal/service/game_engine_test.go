package service

import (
	"testing"
)

// TestAssignWords 测试关键词分配
func TestAssignWords(t *testing.T) {
	players := []string{"player1", "player2", "player3"}
	engine := NewGameEngine("room1", players, 10, 60)
	
	engine.AssignWords()
	
	// 检查每个玩家都有关键词
	if len(engine.Words) != len(players) {
		t.Errorf("Expected %d words, got %d", len(players), len(engine.Words))
	}
	
	// 检查每个玩家的关键词不同
	wordSet := make(map[string]bool)
	for _, word := range engine.Words {
		if wordSet[word] {
			t.Errorf("Duplicate word assigned: %s", word)
		}
		wordSet[word] = true
	}
}

// TestProcessMessage_KeywordDetection 测试关键词检测
func TestProcessMessage_KeywordDetection(t *testing.T) {
	players := []string{"player1", "player2"}
	engine := NewGameEngine("room1", players, 10, 60)
	engine.Start()
	
	// 给 player1 分配"苹果", player2 分配"香蕉"
	engine.Words["player1"] = "苹果"
	engine.Words["player2"] = "香蕉"
	
	// player1 说出 player2 的关键词"香蕉" - 触发关键词
	msg := engine.ProcessMessage("player1", "用户1", "我喜欢吃香蕉")
	if !msg.IsKeyword {
		t.Error("Expected keyword detection for '香蕉' (opponent's word)")
	}
	
	// player1 说出不包含关键词的内容 - 不触发
	msg2 := engine.ProcessMessage("player1", "用户1", "今天天气很好")
	if msg2.IsKeyword {
		t.Error("Should not trigger keyword for normal message")
	}
}

// TestNextTurn 测试回合切换
func TestNextTurn(t *testing.T) {
	players := []string{"player1", "player2", "player3"}
	engine := NewGameEngine("room1", players, 10, 60)
	engine.Start()
	
	// 初始玩家是 player1
	if engine.GetCurrentPlayer() != "player1" {
		t.Errorf("Expected first player to be player1, got %s", engine.GetCurrentPlayer())
	}
	
	// 下一回合
	engine.NextTurn()
	if engine.GetCurrentPlayer() != "player2" {
		t.Errorf("Expected current player to be player2, got %s", engine.GetCurrentPlayer())
	}
	
	// 再下一回合
	engine.NextTurn()
	if engine.GetCurrentPlayer() != "player3" {
		t.Errorf("Expected current player to be player3, got %s", engine.GetCurrentPlayer())
	}
	
	// 回到第一个玩家，回合数+1
	engine.NextTurn()
	if engine.GetCurrentPlayer() != "player1" {
		t.Errorf("Expected current player to be player1, got %s", engine.GetCurrentPlayer())
	}
	if engine.Round != 2 {
		t.Errorf("Expected round to be 2, got %d", engine.Round)
	}
}

// TestGameStatus 游戏状态测试
func TestGameStatus(t *testing.T) {
	players := []string{"player1", "player2"}
	engine := NewGameEngine("room1", players, 2, 60)
	
	// 初始状态是 Waiting
	if engine.Status != GameStatusWaiting {
		t.Errorf("Expected status to be Waiting, got %d", engine.Status)
	}
	
	// 开始游戏
	engine.Start()
	if engine.Status != GameStatusPlaying {
		t.Errorf("Expected status to be Playing, got %d", engine.Status)
	}
	
	// 通过 NextTurn 达到最大回合，结束游戏
	for i := 0; i < 5; i++ {
		engine.NextTurn()
	}
	
	if !engine.IsFinished() {
		t.Error("Expected game to be finished after max rounds")
	}
}

// TestGetWinner 测试获胜者判断
func TestGetWinner(t *testing.T) {
	players := []string{"player1", "player2"}
	engine := NewGameEngine("room1", players, 10, 60)
	
	engine.Scores["player1"] = 100
	engine.Scores["player2"] = 80
	
	winner := engine.GetWinner()
	if winner != "player1" {
		t.Errorf("Expected winner to be player1, got %s", winner)
	}
}

// TestContainsKeyword 测试关键词匹配
func TestContainsKeyword(t *testing.T) {
	tests := []struct {
		content  string
		keyword  string
		expected bool
	}{
		{"我喜欢苹果", "苹果", true},
		{"苹果很好吃", "苹果", true},
		{"这是苹果吗", "苹果", true},
		{"我不喜欢香蕉", "苹果", false},
		{"", "苹果", false},
		{"苹果", "", false},
	}
	
	for _, tt := range tests {
		result := containsKeyword(tt.content, tt.keyword)
		if result != tt.expected {
			t.Errorf("containsKeyword(%q, %q) = %v, expected %v", 
				tt.content, tt.keyword, result, tt.expected)
		}
	}
}
