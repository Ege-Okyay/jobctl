package cmd

import (
	"fmt"
	"runtime"
	"time"

	"github.com/Ege-Okyay/jobctl/internal/db"
	"github.com/Ege-Okyay/jobctl/internal/types"
	"github.com/Ege-Okyay/jobctl/internal/util"
)

var StatusCommand types.Command

func init() {
	StatusCommand = types.Command{
		Name:        "status",
		Description: "Show jobctl system status",
		Usage:       "status",
		Flags:       []types.Flag{},
		Execute:     statusHandler,
	}
}

func statusHandler(args []string) {
	paused, err := db.IsPaused()
	if err != nil {
		util.ErrorMessage(fmt.Sprint("Could not fetch paused status:", err))
		return
	}

	statusStr := "RUNNING"
	if paused {
		statusStr = "PAUSED"
	}

	cfgPath := util.ResolvePaths().ConfigPath
	dbPath := util.ResolvePaths().DBPath
	logDir := util.ResolvePaths().LogDir

	var total, enabled, disabled int

	dueJobs, _ := db.GetDueJobs(time.Now())
	nextJob, err := db.GetNextJob()
	if err != nil {
		util.ErrorMessage(fmt.Sprint("Error fetching job:", err))
		return
	}

	jobs, err := db.GetAllJobs()
	if err != nil {
		util.ErrorMessage(fmt.Sprint("Could not fetch jobs:", err))
		return
	}

	for _, job := range jobs {
		total++
		if job.Enabled {
			enabled++
		} else {
			disabled++
		}
	}

	dueNow := len(dueJobs)

	fmt.Println("\nJobctl System Status")
	fmt.Println("---------------------")
	fmt.Printf("%-22s %s\n", "System Status:", statusStr)
	fmt.Printf("%-22s %d\n", "Jobs (total):", total)
	fmt.Printf("%-22s %d\n", "Jobs (enabled):", enabled)
	fmt.Printf("%-22s %d\n", "Jobs (disabled):", disabled)
	fmt.Printf("%-22s %d\n", "Jobs due now:", dueNow)

	if nextJob != nil {
		delta := nextJob.NextRun.Sub(util.AnchorTime())
		fmt.Printf("%-22s \"%s\" in %s\n", "Next job:", nextJob.Name, util.FormatDuration(delta))
	} else {
		fmt.Printf("%-22s %s\n", "Next job:", "None scheduled")
	}

	fmt.Println()
	fmt.Printf("%-22s %s\n", "Config path:", cfgPath)
	fmt.Printf("%-22s %s\n", "Database path:", dbPath)
	fmt.Printf("%-22s %s\n", "Log directory:", logDir)
	fmt.Printf("%-22s %s/%s\n\n", "OS:", runtime.GOOS, runtime.GOARCH)
}
