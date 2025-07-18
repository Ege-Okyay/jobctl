CREATE TABLE IF NOT EXISTS jobs (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    name TEXT UNIQUE NOT NULL,
    interval_sec INTEGER NOT NULL,
    command TEXT NOT NULL,
    next_run DATETIME NOT NULL,
    retries INTEGER NOT NULL DEFAULT 0,
    timeout_sec INTEGER NOT NULL DEFAULT 0,
    enabled BOOLEAN NOT NULL DEFAULT 1
);

CREATE INDEX IF NOT EXISTS idx_jobs_next_run ON jobs(next_run);

CREATE TABLE IF NOT EXISTS job_runs (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    job_id INTEGER NOT NULL,
    ts DATETIME NOT NULL,
    success BOOLEAN NOT NULL,
    output TEXT,
    FOREIGN KEY(job_id) REFERENCES jobs(id)
);