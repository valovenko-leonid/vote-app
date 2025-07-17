CREATE TABLE IF NOT EXISTS users (
    id SERIAL PRIMARY KEY,
    user_id TEXT NOT NULL UNIQUE,
    fingerprint TEXT NOT NULL
);

CREATE TABLE IF NOT EXISTS options (
    id SERIAL PRIMARY KEY,
    text VARCHAR(40) NOT NULL,
    votes INT DEFAULT 0
);

CREATE TABLE IF NOT EXISTS votes (
    id SERIAL PRIMARY KEY,
    option_id INT NOT NULL REFERENCES options(id) ON DELETE CASCADE,
    user_id TEXT NOT NULL REFERENCES users(user_id),
    created_at TIMESTAMPTZ DEFAULT now()
);

CREATE UNIQUE INDEX IF NOT EXISTS uniq_vote
  ON votes (user_id, option_id);