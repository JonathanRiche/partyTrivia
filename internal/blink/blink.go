package blink

import (
	"encoding/json"
	"fmt"
	"io"
	// "log"
	"net/http"
	"os"

	"github.com/joho/godotenv"
)

type Body struct {
	HTML string `json:"html"`
	CSS  string `json:"css"`
}

type ContentItem struct {
	Title         string `json:"title"`
	ItemID        int    `json:"item_id"`
	Body          Body   `json:"body"`
	Status        string `json:"status"`
	ContentType   string `json:"content_type"`
	Author        string `json:"author"`
	PublishedDate string `json:"published_date"`
	SeoFields     string `json:"seo_fields"`
	PageURL       string `json:"page_url"`
}

// TODO: - Add error handling
func GetBlinkData(item string) (ContentItem, error) {
	envErr := godotenv.Load()
	if envErr != nil {
		return ContentItem{}, fmt.Errorf("error loading .env file: %w", envErr)
	}
	apiKey := os.Getenv("BLINK_API_KEY")
	response, err := http.Get(fmt.Sprintf("https://blinkx.io/api/v1/content-item?dataFilter=html&api_token=%s&page_url=%s", apiKey, item))
	if err != nil {
		return ContentItem{}, fmt.Errorf("error making API request: %w", err)
	}
	defer response.Body.Close()

	body, err := io.ReadAll(response.Body)
	if err != nil {
		return ContentItem{}, fmt.Errorf("error reading response body: %w", err)
	}

	var contentItem ContentItem
	err = json.Unmarshal(body, &contentItem)
	if err != nil {
		return ContentItem{}, fmt.Errorf("error unmarshaling JSON: %w", err)
	}
	return contentItem, nil
}
