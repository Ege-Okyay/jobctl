CREATE TABLE IF NOT EXISTS settings (
    key TEXT PRIMARY KEY,
    value TEXT NOT NULL,
    paused_at TEXT
);

INSERT OR IGNORE INTO settings(key, value) VALUES('paused', 'false');