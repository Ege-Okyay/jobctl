package cmd

import (
	"fmt"

	"github.com/Ege-Okyay/jobctl/internal/logic"
	"github.com/Ege-Okyay/jobctl/internal/types"
	"github.com/Ege-Okyay/jobctl/internal/util"
)

var RunCommand types.Command

func init() {
	RunCommand = types.Command{
		Name:        "run",
		Description: "Run a job immediately (ignores schedule)",
		Usage:       "run <job_name> [--keep-remaining]",
		Flags: []types.Flag{
			{Name: "--keep-remaining", Description: "Do not reset the job's schedule after execution", Required: false},
		},
		Execute: runHandler,
	}
}

func runHandler(args []string) {
	if len(args) < 1 {
		util.ErrorMessage(fmt.Sprint("Usage: ", RunCommand.Usage))
		return
	}

	name := args[0]

	flags := util.ParseFlags(args[1:])
	keep := false
	if _, ok := flags["--keep-remaining"]; ok {
		keep = true
	}

	if err := logic.RunJob(name, keep); err != nil {
		util.ErrorMessage(fmt.Sprintf("Failed to run %q: %v", name, err))
	} else {
		util.SuccessMessage(fmt.Sprintf("Job %q exectued", name))
	}
}
