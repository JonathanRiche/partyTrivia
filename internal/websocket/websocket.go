package websocket

import (
	"fmt"
	"github.com/gorilla/websocket"
	"net/http"
	"richetechguy/internal/game"
	"richetechguy/internal/types"
	"time"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

type Message struct {
	Type    string      `json:"type"`
	Payload interface{} `json:"payload"`
}

func HandleWebSocket(gameManager *game.GameManager) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		conn, err := upgrader.Upgrade(w, r, nil)
		if err != nil {
			fmt.Printf("WebSocket upgrade error: %v\n", err)
			return
		}
		defer conn.Close()

		// Get player name from query parameters
		playerName := r.URL.Query().Get("name")
		if playerName == "" {
			playerName = fmt.Sprintf("Player_%d", time.Now().UnixNano())
		}

		// Create a new player ID
		playerID := fmt.Sprintf("player_%d", time.Now().UnixNano())

		// Find an available game or create a new one
		var activeGame *types.GameState
		games := gameManager.GetAllGames()
		for _, g := range games {
			if !g.IsActive {
				activeGame = g
				break
			}
		}

		if activeGame == nil {
			gameM, err := gameManager.CreateGame("Rookie of the Year")
			if err != nil {
				http.Error(w, "Error creating game", http.StatusInternalServerError)
				return
			}
			activeGame, _ = gameManager.GetGame(gameM.ID)
		}

		// Add player to game
		player := &types.Player{
			ID:      playerID,
			Name:    playerName,
			Score:   0,
			Answers: make(map[int]string),
			WSConn:  conn,
			GameID:  activeGame.ID,
		}

		activeGame.Players[player.ID] = player

		// Broadcast player joined message
		broadcastMessage := Message{
			Type: "playerJoined",
			Payload: map[string]interface{}{
				"players": activeGame.Players,
			},
		}
		broadcastToPlayers(activeGame, broadcastMessage)

		// Handle incoming messages
		for {
			var msg Message
			err := conn.ReadJSON(&msg)
			if err != nil {
				fmt.Printf("WebSocket read error: %v\n", err)
				delete(activeGame.Players, player.ID)
				broadcastPlayerLeft(activeGame, player)
				break
			}

			handlePlayerMessage(msg, player, activeGame)
		}
	}
}

func broadcastToPlayers(gameState *types.GameState, msg Message) {
	for _, player := range gameState.Players {
		if player.WSConn != nil {
			if err := player.WSConn.WriteJSON(msg); err != nil {
				fmt.Printf("Error broadcasting to player %s: %v\n", player.ID, err)
				player.CloseConnection()
				delete(gameState.Players, player.ID)
			}
		}
	}
}
func broadcastPlayerLeft(gameState *types.GameState, player *types.Player) {
	msg := Message{
		Type: "playerLeft",
		Payload: map[string]interface{}{
			"playerID": player.ID,
			"players":  gameState.Players,
		},
	}
	broadcastToPlayers(gameState, msg)
}

func handlePlayerMessage(msg Message, player *types.Player, gameState *types.GameState) {
	switch msg.Type {
	case "answer":
		if payload, ok := msg.Payload.(map[string]interface{}); ok {
			if answer, ok := payload["answer"].(string); ok {
				gameState.SubmitAnswer(player.ID, answer)
			}
		}
	}
}
func HandleAdminWebSocket(gameManager *game.GameManager) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		conn, err := upgrader.Upgrade(w, r, nil)
		if err != nil {
			http.Error(w, "Could not upgrade connection", http.StatusInternalServerError)
			return
		}
		defer conn.Close()

		// Send initial game state
		games := gameManager.GetAllGames()
		for gameID, game := range games {
			msg := Message{
				Type: "gameStatus",
				Payload: map[string]interface{}{
					"gameId": gameID,
					"status": game.GetGameStatus(),
				},
			}
			conn.WriteJSON(msg)
		}

		// Keep connection alive and handle any admin commands
		for {
			var msg Message
			err := conn.ReadJSON(&msg)
			if err != nil {
				break
			}

			// Handle admin messages if needed
		}
	}
}
