CREATE TABLE IF NOT EXISTS options (
    id SERIAL PRIMARY KEY,
    text VARCHAR(40) NOT NULL,
    votes INT DEFAULT 0
);

CREATE TABLE IF NOT EXISTS votes (
    id SERIAL PRIMARY KEY,
    option_id INT NOT NULL REFERENCES options(id) ON DELETE CASCADE,
    fingerprint TEXT NOT NULL,
    created_at TIMESTAMPTZ DEFAULT now()
);

CREATE UNIQUE INDEX IF NOT EXISTS uniq_vote
  ON votes (fingerprint, option_id);