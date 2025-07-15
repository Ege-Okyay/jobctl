package main

import (
	"context"
	"time"

	"github.com/Ege-Okyay/jobctl/internal/initapp"
	"github.com/Ege-Okyay/jobctl/internal/runner"
	"github.com/Ege-Okyay/jobctl/internal/shell"
	"github.com/Ege-Okyay/jobctl/internal/watcher"
)

func main() {
	paths := initapp.SetupApp()

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// TOOD: Watcher still not starting
	go watcher.WatchConfig(ctx, paths.ConfigPath, 5*time.Second)
	go runner.RunScheduler(ctx)

	shell.LaunchInteractiveShell()
}
