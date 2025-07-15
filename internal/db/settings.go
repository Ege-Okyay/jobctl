package db

import (
	"database/sql"
	"fmt"
	"time"
)

func IsPaused() (bool, error) {
	var val string

	err := DB.QueryRow(`SELECT value FROM settings WHERE key = 'paused'`).Scan(&val)
	if err == sql.ErrNoRows {
		return false, nil
	}
	if err != nil {
		return false, err
	}

	return val == "true", nil
}

func SetPaused(state bool) error {
	_, err := DB.Exec(
		`INSERT INTO settings(key, value) VALUES ('paused', ?)
		ON CONFLICT(key) DO UPDATE SET value = excluded.value`,
		fmt.Sprintf("%v", state),
	)

	return err
}

func SetPauseTimestamp(ts time.Time) error {
	_, err := DB.Exec(`
		INSERT INTO settings(key, value, paused_at)
		VALUES('paused', ?, ?)
		ON CONFLICT(key) DO UPDATE
		SET value = excluded.value, paused_at = excluded.paused_at
	`, "true", ts.Format(time.RFC3339))
	return err
}

func GetPauseTimestamp() (time.Time, error) {
	var s string

	err := DB.QueryRow(`SELECT paused_at FROM settings WHERE key = 'paused'`).Scan(&s)
	if err != nil {
		if err == sql.ErrNoRows {
			return time.Time{}, nil
		}
		return time.Time{}, err
	}

	if s == "" {
		return time.Time{}, nil
	}

	return time.Parse(time.RFC3339, s)
}

func ClearPauseTimestamp() error {
	_, err := DB.Exec(`UPDATE settings SET paused_at = NULL WHERE key = 'paused'`)
	return err
}
