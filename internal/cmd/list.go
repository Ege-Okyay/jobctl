package cmd

import (
	"fmt"
	"os"
	"text/tabwriter"

	"github.com/Ege-Okyay/jobctl/internal/db"
	"github.com/Ege-Okyay/jobctl/internal/types"
	"github.com/Ege-Okyay/jobctl/internal/util"
)

var ListCommand types.Command

func init() {
	ListCommand = types.Command{
		Name:        "list",
		Description: "List all scheduled jobs",
		Usage:       "list",
		Flags:       []types.Flag{},
		Execute:     listHandler,
	}
}

func listHandler(args []string) {
	jobs, err := db.GetAllJobs()
	if err != nil {
		util.ErrorMessage(fmt.Sprint("Error fetching jobs:", err))
		return
	}

	w := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)

	fmt.Fprintln(w, "NAME\tINTERVAL(s)\tNEXT RUN\tENABLED")

	for _, j := range jobs {
		rel := "-"
		if j.Enabled {
			delta := j.NextRun.Sub(util.AnchorTime())
			rel = util.FormatDuration(delta)
		}

		fmt.Fprintf(w, "%s\t%d\t%s\t%v\n",
			j.Name,
			j.Interval,
			rel,
			j.Enabled,
		)
	}

	w.Flush()
}
