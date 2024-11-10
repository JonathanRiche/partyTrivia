package game

import (
	"fmt"
	// "github.com/gorilla/websocket"
	"richetechguy/internal/db"
	"richetechguy/internal/types"
	"sync"
	"time"
)

type GameManager struct {
	Games map[string]*types.GameState // Change from 'games' to 'Games'
	mu    sync.RWMutex
	Db    *db.DB
}

// StartGame starts a specific game
func (gm *GameManager) StartGame(gameID string, qm *QuestionManager) error {
	game, err := gm.GetGame(gameID)
	if err != nil {
		return err
	}
	if qm.GetQuestions() == nil {
		return err
	}

	game.Questions = qm.GetQuestions()
	return game.StartGame()
}
func (gm *GameManager) SelectGame(gameID string) (*types.GameState, error) {
	game, err := gm.GetGame(gameID)
	if err != nil {
		return nil, err
	}
	game.StartGame()
	return game, nil
}

// AddPlayer adds a player to a game
func (gm *GameManager) AddPlayer(gameID string, playerName string) (string, error) {
	game, err := gm.GetGame(gameID)
	if err != nil {
		return "", err
	}
	game.Mu.Lock()
	defer game.Mu.Unlock()

	if game.IsActive {
		return "", fmt.Errorf("cannot join active game")
	}

	playerID := fmt.Sprintf("player_%d", len(game.Players)+1)
	game.Players[playerID] = &types.Player{
		ID:      playerID,
		Name:    playerName,
		Score:   0,
		Answers: make(map[int]string),
	}

	return playerID, nil
}

// Add a method to safely get all games
func (gm *GameManager) GetAllGames() map[string]*types.GameState {
	gm.mu.RLock()
	defer gm.mu.RUnlock()

	// Create a copy of the games map to avoid concurrent access issues
	games := make(map[string]*types.GameState)
	for k, v := range gm.Games {
		games[k] = v
	}
	return games
}
func (gm *GameManager) EndGame(gameID string) error {
	gm.mu.Lock()
	defer gm.mu.Unlock()

	if game, exists := gm.Games[gameID]; exists {
		game.Mu.Lock()
		game.IsActive = false
		game.EndTime = time.Now()
		game.Mu.Unlock()

		// Save to database
		return gm.Db.SaveGame(game)
	}
	return nil
}

func (gm *GameManager) ClearAllGames() error {
	gm.mu.Lock()
	defer gm.mu.Unlock()

	// Clear from memory
	gm.Games = make(map[string]*types.GameState)

	// Clear from database
	return gm.Db.ClearAllGames()
}

// Update NewGameManager to initialize Games instead of games

func NewGameManager(dbURL, authToken string) (*GameManager, error) {
	database, err := db.NewDB(dbURL, authToken)
	if err != nil {
		return nil, err
	}

	// Initialize database tables
	if err := database.Initialize(); err != nil {
		return nil, err
	}

	// Load existing games
	games, err := database.LoadGames()
	if err != nil {
		return nil, err
	}

	return &GameManager{
		Games: games,
		Db:    database,
	}, nil
}
func NewGameState(name string) *types.GameState {
	gameID := fmt.Sprintf("game_%d", time.Now().UnixNano())

	return &types.GameState{
		ID:       gameID,
		Name:     name,
		Players:  make(map[string]*types.Player),
		IsActive: false,
		Round:    0,
		Mu:       sync.RWMutex{},
	}
}
func (gm *GameManager) CreateGame(name string) (*types.GameState, error) {
	gm.mu.Lock()
	defer gm.mu.Unlock()

	game := NewGameState(name)
	gm.Games[game.ID] = game

	// Save to database
	if err := gm.Db.SaveGame(game); err != nil {
		return nil, err
	}

	return game, nil
}

// Update CreateGame to use Games instead of games
// func (gm *GameManager) CreateGame() string {
// 	gm.mu.Lock()
// 	defer gm.mu.Unlock()
//
// 	gameID := fmt.Sprintf("game_%d", len(gm.Games)+1)
// 	gm.Games[gameID] = &GameState{
// 		ID:       gameID,
// 		Players:  make(map[string]*Player),
// 		IsActive: false,
// 		Round:    0,
// 	}
// 	fmt.Println(gameID)
// 	return gameID
// }

// Update GetGame to use Games instead of games
func (gm *GameManager) GetGame(gameID string) (*types.GameState, error) {
	gm.mu.RLock()
	defer gm.mu.RUnlock()

	game, exists := gm.Games[gameID]
	if !exists {
		return nil, fmt.Errorf("Game not found please make sure to either create a new game or select an existing one form the dropdown on the right: %s", gameID)
	}
	return game, nil
}

func (gm *GameManager) GetFirstGameID() string {
	gm.mu.RLock()
	defer gm.mu.RUnlock()

	// Get first key from map
	for id := range gm.Games {
		return id
	}
	return ""
}
