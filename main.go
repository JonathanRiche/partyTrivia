package main

import (
	"fmt"
	"net/http"
	"os"
	"richetechguy/internal/game"
	"richetechguy/internal/generate"
	"richetechguy/internal/middleware"
	"richetechguy/internal/template"
	"richetechguy/internal/view"
	"richetechguy/internal/websocket"

	// "strings"

	"github.com/joho/godotenv"
)

func handleJoinGame(gm *game.GameManager) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		name := r.FormValue("name")
		if name == "" {
			http.Error(w, "Name is required", http.StatusBadRequest)
			return
		}

		// Create or join game logic
		// Return game lobby template
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

func main() {

	err := generate.GenerateMain()
	if err != nil {
		panic(err)
	}

	_ = godotenv.Load()
	mux := http.NewServeMux()

	mux.HandleFunc("GET /favicon.ico", view.ServeFavicon)
	mux.HandleFunc("GET /static/", view.ServeStaticFiles)

	gameManager := game.NewGameManager()

	mux.HandleFunc("GET /", func(w http.ResponseWriter, r *http.Request) {
		middleware.Chain(w, r, template.JoinGame())
	})

	mux.HandleFunc("GET /admin", func(w http.ResponseWriter, r *http.Request) {
		// Admin dashboard handler
		// Implement authentication
	})

	mux.HandleFunc("POST /joinGame", handleJoinGame(gameManager))
	mux.HandleFunc("GET /ws/game", websocket.HandleWebSocket(gameManager))

	// mux.HandleFunc("GET /", func(w http.ResponseWriter, r *http.Request) {
	// 	if r.URL.Path == "/" {
	// 		if err != nil {
	// 			http.Error(w, fmt.Sprintf("Error fetching data: %v", err), http.StatusInternalServerError)
	// 			return
	// 		}
	// 		middleware.Chain(w, r, template.Home("Home"))
	// 	} else {
	//
	// 	}
	// })
	//TODO: Add logic for game rooms
	// mux.HandleFunc("POST /joinGame", handleSubmit)
	fmt.Printf("server is running on  http://localhost:%s\n", os.Getenv("PORT"))
	err = http.ListenAndServe(":"+os.Getenv("PORT"), mux)
	if err != nil {
		fmt.Println(err)
	}

}
