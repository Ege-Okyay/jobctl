package db

import (
	"database/sql"
	"fmt"
	"os"
	"path/filepath"
	"sort"

	"github.com/Ege-Okyay/jobctl/internal/logger"
	_ "github.com/glebarez/sqlite"
)

var DB *sql.DB

// Open initializes the SQLite database connection and applies all .sql migrations
// from the specified directory in lexicographical order.
func Open(path, migrationsDir string) error {
	var err error
	DB, err = sql.Open("sqlite", path)
	if err != nil {
		return err
	}

	logger.Log("Looking for migrations in %q\n", migrationsDir)
	files, err := filepath.Glob(filepath.Join(migrationsDir, "*.sql"))
	if err != nil {
		return err
	}
	logger.Log("Found %d migration files: %v\n", len(files), files)
	sort.Strings(files)

	for _, f := range files {
		sqlBytes, err := os.ReadFile(f)
		if err != nil {
			return fmt.Errorf("reading %s: %w", f, err)
		}

		if _, err := DB.Exec(string(sqlBytes)); err != nil {
			return fmt.Errorf("applying %s: %w", f, err)
		}
	}

	return nil
}
