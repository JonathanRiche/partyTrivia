-- Games table
CREATE TABLE IF NOT EXISTS games (
    id TEXT PRIMARY KEY,
    name TEXT NOT NULL,
    is_active BOOLEAN DEFAULT true,
    start_time DATETIME,
    end_time DATETIME,
    questions JSON,
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP
);

-- Questions table (optional, if you want to store questions separately)
CREATE TABLE IF NOT EXISTS questions (
    id TEXT PRIMARY KEY,
    game_id TEXT,
    question TEXT NOT NULL,
    correct_answer TEXT NOT NULL,
    incorrect_answers JSON,
    points INTEGER DEFAULT 1,
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (game_id) REFERENCES games(id)
);

-- Players table (optional, for persistent player stats)
CREATE TABLE IF NOT EXISTS players (
    id TEXT PRIMARY KEY,
    name TEXT NOT NULL,
    score INTEGER DEFAULT 0,
    game_id TEXT,
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (game_id) REFERENCES games(id)
);

-- Game results table (optional, for historical data)
CREATE TABLE IF NOT EXISTS game_results (
    id TEXT PRIMARY KEY,
    game_id TEXT,
    player_id TEXT,
    score INTEGER,
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (game_id) REFERENCES games(id),
    FOREIGN KEY (player_id) REFERENCES players(id)
);

