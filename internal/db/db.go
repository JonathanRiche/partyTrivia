package db

import (
	"context"
	"database/sql"
	"encoding/json"
	_ "github.com/tursodatabase/libsql-client-go/libsql"
	"richetechguy/internal/types"
	"sync"
	// "time"
)

type DB struct {
	db *sql.DB
}

func (d *DB) LoadGames() (map[string]*types.GameState, error) {
	ctx := context.Background()

	rows, err := d.db.QueryContext(ctx, `
        SELECT id, name, is_active, start_time, end_time, questions
        FROM games
    `)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	games := make(map[string]*types.GameState)
	for rows.Next() {
		var game types.GameState
		var questionsJSON string
		var questions []types.Question

		err := rows.Scan(
			&game.ID,
			&game.Name,
			&game.IsActive,
			&game.StartTime,
			&game.EndTime,
			&questionsJSON,
		)
		if err != nil {
			return nil, err
		}

		// Parse questions JSON
		if err := json.Unmarshal([]byte(questionsJSON), &questions); err != nil {
			return nil, err
		}

		game.Questions = questions
		game.Players = make(map[string]*types.Player)

		game.Mu = sync.RWMutex{}

		games[game.ID] = &game
	}

	return games, rows.Err()
}

// func (d *DB) SaveQuestions(gameState *types.GameState, questions types.Question) error {
// 	ctx := context.Background()
//
// 	// Use upsert (INSERT OR REPLACE)
// 	_, err = d.db.ExecContext(ctx, `
//         INSERT OR REPLACE INTO questions (
//             id, text, options, correct
//         ) VALUES (?, ?, ?, ?)
//     `,
// 		questions.ID,
// 		questions.Text,
// 		questions.Options,
// 		questions.Correct,
// 	)
//
// 	return err
// }

func NewDB(url string, authToken string) (*DB, error) {
	db, err := sql.Open("libsql", url+"?authToken="+authToken)
	if err != nil {
		return nil, err
	}

	// Test the connection
	if err := db.Ping(); err != nil {
		return nil, err
	}

	return &DB{db: db}, nil
}

// Initialize creates the necessary tables if they don't exist
func (d *DB) Initialize() error {
	// Create games table
	_, err := d.db.Exec(`
        CREATE TABLE IF NOT EXISTS games (
            id TEXT PRIMARY KEY,
            name TEXT NOT NULL,
            is_active BOOLEAN DEFAULT true,
            start_time DATETIME,
            end_time DATETIME,
            questions JSON,
            created_at DATETIME DEFAULT CURRENT_TIMESTAMP
        )
    `)
	return err
}

// SaveGame saves or updates a game in the database
func (d *DB) SaveGame(game *types.GameState) error {
	ctx := context.Background()

	// Convert questions to JSON
	questionsJSON, err := json.Marshal(game.Questions)
	if err != nil {
		return err
	}

	// Use upsert (INSERT OR REPLACE)
	_, err = d.db.ExecContext(ctx, `
        INSERT OR REPLACE INTO games (
            id, name, is_active, start_time, end_time, questions
        ) VALUES (?, ?, ?, ?, ?, ?)
    `,
		game.ID,
		game.Name,
		game.IsActive,
		game.StartTime,
		game.EndTime,
		string(questionsJSON))

	return err
}

// LoadGames retrieves all games from the database
// func (d *DB) LoadGames() (map[string]*GameState, error) {
// 	ctx := context.Background()
//
// 	rows, err := d.db.QueryContext(ctx, `
//         SELECT id, name, is_active, start_time, end_time, questions
//         FROM games
//     `)
// 	if err != nil {
// 		return nil, err
// 	}
// 	defer rows.Close()
//
// 	games := make(map[string]*GameState)
// 	for rows.Next() {
// 		var game GameState
// 		var questionsJSON string
// 		var questions []Question
//
// 		err := rows.Scan(
// 			&game.ID,
// 			&game.Name,
// 			&game.IsActive,
// 			&game.StartTime,
// 			&game.EndTime,
// 			&questionsJSON,
// 		)
// 		if err != nil {
// 			return nil, err
// 		}
//
// 		// Parse questions JSON
// 		if err := json.Unmarshal([]byte(questionsJSON), &questions); err != nil {
// 			return nil, err
// 		}
//
// 		game.Questions = questions
// 		game.Players = make(map[string]*Player)
// 		game.mu = &sync.Mutex{}
//
// 		games[game.ID] = &game
// 	}
//
// 	return games, rows.Err()
// }

// DeleteGame removes a game from the database
func (d *DB) DeleteGame(gameID string) error {
	ctx := context.Background()
	_, err := d.db.ExecContext(ctx, "DELETE FROM games WHERE id = ?", gameID)
	return err
}

// ClearAllGames removes all games from the database
func (d *DB) ClearAllGames() error {
	ctx := context.Background()
	_, err := d.db.ExecContext(ctx, "DELETE FROM games")
	return err
}
