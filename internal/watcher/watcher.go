package watcher

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/Ege-Okyay/jobctl/internal/logger"
	"github.com/Ege-Okyay/jobctl/internal/logic"
	"github.com/Ege-Okyay/jobctl/internal/util"
)

// WatchConfig watches the configuration file for changes and triggers a synchronization
// with the database when a change is detected.
func WatchConfig(ctx context.Context, path string, pollInterval time.Duration) {
	logger.Log("Watcher: watching %s every %s", path, pollInterval)

	var lastMod time.Time
	if info, err := os.Stat(path); err == nil {
		lastMod = info.ModTime()
	}

	ticker := time.NewTicker(pollInterval)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			logger.Log("Watcher: stopped")
			return

		case <-ticker.C:
			info, err := os.Stat(path)
			if err != nil {
				logger.Log("Watcher: stat error %v", err)
				continue
			}

			// If the file has been modified since the last check, reload it.
			if info.ModTime().After(lastMod) {
				lastMod = info.ModTime()
				logger.Log("Watcher: detected change, reloading...")

				if err := logic.SyncDBWithConfig(path); err != nil {
					logger.Log("Watcher: reload error: %v", err)
					util.ErrorMessage(fmt.Sprint("\nInvalid config:", err))
				} else {
					logger.Log("Watcher: reload successful")
				}
			}
		}
	}
}
