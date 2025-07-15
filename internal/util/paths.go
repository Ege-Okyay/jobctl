package util

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/Ege-Okyay/jobctl/internal/config"
	"github.com/Ege-Okyay/jobctl/internal/types"
)

func ResolvePaths() *types.AppPaths {
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

	return &types.AppPaths{
		DBPath:     dbPath,
		ConfigPath: cfgPath,
		LogDir:     logDir,
		MigDir:     migDir,
	}
}
