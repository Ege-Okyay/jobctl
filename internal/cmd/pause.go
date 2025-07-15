package cmd

import (
	"fmt"
	"time"

	"github.com/Ege-Okyay/jobctl/internal/db"
	"github.com/Ege-Okyay/jobctl/internal/types"
	"github.com/Ege-Okyay/jobctl/internal/util"
)

var PauseCommand types.Command

func init() {
	PauseCommand = types.Command{
		Name:        "pause",
		Description: "Pause all job executions",
		Usage:       "pause",
		Flags:       []types.Flag{},
		Execute:     pauseHandler,
	}
}

func pauseHandler(args []string) {
	isPaused, err := db.IsPaused()
	if err != nil {
		util.ErrorMessage(fmt.Sprintf("Could not check paused state: %v", err))
		return
	}
	if isPaused {
		util.ErrorMessage("System is already paused")
		return
	}

	now := time.Now()

	if err := db.SetPaused(true); err != nil {
		util.ErrorMessage(fmt.Sprint("Failed to pause system:", err))
		return
	}

	// Record the timestamp to calculate the duration of the pause upon resume.
	if err := db.SetPauseTimestamp(now); err != nil {
		util.ErrorMessage(fmt.Sprintf("Paused flag set, but failed to record timestamp: %v", err))
		return
	}

	util.SuccessMessage("System paused: all job countdowns frozen, no executions will occur")
}
