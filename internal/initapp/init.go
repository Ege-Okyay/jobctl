package initapp

import (
	"fmt"
	"log"
	"os"
	"path/filepath"

	"github.com/Ege-Okyay/jobctl/internal/cli"
	"github.com/Ege-Okyay/jobctl/internal/config"
	"github.com/Ege-Okyay/jobctl/internal/db"
	"github.com/Ege-Okyay/jobctl/internal/logger"
	"github.com/Ege-Okyay/jobctl/internal/logic"
	"github.com/Ege-Okyay/jobctl/internal/types"
	"github.com/Ege-Okyay/jobctl/internal/util"
)

// SetupApp orchestrates the application's startup sequence.
func SetupApp() types.AppPaths {
	paths := resolvePaths()

	createDirs(paths)

	initLogger(paths.LogDir)
	initDatabase(paths)

	// Ensure the database state is synchronized with the configuration file.
	syncConfig(paths.ConfigPath)

	util.PrintBanner()
	cli.Setup()

	return paths
}

// resolvePas determines and returns athll critical application paths.
func resolvePaths() types.AppPaths {
	cfgPath, err := config.ConfigPath()
	if err != nil {
		fmt.Fprintf(os.Stderr, "cannot determine config path: %v\n", err)
		os.Exit(1)
	}

	cfgDir := filepath.Dir(cfgPath)
	dbPath := filepath.Join(cfgDir, "jobctl.db")

	cacheDir, _ := os.UserCacheDir()
	logDir := filepath.Join(cacheDir, "jobctl", "logs")

	wd, err := os.Getwd()
	if err != nil {
		fmt.Fprintf(os.Stderr, "cannot determine working directory: %v\n", err)
		os.Exit(1)
	}
	migDir := filepath.Join(wd, "migrations")

	return types.AppPaths{
		DBPath:     dbPath,
		ConfigPath: cfgPath,
		LogDir:     logDir,
		MigDir:     migDir,
	}
}

// createDirs ensures that all necessary application directories exist.
func createDirs(paths types.AppPaths) {
	for _, dir := range []string{
		filepath.Dir(paths.DBPath),
		filepath.Dir(paths.ConfigPath),
		paths.LogDir,
	} {
		if err := os.MkdirAll(dir, 0o755); err != nil {
			fmt.Fprintf(os.Stderr, "Failed to create directory %q: %v\n", dir, err)
			os.Exit(1)
		}
	}
}

// initLogger sets up the application's logging system.
func initLogger(logDir string) {
	if err := logger.Init(logDir); err != nil {
		fmt.Fprintln(os.Stderr, "Failed to init logger:", err)
		os.Exit(1)
	}
}

// initDatabase connects to and migrates the database.
func initDatabase(paths types.AppPaths) {
	if err := db.Open(paths.DBPath, paths.MigDir); err != nil {
		fmt.Fprintln(os.Stderr, "DB/migrations error:", err)
		os.Exit(1)
	}
}

// syncConfig synchronizes the database with the contents of the config file.
func syncConfig(configPath string) {
	if err := logic.SyncDBWithConfig(configPath); err != nil {
		log.Fatalf("Sync error: %v", err)
	}
}
