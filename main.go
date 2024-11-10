package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"richetechguy/internal/game"
	"richetechguy/internal/generate"
	"richetechguy/internal/middleware"
	"richetechguy/internal/template"
	"richetechguy/internal/types"
	"richetechguy/internal/view"
	"richetechguy/internal/websocket"
	"strconv"
	"time"

	"richetechguy/internal/admin"
	// "strings"

	"github.com/joho/godotenv"
)

func handleJoinGame(gm *game.GameManager) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		name := r.FormValue("name")
		gameID := r.FormValue("gameId")

		if name == "" || gameID == "" {
			http.Error(w, "Name and game ID are required", http.StatusBadRequest)
			return
		}

		// Add player to game
		playerID, err := gm.AddPlayer(gameID, name)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		fmt.Printf("Player ID added: %s and name: %s\n", playerID, name)

		// Render game lobby with player info
		template.GameLobby(name).Render(r.Context(), w)
	}
}

func handleSubmit(w http.ResponseWriter, r *http.Request) {

	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	err := r.ParseForm()
	if err != nil {
		http.Error(w, "Error parsing form data", http.StatusBadRequest)
		return
	}
	//Placeholder for the email and name wire this up for correct fields
	name := r.Form.Get("name")

	if name == "" {
		http.Error(w, "Email and name are required", http.StatusBadRequest)
		return
	}

	// Process the data (for now, we'll just print it)
	fmt.Printf("Received submission - Name: %s\n", name)
	// Send a response
	w.Write([]byte("Form submitted successfully"))
}
func handleAdmin(gm *game.GameManager) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Add authentication check here
		middleware.Chain(w, r, admin.Dashboard(gm))
	}
}

// func handleAdminOld(w http.ResponseWriter, r *http.Request) {
// 	// Add authentication check here
// 	middleware.Chain(w, r, admin.Dashboard())
// }

func handleCreateGame(gm *game.GameManager) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		//NOTE: THIS is a placeholder for the game creation logic
		gameID, err := gm.CreateGame("Rookie of the Year")
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		w.Header().Set("HX-Trigger", "gameCreated")
		w.Header().Set("HX-Refresh", "true") // This will refresh the page
		fmt.Fprintf(w, "Game created: %s", gameID.Name)
	}
}
func handleStartGame(gm *game.GameManager, qm *game.QuestionManager) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		gameID := r.FormValue("gameID")
		if err := gm.StartGame(gameID, qm); err != nil {
			w.Header().Set("HX-Trigger", fmt.Sprintf(`{"showMessage": "%s"}`, err))
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		// Get the game state
		game, _ := gm.GetGame(gameID)

		// Broadcast to players
		playerMessage := websocket.Message{
			Type: "gameState",
			Payload: map[string]interface{}{
				"state":   "active",
				"message": "Game has started!",
				"players": game.Players,
			},
		}
		websocket.BroadcastToPlayers(game, playerMessage)

		// Broadcast to admins
		adminMessage := websocket.Message{
			Type: "playerList",
			Payload: map[string]interface{}{
				"gameId":   gameID,
				"players":  game.Players,
				"isActive": true,
			},
		}
		websocket.BroadcastToAdmins(adminMessage)

		w.Header().Set("HX-Trigger", "gameStarted")
		fmt.Fprintf(w, "Game started")
	}
}

func handleSelectGame(gm *game.GameManager) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		gameID := r.FormValue("gameID")
		fmt.Println(gameID)
		gameState, err := gm.SelectGame(gameID)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		w.Header().Set("HX-Trigger", "gameSelected")
		fmt.Fprintf(w, "Game selected %s", gameState.Name)
	}
}
func handleEndGame(gm *game.GameManager) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		gameID := r.FormValue("gameID")
		gm.EndGame(gameID)
		w.Header().Set("HX-Trigger", "gameEnded")
		fmt.Fprintf(w, "Game ended")
	}
}
func handleClearGames(gm *game.GameManager) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		gm.ClearAllGames()
		w.Header().Set("HX-Trigger", "gamesCleared")
		w.Header().Set("HX-Refresh", "true") // This will refresh the page
		fmt.Fprintf(w, "All games have been cleared")
	}
}

func handleAddQuestion(qm *game.QuestionManager) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		q := types.Question{
			Text: r.FormValue("questionText"),
			Options: []string{
				r.FormValue("option1"),
				r.FormValue("option2"),
				r.FormValue("option3"),
				r.FormValue("option4"),
			},
			Correct: r.FormValue("correctAnswer"),
		}

		if err := qm.AddQuestion(q); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		// Return updated question list
		admin.QuestionList(qm.GetQuestions()).Render(r.Context(), w)
	}
}
func handleGameStatus(gm *game.GameManager) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		gameID := r.FormValue("gameID")
		game, _ := gm.GetGame(gameID)
		admin.GameStatus(game).Render(r.Context(), w)
	}
}

func handlePlayerList(gm *game.GameManager) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		gameID := r.FormValue("gameID")
		game, _ := gm.GetGame(gameID)
		if game != nil {
			admin.PlayerList(game.Players, gameID).Render(r.Context(), w)
		}
	}
}
func handleStartQuestions(gm *game.GameManager, qm *game.QuestionManager) http.HandlerFunc {

	return func(w http.ResponseWriter, r *http.Request) {
		gameID := r.FormValue("gameID")
		// if err := gm.StartGame(gameID, qm); err != nil {
		// 	w.Header().Set("HX-Trigger", fmt.Sprintf(`{"showMessage": "%s"}`, err))
		// 	http.Error(w, err.Error(), http.StatusBadRequest)
		// 	return
		// }

		// Get the game state
		game, err := gm.GetGame(gameID)
		if err != nil {
			w.Header().Set("HX-Trigger", fmt.Sprintf(`{"showMessage": "%s"}`, err))
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		// Broadcast to players
		// {
		//   id: string;
		//   text: string;
		//   options: string[];
		// }
		questions := qm.GetQuestions()
		playerMessage := websocket.Message{
			Type: "question",
			Payload: map[string]interface{}{
				"state":     "active",
				"message":   "Questions has started!",
				"questions": questions,
				"gameId":    gameID,
			},
		}
		websocket.BroadcastToPlayers(game, playerMessage)

		// // Broadcast to admins
		// adminMessage := websocket.Message{
		// 	Type: "playerList",
		// 	Payload: map[string]interface{}{
		// 		"gameId":   gameID,
		// 		"players":  game.Players,
		// 		"isActive": true,
		// 	},
		// }
		// websocket.BroadcastToAdmins(adminMessage)

		// w.Header().Set("HX-Trigger", "gameStarted")
		fmt.Fprintf(w, "Questions started")
	}
}
func handleAnswerSubmission(gm *game.GameManager, qm *game.QuestionManager) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		gameID := r.FormValue("gameID")
		questionID := r.FormValue("questionID")
		answer := r.FormValue("answer")
		playerID := r.FormValue("playerID") // You'll need to pass this from the session

		game, err := gm.GetGame(gameID)
		if err != nil {
			http.Error(w, "Game not found", http.StatusBadRequest)
			return
		}

		player, exists := game.Players[playerID]
		if !exists {
			http.Error(w, "Player not found", http.StatusBadRequest)
			return
		}

		qID, _ := strconv.Atoi(questionID)
		player.SubmitAnswer(qID, answer)

		// Broadcast answer submission to admin
		adminMsg := websocket.Message{
			Type: "playerAnswered",
			Payload: map[string]interface{}{
				"gameId":     gameID,
				"playerId":   playerID,
				"questionId": qID,
				"score":      player.Score,
			},
		}
		websocket.BroadcastToAdmins(adminMsg)

		questions := qm.GetQuestions()
		playerMessage := websocket.Message{
			Type: "question",
			Payload: map[string]interface{}{
				"state":     "active",
				"message":   "On the next question!",
				"questions": questions,
				"gameId":    gameID,
			},
		}
		websocket.BroadcastToPlayers(game, playerMessage)

		// Return updated question view or confirmation
		w.Write([]byte("Answer submitted!"))
	}
}
func main() {

	err := generate.GenerateMain()
	if err != nil {
		panic(err)
	}

	_ = godotenv.Load()
	mux := http.NewServeMux()

	dbURL := os.Getenv("TURSO_DATABASE_URL")
	authToken := os.Getenv("TURSO_AUTH_TOKEN")

	questionManager := game.NewQuestionManager()
	gameManager, err := game.NewGameManager(dbURL, authToken)
	if err != nil {
		log.Fatalf("Failed to initialize game manager: %v", err)
	}
	// Add periodic state saving
	go func() {
		ticker := time.NewTicker(5 * time.Minute)
		for range ticker.C {
			// fmt.Println("Starting periodic state saving...")

			for _, game := range gameManager.Games {
				if err := gameManager.Db.SaveGame(game); err != nil {
					log.Printf("Error saving game state: %v", err)
				}
			}
		}
	}()

	mux.HandleFunc("GET /favicon.ico", view.ServeFavicon)
	mux.HandleFunc("GET /static/", view.ServeStaticFiles)

	// gameManager := game.NewGameManager()

	// Load existing questions
	if err := questionManager.LoadQuestions(); err != nil {
		log.Printf("Error loading questions: %v", err)
	}
	//TODO: modify to bd db instead of questions.json
	// questionManager.LoadQuestions()
	// fmt.Println(questionManager.GetQuestions())

	// Admin routes
	mux.HandleFunc("GET /admin", handleAdmin(gameManager))
	mux.HandleFunc("POST /admin/game/create", handleCreateGame(gameManager))
	mux.HandleFunc("POST /admin/game/start", handleStartGame(gameManager, questionManager))
	mux.HandleFunc("POST /admin/game/end", handleEndGame(gameManager))
	mux.HandleFunc("POST /admin/game/clear", handleClearGames(gameManager))
	mux.HandleFunc("POST /admin/game/select", handleSelectGame(gameManager))
	mux.HandleFunc("POST /admin/questions/add", handleAddQuestion(questionManager))
	mux.HandleFunc("GET /", func(w http.ResponseWriter, r *http.Request) {
		middleware.Chain(w, r, template.JoinGame(gameManager))
	})

	mux.HandleFunc("POST /admin/game/startQuestions", handleStartQuestions(gameManager, questionManager))
	mux.HandleFunc("POST /game/submit-answer", handleAnswerSubmission(gameManager, questionManager))

	mux.HandleFunc("GET /admin/game/status", handleGameStatus(gameManager))
	mux.HandleFunc("GET /admin/game/players", handlePlayerList(gameManager))

	mux.HandleFunc("GET /ws/admin", websocket.HandleAdminWebSocket(gameManager))
	mux.HandleFunc("POST /joinGame", handleJoinGame(gameManager))
	mux.HandleFunc("GET /ws/game", websocket.HandleWebSocket(gameManager))

	fmt.Printf("server is running on  http://localhost:%s\n", os.Getenv("PORT"))
	err = http.ListenAndServe(":"+os.Getenv("PORT"), mux)
	if err != nil {
		fmt.Println(err)
	}

}
