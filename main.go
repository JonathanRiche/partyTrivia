package main

import (
	"fmt"
	"net/http"
	"os"
	"richetechguy/internal/blink"
	"richetechguy/internal/generate"
	"richetechguy/internal/middleware"
	"richetechguy/internal/template"
	"richetechguy/internal/view"
	"strings"

	"github.com/joho/godotenv"
)

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
	email := r.Form.Get("email")
	name := r.Form.Get("name")

	if email == "" || name == "" {
		http.Error(w, "Email and name are required", http.StatusBadRequest)
		return
	}

	// Process the data (for now, we'll just print it)
	fmt.Printf("Received submission - Email: %s, Name: %s\n", email, name)
	// Send a response
	w.Write([]byte("Form submitted successfully"))
}
func dynamicPath(w http.ResponseWriter, r *http.Request) {
	item, err := blink.GetBlinkData(r.URL.Path)
	if err != nil {
		http.Error(w, fmt.Sprintf("Error fetching data: %v", err), http.StatusInternalServerError)
		return
	}
	middleware.Chain(w, r, template.Home("Home", item))
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

	mux.HandleFunc("GET /", func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/" {
			item, err := blink.GetBlinkData("home")
			if err != nil {
				http.Error(w, fmt.Sprintf("Error fetching data: %v", err), http.StatusInternalServerError)
				return
			}
			middleware.Chain(w, r, template.Home("Home", item))
		} else {
			item, err := blink.GetBlinkData(strings.TrimPrefix(r.URL.Path, "/"))
			fmt.Println(item.Status)
			if item.Status == "" {
				http.NotFound(w, r)
				return
			}
			if err != nil {
				http.Error(w, fmt.Sprintf("Error fetching data: %v", err), http.StatusInternalServerError)
				return
			}
			middleware.Chain(w, r, template.Home("Page Name", item))
		}
	})
	mux.HandleFunc("POST /submitDealer", handleSubmit)
	fmt.Printf("server is running on port %s\n", os.Getenv("PORT"))
	err = http.ListenAndServe(":"+os.Getenv("PORT"), mux)
	if err != nil {
		fmt.Println(err)
	}

}
