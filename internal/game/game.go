package game

import (
	"sync"
	"time"

	"github.com/gorilla/websocket"
)

type Player struct {
	ID      string
	Name    string
	Score   int
	Answers map[int]string
	WSConn  *websocket.Conn
}

type Question struct {
	ID        int
	Text      string
	Options   []string
	Correct   string
	TimeLimit time.Duration
}

type GameState struct {
	ID              string
	Players         map[string]*Player
	CurrentQuestion *Question
	Questions       []Question
	IsActive        bool
	Round           int
	mu              sync.RWMutex
}

type GameManager struct {
	Games map[string]*GameState
	mu    sync.RWMutex
}

func NewGameManager() *GameManager {
	return &GameManager{
		Games: make(map[string]*GameState),
	}
}
