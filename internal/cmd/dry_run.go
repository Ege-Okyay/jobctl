package cmd

import (
	"fmt"
	"time"

	"github.com/Ege-Okyay/jobctl/internal/db"
	"github.com/Ege-Okyay/jobctl/internal/types"
	"github.com/Ege-Okyay/jobctl/internal/util"
)

var DryRunCommand types.Command

func init() {
	DryRunCommand = types.Command{
		Name:        "dry-run",
		Description: "Show which jobs would run in the next N seconds without executing them",
		Usage:       "dry-run <seconds>",
		Flags:       []types.Flag{},
		Execute:     dryRunHandler,
	}
}

// dryRunHandler simulates future job executions without running them.
func dryRunHandler(args []string) {
	if len(args) != 1 {
		util.ErrorMessage(fmt.Sprint("Usage: ", DryRunCommand.Usage))
		return
	}

	secs := util.AtoiOrDefault(args[0], 0)
	if secs <= 0 {
		util.ErrorMessage(fmt.Sprintf("Invalid --seconds %d; must be positive\n", secs))
		return
	}

	now := time.Now()
	cutoff := now.Add(time.Duration(secs) * time.Second)

	jobs, err := db.GetDueJobs(cutoff)
	if err != nil {
		util.ErrorMessage(fmt.Sprint("Error fetching due jobs:", err))
		return
	}

	if len(jobs) == 0 {
		util.SuccessMessage(fmt.Sprintf("No jobs due in the next %d seconds\n", secs))
		return
	}

	fmt.Printf("Jobs due in the next %d seconds:\n", secs)
	for _, j := range jobs {
		delta := j.NextRun.Sub(util.AnchorTime())
		fmt.Printf("\t- %s (in %s)\n", j.Name, util.FormatDuration(delta))
	}
}
