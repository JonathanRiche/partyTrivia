package game

import (
	"encoding/json"
	"fmt"
	"os"
	"richetechguy/internal/types"
	"sync"
)

type QuestionManager struct {
	questions []types.Question
	mu        sync.RWMutex
}

func NewQuestionManager() *QuestionManager {
	return &QuestionManager{
		questions: make([]types.Question, 0),
	}
}

func (qm *QuestionManager) AddQuestion(q types.Question) error {
	qm.mu.Lock()
	defer qm.mu.Unlock()

	// Validate question
	if q.Text == "" || len(q.Options) != 4 || q.Correct == "" {
		return fmt.Errorf("invalid question format")
	}

	q.ID = len(qm.questions) + 1
	qm.questions = append(qm.questions, q)

	// Save to persistent storage
	return qm.saveQuestions()
}

func (qm *QuestionManager) GetQuestions() []types.Question {
	qm.mu.RLock()
	defer qm.mu.RUnlock()
	return qm.questions
}

func (qm *QuestionManager) saveQuestions() error {
	data, err := json.Marshal(qm.questions)
	if err != nil {
		return err
	}
	return os.WriteFile("questions.json", data, 0644)
}

func (qm *QuestionManager) LoadQuestions() error {
	data, err := os.ReadFile("questions.json")
	if err != nil {
		if os.IsNotExist(err) {
			return nil
		}
		return err
	}

	qm.mu.Lock()
	defer qm.mu.Unlock()
	return json.Unmarshal(data, &qm.questions)
}
