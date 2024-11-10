package types

import (
	"encoding/json"
	"fmt"
	"strings"
	"sync"
	"time"

	"github.com/gorilla/websocket"
)

type QuestionType string

const (
	MultipleChoice QuestionType = "multiple"
	SingleChoice   QuestionType = "single"
)

// IsValid checks if the question type is valid
func (qt QuestionType) IsValid() bool {
	switch qt {
	case MultipleChoice, SingleChoice:
		return true
	default:
		return false
	}
}

// String implements the Stringer interface
func (qt QuestionType) String() string {
	return string(qt)
}

// Question represents a single trivia question
type Question struct {
	ID      int          `json:"id"`
	Text    string       `json:"text"`
	Options []string     `json:"options"`
	Type    QuestionType `json:"type"`
	Correct string       `json:"correct"`
}

// ValidateType ensures the question type is valid
func (q *Question) ValidateType() error {
	if !q.Type.IsValid() {
		return fmt.Errorf("invalid question type: %s", q.Type)
	}
	return nil
}
func (qt *QuestionType) UnmarshalJSON(data []byte) error {
	var s string
	if err := json.Unmarshal(data, &s); err != nil {
		return err
	}

	temp := QuestionType(s)
	if !temp.IsValid() {
		return fmt.Errorf("invalid question type: %s", s)
	}

	*qt = temp
	return nil
}

// Player represents a game participant
type Player struct {
	ID      string          `json:"id"`
	Name    string          `json:"name"`
	Score   int             `json:"score"`
	Answers map[int]string  `json:"answers"` // maps question ID to answer
	WSConn  *websocket.Conn // Add this field
	GameID  string
}

// GameState represents the current state of a trivia game
type GameState struct {
	ID              string
	Players         map[string]*Player
	CurrentQuestion *Question
	Questions       []Question
	IsActive        bool
	Round           int
	Name            string
	StartTime       time.Time
	EndTime         time.Time
	Mu              sync.RWMutex
}

// func (gs *GameState) SetQuestions(questions []Question) {
// 	gs.Questions = questions
// }

// Game state management methods
func (gs *GameState) StartGame() error {
	gs.Mu.Lock()
	defer gs.Mu.Unlock()

	if len(gs.Questions) == 0 {
		return fmt.Errorf("no questions available")
	}

	if len(gs.Players) == 0 {
		return fmt.Errorf("no players joined")
	}

	gs.IsActive = true
	gs.Round = 0
	gs.StartTime = time.Now()
	return nil
}

func (gs *GameState) EndGame() {
	gs.Mu.Lock()
	defer gs.Mu.Unlock()

	gs.IsActive = false
	gs.EndTime = time.Now()
	gs.calculateFinalScores()
}

func (gs *GameState) NextQuestion() (*Question, error) {
	gs.Mu.Lock()
	defer gs.Mu.Unlock()

	if !gs.IsActive {
		return nil, fmt.Errorf("game is not active")
	}

	if gs.Round >= len(gs.Questions) {
		return nil, fmt.Errorf("no more questions")
	}

	gs.Round++
	gs.CurrentQuestion = &gs.Questions[gs.Round-1]
	return gs.CurrentQuestion, nil
}

func (gs *GameState) SubmitAnswer(playerID string, answer string) error {
	gs.Mu.Lock()
	defer gs.Mu.Unlock()

	player, exists := gs.Players[playerID]
	if !exists {
		return fmt.Errorf("player not found")
	}

	if gs.CurrentQuestion == nil {
		return fmt.Errorf("no active question")
	}

	player.Answers[gs.CurrentQuestion.ID] = answer

	// Update score if answer is correct
	if answer == gs.CurrentQuestion.Correct {
		player.Score += 10 // Basic scoring - 10 points per correct answer
	}

	return nil
}

func (gs *GameState) calculateFinalScores() {
	for _, player := range gs.Players {
		// You can implement more complex scoring logic here
		// For example, time bonuses, streaks, etc.

		// Currently using basic scoring from submitted answers
		finalScore := 0
		for questionID, answer := range player.Answers {
			// Find the question and check if the answer was correct
			for _, q := range gs.Questions {
				if q.ID == questionID && answer == q.Correct {
					finalScore += 10
				}
			}
		}
		player.Score = finalScore
	}
}

// GetGameStatus returns a snapshot of the current game state
func (gs *GameState) GetGameStatus() map[string]interface{} {
	gs.Mu.RLock()
	defer gs.Mu.RUnlock()

	return map[string]interface{}{
		"id":        gs.ID,
		"isActive":  gs.IsActive,
		"round":     gs.Round,
		"players":   gs.Players,
		"question":  gs.CurrentQuestion,
		"startTime": gs.StartTime,
		"endTime":   gs.EndTime,
	}
}

// Add methods to safely handle the WebSocket connection
func (p *Player) SetConnection(conn *websocket.Conn) {
	p.WSConn = conn
}

func (p *Player) CloseConnection() {
	if p.WSConn != nil {
		p.WSConn.Close()
		p.WSConn = nil
	}
}

func (p *Player) SubmitAnswer(questionID int, answer string) {
	if p.Answers == nil {
		p.Answers = make(map[int]string)
	}
	p.Answers[questionID] = answer
}

func (p *Player) GetAnswer(questionID int) (string, bool) {
	if p.Answers == nil {
		return "", false
	}
	answer, exists := p.Answers[questionID]
	return answer, exists
}

func (p *Player) GetAllAnswers() map[int]string {
	if p.Answers == nil {
		return make(map[int]string)
	}
	return p.Answers
}
func (q *Question) ValidateAnswer(answer string) bool {
	switch q.Type {
	case QuestionType(SingleChoice):
		return answer == q.Correct
	case QuestionType(MultipleChoice):
		// Split both correct answers and submitted answers
		correctAnswers := strings.Split(q.Correct, ",")
		submittedAnswers := strings.Split(answer, ",")
		// Compare sets of answers
		return compareAnswerSets(correctAnswers, submittedAnswers)
	default:
		return false
	}
}

func compareAnswerSets(correct, submitted []string) bool {
	if len(correct) != len(submitted) {
		return false
	}

	// Create maps for easier comparison
	correctMap := make(map[string]bool)
	submittedMap := make(map[string]bool)

	for _, ans := range correct {
		correctMap[strings.TrimSpace(ans)] = true
	}

	for _, ans := range submitted {
		submittedMap[strings.TrimSpace(ans)] = true
	}

	// Compare maps
	for ans := range correctMap {
		if !submittedMap[ans] {
			return false
		}
	}

	return true
}
