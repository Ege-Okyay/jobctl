package cmd

import (
	"fmt"
	"os"
	"text/tabwriter"
	"time"

	"github.com/Ege-Okyay/jobctl/internal/db"
	"github.com/Ege-Okyay/jobctl/internal/types"
	"github.com/Ege-Okyay/jobctl/internal/util"
)

var InspectCommand types.Command

func init() {
	InspectCommand = types.Command{
		Name:        "inspect",
		Description: "Show a detailed info about one job",
		Usage:       "inspect <job_name>",
		Flags:       []types.Flag{},
		Execute:     inspectHandler,
	}
}

func inspectHandler(args []string) {
	if len(args) != 1 {
		util.ErrorMessage(fmt.Sprint("Usage: ", InspectCommand.Usage))
		return
	}

	name := args[0]

	job, err := db.GetJob(name)
	if err != nil {
		util.ErrorMessage(fmt.Sprintf("fetching job: %v", err))
		return
	}

	if job == nil {
		util.ErrorMessage(fmt.Sprintf("No job named %q", name))
		return
	}

	lastRec, err := db.GetLastRun(name)
	if err != nil {
		util.ErrorMessage(fmt.Sprintf("Error fetching last run: %v", err))
		return
	}

	delta := job.NextRun.Sub(util.AnchorTime())
	nextRunRel := util.FormatDuration(delta)
	nextRunAbs := job.NextRun.Format("2006-01-02 15:04:05")

	w := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)

	fmt.Fprintln(w, "FIELD\tVALUE")
	fmt.Fprintf(w, "Name:\t%s\n", job.Name)
	fmt.Fprintf(w, "Interval:\t%d seconds\n", job.Interval)
	fmt.Fprintf(w, "Command:\t%s\n", job.Command)
	fmt.Fprintf(w, "Retries:\t%d\n", job.Retries)
	fmt.Fprintf(w, "Timeout:\t%d seconds\n", job.Timeout)
	fmt.Fprintf(w, "Enabled:\t%v\n", job.Enabled)
	fmt.Fprintf(w, "Next run:\t%s (in %s)\n", nextRunAbs, nextRunRel)

	if lastRec != nil {
		abs := lastRec.Timestamp.Format("2006-01-02 15:04:05")
		rel := util.FormatDuration(time.Since(lastRec.Timestamp))

		status := "FAILED"
		if lastRec.Success {
			status = "SUCCESS"
		}

		fmt.Fprintf(w, "Last run:\t%s (%s ago) [%s]\n", abs, rel, status)
	} else {
		fmt.Fprintf(w, "Last run:\t%s\n", "never")
	}

	w.Flush()
}
