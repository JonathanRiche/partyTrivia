package websocket

import (
	"net/http"
	"richetechguy/internal/game"

	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		return true // You might want to add more security here
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
			http.Error(w, "Could not upgrade connection", http.StatusInternalServerError)
			return
		}
		defer conn.Close()

		// Get player info from session/query params
		playerID := r.URL.Query().Get("player_id")
		gameID := r.URL.Query().Get("game_id")

		// Add player to game
		game := gameManager.Games[gameID]
		if game == nil {
			return
		}

		player := game.Players[playerID]
		if player == nil {
			return
		}

		player.WSConn = conn

		// Handle WebSocket messages
		for {
			var msg Message
			err := conn.ReadJSON(&msg)
			if err != nil {
				break
			}

			// Handle different message types
			switch msg.Type {
			case "answer":
				// Handle player answer
				// ...
			}
		}
	}
}
