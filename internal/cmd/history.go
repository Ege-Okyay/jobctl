package cmd

import (
	"fmt"
	"strings"

	"github.com/Ege-Okyay/jobctl/internal/db"
	"github.com/Ege-Okyay/jobctl/internal/types"
	"github.com/Ege-Okyay/jobctl/internal/util"
)

var HistoryCommand types.Command

func init() {
	HistoryCommand = types.Command{
		Name:        "history",
		Description: "Show past run records for a job",
		Usage:       "history <job_name> [--limit <n>]",
		Flags: []types.Flag{
			{Name: "--limit", Description: "Max records to show (default 10)", Required: false},
		},
		Execute: historyHandler,
	}
}

func historyHandler(args []string) {
	if len(args) < 1 {
		util.ErrorMessage(fmt.Sprint("Usage: ", HistoryCommand.Usage))
		return
	}

	name := args[0]

	flags := util.ParseFlags(args[1:])

	limit := util.AtoiOrDefault(flags["--limit"], 10)
	if limit <= 0 {
		util.ErrorMessage(fmt.Sprintf("Invalid --limit %d: must be > 0", limit))
		return
	}

	recs, err := db.GetRunHistory(name, limit)
	if err != nil {
		util.ErrorMessage(fmt.Sprintf("Error fetching history for %q: %v", name, err))
		return
	}

	if len(recs) == 0 {
		util.SuccessMessage(fmt.Sprintf("No run records found for %q", name))
		return
	}

	fmt.Printf("Last %d runs for %q:\n", len(recs), name)

	for i, r := range recs {
		status := "FAILED"
		if r.Success {
			status = "SUCCESS"
		}

		ts := r.Timestamp.Format("2006-01-02 15:04:05")

		fmt.Printf("%2d) %s  %s\n", i+1, ts, status)

		// Display command output if it exists.
		if out := strings.TrimSpace(r.Output); out != "" {
			fmt.Println("    ▶ Output:")

			for _, line := range strings.Split(out, "\n") {
				fmt.Printf("       %s\n", line)
			}
		}
	}
}
