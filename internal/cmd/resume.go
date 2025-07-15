package cmd

import (
	"fmt"
	"time"

	"github.com/Ege-Okyay/jobctl/internal/db"
	"github.com/Ege-Okyay/jobctl/internal/types"
	"github.com/Ege-Okyay/jobctl/internal/util"
)

var ResumeCommand types.Command

func init() {
	ResumeCommand = types.Command{
		Name:        "resume",
		Description: "Resume all job executions",
		Usage:       "resume",
		Flags:       []types.Flag{},
		Execute:     resumeHandler,
	}
}

// resumeHandler unfreezes the systeman d adjusts all job schedules to account for the paused duration.
func resumeHandler(args []string) {
	isPaused, err := db.IsPaused()
	if err != nil {
		util.ErrorMessage(fmt.Sprintf("Could not check paused state: %v", err))
		return
	}
	if !isPaused {
		util.ErrorMessage("System is not paused")
		return
	}

	pausedAt, err := db.GetPauseTimestamp()
	if err != nil {
		util.ErrorMessage(fmt.Sprintf("Error reading pause timestamp: %v", err))
		return
	}

	// Calculate the time elapsed since the system was paused.
	now := time.Now()
	delta := now.Sub(pausedAt)

	jobs, err := db.GetAllJobs()
	if err != nil {
		util.ErrorMessage(fmt.Sprintf("Error loading jobs: %v", err))
		return
	}

	// Shift the next_run time for all enabled jobs forward by the pause duration.
	for _, j := range jobs {
		if !j.Enabled {
			continue
		}

		newNext := j.NextRun.Add(delta)
		if err := db.UpdateNextRun(j.Name, newNext); err != nil {
			util.ErrorMessage(fmt.Sprintf("Failed to update next_run for %q: %v", j.Name, err))
			return // Stop on first error to prevent inconsistent state.
		}
	}

	if err := db.SetPaused(false); err != nil {
		util.ErrorMessage(fmt.Sprintf("Failed to clear paused flag: %v", err))
		return
	}

	if err := db.ClearPauseTimestamp(); err != nil {
		util.ErrorMessage(fmt.Sprintf("Failed to clear pause timestamp: %v", err))
		return
	}

	util.SuccessMessage(fmt.Sprintf("System resumed: schedules shifted forward by %s", delta.Round(time.Second)))
}
