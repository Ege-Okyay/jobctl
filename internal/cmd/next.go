package cmd

import (
	"fmt"

	"github.com/Ege-Okyay/jobctl/internal/db"
	"github.com/Ege-Okyay/jobctl/internal/types"
	"github.com/Ege-Okyay/jobctl/internal/util"
)

var NextCommand types.Command

func init() {
	NextCommand = types.Command{
		Name:        "next",
		Description: "Show the next scheduled job to run",
		Usage:       "next",
		Flags:       []types.Flag{},
		Execute:     nextHandler,
	}
}

func nextHandler(args []string) {
	job, err := db.GetNextJob()
	if err != nil {
		util.ErrorMessage(fmt.Sprint("Error fetching job:", err))
		return
	}
	if job == nil {
		util.SuccessMessage("No enabled jobs found")
		return
	}

	delta := job.NextRun.Sub(util.AnchorTime())
	rel := util.FormatDuration(delta)

	fmt.Printf("Next: %q %s\n", job.Name, rel)
}
